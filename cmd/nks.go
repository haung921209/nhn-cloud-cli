package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/nks"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/spf13/cobra"
)

var nksCmd = &cobra.Command{
	Use:     "nks",
	Aliases: []string{"kubernetes"},
	Short:   "Manage NHN Kubernetes Service (NKS) clusters",
	Long:    `Manage NKS clusters including create, delete, scale node groups, and get kubeconfig.`,
}

func init() {
	rootCmd.AddCommand(nksCmd)
}

func getNKSClient() *nks.Client {
	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	return nks.NewClient(getRegion(), creds, nil, debug)
}
