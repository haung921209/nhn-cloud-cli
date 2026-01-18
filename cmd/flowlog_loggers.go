package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/flowlog"
	"github.com/spf13/cobra"
)

func init() {
	flowlogCmd.AddCommand(flDescribeLoggersCmd)
	flowlogCmd.AddCommand(flGetLoggerCmd)
	flowlogCmd.AddCommand(flCreateLoggerCmd)
	flowlogCmd.AddCommand(flUpdateLoggerCmd)
	flowlogCmd.AddCommand(flDeleteLoggerCmd)

	flGetLoggerCmd.Flags().String("logger-id", "", "Logger ID (required)")
	flGetLoggerCmd.MarkFlagRequired("logger-id")

	flCreateLoggerCmd.Flags().String("name", "", "Logger name (required)")
	flCreateLoggerCmd.Flags().String("description", "", "Description")
	flCreateLoggerCmd.Flags().String("resource-type", "", "Resource type: VPC, SUBNET, PORT (required)")
	flCreateLoggerCmd.Flags().String("resource-id", "", "Resource ID (required)")
	flCreateLoggerCmd.Flags().String("filter-type", "ALL", "Filter type: ALL, ACCEPT, DROP")
	flCreateLoggerCmd.Flags().String("connection-action", "", "Connection action: enable, disable")
	flCreateLoggerCmd.Flags().String("storage-type", "OBS", "Storage type: OBS")
	flCreateLoggerCmd.Flags().String("storage-url", "", "Object Storage URL (required)")
	flCreateLoggerCmd.Flags().String("log-format", "CSV", "Log format: CSV, PARQUET")
	flCreateLoggerCmd.Flags().String("compression-type", "RAW", "Compression type: RAW, GZIP")
	flCreateLoggerCmd.Flags().String("partition-period", "HOUR", "Partition period: HOUR, DAY")
	flCreateLoggerCmd.Flags().Bool("admin-state-up", true, "Admin state up")
	flCreateLoggerCmd.MarkFlagRequired("name")
	flCreateLoggerCmd.MarkFlagRequired("resource-type")
	flCreateLoggerCmd.MarkFlagRequired("resource-id")
	flCreateLoggerCmd.MarkFlagRequired("storage-url")

	flUpdateLoggerCmd.Flags().String("logger-id", "", "Logger ID (required)")
	flUpdateLoggerCmd.Flags().String("name", "", "Logger name")
	flUpdateLoggerCmd.Flags().String("description", "", "Description")
	flUpdateLoggerCmd.Flags().String("connection-action", "", "Connection action: enable, disable")
	flUpdateLoggerCmd.Flags().Bool("admin-state-up", true, "Admin state up")
	flUpdateLoggerCmd.Flags().Bool("set-admin-state", false, "Set admin state (use with --admin-state-up)")
	flUpdateLoggerCmd.MarkFlagRequired("logger-id")

	flDeleteLoggerCmd.Flags().String("logger-id", "", "Logger ID (required)")
	flDeleteLoggerCmd.MarkFlagRequired("logger-id")
}

var flDescribeLoggersCmd = &cobra.Command{
	Use:     "describe-loggers",
	Aliases: []string{"list-loggers"},
	Short:   "List all flow log loggers",
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		result, err := client.ListLoggers(context.Background())
		if err != nil {
			exitWithError("Failed to list loggers", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tRESOURCE_TYPE\tFILTER_TYPE\tLOG_FORMAT\tSTATE")
		for _, logger := range result.Loggers {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				logger.ID, logger.Name, logger.ResourceType, logger.FilterType, logger.LogFormat, logger.State)
		}
		w.Flush()
	},
}

var flGetLoggerCmd = &cobra.Command{
	Use:     "describe-logger",
	Aliases: []string{"get-logger"},
	Short:   "Get flow log logger details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("logger-id")

		result, err := client.GetLogger(ctx, id)
		if err != nil {
			exitWithError("Failed to get logger", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		logger := result.Logger
		fmt.Printf("ID:                %s\n", logger.ID)
		fmt.Printf("Name:              %s\n", logger.Name)
		fmt.Printf("Description:       %s\n", logger.Description)
		fmt.Printf("Resource Type:     %s\n", logger.ResourceType)
		fmt.Printf("Resource ID:       %s\n", logger.ResourceID)
		fmt.Printf("Filter Type:       %s\n", logger.FilterType)
		fmt.Printf("Connection Action: %s\n", logger.ConnectionAction)
		fmt.Printf("Storage Type:      %s\n", logger.StorageType)
		fmt.Printf("Storage URL:       %s\n", logger.StorageURL)
		fmt.Printf("Log Format:        %s\n", logger.LogFormat)
		fmt.Printf("Compression Type:  %s\n", logger.CompressionType)
		fmt.Printf("Partition Period:  %s\n", logger.PartitionPeriod)
		fmt.Printf("Admin State Up:    %v\n", logger.AdminStateUp)
		fmt.Printf("State:             %s\n", logger.State)
		fmt.Printf("Created At:        %s\n", logger.CreatedAt)
		fmt.Printf("Updated At:        %s\n", logger.UpdatedAt)
	},
}

var flCreateLoggerCmd = &cobra.Command{
	Use:   "create-logger",
	Short: "Create a new flow log logger",
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		resourceType, _ := cmd.Flags().GetString("resource-type")
		resourceID, _ := cmd.Flags().GetString("resource-id")
		filterType, _ := cmd.Flags().GetString("filter-type")
		connectionAction, _ := cmd.Flags().GetString("connection-action")
		storageType, _ := cmd.Flags().GetString("storage-type")
		storageURL, _ := cmd.Flags().GetString("storage-url")
		logFormat, _ := cmd.Flags().GetString("log-format")
		compressionType, _ := cmd.Flags().GetString("compression-type")
		partitionPeriod, _ := cmd.Flags().GetString("partition-period")
		adminStateUp, _ := cmd.Flags().GetBool("admin-state-up")

		input := &flowlog.CreateLoggerInput{
			Name:             name,
			Description:      description,
			ResourceType:     resourceType,
			ResourceID:       resourceID,
			FilterType:       filterType,
			ConnectionAction: connectionAction,
			StorageType:      storageType,
			StorageURL:       storageURL,
			LogFormat:        logFormat,
			CompressionType:  compressionType,
			PartitionPeriod:  partitionPeriod,
			AdminStateUp:     adminStateUp,
		}

		result, err := client.CreateLogger(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create logger", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Flow log logger created: %s\n", result.Logger.ID)
		fmt.Printf("Name: %s\n", result.Logger.Name)
		fmt.Printf("State: %s\n", result.Logger.State)
	},
}

var flUpdateLoggerCmd = &cobra.Command{
	Use:   "update-logger",
	Short: "Update a flow log logger",
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		id, _ := cmd.Flags().GetString("logger-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		connectionAction, _ := cmd.Flags().GetString("connection-action")
		setAdminState, _ := cmd.Flags().GetBool("set-admin-state")
		adminStateUp, _ := cmd.Flags().GetBool("admin-state-up")

		input := &flowlog.UpdateLoggerInput{
			Name:             name,
			Description:      description,
			ConnectionAction: connectionAction,
		}

		if setAdminState {
			input.AdminStateUp = &adminStateUp
		}

		result, err := client.UpdateLogger(context.Background(), id, input)
		if err != nil {
			exitWithError("Failed to update logger", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Flow log logger updated: %s\n", result.Logger.ID)
		fmt.Printf("Name: %s\n", result.Logger.Name)
	},
}

var flDeleteLoggerCmd = &cobra.Command{
	Use:   "delete-logger",
	Short: "Delete a flow log logger",
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		id, _ := cmd.Flags().GetString("logger-id")
		if err := client.DeleteLogger(context.Background(), id); err != nil {
			exitWithError("Failed to delete logger", err)
		}
		fmt.Printf("Flow log logger %s deleted\n", id)
	},
}
