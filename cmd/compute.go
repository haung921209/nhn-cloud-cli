package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/compute"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/spf13/cobra"
)

var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Manage Compute instances (VMs)",
	Long:  `Manage Compute instances including create, delete, start, stop, and more.`,
}

func init() {
	rootCmd.AddCommand(computeCmd)
}

func getComputeClient() *compute.Client {
	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	return compute.NewClient(getRegion(), creds, nil, debug)
}
