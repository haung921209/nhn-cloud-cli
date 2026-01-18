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
	loadbalancerCmd.AddCommand(lbDescribeListenersCmd)
	loadbalancerCmd.AddCommand(lbCreateListenerCmd)
	loadbalancerCmd.AddCommand(lbDeleteListenerCmd)

	lbDescribeListenersCmd.Flags().String("listener-id", "", "Listener ID")

	lbCreateListenerCmd.Flags().String("name", "", "Listener name (required)")
	lbCreateListenerCmd.Flags().String("lb-id", "", "Load balancer ID (required)")
	lbCreateListenerCmd.Flags().String("protocol", "TCP", "Protocol (TCP/HTTP/HTTPS/TERMINATED_HTTPS)")
	lbCreateListenerCmd.Flags().Int("port", 80, "Protocol port (required)")
	lbCreateListenerCmd.Flags().String("pool-id", "", "Default pool ID")
	lbCreateListenerCmd.Flags().Int("connection-limit", -1, "Connection limit")
	lbCreateListenerCmd.MarkFlagRequired("name")
	lbCreateListenerCmd.MarkFlagRequired("lb-id")
	lbCreateListenerCmd.MarkFlagRequired("port")

	lbDeleteListenerCmd.Flags().String("listener-id", "", "Listener ID (required)")
	lbDeleteListenerCmd.MarkFlagRequired("listener-id")
}

var lbDescribeListenersCmd = &cobra.Command{
	Use:   "describe-listeners",
	Short: "Describe listeners",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("listener-id")

		if id != "" {
			result, err := client.GetListener(ctx, id)
			if err != nil {
				exitWithError("Failed to get listener", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			l := result.Listener
			fmt.Printf("ID:              %s\n", l.ID)
			fmt.Printf("Name:            %s\n", l.Name)
			fmt.Printf("Protocol:        %s\n", l.Protocol)
			fmt.Printf("Port:            %d\n", l.ProtocolPort)
			fmt.Printf("Load Balancer:   %s\n", l.LoadBalancerID)
			fmt.Printf("Default Pool:    %s\n", l.DefaultPoolID)
			fmt.Printf("Connection Limit:%d\n", l.ConnectionLimit)
			fmt.Printf("Status:          %s\n", l.OperatingStatus)
		} else {
			result, err := client.ListListeners(ctx)
			if err != nil {
				exitWithError("Failed to list listeners", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tPROTOCOL\tPORT\tLB_ID\tSTATUS")
			for _, l := range result.Listeners {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n",
					l.ID, l.Name, l.Protocol, l.ProtocolPort, l.LoadBalancerID, l.OperatingStatus)
			}
			w.Flush()
		}
	},
}

var lbCreateListenerCmd = &cobra.Command{
	Use:   "create-listener",
	Short: "Create a new listener",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		name, _ := cmd.Flags().GetString("name")
		lbID, _ := cmd.Flags().GetString("lb-id")
		protocol, _ := cmd.Flags().GetString("protocol")
		port, _ := cmd.Flags().GetInt("port")
		poolID, _ := cmd.Flags().GetString("pool-id")
		connLimit, _ := cmd.Flags().GetInt("connection-limit")

		input := &loadbalancer.CreateListenerInput{
			Name:            name,
			LoadBalancerID:  lbID,
			Protocol:        protocol,
			ProtocolPort:    port,
			DefaultPoolID:   poolID,
			ConnectionLimit: connLimit,
		}

		result, err := client.CreateListener(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create listener", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Listener created: %s\n", result.Listener.ID)
		fmt.Printf("Name: %s\n", result.Listener.Name)
	},
}

var lbDeleteListenerCmd = &cobra.Command{
	Use:   "delete-listener",
	Short: "Delete a listener",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		id, _ := cmd.Flags().GetString("listener-id")
		if err := client.DeleteListener(context.Background(), id); err != nil {
			exitWithError("Failed to delete listener", err)
		}
		fmt.Printf("Listener %s deleted\n", id)
	},
}
