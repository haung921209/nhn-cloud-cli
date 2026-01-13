package cmd

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/haung921209/nhn-cloud-cli/internal/cert"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/rds/mariadb"
	"github.com/spf13/cobra"
)

// ============================================================================
// MariaDB Certificate Management Commands
// ============================================================================

var mariadbCertCmd = &cobra.Command{
	Use:   "cert",
	Short: "Manage SSL certificates for MariaDB",
	Long:  "Manage SSL certificates for MariaDB database connections",
}

var mariadbCertListCmd = &cobra.Command{
	Use:   "list",
	Short: "List MariaDB SSL certificates",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("instance-id")

		store, err := cert.NewCertificateStore()
		if err != nil {
			exitWithError("failed to create certificate store", err)
		}

		certificates, err := store.ListCertificates("mariadb", getRegion(), instanceID)
		if err != nil {
			exitWithError("failed to list certificates", err)
		}

		if len(certificates) == 0 {
			fmt.Println("No MariaDB SSL certificates found.")
			if instanceID != "" {
				fmt.Printf("\nTo import certificates for instance %s:\n", instanceID)
				fmt.Printf("  nhncloud rds-mariadb cert import --instance-id %s --ca-cert ca.pem --client-cert client.pem --client-key client.key\n", instanceID)
			}
			return
		}

		if output == "json" {
			printMariaDBJSON(certificates)
			return
		}

		fmt.Printf("Found %d MariaDB SSL certificate(s):\n\n", len(certificates))
		for _, c := range certificates {
			fmt.Printf("Certificate ID: %s\n", c.ID)
			fmt.Printf("  Service:     %s\n", c.ServiceType)
			fmt.Printf("  Region:      %s\n", c.Region)
			if c.InstanceID != "" {
				fmt.Printf("  Instance ID: %s\n", c.InstanceID)
			}
			if c.Version != "" {
				fmt.Printf("  Version:     %s\n", c.Version)
			}
			fmt.Printf("  Source:      %s\n", c.Source)
			fmt.Printf("  Stored:      %s\n", c.StoredAt.Format("2006-01-02 15:04:05"))
			if c.Description != "" {
				fmt.Printf("  Description: %s\n", c.Description)
			}
			fmt.Println()
		}
	},
}

var mariadbCertImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import MariaDB SSL certificate",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("instance-id")
		version, _ := cmd.Flags().GetString("version")
		caCertPath, _ := cmd.Flags().GetString("ca-cert")
		clientCertPath, _ := cmd.Flags().GetString("client-cert")
		clientKeyPath, _ := cmd.Flags().GetString("client-key")
		description, _ := cmd.Flags().GetString("description")

		if instanceID == "" {
			exitWithError("--instance-id is required", nil)
		}
		if caCertPath == "" {
			exitWithError("--ca-cert is required", nil)
		}

		fmt.Printf("Importing MariaDB SSL certificates for instance: %s\n", instanceID)

		store, err := cert.NewCertificateStore()
		if err != nil {
			exitWithError("failed to create certificate store", err)
		}

		region := getRegion()

		// Import CA certificate
		fmt.Printf("  Importing CA certificate...\n")
		caCertData, err := os.ReadFile(caCertPath)
		if err != nil {
			exitWithError("failed to read CA certificate file", err)
		}

		caCertReq := &cert.CertificateRequest{
			ServiceType: "mariadb",
			Region:      region,
			InstanceID:  instanceID,
			Version:     version,
			Source:      "manual",
			Description: fmt.Sprintf("CA certificate for MariaDB instance %s", instanceID),
			CertData:    caCertData,
		}
		if description != "" {
			caCertReq.Description = description + " (CA)"
		}

		caCertInfo, err := store.StoreCertificate(caCertReq)
		if err != nil {
			exitWithError("failed to import CA certificate", err)
		}
		fmt.Printf("  CA certificate imported with ID: %s\n", caCertInfo.ID)

		// Import client certificate if provided
		if clientCertPath != "" {
			fmt.Printf("  Importing client certificate...\n")
			clientCertData, err := os.ReadFile(clientCertPath)
			if err != nil {
				exitWithError("failed to read client certificate file", err)
			}

			clientCertReq := &cert.CertificateRequest{
				ServiceType: "mariadb",
				Region:      region,
				InstanceID:  instanceID,
				Version:     version,
				Source:      "manual",
				Description: fmt.Sprintf("Client certificate for MariaDB instance %s", instanceID),
				CertData:    clientCertData,
			}
			if description != "" {
				clientCertReq.Description = description + " (Client)"
			}

			clientCertInfo, err := store.StoreCertificate(clientCertReq)
			if err != nil {
				exitWithError("failed to import client certificate", err)
			}
			fmt.Printf("  Client certificate imported with ID: %s\n", clientCertInfo.ID)
		}

		// Import client key if provided
		if clientKeyPath != "" {
			fmt.Printf("  Importing client key...\n")
			clientKeyData, err := os.ReadFile(clientKeyPath)
			if err != nil {
				exitWithError("failed to read client key file", err)
			}

			clientKeyReq := &cert.CertificateRequest{
				ServiceType: "mariadb",
				Region:      region,
				InstanceID:  instanceID,
				Version:     version,
				Source:      "manual",
				Description: fmt.Sprintf("Client key for MariaDB instance %s", instanceID),
				CertData:    clientKeyData,
			}
			if description != "" {
				clientKeyReq.Description = description + " (Key)"
			}

			clientKeyInfo, err := store.StoreCertificate(clientKeyReq)
			if err != nil {
				exitWithError("failed to import client key", err)
			}
			fmt.Printf("  Client key imported with ID: %s\n", clientKeyInfo.ID)
		}

		fmt.Printf("\nSuccessfully imported SSL certificates for MariaDB instance %s\n", instanceID)
		fmt.Printf("\nYou can now use these certificates with:\n")
		fmt.Printf("  nhncloud rds-mariadb connect %s --user <username> --password <password>\n", instanceID)
	},
}

var mariadbCertDeleteCmd = &cobra.Command{
	Use:   "delete [certificate-id]",
	Short: "Delete MariaDB SSL certificate",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		certID := args[0]

		store, err := cert.NewCertificateStore()
		if err != nil {
			exitWithError("failed to create certificate store", err)
		}

		err = store.RemoveCertificate(certID)
		if err != nil {
			exitWithError("failed to delete certificate", err)
		}

		fmt.Printf("Certificate %s deleted successfully\n", certID)
	},
}

// ============================================================================
// MariaDB Connection Commands
// ============================================================================

var mariadbConnectCmd = &cobra.Command{
	Use:   "connect [instance-id]",
	Short: "Connect to MariaDB instance",
	Long:  "Connect to a MariaDB instance using SSL certificates and credentials",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		instanceID := args[0]

		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		database, _ := cmd.Flags().GetString("database")
		disableSSL, _ := cmd.Flags().GetBool("disable-ssl")
		timeout, _ := cmd.Flags().GetInt("timeout")

		if user == "" || password == "" {
			exitWithError("both --user and --password are required", nil)
		}

		fmt.Printf("Connecting to MariaDB instance: %s\n", instanceID)

		// Get instance details
		client := newMariaDBClient()
		result, err := client.GetInstance(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to get instance details", err)
		}

		// Get network info for endpoint
		networkInfo, err := client.GetNetworkInfo(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to get network info", err)
		}

		var endpoint string
		port := 3306

		// Find the endpoint from network info
		if networkInfo != nil && len(networkInfo.EndPoints) > 0 {
			for _, ep := range networkInfo.EndPoints {
				if ep.Domain != "" {
					endpoint = ep.Domain
					break
				}
				if ep.IPAddress != "" {
					endpoint = ep.IPAddress
				}
			}
		}

		if endpoint == "" {
			exitWithError("no network endpoint found for instance", nil)
		}

		if result.DBPort > 0 {
			port = result.DBPort
		}

		fmt.Printf("  Endpoint: %s:%d\n", endpoint, port)
		fmt.Printf("  User: %s\n", user)
		fmt.Printf("  Database: %s\n", database)

		// Prepare TLS config
		var tlsConfigName string
		if !disableSSL {
			fmt.Printf("  SSL: Enabled\n")

			// Try to find certificates
			store, err := cert.NewCertificateStore()
			if err == nil {
				certs, err := store.ListCertificates("mariadb", getRegion(), instanceID)
				if err == nil && len(certs) > 0 {
					fmt.Printf("  SSL Certificates: Found %d certificate(s)\n", len(certs))
				}
			}

			// Register TLS config
			tlsConfig := &tls.Config{
				InsecureSkipVerify: true,
			}
			mysql.RegisterTLSConfig("mariadb-custom", tlsConfig)
			tlsConfigName = "mariadb-custom"
		} else {
			fmt.Printf("  SSL: DISABLED (not recommended for production)\n")
		}

		// Build MySQL DSN (MariaDB uses MySQL protocol)
		cfg := mysql.Config{
			User:                 user,
			Passwd:               password,
			Net:                  "tcp",
			Addr:                 fmt.Sprintf("%s:%d", endpoint, port),
			DBName:               database,
			Timeout:              time.Duration(timeout) * time.Second,
			ReadTimeout:          time.Duration(timeout) * time.Second,
			WriteTimeout:         time.Duration(timeout) * time.Second,
			AllowNativePasswords: true,
			ParseTime:            true,
		}

		if tlsConfigName != "" {
			cfg.TLSConfig = tlsConfigName
		}

		dsn := cfg.FormatDSN()

		// Open connection
		fmt.Printf("\nAttempting connection...\n")
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			exitWithError("failed to open connection", err)
		}
		defer db.Close()

		// Test connection
		err = db.Ping()
		if err != nil {
			exitWithError("failed to connect to MariaDB instance", err)
		}

		fmt.Printf("Successfully connected to MariaDB instance!\n")

		// Execute test query
		var version, currentUser, currentDB string
		err = db.QueryRow("SELECT VERSION(), USER(), DATABASE()").Scan(&version, &currentUser, &currentDB)
		if err != nil {
			fmt.Printf("Warning: failed to execute test query: %v\n", err)
		} else {
			fmt.Printf("\nConnection Test:\n")
			fmt.Printf("  MariaDB Version: %s\n", version)
			fmt.Printf("  Connected User: %s\n", currentUser)
			fmt.Printf("  Current Database: %s\n", currentDB)
		}

		fmt.Printf("\nConnection successful! Use 'nhncloud rds-mariadb query' to execute SQL.\n")
	},
}

// ============================================================================
// MariaDB Query Commands
// ============================================================================

var mariadbQueryCmd = &cobra.Command{
	Use:   "query [instance-id] [sql]",
	Short: "Execute SQL query on MariaDB instance",
	Long:  "Execute a SQL query on a MariaDB instance and return formatted results",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		instanceID := args[0]
		sqlQuery := args[1]

		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		database, _ := cmd.Flags().GetString("database")
		disableSSL, _ := cmd.Flags().GetBool("disable-ssl")
		timeout, _ := cmd.Flags().GetInt("timeout")
		format, _ := cmd.Flags().GetString("format")

		if user == "" || password == "" {
			exitWithError("both --user and --password are required", nil)
		}

		// Get instance details
		client := newMariaDBClient()
		result, err := client.GetInstance(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to get instance details", err)
		}

		// Get network info for endpoint
		networkInfo, err := client.GetNetworkInfo(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to get network info", err)
		}

		var endpoint string
		port := 3306

		if networkInfo != nil && len(networkInfo.EndPoints) > 0 {
			for _, ep := range networkInfo.EndPoints {
				if ep.Domain != "" {
					endpoint = ep.Domain
					break
				}
				if ep.IPAddress != "" {
					endpoint = ep.IPAddress
				}
			}
		}

		if endpoint == "" {
			exitWithError("no network endpoint found for instance", nil)
		}

		if result.DBPort > 0 {
			port = result.DBPort
		}

		// Prepare TLS config
		var tlsConfigName string
		if !disableSSL {
			tlsConfig := &tls.Config{
				InsecureSkipVerify: true,
			}
			mysql.RegisterTLSConfig("mariadb-query", tlsConfig)
			tlsConfigName = "mariadb-query"
		}

		// Build MySQL DSN
		cfg := mysql.Config{
			User:                 user,
			Passwd:               password,
			Net:                  "tcp",
			Addr:                 fmt.Sprintf("%s:%d", endpoint, port),
			DBName:               database,
			Timeout:              time.Duration(timeout) * time.Second,
			ReadTimeout:          time.Duration(timeout) * time.Second,
			WriteTimeout:         time.Duration(timeout) * time.Second,
			AllowNativePasswords: true,
			ParseTime:            true,
		}

		if tlsConfigName != "" {
			cfg.TLSConfig = tlsConfigName
		}

		dsn := cfg.FormatDSN()

		// Open connection
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			exitWithError("failed to open connection", err)
		}
		defer db.Close()

		// Test connection
		err = db.Ping()
		if err != nil {
			exitWithError("failed to connect to MariaDB instance", err)
		}

		// Execute query
		startTime := time.Now()

		lowerQuery := strings.ToLower(strings.TrimSpace(sqlQuery))
		isSelectQuery := strings.HasPrefix(lowerQuery, "select") ||
			strings.HasPrefix(lowerQuery, "show") ||
			strings.HasPrefix(lowerQuery, "describe") ||
			strings.HasPrefix(lowerQuery, "desc") ||
			strings.HasPrefix(lowerQuery, "explain")

		if isSelectQuery {
			rows, err := db.Query(sqlQuery)
			if err != nil {
				exitWithError("query failed", err)
			}
			defer rows.Close()

			columns, err := rows.Columns()
			if err != nil {
				exitWithError("failed to get columns", err)
			}

			var allRows [][]string
			for rows.Next() {
				values := make([]interface{}, len(columns))
				valuePtrs := make([]interface{}, len(columns))
				for i := range values {
					valuePtrs[i] = &values[i]
				}

				if err := rows.Scan(valuePtrs...); err != nil {
					exitWithError("failed to scan row", err)
				}

				row := make([]string, len(columns))
				for i, val := range values {
					if val == nil {
						row[i] = "NULL"
					} else {
						switch v := val.(type) {
						case []byte:
							row[i] = string(v)
						case string:
							row[i] = v
						default:
							row[i] = fmt.Sprintf("%v", v)
						}
					}
				}
				allRows = append(allRows, row)
			}

			duration := time.Since(startTime)

			switch format {
			case "json":
				printMariaDBQueryJSON(columns, allRows)
			case "csv":
				printMariaDBQueryCSV(columns, allRows)
			default:
				printMariaDBQueryTable(columns, allRows, duration)
			}
		} else {
			// Execute non-SELECT query
			execResult, err := db.Exec(sqlQuery)
			if err != nil {
				exitWithError("command failed", err)
			}

			duration := time.Since(startTime)

			rowsAffected, _ := execResult.RowsAffected()
			fmt.Printf("Query executed successfully\n")
			fmt.Printf("  Rows affected: %d\n", rowsAffected)
			fmt.Printf("  Duration: %.2fms\n", float64(duration.Nanoseconds())/1000000)
		}
	},
}

// Helper functions for MariaDB query output
func printMariaDBQueryTable(columns []string, rows [][]string, duration time.Duration) {
	if len(rows) == 0 {
		fmt.Printf("(0 rows) (%.2fms)\n", float64(duration.Nanoseconds())/1000000)
		return
	}

	maxWidth := 50
	colWidths := make([]int, len(columns))

	for i, col := range columns {
		colWidths[i] = min(len(col), maxWidth)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				cellLen := min(len(cell), maxWidth)
				if cellLen > colWidths[i] {
					colWidths[i] = cellLen
				}
			}
		}
	}

	fmt.Print("+")
	for _, width := range colWidths {
		fmt.Print(strings.Repeat("-", width+2) + "+")
	}
	fmt.Println()

	fmt.Print("|")
	for i, col := range columns {
		fmt.Printf(" %-*s |", colWidths[i], col)
	}
	fmt.Println()

	fmt.Print("+")
	for _, width := range colWidths {
		fmt.Print(strings.Repeat("-", width+2) + "+")
	}
	fmt.Println()

	for _, row := range rows {
		fmt.Print("|")
		for i, cell := range row {
			if i < len(colWidths) {
				displayCell := cell
				if len(cell) > maxWidth {
					displayCell = cell[:maxWidth-3] + "..."
				}
				fmt.Printf(" %-*s |", colWidths[i], displayCell)
			}
		}
		fmt.Println()
	}

	fmt.Print("+")
	for _, width := range colWidths {
		fmt.Print(strings.Repeat("-", width+2) + "+")
	}
	fmt.Println()

	fmt.Printf("(%d rows) (%.2fms)\n", len(rows), float64(duration.Nanoseconds())/1000000)
}

func printMariaDBQueryJSON(columns []string, rows [][]string) {
	result := make([]map[string]interface{}, len(rows))

	for i, row := range rows {
		rowMap := make(map[string]interface{})
		for j, col := range columns {
			if j < len(row) {
				rowMap[col] = row[j]
			}
		}
		result[i] = rowMap
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		exitWithError("failed to marshal JSON", err)
	}

	fmt.Printf("%s\n", jsonData)
}

func printMariaDBQueryCSV(columns []string, rows [][]string) {
	fmt.Println(strings.Join(columns, ","))

	for _, row := range rows {
		escapedRow := make([]string, len(row))
		for i, cell := range row {
			if strings.Contains(cell, ",") || strings.Contains(cell, "\"") {
				escapedRow[i] = fmt.Sprintf("\"%s\"", strings.ReplaceAll(cell, "\"", "\"\""))
			} else {
				escapedRow[i] = cell
			}
		}
		fmt.Println(strings.Join(escapedRow, ","))
	}
}

func printMariaDBJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

// ============================================================================
// MariaDB Instance Group Commands
// ============================================================================

var mariadbInstanceGroupCmd = &cobra.Command{
	Use:   "instance-groups",
	Short: "Manage MariaDB instance groups",
}

var mariadbListInstanceGroupsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all MariaDB instance groups",
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.ListInstanceGroups(context.Background())
		if err != nil {
			exitWithError("failed to list instance groups", err)
		}
		printMariaDBInstanceGroups(result)
	},
}

var mariadbGetInstanceGroupCmd = &cobra.Command{
	Use:   "get [group-id]",
	Short: "Get details of a MariaDB instance group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newMariaDBClient()
		result, err := client.GetInstanceGroup(context.Background(), args[0])
		if err != nil {
			exitWithError("failed to get instance group", err)
		}
		printMariaDBInstanceGroupDetail(result)
	},
}

func printMariaDBInstanceGroups(result *mariadb.ListInstanceGroupsOutput) {
	if output == "json" {
		printMariaDBJSON(result)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tREPLICATION TYPE\tCREATED")
	for _, g := range result.DBInstanceGroups {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			g.DBInstanceGroupID, g.ReplicationType, g.CreatedYmdt)
	}
	w.Flush()
}

func printMariaDBInstanceGroupDetail(result *mariadb.InstanceGroupOutput) {
	if output == "json" {
		printMariaDBJSON(result)
		return
	}

	fmt.Printf("ID:               %s\n", result.DBInstanceGroupID)
	fmt.Printf("Replication Type: %s\n", result.ReplicationType)
	fmt.Printf("Created:          %s\n", result.CreatedYmdt)
	fmt.Printf("Updated:          %s\n", result.UpdatedYmdt)
	fmt.Println("\nInstances:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  ID\tTYPE\tSTATUS")
	for _, inst := range result.DBInstances {
		fmt.Fprintf(w, "  %s\t%s\t%s\n",
			inst.DBInstanceID, inst.DBInstanceType, inst.DBInstanceStatus)
	}
	w.Flush()
}

func init() {
	// Certificate commands
	rdsMariaDBCmd.AddCommand(mariadbCertCmd)

	// Instance Group commands
	rdsMariaDBCmd.AddCommand(mariadbInstanceGroupCmd)
	mariadbInstanceGroupCmd.AddCommand(mariadbListInstanceGroupsCmd)
	mariadbInstanceGroupCmd.AddCommand(mariadbGetInstanceGroupCmd)
	mariadbCertCmd.AddCommand(mariadbCertListCmd)
	mariadbCertCmd.AddCommand(mariadbCertImportCmd)
	mariadbCertCmd.AddCommand(mariadbCertDeleteCmd)

	mariadbCertListCmd.Flags().String("instance-id", "", "Filter by database instance ID")

	mariadbCertImportCmd.Flags().String("instance-id", "", "Database instance ID (required)")
	mariadbCertImportCmd.Flags().String("version", "", "Certificate version")
	mariadbCertImportCmd.Flags().String("ca-cert", "", "CA certificate file path (required)")
	mariadbCertImportCmd.Flags().String("client-cert", "", "Client certificate file path")
	mariadbCertImportCmd.Flags().String("client-key", "", "Client key file path")
	mariadbCertImportCmd.Flags().String("description", "", "Certificate description")

	// Connect command
	rdsMariaDBCmd.AddCommand(mariadbConnectCmd)
	mariadbConnectCmd.Flags().String("user", "", "Database username (required)")
	mariadbConnectCmd.Flags().String("password", "", "Database password (required)")
	mariadbConnectCmd.Flags().String("database", "", "Database/schema name to connect to")
	mariadbConnectCmd.Flags().Bool("disable-ssl", false, "Disable SSL connection (not recommended)")
	mariadbConnectCmd.Flags().Int("timeout", 30, "Connection timeout in seconds")

	// Query command
	rdsMariaDBCmd.AddCommand(mariadbQueryCmd)
	mariadbQueryCmd.Flags().String("user", "", "Database username (required)")
	mariadbQueryCmd.Flags().String("password", "", "Database password (required)")
	mariadbQueryCmd.Flags().String("database", "", "Database/schema name to connect to")
	mariadbQueryCmd.Flags().Bool("disable-ssl", false, "Disable SSL connection (not recommended)")
	mariadbQueryCmd.Flags().Int("timeout", 30, "Connection timeout in seconds")
	mariadbQueryCmd.Flags().String("format", "table", "Output format: table, json, csv")
}
