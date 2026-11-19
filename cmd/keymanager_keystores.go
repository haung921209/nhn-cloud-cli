package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	keymanagerCmd.AddCommand(kmDescribeKeyStoresCmd)
	keymanagerCmd.AddCommand(kmGetKeyStoreCmd)

	kmGetKeyStoreCmd.Flags().String("keystore-id", "", "Key Store ID (required)")
	kmGetKeyStoreCmd.MarkFlagRequired("keystore-id")
}

var kmDescribeKeyStoresCmd = &cobra.Command{
	Use:     "describe-key-stores",
	Aliases: []string{"list-key-stores"},
	Short:   "List all key stores",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()

		result, err := client.ListKeyStores(ctx)
		if err != nil {
			exitWithError("Failed to list key stores", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tIPV4_AUTH\tMAC_AUTH\tCERT_AUTH")
		for _, ks := range result.Body.KeyStores {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				ks.KeyStoreID, ks.Name, ks.IP4AuthUse, ks.MACAuthUse, ks.CertificateAuthUse)
		}
		w.Flush()
	},
}

var kmGetKeyStoreCmd = &cobra.Command{
	Use:     "describe-key-store",
	Aliases: []string{"get-key-store"},
	Short:   "Get key store details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("keystore-id")

		result, err := client.GetKeyStore(ctx, id)
		if err != nil {
			exitWithError("Failed to get key store", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		ks := result.Body
		fmt.Printf("ID:          %d\n", ks.KeyStoreID)
		fmt.Printf("Name:        %s\n", ks.Name)
		fmt.Printf("Description: %s\n", ks.Description)
		fmt.Printf("IPv4 Auth:   %s\n", ks.IP4AuthUse)
		fmt.Printf("MAC Auth:    %s\n", ks.MACAuthUse)
		fmt.Printf("Cert Auth:   %s\n", ks.CertificateAuthUse)
		fmt.Printf("Created:     %s\n", ks.CreationDatetime)
		fmt.Printf("Updated:     %s\n", ks.LastChangeDatetime)
	},
}
