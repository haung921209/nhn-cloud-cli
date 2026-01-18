package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/block"
	"github.com/spf13/cobra"
)

var blockStorageCmd = &cobra.Command{
	Use:     "block-storage",
	Aliases: []string{"volume", "bs"},
	Short:   "Manage Block Storage volumes and snapshots",
	Long:    `Manage block storage volumes, snapshots, and volume types.`,
}

func init() {
	rootCmd.AddCommand(blockStorageCmd)
}

func getBlockStorageClient() *block.Client {
	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	return block.NewClient(getRegion(), creds, nil, debug)
}
