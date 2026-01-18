package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/loadbalancer"
	"github.com/spf13/cobra"
)

var loadbalancerCmd = &cobra.Command{
	Use:     "loadbalancer",
	Aliases: []string{"lb"},
	Short:   "Manage Load Balancers",
}

func init() {
	rootCmd.AddCommand(loadbalancerCmd)
}

func newLBClient() *loadbalancer.Client {
	return loadbalancer.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}
