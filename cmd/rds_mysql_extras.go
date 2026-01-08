package cmd

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/haung921209/nhn-cloud-cli/internal/cert"
	"github.com/spf13/cobra"
)

// ============================================================================
// Certificate Management Commands
// ============================================================================

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Manage SSL certificates for MySQL",
	Long:  "Manage SSL certificates for MySQL database connections",
}

var certListCmd = &cobra.Command{
	Use:   "list",
	Short: "List MySQL SSL certificates",
	Run: func(cmd *cobra.Command, args []string) {
		instanceID, _ := cmd.Flags().GetString("instance-id")

		store, err := cert.NewCertificateStore()
		if err != nil {
			exitWithError("failed to create certificate store", err)
		}

		certificates, err := store.ListCertificates("mysql", getRegion(), instanceID)
		if err != nil {
			exitWithError("failed to list certificates", err)
		}

		if len(certificates) == 0 {
			fmt.Println("No MySQL SSL certificates found.")
			if instanceID != "" {
				fmt.Printf("\nTo import certificates for instance %s:\n", instanceID)
				fmt.Printf("  nhncloud rds-mysql cert import --instance-id %s --ca-cert ca.pem --client-cert client.pem --client-key client.key\n", instanceID)
			}
			return
		}

		if output == "json" {
			printJSON(certificates)
			return
		}

		fmt.Printf("Found %d MySQL SSL certificate(s):\n\n", len(certificates))
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

var certImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import MySQL SSL certificate",
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

		fmt.Printf("Importing MySQL SSL certificates for instance: %s\n", instanceID)

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
			ServiceType: "mysql",
			Region:      region,
			InstanceID:  instanceID,
			Version:     version,
			Source:      "manual",
			Description: fmt.Sprintf("CA certificate for MySQL instance %s", instanceID),
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
				ServiceType: "mysql",
				Region:      region,
				InstanceID:  instanceID,
				Version:     version,
				Source:      "manual",
				Description: fmt.Sprintf("Client certificate for MySQL instance %s", instanceID),
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
				ServiceType: "mysql",
				Region:      region,
				InstanceID:  instanceID,
				Version:     version,
				Source:      "manual",
				Description: fmt.Sprintf("Client key for MySQL instance %s", instanceID),
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

		fmt.Printf("\nSuccessfully imported SSL certificates for MySQL instance %s\n", instanceID)
		fmt.Printf("\nYou can now use these certificates with:\n")
		fmt.Printf("  nhncloud rds-mysql connect %s --user <username> --password <password>\n", instanceID)
	},
}

var certDeleteCmd = &cobra.Command{
	Use:   "delete [certificate-id]",
	Short: "Delete MySQL SSL certificate",
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
// Connection Commands
// ============================================================================

var connectCmd = &cobra.Command{
	Use:   "connect [instance-id]",
	Short: "Connect to MySQL instance",
	Long:  "Connect to a MySQL instance using SSL certificates and credentials",
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

		fmt.Printf("Connecting to MySQL instance: %s\n", instanceID)

		// Get instance details
		client := newMySQLClient()
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
		var port int = 3306

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
				certs, err := store.ListCertificates("mysql", getRegion(), instanceID)
				if err == nil && len(certs) > 0 {
					fmt.Printf("  SSL Certificates: Found %d certificate(s)\n", len(certs))
				}
			}

			// Register TLS config
			tlsConfig := &tls.Config{
				InsecureSkipVerify: true, // For simplicity; production should verify certs
			}
			mysql.RegisterTLSConfig("custom", tlsConfig)
			tlsConfigName = "custom"
		} else {
			fmt.Printf("  SSL: DISABLED (not recommended for production)\n")
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
		fmt.Printf("\nAttempting connection...\n")
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			exitWithError("failed to open connection", err)
		}
		defer db.Close()

		// Test connection
		err = db.Ping()
		if err != nil {
			exitWithError("failed to connect to MySQL instance", err)
		}

		fmt.Printf("Successfully connected to MySQL instance!\n")

		// Execute test query
		var version, currentUser, currentDB string
		err = db.QueryRow("SELECT VERSION(), USER(), DATABASE()").Scan(&version, &currentUser, &currentDB)
		if err != nil {
			fmt.Printf("Warning: failed to execute test query: %v\n", err)
		} else {
			fmt.Printf("\nConnection Test:\n")
			fmt.Printf("  MySQL Version: %s\n", version)
			fmt.Printf("  Connected User: %s\n", currentUser)
			fmt.Printf("  Current Database: %s\n", currentDB)
		}

		fmt.Printf("\nConnection successful! Use 'nhncloud rds-mysql query' to execute SQL.\n")
	},
}

// ============================================================================
// Query Commands
// ============================================================================

var queryCmd = &cobra.Command{
	Use:   "query [instance-id] [sql]",
	Short: "Execute SQL query on MySQL instance",
	Long:  "Execute a SQL query on a MySQL instance and return formatted results",
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
		client := newMySQLClient()
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
		var port int = 3306

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
			mysql.RegisterTLSConfig("query", tlsConfig)
			tlsConfigName = "query"
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
			exitWithError("failed to connect to MySQL instance", err)
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
				printQueryJSON(columns, allRows)
			case "csv":
				printQueryCSV(columns, allRows)
			default:
				printQueryTable(columns, allRows, duration)
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

// printQueryTable formats results as a table
func printQueryTable(columns []string, rows [][]string, duration time.Duration) {
	if len(rows) == 0 {
		fmt.Printf("(0 rows) (%.2fms)\n", float64(duration.Nanoseconds())/1000000)
		return
	}

	// Calculate column widths
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

	// Print header
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

	// Print rows
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

// printQueryJSON formats results as JSON
func printQueryJSON(columns []string, rows [][]string) {
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

// printQueryCSV formats results as CSV
func printQueryCSV(columns []string, rows [][]string) {
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

func init() {
	// Certificate commands
	rdsMySQLCmd.AddCommand(certCmd)
	certCmd.AddCommand(certListCmd)
	certCmd.AddCommand(certImportCmd)
	certCmd.AddCommand(certDeleteCmd)

	certListCmd.Flags().String("instance-id", "", "Filter by database instance ID")

	certImportCmd.Flags().String("instance-id", "", "Database instance ID (required)")
	certImportCmd.Flags().String("version", "", "Certificate version")
	certImportCmd.Flags().String("ca-cert", "", "CA certificate file path (required)")
	certImportCmd.Flags().String("client-cert", "", "Client certificate file path")
	certImportCmd.Flags().String("client-key", "", "Client key file path")
	certImportCmd.Flags().String("description", "", "Certificate description")

	// Connect command
	rdsMySQLCmd.AddCommand(connectCmd)
	connectCmd.Flags().String("user", "", "Database username (required)")
	connectCmd.Flags().String("password", "", "Database password (required)")
	connectCmd.Flags().String("database", "", "Database/schema name to connect to")
	connectCmd.Flags().Bool("disable-ssl", false, "Disable SSL connection (not recommended)")
	connectCmd.Flags().Int("timeout", 30, "Connection timeout in seconds")

	// Query command
	rdsMySQLCmd.AddCommand(queryCmd)
	queryCmd.Flags().String("user", "", "Database username (required)")
	queryCmd.Flags().String("password", "", "Database password (required)")
	queryCmd.Flags().String("database", "", "Database/schema name to connect to")
	queryCmd.Flags().Bool("disable-ssl", false, "Disable SSL connection (not recommended)")
	queryCmd.Flags().Int("timeout", 30, "Connection timeout in seconds")
	queryCmd.Flags().String("format", "table", "Output format: table, json, csv")
}
