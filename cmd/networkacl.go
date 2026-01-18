package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/networkacl"
	"github.com/spf13/cobra"
)

var networkACLCmd = &cobra.Command{
	Use:     "network-acl",
	Aliases: []string{"acl", "nacl"},
	Short:   "Manage Network ACLs, rules, and subnet bindings",
}

func init() {
	rootCmd.AddCommand(networkACLCmd)
}

func newNetworkACLClient() *networkacl.Client {
	return networkacl.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}
