package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/nas"
	"github.com/spf13/cobra"
)

var nasCmd = &cobra.Command{
	Use:     "nas",
	Aliases: []string{"nas-storage"},
	Short:   "Manage NAS Storage volumes, snapshots, and interfaces",
}

func init() {
	rootCmd.AddCommand(nasCmd)
}

func newNASClient() *nas.Client {
	return nas.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}
