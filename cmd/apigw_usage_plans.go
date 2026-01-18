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
	apigwCmd.AddCommand(apigwDescribeUsagePlansCmd)
	apigwCmd.AddCommand(apigwCreateUsagePlanCmd)
	apigwCmd.AddCommand(apigwUpdateUsagePlanCmd)
	apigwCmd.AddCommand(apigwDeleteUsagePlanCmd)
	apigwCmd.AddCommand(apigwGetUsagePlanCmd)

	apigwCreateUsagePlanCmd.Flags().String("name", "", "Usage plan name (required)")
	apigwCreateUsagePlanCmd.Flags().String("description", "", "Usage plan description")
	apigwCreateUsagePlanCmd.Flags().Int("rate-limit", 0, "Rate limit (requests per second)")
	apigwCreateUsagePlanCmd.Flags().Int("quota-limit", 0, "Quota limit (request count)")
	apigwCreateUsagePlanCmd.Flags().String("quota-period", "", "Quota period: DAY or MONTH")
	apigwCreateUsagePlanCmd.MarkFlagRequired("name")

	apigwUpdateUsagePlanCmd.Flags().String("plan-id", "", "Usage plan ID (required)")
	apigwUpdateUsagePlanCmd.Flags().String("name", "", "Usage plan name")
	apigwUpdateUsagePlanCmd.Flags().String("description", "", "Usage plan description")
	apigwUpdateUsagePlanCmd.Flags().Int("rate-limit", 0, "Rate limit")
	apigwUpdateUsagePlanCmd.Flags().Int("quota-limit", 0, "Quota limit")
	apigwUpdateUsagePlanCmd.Flags().String("quota-period", "", "Quota period")
	apigwUpdateUsagePlanCmd.MarkFlagRequired("plan-id")

	apigwDeleteUsagePlanCmd.Flags().String("plan-id", "", "Usage plan ID (required)")
	apigwDeleteUsagePlanCmd.MarkFlagRequired("plan-id")

	apigwGetUsagePlanCmd.Flags().String("plan-id", "", "Usage plan ID (required)")
	apigwGetUsagePlanCmd.MarkFlagRequired("plan-id")
}

var apigwDescribeUsagePlansCmd = &cobra.Command{
	Use:     "describe-usage-plans",
	Aliases: []string{"list-usage-plans"},
	Short:   "List all usage plans",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()

		result, err := client.ListUsagePlans(ctx)
		if err != nil {
			exitWithError("Failed to list usage plans", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tRATE_LIMIT\tQUOTA_LIMIT\tQUOTA_PERIOD")
		for _, p := range result.UsagePlans {
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\n", p.ID, p.Name, p.RateLimitRequestPerSecond, p.QuotaLimitRequestCount, p.QuotaPeriodUnitCode)
		}
		w.Flush()
	},
}

var apigwGetUsagePlanCmd = &cobra.Command{
	Use:   "get-usage-plan",
	Short: "Get usage plan details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("plan-id")

		result, err := client.GetUsagePlan(ctx, id)
		if err != nil {
			exitWithError("Failed to get usage plan", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		p := result.UsagePlan
		fmt.Printf("ID:           %s\n", p.ID)
		fmt.Printf("Name:         %s\n", p.Name)
		fmt.Printf("Description:  %s\n", p.Description)
		fmt.Printf("Rate Limit:   %d req/s\n", p.RateLimitRequestPerSecond)
		fmt.Printf("Quota Limit:  %d\n", p.QuotaLimitRequestCount)
		fmt.Printf("Quota Period: %s\n", p.QuotaPeriodUnitCode)
		if p.CreatedAt != nil {
			fmt.Printf("Created:      %s\n", p.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	},
}

var apigwCreateUsagePlanCmd = &cobra.Command{
	Use:   "create-usage-plan",
	Short: "Create a new usage plan",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		rateLimit, _ := cmd.Flags().GetInt("rate-limit")
		quotaLimit, _ := cmd.Flags().GetInt("quota-limit")
		quotaPeriod, _ := cmd.Flags().GetString("quota-period")

		input := &apigw.CreateUsagePlanInput{
			Name:                      name,
			Description:               description,
			RateLimitRequestPerSecond: rateLimit,
			QuotaLimitRequestCount:    quotaLimit,
			QuotaPeriodUnitCode:       quotaPeriod,
		}

		result, err := client.CreateUsagePlan(ctx, input)
		if err != nil {
			exitWithError("Failed to create usage plan", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Usage Plan created: %s\n", result.UsagePlan.ID)
		fmt.Printf("Name: %s\n", result.UsagePlan.Name)
	},
}

var apigwUpdateUsagePlanCmd = &cobra.Command{
	Use:   "update-usage-plan",
	Short: "Update a usage plan",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("plan-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		rateLimit, _ := cmd.Flags().GetInt("rate-limit")
		quotaLimit, _ := cmd.Flags().GetInt("quota-limit")
		quotaPeriod, _ := cmd.Flags().GetString("quota-period")

		input := &apigw.UpdateUsagePlanInput{
			Name:                      name,
			Description:               description,
			RateLimitRequestPerSecond: rateLimit,
			QuotaLimitRequestCount:    quotaLimit,
			QuotaPeriodUnitCode:       quotaPeriod,
		}

		result, err := client.UpdateUsagePlan(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update usage plan", err)
		}

		fmt.Printf("Usage Plan updated: %s\n", result.UsagePlan.ID)
	},
}

var apigwDeleteUsagePlanCmd = &cobra.Command{
	Use:   "delete-usage-plan",
	Short: "Delete a usage plan",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("plan-id")

		if err := client.DeleteUsagePlan(ctx, id); err != nil {
			exitWithError("Failed to delete usage plan", err)
		}

		fmt.Printf("Usage Plan %s deleted\n", id)
	},
}
