package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/haung921209/nhn-cloud-cli/internal/cert"
	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

func init() {
	rdsPostgreSQLCmd.AddCommand(connectPostgreSQLCmd)

	connectPostgreSQLCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	connectPostgreSQLCmd.Flags().String("username", "", "Database username (Required)")
	connectPostgreSQLCmd.Flags().String("password", "", "Database password (Required)")
	connectPostgreSQLCmd.Flags().String("database", "postgres", "Database name (Default: postgres)")
	connectPostgreSQLCmd.Flags().String("region", "kr1", "Region (Default: kr1)")
	connectPostgreSQLCmd.Flags().String("ca-path", "", "Explicit path to CA certificate (Optional)")

	// Native execution flag
	connectPostgreSQLCmd.Flags().StringP("execute", "e", "", "Execute query and exit (Uses built-in driver, no external dependency)")

	connectPostgreSQLCmd.MarkFlagRequired("db-instance-identifier")
	connectPostgreSQLCmd.MarkFlagRequired("username")
	connectPostgreSQLCmd.MarkFlagRequired("password")
}

var connectPostgreSQLCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a PostgreSQL DB instance",
	Long: `Connects to a PostgreSQL DB instance.
Modes:
1. Interactive: Launches 'psql' client (Requires local installation).
2. Execute: Runs a query using the built-in Go driver (No external dependency).

Example:
  Interactive: nhncloud rds-postgresql connect --db-instance-identifier mydb -u user -p pass
  Execute:     nhncloud rds-postgresql connect --db-instance-identifier mydb -u user -p pass -e "SELECT version();"`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		dbInstanceIdentifier, _ := cmd.Flags().GetString("db-instance-identifier")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		database, _ := cmd.Flags().GetString("database")
		region, _ := cmd.Flags().GetString("region")
		caPath, _ := cmd.Flags().GetString("ca-path")
		query, _ := cmd.Flags().GetString("execute")

		// Resolve ID
		dbInstanceID, err := resolvePostgreSQLInstanceIdentifier(client, dbInstanceIdentifier)
		if err != nil {
			exitWithError("failed to resolve instance identifier", err)
		}

		ctx := context.Background()
		inst, err := client.GetInstance(ctx, dbInstanceID)
		if err != nil {
			exitWithError("failed to get instance details", err)
		}

		// Host Detection Logic (Strict)
		host := ""
		// 1. Try standard network info from GetInstance
		if inst.Network.DomainName != "" {
			host = inst.Network.DomainName
		} else if inst.Network.IPAddress != "" {
			host = inst.Network.IPAddress
		}

		// 2. Fallback: GetNetworkInfo (Checking for EXTERNAL endpoint)
		if host == "" {
			netInfo, err := client.GetNetworkInfo(ctx, dbInstanceID)
			if err == nil && netInfo != nil {
				for _, ep := range netInfo.EndPoints {
					if ep.EndPointType == "EXTERNAL" {
						if ep.Domain != "" {
							host = ep.Domain
						} else if ep.IPAddress != "" {
							host = ep.IPAddress
						}
						break
					}
					// Fallback to internal if needed?
					// Generally if user connects via CLI they want public or they have VPN.
					// We'll prioritize EXTERNAL.
				}
			}
		}

		if host == "" {
			exitWithError(fmt.Sprintf("Unable to determine public host for instance '%s'. Ensure Public Access is enabled.", inst.DBInstanceName), nil)
		}

		fmt.Printf("Connecting to %s (%s:%d)...\n", inst.DBInstanceName, host, inst.DBPort)

		// Initialize Helper
		helper, err := cert.NewHelper()
		if err != nil {
			exitWithError("failed to initialize certificate helper", err)
		}

		// Check if we need Native Connection (either -e flag OR psql missing)
		useNative := false
		_, err = exec.LookPath("psql")
		if query != "" {
			useNative = true
		} else if err != nil {
			// psql not found, fallback to native REPL
			useNative = true
			fmt.Println("Notice: 'psql' client not found in PATH. Falling back to built-in native shell.")
		}

		if useNative {
			// Get CA Path for SSL
			caPathReal, err := helper.GetCertificateForDatabase("rds-postgresql", region, dbInstanceID, inst.DBVersion, true, caPath, "CA")
			if err != nil {
				fmt.Printf("Warning: Could not locate CA cert: %v. Connection may fail if SSL required.\n", err)
			}

			// Construct Connection String (lib/pq format)
			connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				host, inst.DBPort, username, password, database)

			if caPathReal != "" {
				// We removed sslrootcert because we use sslmode=disable, but if we revert to prefer/require:
				// connStr += fmt.Sprintf(" sslrootcert=%s", caPathReal)
			}

			db, err := sql.Open("postgres", connStr)
			if err != nil {
				exitWithError("failed to open database connection", err)
			}
			defer db.Close()

			err = db.Ping()
			if err != nil {
				exitWithError("failed to connect/ping database", err)
			}

			if query != "" {
				// Single shot
				executeNativeQuery(db, query)
			} else {
				// REPL
				runNativeREPL(db, database, database)
			}
			return
		}

		// Interactive Mode (using psql)
		cmdArgs, err := helper.GetConnectionCommand(
			"rds-postgresql",
			host,
			strconv.Itoa(inst.DBPort),
			database,
			username,
			password,
			region,
			dbInstanceID,
			inst.DBVersion,
			true,   // Auto-find certs
			caPath, // Explicit CA override
		)
		if err != nil {
			exitWithError("failed to generate connection command", err)
		}

		// Append extra args
		if len(args) > 0 {
			cmdArgs = append(cmdArgs, args...)
		}

		fmt.Printf("Executing: [%s]\n", cmdArgs)

		// Execute
		c := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			fmt.Printf("Command exited: %v\n", err)
			// This fallback is only triggered if exec fails, but usually we check LookPath earlier.
			// However, if psql fails for other reasons (like signal), we just exit.
			os.Exit(1)
		}
	},
}
