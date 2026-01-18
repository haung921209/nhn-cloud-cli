package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/s3credential"
	"github.com/spf13/cobra"
)

var s3credentialCmd = &cobra.Command{
	Use:     "s3-credential",
	Aliases: []string{"s3-cred", "s3cred"},
	Short:   "Manage S3 API credentials",
}

func init() {
	rootCmd.AddCommand(s3credentialCmd)
}

func newS3CredentialClient() *s3credential.Client {
	return s3credential.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}
