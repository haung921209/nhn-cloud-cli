package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/compute"
	"github.com/spf13/cobra"
)

func init() {
	computeCmd.AddCommand(computeDescribeKeyPairsCmd)
	computeCmd.AddCommand(computeCreateKeyPairCmd)
	computeCmd.AddCommand(computeDeleteKeyPairCmd)

	computeCreateKeyPairCmd.Flags().String("key-name", "", "Keypair name (required)")
	computeCreateKeyPairCmd.Flags().String("public-key", "", "Public key content (optional)")
	computeCreateKeyPairCmd.MarkFlagRequired("key-name")

	computeDeleteKeyPairCmd.Flags().String("key-name", "", "Keypair name (required)")
	computeDeleteKeyPairCmd.MarkFlagRequired("key-name")

	computeCmd.AddCommand(computeImportKeyCmd)
	computeImportKeyCmd.Flags().String("key-name", "", "Keypair name (required)")
	computeImportKeyCmd.Flags().String("private-key-file", "", "Path to private key file (required)")
	computeImportKeyCmd.MarkFlagRequired("key-name")
	computeImportKeyCmd.MarkFlagRequired("private-key-file")
}

var computeDescribeKeyPairsCmd = &cobra.Command{
	Use:     "describe-key-pairs",
	Aliases: []string{"keypairs"},
	Short:   "List SSH keypairs",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		result, err := client.ListKeyPairs(ctx)
		if err != nil {
			exitWithError("Failed to list keypairs", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tFINGERPRINT")
		for _, kp := range result.KeyPairs {
			fmt.Fprintf(w, "%s\t%s\n", kp.KeyPair.Name, kp.KeyPair.Fingerprint)
		}
		w.Flush()
	},
}

var computeCreateKeyPairCmd = &cobra.Command{
	Use:     "create-key-pair",
	Aliases: []string{"keypair-create"},
	Short:   "Create a new SSH keypair",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("key-name")
		publicKey, _ := cmd.Flags().GetString("public-key")

		input := &compute.CreateKeyPairInput{
			Name:      name,
			PublicKey: publicKey,
		}

		result, err := client.CreateKeyPair(ctx, input)
		if err != nil {
			exitWithError("Failed to create keypair", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Keypair created: %s\n", result.KeyPair.Name)
		fmt.Printf("Fingerprint: %s\n", result.KeyPair.Fingerprint)
		if result.KeyPair.PrivateKey != "" {
			fmt.Printf("\nPrivate Key (save this - it won't be shown again):\n%s\n", result.KeyPair.PrivateKey)
		}
	},
}

var computeDeleteKeyPairCmd = &cobra.Command{
	Use:     "delete-key-pair",
	Aliases: []string{"keypair-delete"},
	Short:   "Delete an SSH keypair",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()
		keyName, _ := cmd.Flags().GetString("key-name")

		if err := client.DeleteKeyPair(ctx, keyName); err != nil {
			exitWithError("Failed to delete keypair", err)
		}

		fmt.Printf("Keypair %s deleted\n", keyName)
	},
}

var computeImportKeyCmd = &cobra.Command{
	Use:   "import-key",
	Short: "Import a private key to local SSH directory",
	Long:  `Import an existing private key file to ~/.ssh/ for use with 'connect' command.`,
	Run: func(cmd *cobra.Command, args []string) {
		keyName, _ := cmd.Flags().GetString("key-name")
		privateKeyFile, _ := cmd.Flags().GetString("private-key-file")

		// 1. Verify Source
		srcInfo, err := os.Stat(privateKeyFile)
		if err != nil {
			exitWithError("Failed to find private key file", err)
		}
		if srcInfo.IsDir() {
			exitWithError("Private key path is a directory", fmt.Errorf("path: %s", privateKeyFile))
		}

		// 2. Read Source
		content, err := os.ReadFile(privateKeyFile)
		if err != nil {
			exitWithError("Failed to read private key file", err)
		}

		// 3. Determine Destination
		homeDir, err := os.UserHomeDir()
		if err != nil {
			exitWithError("Failed to get user home directory", err)
		}
		sshDir := filepath.Join(homeDir, ".ssh")
		if _, err := os.Stat(sshDir); os.IsNotExist(err) {
			if err := os.MkdirAll(sshDir, 0700); err != nil {
				exitWithError("Failed to create .ssh directory", err)
			}
		}

		destPath := filepath.Join(sshDir, keyName+".pem")

		// 4. Check Overwrite (Simple check)
		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("Warning: Key file '%s' already exists. Overwriting...\n", destPath)
		}

		// 5. Write and Secure
		// os.WriteFile with 0600 creates or truncates
		if err := os.WriteFile(destPath, content, 0600); err != nil {
			exitWithError("Failed to write key file", err)
		}

		fmt.Printf("Success! Key '%s' imported to '%s'\n", keyName, destPath)
		fmt.Println("Permissions set to 0600 (read/write only for owner).")
	},
}
