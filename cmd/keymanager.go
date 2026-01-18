package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/security/keymanager"
	"github.com/spf13/cobra"
)

var keymanagerCmd = &cobra.Command{
	Use:     "key-manager",
	Aliases: []string{"km", "kms", "keymanager"},
	Short:   "Manage Secure Key Manager (KMS) service",
}

func init() {
	rootCmd.AddCommand(keymanagerCmd)
}

func newKeyManagerClient() *keymanager.Client {
	return keymanager.NewClient(getRegion(), getAppKey(), getAccessKey(), getSecretKey(), debug)
}
