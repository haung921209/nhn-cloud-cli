package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	keymanagerCmd.AddCommand(kmGetSecretCmd)

	kmGetSecretCmd.Flags().String("key-id", "", "Secret Key ID (required)")
	kmGetSecretCmd.MarkFlagRequired("key-id")
}

var kmGetSecretCmd = &cobra.Command{
	Use:   "get-secret",
	Short: "Get secret value",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		keyID, _ := cmd.Flags().GetString("key-id")

		result, err := client.GetSecret(ctx, keyID)
		if err != nil {
			exitWithError("Failed to get secret", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Secret: %s\n", result.Body.Secret)
	},
}
