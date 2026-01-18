package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/servicegateway"
	"github.com/spf13/cobra"
)

var serviceGatewayCmd = &cobra.Command{
	Use:     "service-gateway",
	Aliases: []string{"svcgw", "sg-gateway"},
	Short:   "Manage Service Gateways",
}

func init() {
	rootCmd.AddCommand(serviceGatewayCmd)
}

func newServiceGatewayClient() *servicegateway.Client {
	return servicegateway.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}
