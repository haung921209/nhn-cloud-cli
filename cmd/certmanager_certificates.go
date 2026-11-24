package cmd

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
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
		cmd.SilenceUsage = true
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
		fmt.Fprintln(w, "ID\tNAME\tAUTHORITY\tDOMAINS\tEXPIRES")
		for _, cert := range result.Body.Certificates {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				cert.CertificateID,
				cert.CertificateName,
				cert.Authority,
				strings.Join(cert.Domains, ","),
				cert.ExpirationDate,
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
		cmd.SilenceUsage = true
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

				fmt.Printf("Certificate ID:      %d\n", cert.CertificateID)
				fmt.Printf("Certificate Name:    %s\n", cert.CertificateName)
				fmt.Printf("Authority:           %s\n", cert.Authority)
				fmt.Printf("Domains:             %s\n", strings.Join(cert.Domains, ", "))
				fmt.Printf("Signature Algorithm: %s\n", cert.SignatureAlgorithm)
				fmt.Printf("File Created:        %s\n", cert.FileCreationDate)
				fmt.Printf("Expires:             %s\n", cert.ExpirationDate)
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
		cmd.SilenceUsage = true
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
			if err := plainFileName("--certificate-name", certName); err != nil {
				return err
			}
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}

			if !splitPEM {
				// Default: write the raw bundle as a single .pem file
				bundlePath := fmt.Sprintf("%s/%s.pem", outputDir, certName)
				if err := writePrivateFile(bundlePath, result.Data); err != nil {
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
				if err := writePrivateFile(keyPath, keyBlocks); err != nil {
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
