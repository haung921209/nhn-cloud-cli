package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/security/keymanager"
	"github.com/spf13/cobra"
)

var keymanagerCmd = &cobra.Command{
	Use:     "keymanager",
	Aliases: []string{"km", "skm"},
	Short:   "Manage Secure Key Manager",
}

var kmKeyStoreCmd = &cobra.Command{
	Use:   "keystore",
	Short: "Manage key stores",
}

var kmKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Manage keys",
}

var kmSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Manage secrets",
}

var kmSymmetricCmd = &cobra.Command{
	Use:   "symmetric",
	Short: "Symmetric key operations",
}

var kmAsymmetricCmd = &cobra.Command{
	Use:   "asymmetric",
	Short: "Asymmetric key operations",
}

func init() {
	rootCmd.AddCommand(keymanagerCmd)

	// KeyStore commands
	keymanagerCmd.AddCommand(kmKeyStoreCmd)
	kmKeyStoreCmd.AddCommand(kmKeyStoreListCmd)
	kmKeyStoreCmd.AddCommand(kmKeyStoreGetCmd)

	// Key commands
	keymanagerCmd.AddCommand(kmKeyCmd)
	kmKeyCmd.AddCommand(kmKeyListCmd)
	kmKeyCmd.AddCommand(kmKeyGetCmd)
	kmKeyCmd.AddCommand(kmKeyCreateCmd)
	kmKeyCmd.AddCommand(kmKeyDeleteCmd)

	kmKeyListCmd.Flags().String("keystore-id", "", "Key store ID (required)")
	kmKeyListCmd.MarkFlagRequired("keystore-id")

	kmKeyGetCmd.Flags().String("keystore-id", "", "Key store ID (required)")
	kmKeyGetCmd.MarkFlagRequired("keystore-id")

	kmKeyCreateCmd.Flags().String("type", "", "Key type: secrets, symmetric-keys, asymmetric-keys (required)")
	kmKeyCreateCmd.Flags().String("keystore-name", "", "Key store name (required)")
	kmKeyCreateCmd.Flags().String("name", "", "Key name (required)")
	kmKeyCreateCmd.Flags().String("description", "", "Key description")
	kmKeyCreateCmd.Flags().String("algorithm", "", "Algorithm (AES256 for symmetric, RSA2048/RSA4096/EC_P256/EC_P384 for asymmetric)")
	kmKeyCreateCmd.Flags().String("secret", "", "Secret value (for secrets type)")
	kmKeyCreateCmd.Flags().Int("rotation-period", 0, "Rotation period in days")
	kmKeyCreateCmd.MarkFlagRequired("type")
	kmKeyCreateCmd.MarkFlagRequired("keystore-name")
	kmKeyCreateCmd.MarkFlagRequired("name")

	// Secret commands
	keymanagerCmd.AddCommand(kmSecretCmd)
	kmSecretCmd.AddCommand(kmSecretGetCmd)

	// Symmetric key commands
	keymanagerCmd.AddCommand(kmSymmetricCmd)
	kmSymmetricCmd.AddCommand(kmSymmetricGetCmd)
	kmSymmetricCmd.AddCommand(kmSymmetricEncryptCmd)
	kmSymmetricCmd.AddCommand(kmSymmetricDecryptCmd)
	kmSymmetricCmd.AddCommand(kmSymmetricLocalKeyCmd)

	kmSymmetricEncryptCmd.Flags().String("plaintext", "", "Plaintext to encrypt (required)")
	kmSymmetricEncryptCmd.Flags().String("aad", "", "Additional authenticated data")
	kmSymmetricEncryptCmd.MarkFlagRequired("plaintext")

	kmSymmetricDecryptCmd.Flags().String("ciphertext", "", "Ciphertext to decrypt (required)")
	kmSymmetricDecryptCmd.Flags().String("iv", "", "Initialization vector")
	kmSymmetricDecryptCmd.Flags().String("tag", "", "Authentication tag")
	kmSymmetricDecryptCmd.Flags().String("aad", "", "Additional authenticated data")
	kmSymmetricDecryptCmd.MarkFlagRequired("ciphertext")

	// Asymmetric key commands
	keymanagerCmd.AddCommand(kmAsymmetricCmd)
	kmAsymmetricCmd.AddCommand(kmAsymmetricPrivateKeyCmd)
	kmAsymmetricCmd.AddCommand(kmAsymmetricPublicKeyCmd)
	kmAsymmetricCmd.AddCommand(kmAsymmetricSignCmd)
	kmAsymmetricCmd.AddCommand(kmAsymmetricVerifyCmd)

	kmAsymmetricSignCmd.Flags().String("data", "", "Base64 encoded data to sign (required)")
	kmAsymmetricSignCmd.MarkFlagRequired("data")

	kmAsymmetricVerifyCmd.Flags().String("data", "", "Base64 encoded data (required)")
	kmAsymmetricVerifyCmd.Flags().String("signature", "", "Base64 encoded signature (required)")
	kmAsymmetricVerifyCmd.MarkFlagRequired("data")
	kmAsymmetricVerifyCmd.MarkFlagRequired("signature")

	// Client info command
	keymanagerCmd.AddCommand(kmClientInfoCmd)
}

func newKeyManagerClient() *keymanager.Client {
	return keymanager.NewClient(getRegion(), getAppKey(), getAccessKey(), getSecretKey(), debug)
}

// ============== KeyStore Commands ==============

var kmKeyStoreListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all key stores",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		result, err := client.ListKeyStores(context.Background())
		if err != nil {
			exitWithError("Failed to list key stores", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var kmKeyStoreGetCmd = &cobra.Command{
	Use:   "get [keystore-id]",
	Short: "Get key store details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		result, err := client.GetKeyStore(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get key store", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

// ============== Key Commands ==============

var kmKeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List keys in a key store",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		keyStoreID, _ := cmd.Flags().GetString("keystore-id")

		result, err := client.ListKeys(context.Background(), keyStoreID)
		if err != nil {
			exitWithError("Failed to list keys", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var kmKeyGetCmd = &cobra.Command{
	Use:   "get [key-id]",
	Short: "Get key details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		keyStoreID, _ := cmd.Flags().GetString("keystore-id")

		result, err := client.GetKey(context.Background(), keyStoreID, args[0])
		if err != nil {
			exitWithError("Failed to get key", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var kmKeyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
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

		result, err := client.CreateKey(context.Background(), keyType, input)
		if err != nil {
			exitWithError("Failed to create key", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Key created: %s\n", result.Body.KeyID)
	},
}

var kmKeyDeleteCmd = &cobra.Command{
	Use:   "delete [key-id]",
	Short: "Delete a key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		_, err := client.DeleteKeyImmediately(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to delete key", err)
		}
		fmt.Printf("Key %s deleted\n", args[0])
	},
}

// ============== Secret Commands ==============

var kmSecretGetCmd = &cobra.Command{
	Use:   "get [key-id]",
	Short: "Get secret value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		result, err := client.GetSecret(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get secret", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Secret: %s\n", result.Body.Secret)
	},
}

// ============== Symmetric Key Commands ==============

var kmSymmetricGetCmd = &cobra.Command{
	Use:   "get [key-id]",
	Short: "Get symmetric key value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		result, err := client.GetSymmetricKey(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get symmetric key", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Key Value: %s\n", result.Body.KeyValue)
	},
}

var kmSymmetricEncryptCmd = &cobra.Command{
	Use:   "encrypt [key-id]",
	Short: "Encrypt data with symmetric key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		plaintext, _ := cmd.Flags().GetString("plaintext")
		aad, _ := cmd.Flags().GetString("aad")

		input := &keymanager.EncryptInput{
			Plaintext: plaintext,
			AAD:       aad,
		}

		result, err := client.Encrypt(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to encrypt", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var kmSymmetricDecryptCmd = &cobra.Command{
	Use:   "decrypt [key-id]",
	Short: "Decrypt data with symmetric key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
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

		result, err := client.Decrypt(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to decrypt", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Plaintext: %s\n", result.Body.Plaintext)
	},
}

var kmSymmetricLocalKeyCmd = &cobra.Command{
	Use:   "create-local-key [key-id]",
	Short: "Create a local data key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		result, err := client.CreateLocalKey(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to create local key", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Plain Data Key:     %s\n", result.Body.PlainDataKey)
		fmt.Printf("Encrypted Data Key: %s\n", result.Body.EncryptedDataKey)
	},
}

// ============== Asymmetric Key Commands ==============

var kmAsymmetricPrivateKeyCmd = &cobra.Command{
	Use:   "private-key [key-id]",
	Short: "Get private key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		result, err := client.GetPrivateKey(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get private key", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("%s\n", result.Body.PrivateKey)
	},
}

var kmAsymmetricPublicKeyCmd = &cobra.Command{
	Use:   "public-key [key-id]",
	Short: "Get public key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		result, err := client.GetPublicKey(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get public key", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("%s\n", result.Body.PublicKey)
	},
}

var kmAsymmetricSignCmd = &cobra.Command{
	Use:   "sign [key-id]",
	Short: "Sign data with asymmetric key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		data, _ := cmd.Flags().GetString("data")

		input := &keymanager.SignInput{
			Data: data,
		}

		result, err := client.Sign(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to sign", err)
		}

		if output == "json" {
			out, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(out))
			return
		}

		fmt.Printf("Signature: %s\n", result.Body.Signature)
	},
}

var kmAsymmetricVerifyCmd = &cobra.Command{
	Use:   "verify [key-id]",
	Short: "Verify signature with asymmetric key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		data, _ := cmd.Flags().GetString("data")
		signature, _ := cmd.Flags().GetString("signature")

		input := &keymanager.VerifyInput{
			Data:      data,
			Signature: signature,
		}

		result, err := client.Verify(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to verify", err)
		}

		if output == "json" {
			out, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(out))
			return
		}

		if result.Body.Result {
			fmt.Println("Signature is valid")
		} else {
			fmt.Println("Signature is invalid")
		}
	},
}

// ============== Client Info Command ==============

var kmClientInfoCmd = &cobra.Command{
	Use:   "client-info",
	Short: "Get client information",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		result, err := client.GetClientInfo(context.Background())
		if err != nil {
			exitWithError("Failed to get client info", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("App Key:     %s\n", result.Body.AppKey)
		fmt.Printf("IP Address:  %s\n", result.Body.IPAddress)
		fmt.Printf("MAC Address: %s\n", result.Body.MACAddress)
	},
}
