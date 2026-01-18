package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/transithub"
	"github.com/spf13/cobra"
)

var transitHubCmd = &cobra.Command{
	Use:     "transit-hub",
	Aliases: []string{"th", "transithub"},
	Short:   "Manage Transit Hubs for multi-VPC networking",
}

func init() {
	rootCmd.AddCommand(transitHubCmd)
}

func newTransitHubClient() *transithub.Client {
	return transithub.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}
