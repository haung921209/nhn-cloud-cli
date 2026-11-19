package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/dnsplus"
	"github.com/spf13/cobra"
)

// dnsplusCmd represents the dnsplus command
var dnsplusCmd = &cobra.Command{
	Use:     "dnsplus",
	Aliases: []string{"dns", "dns-plus"},
	Short:   "Manage DNS Plus service",
	Long:    `Manage DNS zones, record sets, GSLB, and health checks in NHN Cloud DNS Plus.`,
}

func init() {
	rootCmd.AddCommand(dnsplusCmd)
}

func newDNSPlusClient() *dnsplus.Client {
	return dnsplus.NewClient(serviceAppKey("dnsplus"), nil, debug)
}
