package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/flowlog"
	"github.com/spf13/cobra"
)

var flowlogCmd = &cobra.Command{
	Use:     "flow-log",
	Aliases: []string{"fl", "flowlog"},
	Short:   "Manage Flow Logs",
}

func init() {
	rootCmd.AddCommand(flowlogCmd)
}

func newFlowlogClient() *flowlog.Client {
	return flowlog.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}
