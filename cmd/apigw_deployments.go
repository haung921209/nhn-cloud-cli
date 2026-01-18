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
	apigwCmd.AddCommand(apigwDescribeDeploymentsCmd)
	apigwCmd.AddCommand(apigwCreateDeploymentCmd)
	apigwCmd.AddCommand(apigwDeleteDeploymentCmd)
	apigwCmd.AddCommand(apigwRollbackDeploymentCmd)
	apigwCmd.AddCommand(apigwGetDeploymentCmd) // for 'latest' or specific

	apigwDescribeDeploymentsCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwDescribeDeploymentsCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwDescribeDeploymentsCmd.MarkFlagRequired("service-id")
	apigwDescribeDeploymentsCmd.MarkFlagRequired("stage-id")

	apigwCreateDeploymentCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwCreateDeploymentCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwCreateDeploymentCmd.Flags().String("description", "", "Deploy description")
	apigwCreateDeploymentCmd.MarkFlagRequired("service-id")
	apigwCreateDeploymentCmd.MarkFlagRequired("stage-id")

	apigwGetDeploymentCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwGetDeploymentCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwGetDeploymentCmd.Flags().Bool("latest", false, "Get latest deployment")
	apigwGetDeploymentCmd.MarkFlagRequired("service-id")
	apigwGetDeploymentCmd.MarkFlagRequired("stage-id")

	apigwDeleteDeploymentCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwDeleteDeploymentCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwDeleteDeploymentCmd.Flags().String("deployment-id", "", "Deployment ID (required)")
	apigwDeleteDeploymentCmd.MarkFlagRequired("service-id")
	apigwDeleteDeploymentCmd.MarkFlagRequired("stage-id")
	apigwDeleteDeploymentCmd.MarkFlagRequired("deployment-id")

	apigwRollbackDeploymentCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwRollbackDeploymentCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwRollbackDeploymentCmd.Flags().String("deployment-id", "", "Deployment ID to rollback TO (required)")
	apigwRollbackDeploymentCmd.MarkFlagRequired("service-id")
	apigwRollbackDeploymentCmd.MarkFlagRequired("stage-id")
	apigwRollbackDeploymentCmd.MarkFlagRequired("deployment-id")
}

var apigwDescribeDeploymentsCmd = &cobra.Command{
	Use:     "describe-deployments",
	Aliases: []string{"list-deployments", "list-deploys"},
	Short:   "List deployments for a stage",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")

		result, err := client.ListDeploys(ctx, serviceID, stageID)
		if err != nil {
			exitWithError("Failed to list deployments", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATUS\tDESCRIPTION\tCREATED")
		for _, d := range result.Deploys {
			created := ""
			if d.CreatedAt != nil {
				created = d.CreatedAt.Format("2006-01-02 15:04:05")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", d.ID, d.StatusCode, d.Description, created)
		}
		w.Flush()
	},
}

var apigwGetDeploymentCmd = &cobra.Command{
	Use:   "get-deployment",
	Short: "Get details of a deployment",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")
		latest, _ := cmd.Flags().GetBool("latest")

		if latest {
			result, err := client.GetLatestDeploy(ctx, serviceID, stageID)
			if err != nil {
				exitWithError("Failed to get latest deployment", err)
			}
			printDeploymentDetails(result.Deploy)
		} else {
			// Currently SDK might not have GetDeploy by ID separate from List?
			// Looking at original apigw.go, `apigwDeployLatestCmd` calls `GetLatestDeploy`.
			// There doesn't seem to be a `GetDeploy(id)` in the original file, only `GetLatestDeploy`.
			// So `get-deployment` serves mostly for latest or would need SDK update for specific ID get if supported.
			// Let's support `latest` primarily for now.
			fmt.Println("Please use --latest flag to get the latest deployment.")
		}
	},
}

func printDeploymentDetails(d apigw.Deploy) {
	if output == "json" {
		printJSON(d)
		return
	}
	fmt.Printf("ID:          %s\n", d.ID)
	fmt.Printf("Status:      %s\n", d.StatusCode)
	fmt.Printf("Description: %s\n", d.Description)
	if d.CreatedAt != nil {
		fmt.Printf("Created:     %s\n", d.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

var apigwCreateDeploymentCmd = &cobra.Command{
	Use:   "create-deployment",
	Short: "Create a new deployment",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")
		description, _ := cmd.Flags().GetString("description")

		input := &apigw.CreateDeployInput{
			Description: description,
		}

		result, err := client.DeployStage(ctx, serviceID, stageID, input)
		if err != nil {
			exitWithError("Failed to create deployment", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Deployment created: %s\n", result.Deploy.ID)
		fmt.Printf("Status: %s\n", result.Deploy.StatusCode)
	},
}

var apigwDeleteDeploymentCmd = &cobra.Command{
	Use:   "delete-deployment",
	Short: "Delete a deployment",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")
		deployID, _ := cmd.Flags().GetString("deployment-id")

		if err := client.DeleteDeploy(ctx, serviceID, stageID, deployID); err != nil {
			exitWithError("Failed to delete deployment", err)
		}

		fmt.Printf("Deployment %s deleted\n", deployID)
	},
}

var apigwRollbackDeploymentCmd = &cobra.Command{
	Use:   "rollback-deployment",
	Short: "Rollback to a specific deployment",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")
		deployID, _ := cmd.Flags().GetString("deployment-id")

		result, err := client.RollbackDeploy(ctx, serviceID, stageID, deployID)
		if err != nil {
			exitWithError("Failed to rollback deployment", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Rolled back to deployment: %s\n", result.Deploy.ID)
	},
}
