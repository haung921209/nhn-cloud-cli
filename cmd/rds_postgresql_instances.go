package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
	"github.com/spf13/cobra"
)

var createPostgreSQLInstanceCmd = &cobra.Command{
	Use:   "create-db-instance",
	Short: "Create a PostgreSQL DB instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()

		port, _ := cmd.Flags().GetInt("port")
		name, _ := cmd.Flags().GetString("db-instance-name")
		candidateName, _ := cmd.Flags().GetString("db-instance-candidate-name")
		dbName, _ := cmd.Flags().GetString("database-name")
		flavorID, _ := cmd.Flags().GetString("db-flavor-id")
		version, _ := cmd.Flags().GetString("db-version")
		username, _ := cmd.Flags().GetString("db-user-name")
		password, _ := cmd.Flags().GetString("db-password")
		description, _ := cmd.Flags().GetString("description")
		paramGroupID, _ := cmd.Flags().GetString("db-parameter-group-id")
		securityGroupIDs, _ := cmd.Flags().GetStringSlice("db-security-group-ids")
		userGroupIDs, _ := cmd.Flags().GetStringSlice("user-group-ids")
		notiGroupIDs, _ := cmd.Flags().GetStringSlice("notification-group-ids")

		subnetID, _ := cmd.Flags().GetString("subnet-id")
		az, _ := cmd.Flags().GetString("availability-zone")
		publicAccess, _ := cmd.Flags().GetBool("public-access")

		storageType, _ := cmd.Flags().GetString("storage-type")
		storageSize, _ := cmd.Flags().GetInt("allocated-storage")

		backupPeriod, _ := cmd.Flags().GetInt("backup-retention-period")

		backupStartTime, _ := cmd.Flags().GetString("backup-start-time")
		backupDuration, _ := cmd.Flags().GetString("backup-duration")

		if name == "" {
			exitWithError("--db-instance-name is required", nil)
		}
		if dbName == "" {
			exitWithError("--database-name is required", nil)
		}
		if flavorID == "" {
			exitWithError("--db-flavor-id is required", nil)
		}
		if version == "" {
			exitWithError("--db-version is required", nil)
		}
		if username == "" {
			exitWithError("--db-user-name is required", nil)
		}
		if password == "" {
			exitWithError("--db-password is required", nil)
		}
		if paramGroupID == "" {
			exitWithError("--db-parameter-group-id is required", nil)
		}
		if subnetID == "" {
			exitWithError("--subnet-id is required", nil)
		}

		req := &postgresql.CreateInstanceRequest{
			DBInstanceName:          name,
			DBInstanceCandidateName: candidateName,
			DatabaseName:            dbName,
			Description:             description,
			DBFlavorID:              flavorID,
			DBVersion:               version,
			DBUserName:              username,
			DBPassword:              password,
			DBPort:                  &port,
			ParameterGroupID:        paramGroupID,
			DBSecurityGroupIDs:      securityGroupIDs,
			UserGroupIDs:            userGroupIDs,
			NotificationGroupIDs:    notiGroupIDs,
			Network: postgresql.CreateInstanceNetworkConfig{
				SubnetID:         subnetID,
				AvailabilityZone: az,
				UsePublicAccess:  &publicAccess,
			},
			Storage: postgresql.CreateInstanceStorageConfig{
				StorageType: storageType,
				StorageSize: storageSize,
			},
			Backup: postgresql.CreateInstanceBackupConfig{
				BackupPeriod: backupPeriod,
				BackupSchedules: []postgresql.CreateInstanceBackupSchedule{
					{
						BackupWndBgnTime:  backupStartTime,
						BackupWndDuration: backupDuration,
					},
				},
			},
		}

		result, err := client.CreateInstance(context.Background(), req)
		if err != nil {
			exitWithError("failed to create instance", err)
		}

		fmt.Printf("Instance creation initiated.\n")
		if result.JobID != "" {
			fmt.Printf("Job ID: %s\n", result.JobID)
		}
	},
}

func init() {
	rdsPostgreSQLCmd.AddCommand(createPostgreSQLInstanceCmd)

	createPostgreSQLInstanceCmd.Flags().String("db-instance-name", "", "DB instance name (required)")
	createPostgreSQLInstanceCmd.Flags().String("db-instance-candidate-name", "", "Candidate instance name (HA only)")
	createPostgreSQLInstanceCmd.Flags().String("database-name", "", "Initial database name (required)")
	createPostgreSQLInstanceCmd.Flags().String("db-flavor-id", "", "Flavor ID (required)")
	createPostgreSQLInstanceCmd.Flags().String("db-version", "", "DB engine version (required)")
	createPostgreSQLInstanceCmd.Flags().String("db-user-name", "", "Master username (required)")
	createPostgreSQLInstanceCmd.Flags().String("db-password", "", "Master password (required)")
	createPostgreSQLInstanceCmd.Flags().String("description", "", "Description")
	createPostgreSQLInstanceCmd.Flags().String("db-parameter-group-id", "", "Parameter group ID (required)")
	createPostgreSQLInstanceCmd.Flags().StringSlice("db-security-group-ids", []string{}, "Security group IDs")
	createPostgreSQLInstanceCmd.Flags().StringSlice("user-group-ids", []string{}, "User group IDs")
	createPostgreSQLInstanceCmd.Flags().StringSlice("notification-group-ids", []string{}, "Notification group IDs")

	createPostgreSQLInstanceCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	createPostgreSQLInstanceCmd.Flags().String("availability-zone", "", "Availability Zone")
	createPostgreSQLInstanceCmd.Flags().Bool("public-access", false, "Enable public access")

	createPostgreSQLInstanceCmd.Flags().String("storage-type", "SSD", "Storage type")
	createPostgreSQLInstanceCmd.Flags().Int("allocated-storage", 20, "Storage size in GB")

	createPostgreSQLInstanceCmd.Flags().Int("port", 5432, "DB port (default 5432)")

	createPostgreSQLInstanceCmd.Flags().Int("backup-retention-period", 0, "Backup retention period (days)")
	createPostgreSQLInstanceCmd.Flags().String("backup-start-time", "01:00", "Backup window start time (HH:MM)")
	createPostgreSQLInstanceCmd.Flags().String("backup-duration", "ONE_HOUR", "Backup window duration (ONE_HOUR, TWO_HOURS, etc.)")
}
