package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/privatedns"
	"github.com/spf13/cobra"
)

var privateDNSCmd = &cobra.Command{
	Use:     "private-dns",
	Aliases: []string{"pdns", "privatedns"},
	Short:   "Manage Private DNS zones and records",
}

func init() {
	rootCmd.AddCommand(privateDNSCmd)
}

func newPrivateDNSClient() *privatedns.Client {
	return privatedns.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}
