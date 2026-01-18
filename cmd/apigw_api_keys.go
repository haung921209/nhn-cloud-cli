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
	apigwCmd.AddCommand(apigwDescribeAPIKeysCmd)
	apigwCmd.AddCommand(apigwCreateAPIKeyCmd)
	apigwCmd.AddCommand(apigwUpdateAPIKeyCmd)
	apigwCmd.AddCommand(apigwDeleteAPIKeyCmd)
	apigwCmd.AddCommand(apigwRegenerateAPIKeyCmd)

	apigwCreateAPIKeyCmd.Flags().String("name", "", "API key name (required)")
	apigwCreateAPIKeyCmd.Flags().String("description", "", "API key description")
	apigwCreateAPIKeyCmd.MarkFlagRequired("name")

	apigwUpdateAPIKeyCmd.Flags().String("apikey-id", "", "API Key ID (required)")
	apigwUpdateAPIKeyCmd.Flags().String("name", "", "API key name")
	apigwUpdateAPIKeyCmd.Flags().String("description", "", "API key description")
	apigwUpdateAPIKeyCmd.Flags().String("status", "", "API key status")
	apigwUpdateAPIKeyCmd.MarkFlagRequired("apikey-id")

	apigwDeleteAPIKeyCmd.Flags().String("apikey-id", "", "API Key ID (required)")
	apigwDeleteAPIKeyCmd.MarkFlagRequired("apikey-id")

	apigwRegenerateAPIKeyCmd.Flags().String("apikey-id", "", "API Key ID (required)")
	apigwRegenerateAPIKeyCmd.Flags().String("key-type", "PRIMARY", "Key type: PRIMARY or SECONDARY")
	apigwRegenerateAPIKeyCmd.MarkFlagRequired("apikey-id")
}

var apigwDescribeAPIKeysCmd = &cobra.Command{
	Use:     "describe-api-keys",
	Aliases: []string{"list-api-keys"},
	Short:   "List API keys",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()

		result, err := client.ListAPIKeys(ctx)
		if err != nil {
			exitWithError("Failed to list API keys", err)
		}

		if output == "json" {
			printJSON(result)
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

var apigwCreateAPIKeyCmd = &cobra.Command{
	Use:   "create-api-key",
	Short: "Create a new API key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &apigw.CreateAPIKeyInput{
			Name:        name,
			Description: description,
		}

		result, err := client.CreateAPIKey(ctx, input)
		if err != nil {
			exitWithError("Failed to create API key", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("API Key created: %s\n", result.APIKey.ID)
		fmt.Printf("Name: %s\n", result.APIKey.Name)
		fmt.Printf("Primary Key: %s\n", result.APIKey.PrimaryKey)
		fmt.Printf("Secondary Key: %s\n", result.APIKey.SecondaryKey)
	},
}

var apigwUpdateAPIKeyCmd = &cobra.Command{
	Use:   "update-api-key",
	Short: "Update an API key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("apikey-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")

		input := &apigw.UpdateAPIKeyInput{
			Name:        name,
			Description: description,
			StatusCode:  status,
		}

		result, err := client.UpdateAPIKey(ctx, id, input)
		if err != nil {
			exitWithError("Failed to update API key", err)
		}

		fmt.Printf("API Key updated: %s\n", result.APIKey.ID)
	},
}

var apigwDeleteAPIKeyCmd = &cobra.Command{
	Use:   "delete-api-key",
	Short: "Delete an API key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("apikey-id")

		if err := client.DeleteAPIKey(ctx, id); err != nil {
			exitWithError("Failed to delete API key", err)
		}

		fmt.Printf("API Key %s deleted\n", id)
	},
}

var apigwRegenerateAPIKeyCmd = &cobra.Command{
	Use:   "regenerate-api-key",
	Short: "Regenerate an API key",
	Run: func(cmd *cobra.Command, args []string) {
		client := newAPIGWClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("apikey-id")
		keyType, _ := cmd.Flags().GetString("key-type")

		input := &apigw.RegenerateAPIKeyInput{
			KeyType: keyType,
		}

		result, err := client.RegenerateAPIKey(ctx, id, input)
		if err != nil {
			exitWithError("Failed to regenerate API key", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("API Key regenerated: %s\n", result.APIKey.ID)
		fmt.Printf("Primary Key: %s\n", result.APIKey.PrimaryKey)
		fmt.Printf("Secondary Key: %s\n", result.APIKey.SecondaryKey)
	},
}
