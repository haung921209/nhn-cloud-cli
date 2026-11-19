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

var kmClientMAC string

func init() {
	rootCmd.AddCommand(keymanagerCmd)
}

// addClientMACFlag gives the commands that read or use keys a --client-mac
// flag, for key stores with MAC address authentication. (Not a persistent
// flag: the catalog ops of this service define their own --client-mac.)
func addClientMACFlag(cmds ...*cobra.Command) {
	for _, c := range cmds {
		c.Flags().StringVar(&kmClientMAC, "client-mac", "", "Client MAC address, for key stores with MAC address authentication (X-TOAST-CLIENT-MAC-ADDR)")
	}
}

func newKeyManagerClient() *keymanager.Client {
	c := keymanager.NewClient(getRegion(), serviceAppKey("key-manager"), getAccessKey(), getSecretKey(), debug)
	c.ClientMAC = kmClientMAC
	return c
}
