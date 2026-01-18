package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/security/keymanager"
	"github.com/spf13/cobra"
)

func init() {
	keymanagerCmd.AddCommand(kmGetPrivateKeyCmd)
	keymanagerCmd.AddCommand(kmGetPublicKeyCmd)
	keymanagerCmd.AddCommand(kmSignCmd)
	keymanagerCmd.AddCommand(kmVerifyCmd)

	kmGetPrivateKeyCmd.Flags().String("key-id", "", "Key ID (required)")
	kmGetPrivateKeyCmd.MarkFlagRequired("key-id")

	kmGetPublicKeyCmd.Flags().String("key-id", "", "Key ID (required)")
	kmGetPublicKeyCmd.MarkFlagRequired("key-id")

	kmSignCmd.Flags().String("key-id", "", "Key ID (required)")
	kmSignCmd.Flags().String("data", "", "Base64 encoded data to sign (required)")
	kmSignCmd.MarkFlagRequired("key-id")
	kmSignCmd.MarkFlagRequired("data")

	kmVerifyCmd.Flags().String("key-id", "", "Key ID (required)")
	kmVerifyCmd.Flags().String("data", "", "Base64 encoded data (required)")
	kmVerifyCmd.Flags().String("signature", "", "Base64 encoded signature (required)")
	kmVerifyCmd.MarkFlagRequired("key-id")
	kmVerifyCmd.MarkFlagRequired("data")
	kmVerifyCmd.MarkFlagRequired("signature")
}

var kmGetPrivateKeyCmd = &cobra.Command{
	Use:   "get-private-key",
	Short: "Get private key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		keyID, _ := cmd.Flags().GetString("key-id")

		result, err := client.GetPrivateKey(ctx, keyID)
		if err != nil {
			exitWithError("Failed to get private key", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("%s\n", result.Body.PrivateKey)
	},
}

var kmGetPublicKeyCmd = &cobra.Command{
	Use:   "get-public-key",
	Short: "Get public key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		keyID, _ := cmd.Flags().GetString("key-id")

		result, err := client.GetPublicKey(ctx, keyID)
		if err != nil {
			exitWithError("Failed to get public key", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("%s\n", result.Body.PublicKey)
	},
}

var kmSignCmd = &cobra.Command{
	Use:   "sign-data",
	Short: "Sign data with asymmetric key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		keyID, _ := cmd.Flags().GetString("key-id")
		data, _ := cmd.Flags().GetString("data")

		input := &keymanager.SignInput{
			Data: data,
		}

		result, err := client.Sign(ctx, keyID, input)
		if err != nil {
			exitWithError("Failed to sign", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Signature: %s\n", result.Body.Signature)
	},
}

var kmVerifyCmd = &cobra.Command{
	Use:   "verify-signature",
	Short: "Verify signature with asymmetric key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()
		keyID, _ := cmd.Flags().GetString("key-id")
		data, _ := cmd.Flags().GetString("data")
		signature, _ := cmd.Flags().GetString("signature")

		input := &keymanager.VerifyInput{
			Data:      data,
			Signature: signature,
		}

		result, err := client.Verify(ctx, keyID, input)
		if err != nil {
			exitWithError("Failed to verify", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		if result.Body.Result {
			fmt.Println("Signature is valid")
		} else {
			fmt.Println("Signature is invalid")
		}
	},
}
