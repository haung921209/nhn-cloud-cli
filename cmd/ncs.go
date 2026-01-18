package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncs"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/spf13/cobra"
)

var ncsCmd = &cobra.Command{
	Use:     "ncs",
	Aliases: []string{"container-service"},
	Short:   "Manage NHN Container Service (NCS)",
	Long:    `Manage serverless container workloads and services.`,
}

func init() {
	rootCmd.AddCommand(ncsCmd)
}

func getNCSClient() *ncs.Client {
	creds := credentials.NewStatic(getAccessKey(), getSecretKey())
	return ncs.NewClient(getRegion(), getNCSAppKey(), creds, nil, debug)
}
