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
	appKey := getResourceWatcherAppKey()
	if appKey == "" {
		exitWithError("Resource Watcher appkey is not set: export NHN_CLOUD_RESOURCE_WATCHER_APPKEY or set app_key in the profile", nil)
	}
	return resourcewatcher.NewClient(appKey, getAccessKey(), getSecretKey(), nil, debug)
}
