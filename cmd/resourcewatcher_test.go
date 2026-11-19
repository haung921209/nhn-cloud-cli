package cmd

import "testing"

// The Resource Watcher API has no resource-group or resource-tag listing and
// no single-event read (neither the v1.0 nor the v2.0 guide); those commands
// only ever got 404 URL Not Found.
func TestResourceWatcherHasNoCommandsForMissingAPIs(t *testing.T) {
	for _, c := range resourceWatcherCmd.Commands() {
		if c.Name() == "describe-resource-groups" || c.Name() == "describe-resource-tags" || c.Name() == "describe-event" {
			t.Errorf("%s calls an API that does not exist", c.Name())
		}
	}
}
