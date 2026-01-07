package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/servicegateway"
	"github.com/spf13/cobra"
)

var serviceGatewayCmd = &cobra.Command{
	Use:     "service-gateway",
	Aliases: []string{"svcgw", "sg-gateway"},
	Short:   "Manage Service Gateways",
}

// Endpoint subcommand
var serviceEndpointCmd = &cobra.Command{
	Use:     "endpoint",
	Aliases: []string{"endpoints", "ep"},
	Short:   "Manage Service Endpoints (predefined NHN Cloud endpoints)",
}

func init() {
	rootCmd.AddCommand(serviceGatewayCmd)

	// Gateway subcommands
	serviceGatewayCmd.AddCommand(serviceGatewayListCmd)
	serviceGatewayCmd.AddCommand(serviceGatewayGetCmd)
	serviceGatewayCmd.AddCommand(serviceGatewayCreateCmd)
	serviceGatewayCmd.AddCommand(serviceGatewayUpdateCmd)
	serviceGatewayCmd.AddCommand(serviceGatewayDeleteCmd)

	// Endpoint subcommands
	serviceGatewayCmd.AddCommand(serviceEndpointCmd)
	serviceEndpointCmd.AddCommand(serviceEndpointListCmd)
	serviceEndpointCmd.AddCommand(serviceEndpointGetCmd)

	// Create gateway flags
	serviceGatewayCreateCmd.Flags().String("name", "", "Service gateway name (required)")
	serviceGatewayCreateCmd.Flags().String("description", "", "Description")
	serviceGatewayCreateCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	serviceGatewayCreateCmd.Flags().String("service-endpoint-id", "", "Service endpoint ID (required)")
	serviceGatewayCreateCmd.MarkFlagRequired("name")
	serviceGatewayCreateCmd.MarkFlagRequired("subnet-id")
	serviceGatewayCreateCmd.MarkFlagRequired("service-endpoint-id")

	// Update gateway flags
	serviceGatewayUpdateCmd.Flags().String("name", "", "Service gateway name")
	serviceGatewayUpdateCmd.Flags().String("description", "", "Description")
}

func newServiceGatewayClient() *servicegateway.Client {
	return servicegateway.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

// ================================
// Service Gateway Commands
// ================================

var serviceGatewayListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all service gateways",
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		result, err := client.ListServiceGateways(context.Background())
		if err != nil {
			exitWithError("Failed to list service gateways", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tIP_ADDRESS\tSUBNET_ID\tSTATUS")
		for _, gw := range result.ServiceGateways {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				gw.ID, gw.Name, gw.IPAddress, gw.SubnetID, gw.Status)
		}
		w.Flush()
	},
}

var serviceGatewayGetCmd = &cobra.Command{
	Use:   "get [gateway-id]",
	Short: "Get service gateway details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		result, err := client.GetServiceGateway(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get service gateway", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		gw := result.ServiceGateway
		fmt.Printf("ID:                  %s\n", gw.ID)
		fmt.Printf("Name:                %s\n", gw.Name)
		fmt.Printf("Description:         %s\n", gw.Description)
		fmt.Printf("Subnet ID:           %s\n", gw.SubnetID)
		fmt.Printf("Service Endpoint ID: %s\n", gw.ServiceEndpointID)
		fmt.Printf("IP Address:          %s\n", gw.IPAddress)
		fmt.Printf("Status:              %s\n", gw.Status)
		fmt.Printf("Tenant ID:           %s\n", gw.TenantID)
		fmt.Printf("Created:             %s\n", gw.CreateTime)
		fmt.Printf("Updated:             %s\n", gw.UpdateTime)
	},
}

var serviceGatewayCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new service gateway",
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		serviceEndpointID, _ := cmd.Flags().GetString("service-endpoint-id")

		input := &servicegateway.CreateServiceGatewayInput{
			Name:              name,
			Description:       description,
			SubnetID:          subnetID,
			ServiceEndpointID: serviceEndpointID,
		}

		result, err := client.CreateServiceGateway(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create service gateway", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Service gateway created: %s\n", result.ServiceGateway.ID)
		fmt.Printf("Name: %s\n", result.ServiceGateway.Name)
		fmt.Printf("IP Address: %s\n", result.ServiceGateway.IPAddress)
		fmt.Printf("Status: %s\n", result.ServiceGateway.Status)
	},
}

var serviceGatewayUpdateCmd = &cobra.Command{
	Use:   "update [gateway-id]",
	Short: "Update a service gateway",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &servicegateway.UpdateServiceGatewayInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateServiceGateway(context.Background(), args[0], input)
		if err != nil {
			exitWithError("Failed to update service gateway", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Service gateway updated: %s\n", result.ServiceGateway.ID)
		fmt.Printf("Name: %s\n", result.ServiceGateway.Name)
	},
}

var serviceGatewayDeleteCmd = &cobra.Command{
	Use:   "delete [gateway-id]",
	Short: "Delete a service gateway",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		if err := client.DeleteServiceGateway(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete service gateway", err)
		}
		fmt.Printf("Service gateway %s deleted\n", args[0])
	},
}

// ================================
// Service Endpoint Commands
// ================================

var serviceEndpointListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available service endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		result, err := client.ListServiceEndpoints(context.Background())
		if err != nil {
			exitWithError("Failed to list service endpoints", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSERVICE_NAME\tREGION\tTYPE")
		for _, ep := range result.ServiceEndpoints {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				ep.ID, ep.Name, ep.ServiceName, ep.Region, ep.EndpointType)
		}
		w.Flush()
	},
}

var serviceEndpointGetCmd = &cobra.Command{
	Use:   "get [endpoint-id]",
	Short: "Get service endpoint details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		result, err := client.GetServiceEndpoint(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get service endpoint", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		ep := result.ServiceEndpoint
		fmt.Printf("ID:            %s\n", ep.ID)
		fmt.Printf("Name:          %s\n", ep.Name)
		fmt.Printf("Service Name:  %s\n", ep.ServiceName)
		fmt.Printf("Description:   %s\n", ep.Description)
		fmt.Printf("Region:        %s\n", ep.Region)
		fmt.Printf("Endpoint Type: %s\n", ep.EndpointType)
		fmt.Printf("Created:       %s\n", ep.CreateTime)
	},
}
