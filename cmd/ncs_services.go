package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncs"
	"github.com/spf13/cobra"
)

func init() {
	ncsCmd.AddCommand(ncsDescribeServicesCmd)
	ncsCmd.AddCommand(ncsGetServiceCmd)
	ncsCmd.AddCommand(ncsCreateServiceCmd)
	ncsCmd.AddCommand(ncsDeleteServiceCmd)

	ncsDescribeServicesCmd.Flags().String("namespace", "", "Filter by namespace")

	ncsGetServiceCmd.Flags().String("service-id", "", "Service ID (required)")
	ncsGetServiceCmd.MarkFlagRequired("service-id")

	ncsCreateServiceCmd.Flags().String("name", "", "Service name (required)")
	ncsCreateServiceCmd.Flags().String("namespace", "default", "Namespace")
	ncsCreateServiceCmd.Flags().String("selector", "", "Label selector (key=value)")
	ncsCreateServiceCmd.Flags().Int("port", 80, "Service port")
	ncsCreateServiceCmd.Flags().Int("target-port", 0, "Target container port")
	ncsCreateServiceCmd.Flags().String("type", "LoadBalancer", "Service type (ClusterIP, LoadBalancer)")
	ncsCreateServiceCmd.MarkFlagRequired("name")

	ncsDeleteServiceCmd.Flags().String("service-id", "", "Service ID (required)")
	ncsDeleteServiceCmd.MarkFlagRequired("service-id")
}

var ncsDescribeServicesCmd = &cobra.Command{
	Use:     "describe-services",
	Aliases: []string{"list-services", "services"},
	Short:   "List all services",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		namespace, _ := cmd.Flags().GetString("namespace")

		result, err := client.ListServices(ctx, namespace)
		if err != nil {
			exitWithError("Failed to list services", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tNAMESPACE\tTYPE\tCLUSTER-IP\tEXTERNAL-IP\tPORTS")
		for _, svc := range result.Services {
			ports := formatServicePorts(svc.Ports)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				svc.ID, svc.Name, svc.Namespace, svc.Type, svc.ClusterIP, svc.ExternalIP, ports)
		}
		w.Flush()
	},
}

func formatServicePorts(ports []ncs.ServicePort) string {
	if len(ports) == 0 {
		return "-"
	}
	var parts []string
	for _, p := range ports {
		parts = append(parts, fmt.Sprintf("%d:%d/%s", p.Port, p.TargetPort, p.Protocol))
	}
	return strings.Join(parts, ", ")
}

var ncsGetServiceCmd = &cobra.Command{
	Use:     "describe-service",
	Aliases: []string{"get-service", "service-get"},
	Short:   "Get service details",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		serviceID, _ := cmd.Flags().GetString("service-id")
		if serviceID == "" && len(args) > 0 {
			serviceID = args[0]
		}
		if serviceID == "" {
			exitWithError("Service ID required", nil)
		}

		result, err := client.GetService(ctx, serviceID)
		if err != nil {
			exitWithError("Failed to get service", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("ID:          %s\n", result.ID)
		fmt.Printf("Name:        %s\n", result.Name)
		fmt.Printf("Namespace:   %s\n", result.Namespace)
		fmt.Printf("Type:        %s\n", result.Type)
		fmt.Printf("Cluster IP:  %s\n", result.ClusterIP)
		fmt.Printf("External IP: %s\n", result.ExternalIP)
		if len(result.Ports) > 0 {
			fmt.Printf("Ports:\n")
			for _, p := range result.Ports {
				fmt.Printf("  - %d:%d/%s\n", p.Port, p.TargetPort, p.Protocol)
			}
		}
		fmt.Printf("Created:     %s\n", result.CreatedAt)
	},
}

var ncsCreateServiceCmd = &cobra.Command{
	Use:   "create-service",
	Short: "Create a new service",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		namespace, _ := cmd.Flags().GetString("namespace")
		selectorStr, _ := cmd.Flags().GetString("selector")
		port, _ := cmd.Flags().GetInt("port")
		targetPort, _ := cmd.Flags().GetInt("target-port")
		svcType, _ := cmd.Flags().GetString("type")

		if targetPort == 0 {
			targetPort = port
		}

		// Parse selector string "key=value" into map
		selector := make(map[string]string)
		if selectorStr != "" {
			parts := strings.SplitN(selectorStr, "=", 2)
			if len(parts) == 2 {
				selector[parts[0]] = parts[1]
			}
		}

		input := &ncs.CreateServiceInput{
			Name:      name,
			Namespace: namespace,
			Type:      svcType,
			Selector:  selector,
			Ports: []ncs.ServicePort{
				{
					Port:       port,
					TargetPort: targetPort,
					Protocol:   "TCP",
				},
			},
		}

		result, err := client.CreateService(ctx, input)
		if err != nil {
			exitWithError("Failed to create service", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Service created successfully!\n")
		fmt.Printf("ID:   %s\n", result.ID)
		fmt.Printf("Name: %s\n", result.Name)
	},
}

var ncsDeleteServiceCmd = &cobra.Command{
	Use:   "delete-service",
	Short: "Delete a service",
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCSClient()
		ctx := context.Background()
		serviceID, _ := cmd.Flags().GetString("service-id")
		if serviceID == "" && len(args) > 0 {
			serviceID = args[0]
		}
		if serviceID == "" {
			exitWithError("Service ID required", nil)
		}

		if err := client.DeleteService(ctx, serviceID); err != nil {
			exitWithError("Failed to delete service", err)
		}

		fmt.Printf("Service %s deleted successfully\n", serviceID)
	},
}
