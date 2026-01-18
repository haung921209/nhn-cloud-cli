package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/loadbalancer"
	"github.com/spf13/cobra"
)

func init() {
	loadbalancerCmd.AddCommand(lbDescribeLoadBalancersCmd)
	loadbalancerCmd.AddCommand(lbGetLoadBalancerCmd) // Legacy or specific get? describe handles both usually.
	loadbalancerCmd.AddCommand(lbCreateLoadBalancerCmd)
	loadbalancerCmd.AddCommand(lbDeleteLoadBalancerCmd)

	lbDescribeLoadBalancersCmd.Flags().String("lb-id", "", "Load Balancer ID")

	lbCreateLoadBalancerCmd.Flags().String("name", "", "Load balancer name (required)")
	lbCreateLoadBalancerCmd.Flags().String("description", "", "Description")
	lbCreateLoadBalancerCmd.Flags().String("subnet-id", "", "VIP subnet ID (required)")
	lbCreateLoadBalancerCmd.Flags().String("vip-address", "", "VIP address (optional)")
	lbCreateLoadBalancerCmd.Flags().String("provider", "", "Provider (optional)")
	lbCreateLoadBalancerCmd.MarkFlagRequired("name")
	lbCreateLoadBalancerCmd.MarkFlagRequired("subnet-id")

	lbDeleteLoadBalancerCmd.Flags().String("lb-id", "", "Load balancer ID (required)")
	lbDeleteLoadBalancerCmd.MarkFlagRequired("lb-id")
}

var lbDescribeLoadBalancersCmd = &cobra.Command{
	Use:     "describe-load-balancers",
	Aliases: []string{"describe-lbs"},
	Short:   "Describe load balancers",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("lb-id")

		if id != "" {
			// Get Single
			result, err := client.GetLoadBalancer(ctx, id)
			if err != nil {
				exitWithError("Failed to get load balancer", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			lb := result.LoadBalancer
			fmt.Printf("ID:           %s\n", lb.ID)
			fmt.Printf("Name:         %s\n", lb.Name)
			fmt.Printf("Description:  %s\n", lb.Description)
			fmt.Printf("VIP Address:  %s\n", lb.VIPAddress)
			fmt.Printf("VIP Subnet:   %s\n", lb.VIPSubnetID)
			fmt.Printf("Status:       %s\n", lb.OperatingStatus)
			fmt.Printf("Provisioning: %s\n", lb.ProvisioningStatus)
			fmt.Printf("Provider:     %s\n", lb.Provider)
			fmt.Printf("Created:      %s\n", lb.CreatedAt)
		} else {
			// List All
			result, err := client.ListLoadBalancers(ctx)
			if err != nil {
				exitWithError("Failed to list load balancers", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tVIP_ADDRESS\tSTATUS\tPROVISIONING")
			for _, lb := range result.LoadBalancers {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					lb.ID, lb.Name, lb.VIPAddress, lb.OperatingStatus, lb.ProvisioningStatus)
			}
			w.Flush()
		}
	},
}

var lbGetLoadBalancerCmd = &cobra.Command{
	Use:   "get-load-balancer", // Alias or deprecated?
	Short: "Get load balancer details (use describe-load-balancers --lb-id)",
	Run: func(cmd *cobra.Command, args []string) {
		// Just redirect logic?
		fmt.Println("Please use: describe-load-balancers --lb-id <id>")
	},
}

var lbCreateLoadBalancerCmd = &cobra.Command{
	Use:   "create-load-balancer",
	Short: "Create a new load balancer",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		vipAddress, _ := cmd.Flags().GetString("vip-address")
		provider, _ := cmd.Flags().GetString("provider")

		input := &loadbalancer.CreateLoadBalancerInput{
			Name:        name,
			Description: description,
			VIPSubnetID: subnetID,
			VIPAddress:  vipAddress,
			Provider:    provider,
		}

		result, err := client.CreateLoadBalancer(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create load balancer", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Load balancer created: %s\n", result.LoadBalancer.ID)
		fmt.Printf("Name: %s\n", result.LoadBalancer.Name)
		fmt.Printf("VIP:  %s\n", result.LoadBalancer.VIPAddress)
	},
}

var lbDeleteLoadBalancerCmd = &cobra.Command{
	Use:   "delete-load-balancer",
	Short: "Delete a load balancer",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		id, _ := cmd.Flags().GetString("lb-id")
		if err := client.DeleteLoadBalancer(context.Background(), id); err != nil {
			exitWithError("Failed to delete load balancer", err)
		}
		fmt.Printf("Load balancer %s deleted\n", id)
	},
}
