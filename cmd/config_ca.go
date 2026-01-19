package cmd

import (
	"fmt"
	"os"

	"github.com/haung921209/nhn-cloud-cli/internal/cert"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var configCACmd = &cobra.Command{
	Use:   "ca",
	Short: "Manage CA certificates for database connections",
	Long:  `Import and manage Certificate Authority (CA) certificates for secure database connections.`,
}

var importCACmd = &cobra.Command{
	Use:   "import",
	Short: "Import a CA certificate",
	Long: `Import a CA certificate file for a specific service and region.
Optionally bind it to a specific database instance.

Example:
  nhncloud config ca import --service rds-mysql --region kr1 --file ./ca.pem
  nhncloud config ca import --service rds-postgresql --region kr1 --instance-id <uuid> --file ./root.crt`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		service, _ := cmd.Flags().GetString("service")
		region, _ := cmd.Flags().GetString("region")
		instanceID, _ := cmd.Flags().GetString("instance-id")
		description, _ := cmd.Flags().GetString("description")
		certType, _ := cmd.Flags().GetString("type")

		if filePath == "" || service == "" || region == "" {
			exitWithError("Flags --file, --service, and --region are required", nil)
		}

		// Read file
		data, err := os.ReadFile(filePath)
		if err != nil {
			exitWithError("Failed to read certificate file", err)
		}

		// Initialize Store
		store, err := cert.NewCertificateStore()
		if err != nil {
			exitWithError("Failed to initialize certificate store", err)
		}

		// Store Certificate
		req := &cert.CertificateRequest{
			Type:        certType,
			ServiceType: service,
			Region:      region,
			InstanceID:  instanceID,
			CertData:    data,
			Description: description,
			Source:      "import",
		}

		info, err := store.StoreCertificate(req)
		if err != nil {
			exitWithError("Failed to store certificate", err)
		}

		fmt.Printf("Certificate imported successfully.\n")
		fmt.Printf("ID: %s\n", info.ID)
		fmt.Printf("Type: %s\n", info.Type)
		fmt.Printf("Path: %s\n", info.FilePath)
	},
}

var listCACmd = &cobra.Command{
	Use:   "list",
	Short: "List imported CA certificates",
	Run: func(cmd *cobra.Command, args []string) {
		service, _ := cmd.Flags().GetString("service")
		region, _ := cmd.Flags().GetString("region")
		instanceID, _ := cmd.Flags().GetString("instance-id")

		store, err := cert.NewCertificateStore()
		if err != nil {
			exitWithError("Failed to initialize certificate store", err)
		}

		certs, err := store.ListCertificates(service, region, instanceID)
		if err != nil {
			exitWithError("Failed to list certificates", err)
		}

		if len(certs) == 0 {
			fmt.Println("No certificates found matching criteria.")
			return
		}

		fmt.Printf("%-10s %-12s %-11s %-10s %-36s %s\n", "ID", "TYPE", "SERVICE", "REGION", "INSTANCE ID", "DESCRIPTION")
		for _, c := range certs {
			fmt.Printf("%-10s %-12s %-11s %-10s %-36s %s\n", c.ID, c.Type, c.ServiceType, c.Region, c.InstanceID, c.Description)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configCACmd)
	configCACmd.AddCommand(importCACmd)
	configCACmd.AddCommand(listCACmd)

	importCACmd.Flags().String("type", "CA", "Certificate type (CA, CLIENT-CERT, CLIENT-KEY)")
	importCACmd.Flags().String("file", "", "Certificate file path (required)")
	importCACmd.Flags().String("service", "", "Service type (rds-mysql, rds-postgresql, etc.) (required)")
	importCACmd.Flags().String("region", "", "Region (required)")
	importCACmd.Flags().String("instance-id", "", "Bind to specific instance ID")
	importCACmd.Flags().String("description", "", "Description")

	listCACmd.Flags().String("service", "", "Filter by service")
	listCACmd.Flags().String("region", "", "Filter by region")
	listCACmd.Flags().String("instance-id", "", "Filter by instance ID")
}
