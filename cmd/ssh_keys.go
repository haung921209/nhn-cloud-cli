package cmd

import (
	"github.com/spf13/cobra"
)

var sshKeysCmd = &cobra.Command{
	Use:   "ssh-keys",
	Short: "Manage local SSH key pairs",
	Long:  `Manage SSH private keys for secure access to compute instances.`,
}

func init() {
	rootCmd.AddCommand(sshKeysCmd)
}
