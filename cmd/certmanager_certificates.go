package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	certmanagerCmd.AddCommand(cmDescribeCertificatesCmd)
	certmanagerCmd.AddCommand(cmDescribeCertificateCmd)
	certmanagerCmd.AddCommand(cmDownloadCertificateCmd)

	cmDescribeCertificateCmd.Flags().String("certificate-name", "", "Certificate Name (required)")
	cmDescribeCertificateCmd.MarkFlagRequired("certificate-name")

	cmDownloadCertificateCmd.Flags().String("certificate-name", "", "Certificate Name (required)")
	cmDownloadCertificateCmd.Flags().String("output-dir", "", "Directory to save the PEM bundle (default: writes <name>.pem)")
	cmDownloadCertificateCmd.Flags().Bool("split-pem", false, "When set with --output-dir, parse the PEM bundle and write CERTIFICATE blocks to <name>.crt and PRIVATE KEY blocks to <name>.key")
	cmDownloadCertificateCmd.MarkFlagRequired("certificate-name")
}

var cmDescribeCertificatesCmd = &cobra.Command{
	Use:     "describe-certificates",
	Aliases: []string{"list", "list-certificates"},
	Short:   "List all certificates",
	Long:    `List all SSL/TLS certificates in Certificate Manager.`,
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

var cmDescribeCertificateCmd = &cobra.Command{
	Use:     "describe-certificate",
	Aliases: []string{"get", "get-certificate"},
	Short:   "Get certificate details",
	Long:    `Get detailed information about a specific certificate.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		certName, _ := cmd.Flags().GetString("certificate-name")
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

var cmDownloadCertificateCmd = &cobra.Command{
	Use:     "download-certificate",
	Aliases: []string{"download"},
	Short:   "Download certificate files",
	Long: `Download certificate bundle as raw PEM.

SDK v0.2.0 changed DownloadCertificateFiles to return a raw PEM
bundle (cert + private key + chain concatenated) rather than three
separate fields. This command writes the bundle as-is by default;
pass --split-pem to use encoding/pem to extract CERTIFICATE blocks
into <name>.crt and PRIVATE KEY blocks into <name>.key (with restricted
permissions).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		certName, _ := cmd.Flags().GetString("certificate-name")
		client := newCertManagerClient()
		ctx := context.Background()

		result, err := client.DownloadCertificateFiles(ctx, certName)
		if err != nil {
			return fmt.Errorf("failed to download certificate files: %w", err)
		}
		if len(result.Data) == 0 {
			return fmt.Errorf("certificate bundle is empty (server returned 0 bytes)")
		}

		outputDir, _ := cmd.Flags().GetString("output-dir")
		splitPEM, _ := cmd.Flags().GetBool("split-pem")

		if outputDir != "" {
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			if !splitPEM {
				// Default: write the raw bundle as a single .pem file
				bundlePath := fmt.Sprintf("%s/%s.pem", outputDir, certName)
				if err := os.WriteFile(bundlePath, result.Data, 0600); err != nil {
					return fmt.Errorf("failed to write PEM bundle: %w", err)
				}
				fmt.Printf("PEM bundle written to: %s\n", bundlePath)
				return nil
			}

			// Split the PEM bundle by block type
			certPath := fmt.Sprintf("%s/%s.crt", outputDir, certName)
			keyPath := fmt.Sprintf("%s/%s.key", outputDir, certName)
			var certBlocks, keyBlocks []byte
			rest := result.Data
			for {
				var block *pem.Block
				block, rest = pem.Decode(rest)
				if block == nil {
					break
				}
				encoded := pem.EncodeToMemory(block)
				switch block.Type {
				case "CERTIFICATE":
					certBlocks = append(certBlocks, encoded...)
				case "PRIVATE KEY", "RSA PRIVATE KEY", "EC PRIVATE KEY", "ENCRYPTED PRIVATE KEY":
					keyBlocks = append(keyBlocks, encoded...)
				}
			}
			if len(certBlocks) > 0 {
				if err := os.WriteFile(certPath, certBlocks, 0644); err != nil {
					return fmt.Errorf("failed to write certificate file: %w", err)
				}
				fmt.Printf("Certificate written to: %s\n", certPath)
			}
			if len(keyBlocks) > 0 {
				if err := os.WriteFile(keyPath, keyBlocks, 0600); err != nil {
					return fmt.Errorf("failed to write private key file: %w", err)
				}
				fmt.Printf("Private key written to: %s\n", keyPath)
			}
			return nil
		}

		// JSON: base64-encode the bundle (raw bytes don't roundtrip cleanly through JSON)
		if output == "json" {
			payload := map[string]string{
				"certificateName": certName,
				"pemBundleBase64": base64.StdEncoding.EncodeToString(result.Data),
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(payload)
		}

		// Text: dump the PEM bundle as-is (already human-readable)
		fmt.Print(string(result.Data))
		return nil
	},
}
