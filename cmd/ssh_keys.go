package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-cli/internal/sshkeys"
	"github.com/spf13/cobra"
)

var sshKeysCmd = &cobra.Command{
	Use:   "ssh-keys",
	Short: "Manage local SSH key pairs",
	Long:  `Manage SSH private keys for secure access to compute instances.`,
}

var sshKeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List stored SSH keys",
	Run: func(cmd *cobra.Command, args []string) {
		manager := sshkeys.NewManager()
		keys, err := manager.List()
		if err != nil {
			exitWithError("Failed to list SSH keys", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(keys)
			return
		}

		if len(keys) == 0 {
			fmt.Println("No SSH keys stored.")
			fmt.Println("\nTo import a key:")
			fmt.Println("  nhncloud ssh-keys import <key-name> <path-to-pem-file>")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tFINGERPRINT\tCREATED")
		for _, key := range keys {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				key.Name, key.Type, key.Fingerprint, key.CreatedAt.Format("2006-01-02 15:04:05"))
		}
		w.Flush()
	},
}

var sshKeysImportCmd = &cobra.Command{
	Use:   "import <key-name> <file-path>",
	Short: "Import an SSH private key",
	Long: `Import an SSH private key file into the local key store.

Examples:
  nhncloud ssh-keys import my-key ~/Downloads/my-key.pem
  nhncloud ssh-keys import prod-key /path/to/private-key`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		keyName := args[0]
		filePath := args[1]

		if len(filePath) >= 2 && filePath[:2] == "~/" {
			home, err := os.UserHomeDir()
			if err != nil {
				exitWithError("Failed to get home directory", err)
			}
			filePath = filepath.Join(home, filePath[2:])
		}

		manager := sshkeys.NewManager()
		keyInfo, err := manager.Import(keyName, filePath)
		if err != nil {
			exitWithError("Failed to import SSH key", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(keyInfo)
			return
		}

		fmt.Printf("SSH key '%s' imported successfully\n", keyInfo.Name)
		fmt.Printf("  Path: %s\n", keyInfo.Path)
		fmt.Printf("  Fingerprint: %s\n", keyInfo.Fingerprint)
		fmt.Printf("  Type: %s\n", keyInfo.Type)
	},
}

var sshKeysGetCmd = &cobra.Command{
	Use:   "get <key-name>",
	Short: "Get SSH key details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := sshkeys.NewManager()
		keyInfo, err := manager.Get(args[0])
		if err != nil {
			exitWithError("Failed to get SSH key", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(keyInfo)
			return
		}

		fmt.Printf("Name: %s\n", keyInfo.Name)
		fmt.Printf("Type: %s\n", keyInfo.Type)
		fmt.Printf("Path: %s\n", keyInfo.Path)
		fmt.Printf("Fingerprint: %s\n", keyInfo.Fingerprint)
		fmt.Printf("Created: %s\n", keyInfo.CreatedAt.Format("2006-01-02 15:04:05"))
		if !keyInfo.LastUsed.IsZero() {
			fmt.Printf("Last Used: %s\n", keyInfo.LastUsed.Format("2006-01-02 15:04:05"))
		}
		if keyInfo.PublicKey != "" {
			fmt.Printf("\nPublic Key:\n%s\n", keyInfo.PublicKey)
		}
	},
}

var sshKeysRemoveCmd = &cobra.Command{
	Use:   "remove <key-name>",
	Short: "Remove an SSH key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := sshkeys.NewManager()
		if err := manager.Remove(args[0]); err != nil {
			exitWithError("Failed to remove SSH key", err)
		}

		fmt.Printf("SSH key '%s' removed successfully\n", args[0])
	},
}

var sshKeysExportCmd = &cobra.Command{
	Use:   "export <key-name> <destination-path>",
	Short: "Export SSH key to a file",
	Long: `Export a stored SSH key to a specified location.

Examples:
  nhncloud ssh-keys export my-key ~/.ssh/my-key.pem
  nhncloud ssh-keys export prod-key /tmp/backup-key.pem`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		keyName := args[0]
		destPath := args[1]

		if len(destPath) >= 2 && destPath[:2] == "~/" {
			home, err := os.UserHomeDir()
			if err != nil {
				exitWithError("Failed to get home directory", err)
			}
			destPath = filepath.Join(home, destPath[2:])
		}

		manager := sshkeys.NewManager()
		if err := manager.Export(keyName, destPath); err != nil {
			exitWithError("Failed to export SSH key", err)
		}

		fmt.Printf("SSH key '%s' exported to %s\n", keyName, destPath)
		fmt.Printf("Note: Set permissions with: chmod 600 %s\n", destPath)
	},
}

var sshKeysUseCmd = &cobra.Command{
	Use:   "use <key-name> <user@host>",
	Short: "SSH to instance using stored key",
	Long: `Use a stored SSH key to connect to a compute instance.

Examples:
  nhncloud ssh-keys use my-key ubuntu@192.168.1.100
  nhncloud ssh-keys use prod-key root@10.0.0.1`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		keyName := args[0]
		target := args[1]

		manager := sshkeys.NewManager()
		if err := manager.Connect(keyName, target); err != nil {
			exitWithError("Failed to connect", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sshKeysCmd)

	sshKeysCmd.AddCommand(sshKeysListCmd)
	sshKeysCmd.AddCommand(sshKeysImportCmd)
	sshKeysCmd.AddCommand(sshKeysGetCmd)
	sshKeysCmd.AddCommand(sshKeysRemoveCmd)
	sshKeysCmd.AddCommand(sshKeysExportCmd)
	sshKeysCmd.AddCommand(sshKeysUseCmd)
}
