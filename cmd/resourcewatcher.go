package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/resourcewatcher"
	"github.com/spf13/cobra"
)

var resourceWatcherCmd = &cobra.Command{
	Use:     "resource-watcher",
	Aliases: []string{"rw", "watcher"},
	Short:   "Manage Resource Watcher (Governance) alarms and events",
	Long:    `Manage event alarms, alarm history, events, resource groups, and resource tags.`,
}

func init() {
	rootCmd.AddCommand(resourceWatcherCmd)
}

func getResourceWatcherClient() *resourcewatcher.Client {
	return resourcewatcher.NewClient(getResourceWatcherAppKey(), getAccessKey(), getSecretKey(), nil, debug)
}
