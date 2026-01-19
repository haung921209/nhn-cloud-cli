package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"

	"github.com/haung921209/nhn-cloud-cli/internal/cert"
	"github.com/spf13/cobra"
)

var connectMySQLCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a MySQL DB instance using mysql client",
	Long: `Connect to a MySQL DB instance using the local 'mysql' client.
Automatically configures SSL/TLS certificates if available in the managed store.

Prerequisites:
- 'mysql' client must be installed in PATH.
- Certificates imported via 'nhncloud config ca import'.
- Instance must be accessible (Public Access or VPN).

Example:
  nhncloud rds-mysql connect --db-instance-identifier <uuid> --username <user> --password <pass> --database <db>
  nhncloud rds-mysql connect ... -- -e "SELECT 1;"`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, err := getResolvedInstanceID(cmd, newMySQLClient())
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		database, _ := cmd.Flags().GetString("database")

		// Fetch Instance Details to get Host/Port
		client := newMySQLClient()
		inst, err := client.GetInstance(context.Background(), dbInstanceID)
		if err != nil {
			exitWithError("failed to get instance details", err)
		}

		host := inst.Network.DomainName
		if host == "" {
			host = inst.Network.FloatingIP
		}
		if host == "" {
			host = inst.Network.PublicIP
		}
		if host == "" {
			host = inst.Network.IPAddress
		}

		// Fallback: Try GetNetworkInfo if host is still missing
		if host == "" {
			netInfo, err := client.GetNetworkInfo(context.Background(), dbInstanceID)
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
					// If no EXTERNAL found, take any non-empty (fallback)
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

		port := fmt.Sprintf("%d", inst.DBPort)

		fmt.Printf("Connecting to %s (%s:%s)...\n", inst.DBInstanceName, host, port)

		authPlugin, _ := cmd.Flags().GetString("auth-plugin")

		// Setup Cert Helper
		helper, err := cert.NewHelper()
		if err != nil {
			exitWithError("failed to init cert helper", err)
		}

		// Generate Connection Command
		commandArgs, err := helper.GetConnectionCommand("rds-mysql", host, port, database, username, password, getRegion(), dbInstanceID, "", true, "")
		if err != nil {
			exitWithError("failed to generate connection command", err)
		}
		fmt.Printf("Executing: %v\n", commandArgs)

		// Check for SSL flags in generated command
		hasSSL := false
		for _, arg := range commandArgs {
			if len(arg) >= 8 && arg[:8] == "--ssl-ca" {
				hasSSL = true
				break
			}
		}

		// Enforce SSL for caching_sha2_password
		if authPlugin == "caching_sha2_password" && !hasSSL {
			fmt.Println()
			fmt.Println("Error: 'caching_sha2_password' requires SSL/TLS certificates, but none were found.")
			fmt.Println("Please import certificates using:")
			fmt.Printf("  nhncloud config ca import --service rds-mysql --region %s --file <ca.pem>\n", getRegion())
			fmt.Println("  (Optional) Import client cert/key for mutual TLS if needed.")
			fmt.Println()
			fmt.Println("To bypass (insecure), use: --auth-plugin mysql_native_password (if supported by server)")
			os.Exit(1)
		}

		// Append auth plugin instruction if not default or if needed?
		// Actually mysql client doesn't need --default-auth unless we want to force it.
		// But usually it negotiates.
		// However, if we want to BE explicit, we can add --default-auth=...
		if authPlugin != "" {
			commandArgs = append(commandArgs, "--default-auth="+authPlugin)
		}

		// Append extra args provided by user (e.g. -e "SELECT 1")
		commandArgs = append(commandArgs, args...)

		fmt.Printf("Executing: %v\n", commandArgs)

		// Exec
		c := exec.Command(commandArgs[0], commandArgs[1:]...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		// Check for missing mysql client
		if _, err := exec.LookPath("mysql"); err != nil && len(args) == 0 {
			// Fallback to Native REPL
			fmt.Println("Notice: 'mysql' client not found in PATH. Falling back to built-in native shell.")

			// Construct DSN
			// TODO: Add SSL support if we want parity with strict mode,
			// but for REPL fallback we'll try basic connection first or skip-verify.

			// We need to register the CA if we want strict SSL.
			// Currently simplified to skip-verify for fallback functionality to work.
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=skip-verify&allowNativePasswords=true",
				username, password, host, port, database)

			db, err := sql.Open("mysql", dsn)
			if err != nil {
				exitWithError("failed to open database connection", err)
			}
			defer db.Close()

			if err := db.Ping(); err != nil {
				exitWithError("failed to connect/ping database", err)
			}

			runNativeREPL(db, database, "mysql")
			return
		}

		if err := c.Run(); err != nil {
			exitWithError("connection failed", err)
		}
	},
}

func init() {
	rdsMySQLCmd.AddCommand(connectMySQLCmd)
	connectMySQLCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	connectMySQLCmd.Flags().String("username", "", "Database username")
	connectMySQLCmd.Flags().String("password", "", "Database password")
	connectMySQLCmd.Flags().String("database", "", "Database name")
	connectMySQLCmd.Flags().String("auth-plugin", "caching_sha2_password", "Authentication plugin (caching_sha2_password, mysql_native_password)")
}
