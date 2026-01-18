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
	apigwCmd.AddCommand(apigwDescribeServicesCmd)
	apigwCmd.AddCommand(apigwCreateServiceCmd)
	apigwCmd.AddCommand(apigwUpdateServiceCmd)
	apigwCmd.AddCommand(apigwDeleteServiceCmd)

	apigwCreateServiceCmd.Flags().String("name", "", "Service name (required)")
	apigwCreateServiceCmd.Flags().String("description", "", "Service description")
	apigwCreateServiceCmd.MarkFlagRequired("name")

	apigwUpdateServiceCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwUpdateServiceCmd.Flags().String("name", "", "Service name")
	apigwUpdateServiceCmd.Flags().String("description", "", "Service description")
	apigwUpdateServiceCmd.MarkFlagRequired("service-id")

	apigwDeleteServiceCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwDeleteServiceCmd.MarkFlagRequired("service-id")
}

var apigwDescribeServicesCmd = &cobra.Command{
	Use:     "describe-services",
	Aliases: []string{"list-services"},
	Short:   "List or describe API Gateway services",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()

		// If args provided, could be get single, but standard list approach for now
		result, err := client.ListServices(ctx)
		if err != nil {
			exitWithError("Failed to list services", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tREGION\tCREATED")
		for _, s := range result.Services {
			created := ""
			if s.CreatedAt != nil {
				created = s.CreatedAt.Format("2006-01-02")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				s.ID, s.Name, s.TypeCode, s.RegionCode, created)
		}
		w.Flush()
	},
}

var apigwCreateServiceCmd = &cobra.Command{
	Use:   "create-service",
	Short: "Create a new API Gateway service",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &apigw.CreateServiceInput{
			Name:        name,
			Description: description,
		}

		result, err := client.CreateService(ctx, input)
		if err != nil {
			exitWithError("Failed to create service", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Service created: %s (%s)\n", result.Service.Name, result.Service.ID)
	},
}

var apigwUpdateServiceCmd = &cobra.Command{
	Use:   "update-service",
	Short: "Update an API Gateway service",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("service-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &apigw.UpdateServiceInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateService(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update service", err)
		}

		fmt.Printf("Service updated: %s\n", result.Service.ID)
	},
}

var apigwDeleteServiceCmd = &cobra.Command{
	Use:   "delete-service",
	Short: "Delete an API Gateway service",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("service-id")

		if err := client.DeleteService(ctx, id); err != nil {
			exitWithError("Failed to delete service", err)
		}

		fmt.Printf("Service %s deleted\n", id)
	},
}
