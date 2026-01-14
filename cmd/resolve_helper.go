package cmd

import (
"context"
"fmt"

"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/mysql"
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
