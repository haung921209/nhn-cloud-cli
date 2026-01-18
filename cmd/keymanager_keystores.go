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
		fmt.Fprintln(w, "ID\tNAME\tKEY_COUNT\tSTATUS")
		for _, ks := range result.Body.KeyStores {
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
				ks.KeyStoreID, ks.Name, ks.KeyCount, ks.Status)
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

		ks := result.Body.KeyStore
		fmt.Printf("ID:          %s\n", ks.KeyStoreID)
		fmt.Printf("Name:        %s\n", ks.Name)
		fmt.Printf("Description: %s\n", ks.Description)
		fmt.Printf("Key Count:   %d\n", ks.KeyCount)
		fmt.Printf("Status:      %s\n", ks.Status)
		fmt.Printf("Created:     %s\n", ks.CreatedAt)
		fmt.Printf("Updated:     %s\n", ks.UpdatedAt)
	},
}
