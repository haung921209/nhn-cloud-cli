package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/flowlog"
	"github.com/spf13/cobra"
)

var flowlogCmd = &cobra.Command{
	Use:     "flowlog",
	Aliases: []string{"fl", "flow-log"},
	Short:   "Manage Flow Logs",
}

// Logger subcommand
var flowlogLoggerCmd = &cobra.Command{
	Use:     "logger",
	Aliases: []string{"loggers"},
	Short:   "Manage Flow Log loggers",
}

// Logging Port subcommand
var flowlogPortCmd = &cobra.Command{
	Use:     "port",
	Aliases: []string{"ports", "logging-port", "logging-ports"},
	Short:   "Manage Flow Log logging ports",
}

func init() {
	rootCmd.AddCommand(flowlogCmd)

	// Logger subcommands
	flowlogCmd.AddCommand(flowlogLoggerCmd)
	flowlogLoggerCmd.AddCommand(flowlogLoggerListCmd)
	flowlogLoggerCmd.AddCommand(flowlogLoggerGetCmd)
	flowlogLoggerCmd.AddCommand(flowlogLoggerCreateCmd)
	flowlogLoggerCmd.AddCommand(flowlogLoggerUpdateCmd)
	flowlogLoggerCmd.AddCommand(flowlogLoggerDeleteCmd)

	// Logging Port subcommands
	flowlogCmd.AddCommand(flowlogPortCmd)
	flowlogPortCmd.AddCommand(flowlogPortListCmd)
	flowlogPortCmd.AddCommand(flowlogPortGetCmd)

	// Create logger flags
	flowlogLoggerCreateCmd.Flags().String("name", "", "Logger name (required)")
	flowlogLoggerCreateCmd.Flags().String("description", "", "Description")
	flowlogLoggerCreateCmd.Flags().String("resource-type", "", "Resource type: VPC, SUBNET, PORT (required)")
	flowlogLoggerCreateCmd.Flags().String("resource-id", "", "Resource ID (required)")
	flowlogLoggerCreateCmd.Flags().String("filter-type", "ALL", "Filter type: ALL, ACCEPT, DROP")
	flowlogLoggerCreateCmd.Flags().String("connection-action", "", "Connection action: enable, disable")
	flowlogLoggerCreateCmd.Flags().String("storage-type", "OBS", "Storage type: OBS")
	flowlogLoggerCreateCmd.Flags().String("storage-url", "", "Object Storage URL (required)")
	flowlogLoggerCreateCmd.Flags().String("log-format", "CSV", "Log format: CSV, PARQUET")
	flowlogLoggerCreateCmd.Flags().String("compression-type", "RAW", "Compression type: RAW, GZIP")
	flowlogLoggerCreateCmd.Flags().String("partition-period", "HOUR", "Partition period: HOUR, DAY")
	flowlogLoggerCreateCmd.Flags().Bool("admin-state-up", true, "Admin state up")
	flowlogLoggerCreateCmd.MarkFlagRequired("name")
	flowlogLoggerCreateCmd.MarkFlagRequired("resource-type")
	flowlogLoggerCreateCmd.MarkFlagRequired("resource-id")
	flowlogLoggerCreateCmd.MarkFlagRequired("storage-url")

	// Update logger flags
	flowlogLoggerUpdateCmd.Flags().String("name", "", "Logger name")
	flowlogLoggerUpdateCmd.Flags().String("description", "", "Description")
	flowlogLoggerUpdateCmd.Flags().String("connection-action", "", "Connection action: enable, disable")
	flowlogLoggerUpdateCmd.Flags().Bool("admin-state-up", true, "Admin state up")
	flowlogLoggerUpdateCmd.Flags().Bool("set-admin-state", false, "Set admin state (use with --admin-state-up)")
}

func newFlowlogClient() *flowlog.Client {
	return flowlog.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

// ================================
// Logger Commands
// ================================

var flowlogLoggerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all flow log loggers",
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		result, err := client.ListLoggers(context.Background())
		if err != nil {
			exitWithError("Failed to list loggers", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var flowlogLoggerGetCmd = &cobra.Command{
	Use:   "get [logger-id]",
	Short: "Get flow log logger details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		result, err := client.GetLogger(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get logger", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var flowlogLoggerCreateCmd = &cobra.Command{
	Use:   "create",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Flow log logger created: %s\n", result.Logger.ID)
		fmt.Printf("Name: %s\n", result.Logger.Name)
		fmt.Printf("State: %s\n", result.Logger.State)
	},
}

var flowlogLoggerUpdateCmd = &cobra.Command{
	Use:   "update [logger-id]",
	Short: "Update a flow log logger",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
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

		result, err := client.UpdateLogger(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update logger", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Flow log logger updated: %s\n", result.Logger.ID)
		fmt.Printf("Name: %s\n", result.Logger.Name)
	},
}

var flowlogLoggerDeleteCmd = &cobra.Command{
	Use:   "delete [logger-id]",
	Short: "Delete a flow log logger",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		if err := client.DeleteLogger(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete logger", err)
		}
		fmt.Printf("Flow log logger %s deleted\n", args[0])
	},
}

// ================================
// Logging Port Commands
// ================================

var flowlogPortListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all flow log logging ports",
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		result, err := client.ListLoggingPorts(context.Background())
		if err != nil {
			exitWithError("Failed to list logging ports", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tLOGGER_ID\tPORT_ID\tSTATE")
		for _, port := range result.LoggingPorts {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				port.ID, port.LoggerID, port.PortID, port.State)
		}
		w.Flush()
	},
}

var flowlogPortGetCmd = &cobra.Command{
	Use:   "get [logging-port-id]",
	Short: "Get flow log logging port details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newFlowlogClient()
		result, err := client.GetLoggingPort(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get logging port", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		port := result.LoggingPort
		fmt.Printf("ID:         %s\n", port.ID)
		fmt.Printf("Logger ID:  %s\n", port.LoggerID)
		fmt.Printf("Port ID:    %s\n", port.PortID)
		fmt.Printf("State:      %s\n", port.State)
		fmt.Printf("Created At: %s\n", port.CreatedAt)
	},
}
