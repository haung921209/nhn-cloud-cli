package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/apigw"
	"github.com/spf13/cobra"
)

func init() {
	apigwCmd.AddCommand(apigwDescribeStagesCmd)
	apigwCmd.AddCommand(apigwCreateStageCmd)
	apigwCmd.AddCommand(apigwUpdateStageCmd)
	apigwCmd.AddCommand(apigwDeleteStageCmd)

	apigwDescribeStagesCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwDescribeStagesCmd.MarkFlagRequired("service-id")

	apigwCreateStageCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwCreateStageCmd.Flags().String("name", "", "Stage name (required)")
	apigwCreateStageCmd.Flags().String("description", "", "Stage description")
	apigwCreateStageCmd.Flags().String("backend-url", "", "Backend endpoint URL")
	apigwCreateStageCmd.MarkFlagRequired("service-id")
	apigwCreateStageCmd.MarkFlagRequired("name")

	apigwUpdateStageCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwUpdateStageCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwUpdateStageCmd.Flags().String("name", "", "Stage name")
	apigwUpdateStageCmd.Flags().String("description", "", "Stage description")
	apigwUpdateStageCmd.Flags().String("backend-url", "", "Backend endpoint URL")
	apigwUpdateStageCmd.MarkFlagRequired("service-id")
	apigwUpdateStageCmd.MarkFlagRequired("stage-id")

	apigwDeleteStageCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwDeleteStageCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwDeleteStageCmd.MarkFlagRequired("service-id")
	apigwDeleteStageCmd.MarkFlagRequired("stage-id")
}

var apigwDescribeStagesCmd = &cobra.Command{
	Use:     "describe-stages",
	Aliases: []string{"list-stages"},
	Short:   "List stages for a service",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		serviceID, _ := cmd.Flags().GetString("service-id")

		result, err := client.ListStages(ctx, serviceID)
		if err != nil {
			exitWithError("Failed to list stages", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tURL\tBACKEND")
		for _, s := range result.Stages {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.ID, s.Name, s.URL, s.BackendEndpointURL)
		}
		w.Flush()
	},
}

var apigwCreateStageCmd = &cobra.Command{
	Use:   "create-stage",
	Short: "Create a new stage",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		serviceID, _ := cmd.Flags().GetString("service-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		backendURL, _ := cmd.Flags().GetString("backend-url")

		input := &apigw.CreateStageInput{
			Name:               name,
			Description:        description,
			BackendEndpointURL: backendURL,
		}

		result, err := client.CreateStage(ctx, serviceID, input)
		if err != nil {
			exitWithError("Failed to create stage", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Stage created: %s (%s)\n", result.Stage.Name, result.Stage.ID)
	},
}

var apigwUpdateStageCmd = &cobra.Command{
	Use:   "update-stage",
	Short: "Update a stage",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		backendURL, _ := cmd.Flags().GetString("backend-url")

		input := &apigw.UpdateStageInput{
			Name:               name,
			Description:        description,
			BackendEndpointURL: backendURL,
		}

		result, err := client.UpdateStage(ctx, serviceID, stageID, input)
		if err != nil {
			exitWithError("Failed to update stage", err)
		}

		fmt.Printf("Stage updated: %s\n", result.Stage.ID)
	},
}

var apigwDeleteStageCmd = &cobra.Command{
	Use:   "delete-stage",
	Short: "Delete a stage",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")

		if err := client.DeleteStage(ctx, serviceID, stageID); err != nil {
			exitWithError("Failed to delete stage", err)
		}

		fmt.Printf("Stage %s deleted\n", stageID)
	},
}
