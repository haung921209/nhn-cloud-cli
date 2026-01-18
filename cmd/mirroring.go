package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/mirroring"
	"github.com/spf13/cobra"
)

var mirroringCmd = &cobra.Command{
	Use:     "mirroring",
	Aliases: []string{"traffic-mirroring", "mirror"},
	Short:   "Manage Traffic Mirroring",
}

func init() {
	rootCmd.AddCommand(mirroringCmd)
}

func newMirroringClient() *mirroring.Client {
	return mirroring.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}
