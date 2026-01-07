package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/apigw"
	"github.com/spf13/cobra"
)

var apigwCmd = &cobra.Command{
	Use:     "apigw",
	Aliases: []string{"api-gateway", "apigateway"},
	Short:   "Manage API Gateway services, stages, deployments, and API keys",
}

var apigwServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services", "svc"},
	Short:   "Manage API Gateway services",
}

var apigwStageCmd = &cobra.Command{
	Use:     "stage",
	Aliases: []string{"stages"},
	Short:   "Manage deployment stages",
}

var apigwDeployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"deploys", "deployment"},
	Short:   "Manage stage deployments",
}

var apigwApikeyCmd = &cobra.Command{
	Use:     "apikey",
	Aliases: []string{"apikeys", "key"},
	Short:   "Manage API keys",
}

var apigwUsagePlanCmd = &cobra.Command{
	Use:     "usage-plan",
	Aliases: []string{"usage-plans", "plan"},
	Short:   "Manage usage plans",
}

func init() {
	rootCmd.AddCommand(apigwCmd)

	// Service commands
	apigwCmd.AddCommand(apigwServiceCmd)
	apigwServiceCmd.AddCommand(apigwServiceListCmd)
	apigwServiceCmd.AddCommand(apigwServiceGetCmd)
	apigwServiceCmd.AddCommand(apigwServiceCreateCmd)
	apigwServiceCmd.AddCommand(apigwServiceUpdateCmd)
	apigwServiceCmd.AddCommand(apigwServiceDeleteCmd)

	// Stage commands
	apigwCmd.AddCommand(apigwStageCmd)
	apigwStageCmd.AddCommand(apigwStageListCmd)
	apigwStageCmd.AddCommand(apigwStageCreateCmd)
	apigwStageCmd.AddCommand(apigwStageUpdateCmd)
	apigwStageCmd.AddCommand(apigwStageDeleteCmd)

	// Deploy commands
	apigwCmd.AddCommand(apigwDeployCmd)
	apigwDeployCmd.AddCommand(apigwDeployListCmd)
	apigwDeployCmd.AddCommand(apigwDeployCreateCmd)
	apigwDeployCmd.AddCommand(apigwDeployLatestCmd)
	apigwDeployCmd.AddCommand(apigwDeployDeleteCmd)
	apigwDeployCmd.AddCommand(apigwDeployRollbackCmd)

	// API Key commands
	apigwCmd.AddCommand(apigwApikeyCmd)
	apigwApikeyCmd.AddCommand(apigwApikeyListCmd)
	apigwApikeyCmd.AddCommand(apigwApikeyCreateCmd)
	apigwApikeyCmd.AddCommand(apigwApikeyUpdateCmd)
	apigwApikeyCmd.AddCommand(apigwApikeyDeleteCmd)
	apigwApikeyCmd.AddCommand(apigwApikeyRegenerateCmd)

	// Usage Plan commands
	apigwCmd.AddCommand(apigwUsagePlanCmd)
	apigwUsagePlanCmd.AddCommand(apigwUsagePlanListCmd)
	apigwUsagePlanCmd.AddCommand(apigwUsagePlanGetCmd)
	apigwUsagePlanCmd.AddCommand(apigwUsagePlanCreateCmd)
	apigwUsagePlanCmd.AddCommand(apigwUsagePlanUpdateCmd)
	apigwUsagePlanCmd.AddCommand(apigwUsagePlanDeleteCmd)

	// Service flags
	apigwServiceCreateCmd.Flags().String("name", "", "Service name (required)")
	apigwServiceCreateCmd.Flags().String("description", "", "Service description")
	apigwServiceCreateCmd.MarkFlagRequired("name")

	apigwServiceUpdateCmd.Flags().String("name", "", "Service name")
	apigwServiceUpdateCmd.Flags().String("description", "", "Service description")

	// Stage flags
	apigwStageListCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwStageListCmd.MarkFlagRequired("service-id")

	apigwStageCreateCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwStageCreateCmd.Flags().String("name", "", "Stage name (required)")
	apigwStageCreateCmd.Flags().String("description", "", "Stage description")
	apigwStageCreateCmd.Flags().String("backend-url", "", "Backend endpoint URL")
	apigwStageCreateCmd.MarkFlagRequired("service-id")
	apigwStageCreateCmd.MarkFlagRequired("name")

	apigwStageUpdateCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwStageUpdateCmd.Flags().String("name", "", "Stage name")
	apigwStageUpdateCmd.Flags().String("description", "", "Stage description")
	apigwStageUpdateCmd.Flags().String("backend-url", "", "Backend endpoint URL")
	apigwStageUpdateCmd.MarkFlagRequired("service-id")

	apigwStageDeleteCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwStageDeleteCmd.MarkFlagRequired("service-id")

	// Deploy flags
	apigwDeployListCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwDeployListCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwDeployListCmd.MarkFlagRequired("service-id")
	apigwDeployListCmd.MarkFlagRequired("stage-id")

	apigwDeployCreateCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwDeployCreateCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwDeployCreateCmd.Flags().String("description", "", "Deploy description")
	apigwDeployCreateCmd.MarkFlagRequired("service-id")
	apigwDeployCreateCmd.MarkFlagRequired("stage-id")

	apigwDeployLatestCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwDeployLatestCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwDeployLatestCmd.MarkFlagRequired("service-id")
	apigwDeployLatestCmd.MarkFlagRequired("stage-id")

	apigwDeployDeleteCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwDeployDeleteCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwDeployDeleteCmd.MarkFlagRequired("service-id")
	apigwDeployDeleteCmd.MarkFlagRequired("stage-id")

	apigwDeployRollbackCmd.Flags().String("service-id", "", "Service ID (required)")
	apigwDeployRollbackCmd.Flags().String("stage-id", "", "Stage ID (required)")
	apigwDeployRollbackCmd.MarkFlagRequired("service-id")
	apigwDeployRollbackCmd.MarkFlagRequired("stage-id")

	// API Key flags
	apigwApikeyCreateCmd.Flags().String("name", "", "API key name (required)")
	apigwApikeyCreateCmd.Flags().String("description", "", "API key description")
	apigwApikeyCreateCmd.MarkFlagRequired("name")

	apigwApikeyUpdateCmd.Flags().String("name", "", "API key name")
	apigwApikeyUpdateCmd.Flags().String("description", "", "API key description")
	apigwApikeyUpdateCmd.Flags().String("status", "", "API key status")

	apigwApikeyRegenerateCmd.Flags().String("key-type", "PRIMARY", "Key type: PRIMARY or SECONDARY")

	// Usage Plan flags
	apigwUsagePlanCreateCmd.Flags().String("name", "", "Usage plan name (required)")
	apigwUsagePlanCreateCmd.Flags().String("description", "", "Usage plan description")
	apigwUsagePlanCreateCmd.Flags().Int("rate-limit", 0, "Rate limit (requests per second)")
	apigwUsagePlanCreateCmd.Flags().Int("quota-limit", 0, "Quota limit (request count)")
	apigwUsagePlanCreateCmd.Flags().String("quota-period", "", "Quota period: DAY or MONTH")
	apigwUsagePlanCreateCmd.MarkFlagRequired("name")

	apigwUsagePlanUpdateCmd.Flags().String("name", "", "Usage plan name")
	apigwUsagePlanUpdateCmd.Flags().String("description", "", "Usage plan description")
	apigwUsagePlanUpdateCmd.Flags().Int("rate-limit", 0, "Rate limit (requests per second)")
	apigwUsagePlanUpdateCmd.Flags().Int("quota-limit", 0, "Quota limit (request count)")
	apigwUsagePlanUpdateCmd.Flags().String("quota-period", "", "Quota period: DAY or MONTH")
}

func newAPIGWClient() *apigw.Client {
	appKey := getAppKey()
	accessKey := getAccessKey()
	secretKey := getSecretKey()
	return apigw.NewClient(getRegion(), appKey, accessKey, secretKey, nil, debug)
}

// ================================
// Service Commands
// ================================

var apigwServiceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API Gateway services",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		result, err := client.ListServices(context.Background())
		if err != nil {
			exitWithError("Failed to list services", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tREGION")
		for _, s := range result.Services {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.ID, s.Name, s.TypeCode, s.RegionCode)
		}
		w.Flush()
	},
}

var apigwServiceGetCmd = &cobra.Command{
	Use:   "get [service-id]",
	Short: "Get service details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		result, err := client.GetService(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get service", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		s := result.Service
		fmt.Printf("ID:          %s\n", s.ID)
		fmt.Printf("Name:        %s\n", s.Name)
		fmt.Printf("Description: %s\n", s.Description)
		fmt.Printf("Type:        %s\n", s.TypeCode)
		fmt.Printf("AppKey:      %s\n", s.AppKey)
		fmt.Printf("Region:      %s\n", s.RegionCode)
		if s.CreatedAt != nil {
			fmt.Printf("Created:     %s\n", s.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	},
}

var apigwServiceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API Gateway service",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &apigw.CreateServiceInput{
			Name:        name,
			Description: description,
		}

		result, err := client.CreateService(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create service", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Service created: %s\n", result.Service.ID)
		fmt.Printf("Name: %s\n", result.Service.Name)
	},
}

var apigwServiceUpdateCmd = &cobra.Command{
	Use:   "update [service-id]",
	Short: "Update a service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &apigw.UpdateServiceInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateService(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update service", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Service updated: %s\n", result.Service.ID)
	},
}

var apigwServiceDeleteCmd = &cobra.Command{
	Use:   "delete [service-id]",
	Short: "Delete a service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		if err := client.DeleteService(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete service", err)
		}
		fmt.Printf("Service %s deleted\n", args[0])
	},
}

// ================================
// Stage Commands
// ================================

var apigwStageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List stages for a service",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		serviceID, _ := cmd.Flags().GetString("service-id")

		result, err := client.ListStages(context.Background(), serviceID)
		if err != nil {
			exitWithError("Failed to list stages", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var apigwStageCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new stage",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		serviceID, _ := cmd.Flags().GetString("service-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		backendURL, _ := cmd.Flags().GetString("backend-url")

		input := &apigw.CreateStageInput{
			Name:               name,
			Description:        description,
			BackendEndpointURL: backendURL,
		}

		result, err := client.CreateStage(context.Background(), serviceID, input)
		if err != nil {
			exitWithError("Failed to create stage", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Stage created: %s\n", result.Stage.ID)
		fmt.Printf("Name: %s\n", result.Stage.Name)
		fmt.Printf("URL: %s\n", result.Stage.URL)
	},
}

var apigwStageUpdateCmd = &cobra.Command{
	Use:   "update [stage-id]",
	Short: "Update a stage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		serviceID, _ := cmd.Flags().GetString("service-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		backendURL, _ := cmd.Flags().GetString("backend-url")

		input := &apigw.UpdateStageInput{
			Name:               name,
			Description:        description,
			BackendEndpointURL: backendURL,
		}

		result, err := client.UpdateStage(context.Background(), serviceID, args[0], input)
		if err != nil {
			exitWithError("Failed to update stage", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Stage updated: %s\n", result.Stage.ID)
	},
}

var apigwStageDeleteCmd = &cobra.Command{
	Use:   "delete [stage-id]",
	Short: "Delete a stage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		serviceID, _ := cmd.Flags().GetString("service-id")

		if err := client.DeleteStage(context.Background(), serviceID, args[0]); err != nil {
			exitWithError("Failed to delete stage", err)
		}
		fmt.Printf("Stage %s deleted\n", args[0])
	},
}

// ================================
// Deploy Commands
// ================================

var apigwDeployListCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployments for a stage",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")

		result, err := client.ListDeploys(context.Background(), serviceID, stageID)
		if err != nil {
			exitWithError("Failed to list deploys", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var apigwDeployCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Deploy a stage",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")
		description, _ := cmd.Flags().GetString("description")

		input := &apigw.CreateDeployInput{
			Description: description,
		}

		result, err := client.DeployStage(context.Background(), serviceID, stageID, input)
		if err != nil {
			exitWithError("Failed to deploy stage", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Deployment created: %s\n", result.Deploy.ID)
		fmt.Printf("Status: %s\n", result.Deploy.StatusCode)
	},
}

var apigwDeployLatestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Get the latest deployment for a stage",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")

		result, err := client.GetLatestDeploy(context.Background(), serviceID, stageID)
		if err != nil {
			exitWithError("Failed to get latest deploy", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		d := result.Deploy
		fmt.Printf("ID:          %s\n", d.ID)
		fmt.Printf("Status:      %s\n", d.StatusCode)
		fmt.Printf("Description: %s\n", d.Description)
		if d.CreatedAt != nil {
			fmt.Printf("Created:     %s\n", d.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	},
}

var apigwDeployDeleteCmd = &cobra.Command{
	Use:   "delete [deploy-id]",
	Short: "Delete a deployment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")

		if err := client.DeleteDeploy(context.Background(), serviceID, stageID, args[0]); err != nil {
			exitWithError("Failed to delete deploy", err)
		}
		fmt.Printf("Deployment %s deleted\n", args[0])
	},
}

var apigwDeployRollbackCmd = &cobra.Command{
	Use:   "rollback [deploy-id]",
	Short: "Rollback to a specific deployment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		serviceID, _ := cmd.Flags().GetString("service-id")
		stageID, _ := cmd.Flags().GetString("stage-id")

		result, err := client.RollbackDeploy(context.Background(), serviceID, stageID, args[0])
		if err != nil {
			exitWithError("Failed to rollback deploy", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Rolled back to deployment: %s\n", result.Deploy.ID)
	},
}

// ================================
// API Key Commands
// ================================

var apigwApikeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		result, err := client.ListAPIKeys(context.Background())
		if err != nil {
			exitWithError("Failed to list API keys", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tCREATED")
		for _, k := range result.APIKeys {
			created := ""
			if k.CreatedAt != nil {
				created = k.CreatedAt.Format("2006-01-02 15:04:05")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", k.ID, k.Name, k.StatusCode, created)
		}
		w.Flush()
	},
}

var apigwApikeyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &apigw.CreateAPIKeyInput{
			Name:        name,
			Description: description,
		}

		result, err := client.CreateAPIKey(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create API key", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("API Key created: %s\n", result.APIKey.ID)
		fmt.Printf("Name: %s\n", result.APIKey.Name)
		fmt.Printf("Primary Key: %s\n", result.APIKey.PrimaryKey)
		fmt.Printf("Secondary Key: %s\n", result.APIKey.SecondaryKey)
	},
}

var apigwApikeyUpdateCmd = &cobra.Command{
	Use:   "update [apikey-id]",
	Short: "Update an API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")

		input := &apigw.UpdateAPIKeyInput{
			Name:        name,
			Description: description,
			StatusCode:  status,
		}

		result, err := client.UpdateAPIKey(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update API key", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("API Key updated: %s\n", result.APIKey.ID)
	},
}

var apigwApikeyDeleteCmd = &cobra.Command{
	Use:   "delete [apikey-id]",
	Short: "Delete an API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		if err := client.DeleteAPIKey(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete API key", err)
		}
		fmt.Printf("API Key %s deleted\n", args[0])
	},
}

var apigwApikeyRegenerateCmd = &cobra.Command{
	Use:   "regenerate [apikey-id]",
	Short: "Regenerate an API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		keyType, _ := cmd.Flags().GetString("key-type")

		input := &apigw.RegenerateAPIKeyInput{
			KeyType: keyType,
		}

		result, err := client.RegenerateAPIKey(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to regenerate API key", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("API Key regenerated: %s\n", result.APIKey.ID)
		fmt.Printf("Primary Key: %s\n", result.APIKey.PrimaryKey)
		fmt.Printf("Secondary Key: %s\n", result.APIKey.SecondaryKey)
	},
}

// ================================
// Usage Plan Commands
// ================================

var apigwUsagePlanListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all usage plans",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		result, err := client.ListUsagePlans(context.Background())
		if err != nil {
			exitWithError("Failed to list usage plans", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var apigwUsagePlanGetCmd = &cobra.Command{
	Use:   "get [usage-plan-id]",
	Short: "Get usage plan details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		result, err := client.GetUsagePlan(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get usage plan", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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

var apigwUsagePlanCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new usage plan",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
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

		result, err := client.CreateUsagePlan(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create usage plan", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Usage Plan created: %s\n", result.UsagePlan.ID)
		fmt.Printf("Name: %s\n", result.UsagePlan.Name)
	},
}

var apigwUsagePlanUpdateCmd = &cobra.Command{
	Use:   "update [usage-plan-id]",
	Short: "Update a usage plan",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
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

		result, err := client.UpdateUsagePlan(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update usage plan", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Usage Plan updated: %s\n", result.UsagePlan.ID)
	},
}

var apigwUsagePlanDeleteCmd = &cobra.Command{
	Use:   "delete [usage-plan-id]",
	Short: "Delete a usage plan",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		if err := client.DeleteUsagePlan(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete usage plan", err)
		}
		fmt.Printf("Usage Plan %s deleted\n", args[0])
	},
}
