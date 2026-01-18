package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/image"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"images", "img"},
	Short:   "Manage Glance images",
	Long:    `Manage images including list, get, create, delete, and tag management.`,
}

func init() {
	rootCmd.AddCommand(imageCmd)
}

func getImageClient() *image.Client {
	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	return image.NewClient(getRegion(), creds, nil, debug)
}
