package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
	"github.com/spf13/cobra"
)

// resolveInstanceIdentifier resolves an instance identifier (name or ID) to an ID
func resolveInstanceIdentifier(client *mysql.Client, identifier string) (string, error) {
	// If identifier looks like a UUID, return it as-is
	if len(identifier) == 36 && identifier[8] == '-' && identifier[13] == '-' {
		return identifier, nil
	}

	// Otherwise, treat as name and look it up
	result, err := client.ListInstances(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to list instances: %w", err)
	}

	for _, inst := range result.DBInstances {
		if inst.DBInstanceName == identifier {
			return inst.DBInstanceID, nil
		}
	}

	return "", fmt.Errorf("instance not found: %s", identifier)
}

// getResolvedInstanceID is a helper that gets and resolves instance ID from command flags
func getResolvedInstanceID(cmd *cobra.Command, client *mysql.Client) (string, error) {
	identifier, _ := cmd.Flags().GetString("db-instance-identifier")
	if identifier == "" {
		return "", fmt.Errorf("--db-instance-identifier is required")
	}
	return resolveInstanceIdentifier(client, identifier)
}
