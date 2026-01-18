package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/cloudtrail"
	"github.com/spf13/cobra"
)

var cloudtrailCmd = &cobra.Command{
	Use:     "cloudtrail",
	Aliases: []string{"trail", "audit"},
	Short:   "Manage CloudTrail audit events",
	Long:    `View and search CloudTrail audit events for your NHN Cloud account.`,
}

func init() {
	rootCmd.AddCommand(cloudtrailCmd)
}

func newCloudTrailClient() *cloudtrail.Client {
	return cloudtrail.NewClient(getAppKey(), getAccessKey(), getSecretKey(), nil, debug)
}
