package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/certmanager"
	"github.com/spf13/cobra"
)

// certmanagerCmd represents the certmanager command
var certmanagerCmd = &cobra.Command{
	Use:     "certmanager",
	Aliases: []string{"cert-manager", "cert"},
	Short:   "Manage SSL/TLS certificates",
	Long:    `Manage SSL/TLS certificates in NHN Cloud Certificate Manager.`,
}

// certmanagerListCmd lists all certificates
var certmanagerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all certificates",
	Long:  `List all SSL/TLS certificates in Certificate Manager.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newCertManagerClient()
		ctx := context.Background()

		result, err := client.ListCertificates(ctx)
		if err != nil {
			return fmt.Errorf("failed to list certificates: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Body.Certificates)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tSTATUS\tDOMAIN\tISSUER\tEXPIRES")
		for _, cert := range result.Body.Certificates {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				cert.CertificateName,
				cert.CertificateType,
				cert.Status,
				cert.DomainName,
				cert.Issuer,
				cert.NotAfter.Format("2006-01-02"),
			)
		}
		return w.Flush()
	},
}

// certmanagerGetCmd gets certificate details
var certmanagerGetCmd = &cobra.Command{
	Use:   "get <certificate-name>",
	Short: "Get certificate details",
	Long:  `Get detailed information about a specific certificate.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		certName := args[0]
		client := newCertManagerClient()
		ctx := context.Background()

		result, err := client.ListCertificates(ctx)
		if err != nil {
			return fmt.Errorf("failed to get certificates: %w", err)
		}

		// Find the specific certificate
		for _, cert := range result.Body.Certificates {
			if cert.CertificateName == certName {
				if output == "json" {
					enc := json.NewEncoder(os.Stdout)
					enc.SetIndent("", "  ")
					return enc.Encode(cert)
				}

				fmt.Printf("Certificate Name:    %s\n", cert.CertificateName)
				fmt.Printf("Type:                %s\n", cert.CertificateType)
				fmt.Printf("Status:              %s\n", cert.Status)
				fmt.Printf("Domain Name:         %s\n", cert.DomainName)
				if len(cert.SubjectAlternativeNames) > 0 {
					fmt.Printf("SANs:                %v\n", cert.SubjectAlternativeNames)
				}
				fmt.Printf("Issuer:              %s\n", cert.Issuer)
				fmt.Printf("Serial Number:       %s\n", cert.SerialNumber)
				fmt.Printf("Valid From:          %s\n", cert.NotBefore.Format("2006-01-02 15:04:05"))
				fmt.Printf("Valid Until:         %s\n", cert.NotAfter.Format("2006-01-02 15:04:05"))
				if cert.KeyAlgorithm != "" {
					fmt.Printf("Key Algorithm:       %s\n", cert.KeyAlgorithm)
				}
				if cert.KeySize > 0 {
					fmt.Printf("Key Size:            %d\n", cert.KeySize)
				}
				if cert.SignatureAlgorithm != "" {
					fmt.Printf("Signature Algorithm: %s\n", cert.SignatureAlgorithm)
				}
				fmt.Printf("Created At:          %s\n", cert.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("Updated At:          %s\n", cert.UpdatedAt.Format("2006-01-02 15:04:05"))
				return nil
			}
		}

		return fmt.Errorf("certificate not found: %s", certName)
	},
}

// certmanagerDownloadCmd downloads certificate files
var certmanagerDownloadCmd = &cobra.Command{
	Use:   "download <certificate-name>",
	Short: "Download certificate files",
	Long:  `Download certificate, private key, and certificate chain files for a certificate.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		certName := args[0]
		client := newCertManagerClient()
		ctx := context.Background()

		result, err := client.DownloadCertificateFiles(ctx, certName)
		if err != nil {
			return fmt.Errorf("failed to download certificate files: %w", err)
		}

		// Check if output to files or display
		outputDir, _ := cmd.Flags().GetString("output-dir")
		if outputDir != "" {
			// Create output directory if it doesn't exist
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			// Write certificate file
			if result.Body.Certificate != "" {
				certPath := fmt.Sprintf("%s/%s.crt", outputDir, certName)
				if err := os.WriteFile(certPath, []byte(result.Body.Certificate), 0644); err != nil {
					return fmt.Errorf("failed to write certificate file: %w", err)
				}
				fmt.Printf("Certificate written to: %s\n", certPath)
			}

			// Write private key file (with restricted permissions)
			if result.Body.PrivateKey != "" {
				keyPath := fmt.Sprintf("%s/%s.key", outputDir, certName)
				if err := os.WriteFile(keyPath, []byte(result.Body.PrivateKey), 0600); err != nil {
					return fmt.Errorf("failed to write private key file: %w", err)
				}
				fmt.Printf("Private key written to: %s\n", keyPath)
			}

			// Write certificate chain file
			if result.Body.CertificateChain != "" {
				chainPath := fmt.Sprintf("%s/%s-chain.crt", outputDir, certName)
				if err := os.WriteFile(chainPath, []byte(result.Body.CertificateChain), 0644); err != nil {
					return fmt.Errorf("failed to write certificate chain file: %w", err)
				}
				fmt.Printf("Certificate chain written to: %s\n", chainPath)
			}

			return nil
		}

		// Display as JSON if no output directory specified
		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Body)
		}

		// Display as text
		fmt.Println("=== Certificate ===")
		if result.Body.Certificate != "" {
			fmt.Println(result.Body.Certificate)
		} else {
			fmt.Println("(not available)")
		}

		fmt.Println("\n=== Private Key ===")
		if result.Body.PrivateKey != "" {
			fmt.Println(result.Body.PrivateKey)
		} else {
			fmt.Println("(not available)")
		}

		fmt.Println("\n=== Certificate Chain ===")
		if result.Body.CertificateChain != "" {
			fmt.Println(result.Body.CertificateChain)
		} else {
			fmt.Println("(not available)")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(certmanagerCmd)
	certmanagerCmd.AddCommand(certmanagerListCmd)
	certmanagerCmd.AddCommand(certmanagerGetCmd)
	certmanagerCmd.AddCommand(certmanagerDownloadCmd)

	// Download flags
	certmanagerDownloadCmd.Flags().String("output-dir", "", "Directory to save certificate files")
}

func newCertManagerClient() *certmanager.Client {
	return certmanager.NewClient(getAppKey(), getAccessKey(), getSecretKey(), nil, debug)
}
