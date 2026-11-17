package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/servicegateway"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/vpc"
	"github.com/spf13/cobra"
)

func init() {
	serviceGatewayCmd.AddCommand(sgDescribeCmd)
	serviceGatewayCmd.AddCommand(sgGetCmd)
	serviceGatewayCmd.AddCommand(sgCreateCmd)
	serviceGatewayCmd.AddCommand(sgUpdateCmd)
	serviceGatewayCmd.AddCommand(sgDeleteCmd)

	sgGetCmd.Flags().String("gateway-id", "", requiredHint("Service Gateway ID", hintSvcGatewayID))
	sgGetCmd.MarkFlagRequired("gateway-id")

	sgCreateCmd.Flags().String("name", "", "Service gateway name (required)")
	sgCreateCmd.Flags().String("description", "", "Description")
	sgCreateCmd.Flags().String("subnet-id", "", requiredHint("Subnet ID", hintSubnetID))
	sgCreateCmd.Flags().String("service-endpoint-id", "", requiredHint("Service endpoint ID", hintSvcEndpointID))
	sgCreateCmd.Flags().String("network-id", "", "VPC ID of the subnet (looked up from --subnet-id when omitted)")
	sgCreateCmd.Flags().String("fixed-ip", "", "IP address for the gateway inside the subnet (assigned automatically when omitted)")
	sgCreateCmd.Flags().Bool("include-gateway-identity", false, "Use a fixed NAT IP address (only for endpoints that support it)")
	sgCreateCmd.MarkFlagRequired("name")
	sgCreateCmd.MarkFlagRequired("subnet-id")
	sgCreateCmd.MarkFlagRequired("service-endpoint-id")

	sgUpdateCmd.Flags().String("gateway-id", "", requiredHint("Service Gateway ID", hintSvcGatewayID))
	sgUpdateCmd.Flags().String("name", "", "Service gateway name")
	sgUpdateCmd.Flags().String("description", "", "Description")
	sgUpdateCmd.MarkFlagRequired("gateway-id")

	sgDeleteCmd.Flags().String("gateway-id", "", requiredHint("Service Gateway ID", hintSvcGatewayID))
	sgDeleteCmd.MarkFlagRequired("gateway-id")
}

var sgDescribeCmd = &cobra.Command{
	Use:     "describe-service-gateways",
	Aliases: []string{"list-service-gateways"},
	Short:   "List all service gateways",
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		result, err := client.ListServiceGateways(context.Background())
		if err != nil {
			exitWithError("Failed to list service gateways", err)
		}

		if output == "json" {
			printJSON(result)
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

var sgGetCmd = &cobra.Command{
	Use:     "describe-service-gateway",
	Aliases: []string{"get-service-gateway"},
	Short:   "Get service gateway details",
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("gateway-id")

		result, err := client.GetServiceGateway(ctx, id)
		if err != nil {
			exitWithError("Failed to get service gateway", err)
		}

		if output == "json" {
			printJSON(result)
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

var sgCreateCmd = &cobra.Command{
	Use:   "create-service-gateway",
	Short: "Create a new service gateway",
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		serviceEndpointID, _ := cmd.Flags().GetString("service-endpoint-id")

		networkID, _ := cmd.Flags().GetString("network-id")
		if networkID == "" {
			// The API requires the VPC ID; derive it from the subnet.
			sub, err := vpc.NewClient(getRegion(), getIdentityCreds(), nil, debug).GetSubnet(context.Background(), subnetID)
			if err != nil {
				exitWithError("Failed to look up the subnet's VPC (pass --network-id)", err)
			}
			networkID = sub.Subnet.VPCID
		}

		input := &servicegateway.CreateServiceGatewayInput{
			Name:              name,
			Description:       description,
			SubnetID:          subnetID,
			NetworkID:         networkID,
			ServiceEndpointID: serviceEndpointID,
		}

		// Options the SDK input does not carry go through the API directly.
		extra := map[string]interface{}{}
		if v, _ := cmd.Flags().GetString("fixed-ip"); v != "" {
			extra["fixed_ip"] = v
		}
		if cmd.Flags().Changed("include-gateway-identity") {
			extra["include_gateway_identity"], _ = cmd.Flags().GetBool("include-gateway-identity")
		}
		if len(extra) > 0 {
			body := map[string]interface{}{}
			if err := remarshal(input, &body); err != nil {
				exitWithError("Failed to build request", err)
			}
			for k, v := range extra {
				body[k] = v
			}
			var res struct {
				ServiceGateway struct {
					ID      string `json:"id"`
					Name    string `json:"name"`
					FixedIP string `json:"fixed_ip"`
					Status  string `json:"status"`
				} `json:"servicegateway"`
			}
			raw := apiCall(http.MethodPost, "/v2.0/gateways/servicegateways", map[string]interface{}{"servicegateway": body}, &res)
			emit(raw, func() {
				fmt.Printf("Service gateway created: %s\n", res.ServiceGateway.ID)
				fmt.Printf("Name: %s\n", res.ServiceGateway.Name)
				fmt.Printf("IP Address: %s\n", res.ServiceGateway.FixedIP)
				fmt.Printf("Status: %s\n", res.ServiceGateway.Status)
			})
			return
		}

		result, err := client.CreateServiceGateway(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create service gateway", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Service gateway created: %s\n", result.ServiceGateway.ID)
		fmt.Printf("Name: %s\n", result.ServiceGateway.Name)
		fmt.Printf("IP Address: %s\n", result.ServiceGateway.IPAddress)
		fmt.Printf("Status: %s\n", result.ServiceGateway.Status)
	},
}

var sgUpdateCmd = &cobra.Command{
	Use:   "update-service-gateway",
	Short: "Update a service gateway",
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		id, _ := cmd.Flags().GetString("gateway-id")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		input := &servicegateway.UpdateServiceGatewayInput{
			Name:        name,
			Description: description,
		}

		result, err := client.UpdateServiceGateway(context.Background(), id, input)
		if err != nil {
			exitWithError("Failed to update service gateway", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Service gateway updated: %s\n", result.ServiceGateway.ID)
		fmt.Printf("Name: %s\n", result.ServiceGateway.Name)
	},
}

var sgDeleteCmd = &cobra.Command{
	Use:   "delete-service-gateway",
	Short: "Delete a service gateway",
	Run: func(cmd *cobra.Command, args []string) {
		client := newServiceGatewayClient()
		id, _ := cmd.Flags().GetString("gateway-id")
		if err := client.DeleteServiceGateway(context.Background(), id); err != nil {
			exitWithError("Failed to delete service gateway", err)
		}
		fmt.Printf("Service gateway %s deleted\n", id)
	},
}
