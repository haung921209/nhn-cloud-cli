package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/internetgateway"
	"github.com/spf13/cobra"
)

var internetGatewayCmd = &cobra.Command{
	Use:     "internet-gateway",
	Aliases: []string{"igw", "gateway"},
	Short:   "Manage Internet Gateways",
}

func init() {
	rootCmd.AddCommand(internetGatewayCmd)
}

func newInternetGatewayClient() *internetgateway.Client {
	return internetgateway.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}
