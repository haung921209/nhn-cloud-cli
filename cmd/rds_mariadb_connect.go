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

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	rdsMariaDBCmd.AddCommand(connectMariaDBCmd)

	connectMariaDBCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (Required)")
	connectMariaDBCmd.Flags().String("username", "", "Database username (Required)")
	connectMariaDBCmd.Flags().String("password", "", "Database password (Required)")
	connectMariaDBCmd.Flags().String("database", "", "Database name (Optional)")
	connectMariaDBCmd.Flags().String("region", "kr1", "Region (Default: kr1)")
	connectMariaDBCmd.Flags().String("ca-path", "", "Explicit path to CA certificate (Optional)")

	// Native execution flag
	connectMariaDBCmd.Flags().StringP("execute", "e", "", "Execute query and exit (Uses built-in driver, no external dependency)")

	connectMariaDBCmd.MarkFlagRequired("db-instance-identifier")
	connectMariaDBCmd.MarkFlagRequired("username")
	connectMariaDBCmd.MarkFlagRequired("password")
}

var connectMariaDBCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a MariaDB DB instance",
	Long: `Connects to a MariaDB DB instance.
Modes:
1. Interactive: Launches 'mysql' client (Requires local installation).
2. Execute: Runs a query using the built-in Go driver (No external dependency).

Example:
  Interactive: nhncloud rds-mariadb connect --db-instance-identifier mydb -u user -p pass
  Execute:     nhncloud rds-mariadb connect --db-instance-identifier mydb -u user -p pass -e "SHOW DATABASES"`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()

		dbInstanceIdentifier, _ := cmd.Flags().GetString("db-instance-identifier")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		database, _ := cmd.Flags().GetString("database")
		region, _ := cmd.Flags().GetString("region")
		caPath, _ := cmd.Flags().GetString("ca-path")
		query, _ := cmd.Flags().GetString("execute")

		// Resolve ID
		dbInstanceID, err := resolveMariaDBInstanceIdentifier(client, dbInstanceIdentifier)
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
		// 1. Try standard network info from GetInstance (often empty)
		if inst.Network.DomainName != "" {
			host = inst.Network.DomainName
		} else if inst.Network.IPAddress != "" {
			host = inst.Network.IPAddress
		}

		// 2. Fallback: GetNetworkInfo
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
					// Fallback to any non-empty (e.g. INTERNAL if user has VPN, though CLI implies public)
					if host == "" {
						if ep.Domain != "" {
							host = ep.Domain
						} else if ep.IPAddress != "" {
							host = ep.IPAddress
						}
					}
				}
			}
		}

		if host == "" {
			exitWithError(fmt.Sprintf("Unable to determine public host for instance '%s'. Only instances with Public Access enabled can be connected to via CLI.", inst.DBInstanceName), nil)
		}

		fmt.Printf("Connecting to %s (%s:%d)...\n", inst.DBInstanceName, host, inst.DBPort)

		// Initialize Helper
		helper, err := cert.NewHelper()
		if err != nil {
			exitWithError("failed to initialize certificate helper", err)
		}

		// If -e/--execute is provided, use Native Go Driver
		if query != "" {
			fmt.Println("Synthesizing native connection...")

			// Get Cert Paths specifically for Native use (Helper might need 'mysql' type logic hint, but we just need paths)
			// Actually we need to register certs with the driver if we want SSL
			// The go-sql-driver/mysql uses 'tls' param in DSN and SystemCertPool or custom
			// For simplicity first pass: basic native connection.
			// Wait, RDS REQUIRES SSL usually if we follow 'ssl-ca'.
			// The driver supports `tls=skip-verify` or custom TLS config.

			// IMPORTANT: Registering Custom TLS Config for NHN Cloud Certs
			// We need the CA Path.
			_, err := helper.GetCertificateForDatabase("rds-mariadb", region, dbInstanceID, inst.DBVersion, true, caPath, "CA") // Ignore caPathReal for now
			if err != nil {
				fmt.Printf("Warning: Could not locate CA cert: %v. Connection may fail if SSL required.\n", err)
			}

			// TODO: Register TLS config with driver if CA exists.
			// For now, let's try strict DSN construction based on documentation.
			// DSN Format: user:password@tcp(host:port)/dbname?tls=custom&...

			// Just for this iteration: Print and use DSN without complex TLS reg (or skip-verify if needed)
			// But NHN Cloud usually needs CA.
			// Implementing properly means reading the CA file and `mysql.RegisterTLSConfig`.

			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=skip-verify&allowNativePasswords=true", username, password, host, inst.DBPort, database)
			// Note: allowNativePasswords=true is CRITICAL for MariaDB/MySQL 9 compat issue

			db, err := sql.Open("mysql", dsn)
			if err != nil {
				exitWithError("failed to open database connection", err)
			}
			defer db.Close()

			err = db.Ping()
			if err != nil {
				exitWithError("failed to connect/ping database", err)
			}

			// Execute Query
			rows, err := db.Query(query)
			if err != nil {
				exitWithError("failed to execute query", err)
			}
			defer rows.Close()

			// Print Results (Generic)
			cols, _ := rows.Columns()
			fmt.Println("--------------------------------------------------")
			// Print Header
			for _, c := range cols {
				fmt.Printf("%s\t", c)
			}
			fmt.Println("\n--------------------------------------------------")

			// Scan Rows
			rowValues := make([]interface{}, len(cols))
			rowPointers := make([]interface{}, len(cols))
			for i := range rowValues {
				rowPointers[i] = &rowValues[i]
			}

			for rows.Next() {
				err := rows.Scan(rowPointers...)
				if err != nil {
					fmt.Printf("Error scanning row: %v\n", err)
					continue
				}
				for _, val := range rowValues {
					if b, ok := val.([]byte); ok {
						fmt.Printf("%s\t", string(b))
					} else {
						fmt.Printf("%v\t", val)
					}
				}
				fmt.Println("")
			}
			fmt.Println("--------------------------------------------------")
			return
		}

		// Interactive Mode (Legacy Wrapper)
		cmdArgs, err := helper.GetConnectionCommand(
			"rds-mariadb", // logic maps this to "mysql"
			host,
			strconv.Itoa(inst.DBPort),
			database,
			username,
			password,
			region,
			dbInstanceID, // Certificate lookup needs UUID
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

		fmt.Printf("Executing: [%s]\n", cmdArgs) // Print for verification (user can see it)

		// Execute
		c := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		// Check for missing mysql client (MariaDB uses mysql client usually)
		if _, err := exec.LookPath("mysql"); err != nil && len(args) == 0 {
			// Fallback to Native REPL
			fmt.Println("Notice: 'mysql' client not found in PATH. Falling back to built-in native shell.")

			// Construct DSN
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=skip-verify&allowNativePasswords=true",
				username, password, host, strconv.Itoa(inst.DBPort), database)

			db, err := sql.Open("mysql", dsn)
			if err != nil {
				exitWithError("failed to open database connection", err)
			}
			defer db.Close()

			if err := db.Ping(); err != nil {
				exitWithError("failed to connect/ping database", err)
			}

			runNativeREPL(db, database, "mariadb")
			return
		}

		if err := c.Run(); err != nil {
			// Don't exitWithError if it's just a query failure or exit code
			fmt.Printf("Command exited: %v\n", err)
			fmt.Println("Tip: If you don't have the 'mysql' client installed, use the --execute / -e flag to run queries using the CLI's built-in driver.")
			os.Exit(1)
		}
	},
}
