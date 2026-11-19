package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncr"
	"github.com/spf13/cobra"
)

func init() {
	ncrCmd.AddCommand(ncrDescribeRegistriesCmd)
	ncrCmd.AddCommand(ncrCreateRegistryCmd)
	ncrCmd.AddCommand(ncrDeleteRegistryCmd)

	ncrDescribeRegistriesCmd.Flags().String("registry-id", "", "Registry ID")

	ncrCreateRegistryCmd.Flags().String("name", "", "Registry name (required)")
	ncrCreateRegistryCmd.Flags().String("description", "", "Registry description")
	ncrCreateRegistryCmd.Flags().Bool("public", false, "Make registry public")
	ncrCreateRegistryCmd.MarkFlagRequired("name")

	ncrDeleteRegistryCmd.Flags().String("registry-id", "", "Registry ID (required)")
	ncrDeleteRegistryCmd.MarkFlagRequired("registry-id")
}

var ncrDescribeRegistriesCmd = &cobra.Command{
	Use:     "describe-registries",
	Aliases: []string{"list-registries"},
	Short:   "Describe container registries",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("registry-id")

		if id != "" {
			result, err := client.GetRegistry(ctx, id)
			if err != nil {
				exitWithError("Failed to get registry", err)
			}
			if structuredOutput() {
				printOutput(result)
				return
			}
			r := result.Registry
			fmt.Printf("ID:      %d\n", r.ID)
			fmt.Printf("Name:    %s\n", r.Name)
			fmt.Printf("URI:     %s\n", r.URI)
			fmt.Printf("Public:  %s\n", r.Metadata["public"])
			fmt.Printf("Images:  %d\n", r.RepoCount)
			fmt.Printf("Created: %s\n", r.CreatedAt)
		} else {
			result, err := client.ListRegistries(ctx)
			if err != nil {
				exitWithError("Failed to list registries", err)
			}
			if structuredOutput() {
				printOutput(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tURI\tPUBLIC\tIMAGES\tCREATED")
			for _, r := range result.Registries {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%s\n",
					r.ID, r.Name, r.URI, r.Metadata["public"], r.RepoCount, r.CreatedAt)
			}
			w.Flush()
		}
	},
}

var ncrCreateRegistryCmd = &cobra.Command{
	Use:   "create-registry",
	Short: "Create a new container registry",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		isPublic, _ := cmd.Flags().GetBool("public")

		input := &ncr.CreateRegistryInput{
			Name:        name,
			Description: description,
			IsPublic:    isPublic,
		}

		if _, err := client.CreateRegistry(ctx, input); err != nil {
			exitWithError("Failed to create registry", err)
		}
		// The create response is the common header only; read the registry back.
		result, err := client.GetRegistry(ctx, name)
		if err != nil {
			exitWithError("Registry created, but reading it back failed", err)
		}

		r := result.Registry
		printWriteResult(result, "Registry created successfully!\nID:   %d\nName: %s\nURI:  %s\n", r.ID, r.Name, r.URI)
	},
}

var ncrDeleteRegistryCmd = &cobra.Command{
	Use:   "delete-registry",
	Short: "Delete a container registry",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("registry-id")

		if err := client.DeleteRegistry(ctx, id); err != nil {
			exitWithError("Failed to delete registry", err)
		}

		fmt.Printf("Registry %s deleted successfully\n", id)
	},
}
