package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/certmanager"
	"github.com/spf13/cobra"
)

var certmanagerCmd = &cobra.Command{
	Use:     "certmanager",
	Aliases: []string{"cert-manager", "cert"},
	Short:   "Manage SSL/TLS certificates",
	Long:    `Manage SSL/TLS certificates in NHN Cloud Certificate Manager.`,
}

func init() {
	rootCmd.AddCommand(certmanagerCmd)
}

func newCertManagerClient() *certmanager.Client {
	return certmanager.NewClient(getAppKey(), getAccessKey(), getSecretKey(), nil, debug)
}
