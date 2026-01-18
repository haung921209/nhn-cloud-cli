package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/security/keymanager"
	"github.com/spf13/cobra"
)

func init() {
	keymanagerCmd.AddCommand(kmGetSymmetricKeyCmd)
	keymanagerCmd.AddCommand(kmEncryptCmd)
	keymanagerCmd.AddCommand(kmDecryptCmd)
	keymanagerCmd.AddCommand(kmCreateLocalKeyCmd)

	kmGetSymmetricKeyCmd.Flags().String("key-id", "", "Key ID (required)")
	kmGetSymmetricKeyCmd.MarkFlagRequired("key-id")

	kmEncryptCmd.Flags().String("key-id", "", "Key ID (required)")
	kmEncryptCmd.Flags().String("plaintext", "", "Plaintext to encrypt (required)")
	kmEncryptCmd.Flags().String("aad", "", "Additional authenticated data")
	kmEncryptCmd.MarkFlagRequired("key-id")
	kmEncryptCmd.MarkFlagRequired("plaintext")

	kmDecryptCmd.Flags().String("key-id", "", "Key ID (required)")
	kmDecryptCmd.Flags().String("ciphertext", "", "Ciphertext to decrypt (required)")
	kmDecryptCmd.Flags().String("iv", "", "Initialization vector")
	kmDecryptCmd.Flags().String("tag", "", "Authentication tag")
	kmDecryptCmd.Flags().String("aad", "", "Additional authenticated data")
	kmDecryptCmd.MarkFlagRequired("key-id")
	kmDecryptCmd.MarkFlagRequired("ciphertext")

	kmCreateLocalKeyCmd.Flags().String("key-id", "", "Key ID (required)")
	kmCreateLocalKeyCmd.MarkFlagRequired("key-id")
}

var kmGetSymmetricKeyCmd = &cobra.Command{
	Use:   "get-symmetric-key",
	Short: "Get symmetric key value",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		keyID, _ := cmd.Flags().GetString("key-id")

		result, err := client.GetSymmetricKey(ctx, keyID)
		if err != nil {
			exitWithError("Failed to get symmetric key", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Key Value: %s\n", result.Body.KeyValue)
	},
}

var kmEncryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt data with symmetric key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		keyID, _ := cmd.Flags().GetString("key-id")
		plaintext, _ := cmd.Flags().GetString("plaintext")
		aad, _ := cmd.Flags().GetString("aad")

		input := &keymanager.EncryptInput{
			Plaintext: plaintext,
			AAD:       aad,
		}

		result, err := client.Encrypt(ctx, keyID, input)
		if err != nil {
			exitWithError("Failed to encrypt", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Ciphertext: %s\n", result.Body.Ciphertext)
		if result.Body.IV != "" {
			fmt.Printf("IV:         %s\n", result.Body.IV)
		}
		if result.Body.Tag != "" {
			fmt.Printf("Tag:        %s\n", result.Body.Tag)
		}
	},
}

var kmDecryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt data with symmetric key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		keyID, _ := cmd.Flags().GetString("key-id")
		ciphertext, _ := cmd.Flags().GetString("ciphertext")
		iv, _ := cmd.Flags().GetString("iv")
		tag, _ := cmd.Flags().GetString("tag")
		aad, _ := cmd.Flags().GetString("aad")

		input := &keymanager.DecryptInput{
			Ciphertext: ciphertext,
			IV:         iv,
			Tag:        tag,
			AAD:        aad,
		}

		result, err := client.Decrypt(ctx, keyID, input)
		if err != nil {
			exitWithError("Failed to decrypt", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Plaintext: %s\n", result.Body.Plaintext)
	},
}

var kmCreateLocalKeyCmd = &cobra.Command{
	Use:   "create-local-key",
	Short: "Create a local data key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		keyID, _ := cmd.Flags().GetString("key-id")

		result, err := client.CreateLocalKey(ctx, keyID)
		if err != nil {
			exitWithError("Failed to create local key", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Plain Data Key:     %s\n", result.Body.PlainDataKey)
		fmt.Printf("Encrypted Data Key: %s\n", result.Body.EncryptedDataKey)
	},
}
