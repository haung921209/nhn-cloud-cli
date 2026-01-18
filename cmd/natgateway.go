package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/natgateway"
	"github.com/spf13/cobra"
)

var natGatewayCmd = &cobra.Command{
	Use:     "nat-gateway",
	Aliases: []string{"nat", "natgw"},
	Short:   "Manage NAT Gateways",
}

func init() {
	rootCmd.AddCommand(natGatewayCmd)
}

func newNATGatewayClient() *natgateway.Client {
	return natgateway.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}
