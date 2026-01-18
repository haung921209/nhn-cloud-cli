package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/security/keymanager"
	"github.com/spf13/cobra"
)

func init() {
	keymanagerCmd.AddCommand(kmDescribeKeysCmd)
	keymanagerCmd.AddCommand(kmGetKeyCmd)
	keymanagerCmd.AddCommand(kmCreateKeyCmd)
	keymanagerCmd.AddCommand(kmDeleteKeyCmd)

	kmDescribeKeysCmd.Flags().String("keystore-id", "", "Key Store ID (required)")
	kmDescribeKeysCmd.MarkFlagRequired("keystore-id")

	kmGetKeyCmd.Flags().String("keystore-id", "", "Key Store ID (required)")
	kmGetKeyCmd.Flags().String("key-id", "", "Key ID (required)")
	kmGetKeyCmd.MarkFlagRequired("keystore-id")
	kmGetKeyCmd.MarkFlagRequired("key-id")

	kmCreateKeyCmd.Flags().String("name", "", "Key name (required)")
	kmCreateKeyCmd.Flags().String("keystore-name", "", "Key Store Name (required)")
	kmCreateKeyCmd.Flags().String("type", "", "Key type: secrets, symmetric-keys, asymmetric-keys (required)")
	kmCreateKeyCmd.Flags().String("description", "", "Description")
	kmCreateKeyCmd.Flags().String("algorithm", "", "Algorithm (AES256, RSA2048, etc.)")
	kmCreateKeyCmd.Flags().String("secret", "", "Secret value (for secrets type)")
	kmCreateKeyCmd.Flags().Int("rotation-period", 0, "Rotation period in days")
	kmCreateKeyCmd.MarkFlagRequired("name")
	kmCreateKeyCmd.MarkFlagRequired("keystore-name")
	kmCreateKeyCmd.MarkFlagRequired("type")

	kmDeleteKeyCmd.Flags().String("key-id", "", "Key ID (required)")
	kmDeleteKeyCmd.MarkFlagRequired("key-id")
}

var kmDescribeKeysCmd = &cobra.Command{
	Use:     "describe-keys",
	Aliases: []string{"list-keys"},
	Short:   "List keys in a key store",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		ksID, _ := cmd.Flags().GetString("keystore-id")

		result, err := client.ListKeys(ctx, ksID)
		if err != nil {
			exitWithError("Failed to list keys", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tALGORITHM\tSTATUS")
		for _, key := range result.Body.Keys {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				key.KeyID, key.Name, key.KeyType, key.KeyAlgorithm, key.Status)
		}
		w.Flush()
	},
}

var kmGetKeyCmd = &cobra.Command{
	Use:     "describe-key",
	Aliases: []string{"get-key"},
	Short:   "Get key details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		ksID, _ := cmd.Flags().GetString("keystore-id")
		keyID, _ := cmd.Flags().GetString("key-id")

		result, err := client.GetKey(ctx, ksID, keyID)
		if err != nil {
			exitWithError("Failed to get key", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		key := result.Body.Key
		fmt.Printf("ID:          %s\n", key.KeyID)
		fmt.Printf("Name:        %s\n", key.Name)
		fmt.Printf("Description: %s\n", key.Description)
		fmt.Printf("Type:        %s\n", key.KeyType)
		fmt.Printf("Algorithm:   %s\n", key.KeyAlgorithm)
		fmt.Printf("Size:        %d\n", key.KeySize)
		fmt.Printf("Status:      %s\n", key.Status)
		fmt.Printf("Rotation:    %d days\n", key.RotationPeriod)
		fmt.Printf("Created:     %s\n", key.CreatedAt)
		fmt.Printf("Updated:     %s\n", key.UpdatedAt)
	},
}

var kmCreateKeyCmd = &cobra.Command{
	Use:   "create-key",
	Short: "Create a new key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()

		keyType, _ := cmd.Flags().GetString("type")
		keyStoreName, _ := cmd.Flags().GetString("keystore-name")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		algorithm, _ := cmd.Flags().GetString("algorithm")
		secret, _ := cmd.Flags().GetString("secret")
		rotationPeriod, _ := cmd.Flags().GetInt("rotation-period")

		input := &keymanager.CreateKeyInput{
			KeyStoreName:   keyStoreName,
			Name:           name,
			Description:    description,
			RotationPeriod: rotationPeriod,
		}

		switch keyType {
		case "secrets":
			input.Secret = secret
		case "symmetric-keys":
			input.KeyAlgorithm = algorithm
		case "asymmetric-keys":
			input.Algorithm = algorithm
		}

		result, err := client.CreateKey(ctx, keyType, input)
		if err != nil {
			exitWithError("Failed to create key", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Key created: %s\n", result.Body.KeyID)
	},
}

var kmDeleteKeyCmd = &cobra.Command{
	Use:   "delete-key",
	Short: "Delete a key immediately",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		keyID, _ := cmd.Flags().GetString("key-id")

		_, err := client.DeleteKeyImmediately(ctx, keyID)
		if err != nil {
			exitWithError("Failed to delete key", err)
		}
		fmt.Printf("Key %s deleted\n", keyID)
	},
}
