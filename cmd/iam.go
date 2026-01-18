package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/iam"
	"github.com/spf13/cobra"
)

var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Manage IAM organizations, projects, and members",
	Long:  `Manage Identity and Access Management resources including organizations, projects, and members.`,
}

func init() {
	rootCmd.AddCommand(iamCmd)
}

func getIAMClient() *iam.Client {
	creds := credentials.NewStatic(getAccessKey(), getSecretKey())
	return iam.NewClient(getRegion(), creds, nil, debug)
}
