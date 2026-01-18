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
	loadbalancerCmd.AddCommand(lbDescribePoolsCmd)
	loadbalancerCmd.AddCommand(lbCreatePoolCmd)
	loadbalancerCmd.AddCommand(lbDeletePoolCmd)
	loadbalancerCmd.AddCommand(lbDescribeMembersCmd)
	loadbalancerCmd.AddCommand(lbCreateMemberCmd)
	loadbalancerCmd.AddCommand(lbDeleteMemberCmd)

	lbDescribePoolsCmd.Flags().String("pool-id", "", "Pool ID")

	lbCreatePoolCmd.Flags().String("name", "", "Pool name (required)")
	lbCreatePoolCmd.Flags().String("protocol", "TCP", "Protocol (TCP/HTTP/HTTPS/PROXY)")
	lbCreatePoolCmd.Flags().String("algorithm", "ROUND_ROBIN", "LB algorithm (ROUND_ROBIN/LEAST_CONNECTIONS/SOURCE_IP)")
	lbCreatePoolCmd.Flags().String("lb-id", "", "Load balancer ID")
	lbCreatePoolCmd.Flags().String("listener-id", "", "Listener ID")
	lbCreatePoolCmd.MarkFlagRequired("name")

	lbDeletePoolCmd.Flags().String("pool-id", "", "Pool ID (required)")
	lbDeletePoolCmd.MarkFlagRequired("pool-id")

	lbDescribeMembersCmd.Flags().String("pool-id", "", "Pool ID (required)")
	lbDescribeMembersCmd.Flags().String("member-id", "", "Member ID (optional, to describe one)")
	lbDescribeMembersCmd.MarkFlagRequired("pool-id")

	lbCreateMemberCmd.Flags().String("pool-id", "", "Pool ID (required)")
	lbCreateMemberCmd.Flags().String("address", "", "Member IP address (required)")
	lbCreateMemberCmd.Flags().Int("port", 80, "Member port (required)")
	lbCreateMemberCmd.Flags().Int("weight", 1, "Member weight")
	lbCreateMemberCmd.Flags().String("subnet-id", "", "Subnet ID")
	lbCreateMemberCmd.MarkFlagRequired("pool-id")
	lbCreateMemberCmd.MarkFlagRequired("address")
	lbCreateMemberCmd.MarkFlagRequired("port")

	lbDeleteMemberCmd.Flags().String("pool-id", "", "Pool ID (required)")
	lbDeleteMemberCmd.Flags().String("member-id", "", "Member ID (required)")
	lbDeleteMemberCmd.MarkFlagRequired("pool-id")
	lbDeleteMemberCmd.MarkFlagRequired("member-id")
}

var lbDescribePoolsCmd = &cobra.Command{
	Use:   "describe-pools",
	Short: "Describe pools",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("pool-id")

		if id != "" {
			result, err := client.GetPool(ctx, id)
			if err != nil {
				exitWithError("Failed to get pool", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			p := result.Pool
			fmt.Printf("ID:            %s\n", p.ID)
			fmt.Printf("Name:          %s\n", p.Name)
			fmt.Printf("Protocol:      %s\n", p.Protocol)
			fmt.Printf("Algorithm:     %s\n", p.LBAlgorithm)
			fmt.Printf("Load Balancer: %s\n", p.LoadBalancerID)
			fmt.Printf("Listener:      %s\n", p.ListenerID)
			fmt.Printf("Health Monitor:%s\n", p.HealthMonitorID)
			fmt.Printf("Status:        %s\n", p.OperatingStatus)
			fmt.Printf("Members:       %d\n", len(p.Members))
		} else {
			result, err := client.ListPools(ctx)
			if err != nil {
				exitWithError("Failed to list pools", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tPROTOCOL\tALGORITHM\tSTATUS")
			for _, p := range result.Pools {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					p.ID, p.Name, p.Protocol, p.LBAlgorithm, p.OperatingStatus)
			}
			w.Flush()
		}
	},
}

var lbCreatePoolCmd = &cobra.Command{
	Use:   "create-pool",
	Short: "Create a new pool",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		name, _ := cmd.Flags().GetString("name")
		protocol, _ := cmd.Flags().GetString("protocol")
		algorithm, _ := cmd.Flags().GetString("algorithm")
		lbID, _ := cmd.Flags().GetString("lb-id")
		listenerID, _ := cmd.Flags().GetString("listener-id")

		input := &loadbalancer.CreatePoolInput{
			Name:           name,
			Protocol:       protocol,
			LBAlgorithm:    algorithm,
			LoadBalancerID: lbID,
			ListenerID:     listenerID,
		}

		result, err := client.CreatePool(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create pool", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Pool created: %s\n", result.Pool.ID)
		fmt.Printf("Name: %s\n", result.Pool.Name)
	},
}

var lbDeletePoolCmd = &cobra.Command{
	Use:   "delete-pool",
	Short: "Delete a pool",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		id, _ := cmd.Flags().GetString("pool-id")
		if err := client.DeletePool(context.Background(), id); err != nil {
			exitWithError("Failed to delete pool", err)
		}
		fmt.Printf("Pool %s deleted\n", id)
	},
}

var lbDescribeMembersCmd = &cobra.Command{
	Use:   "describe-members",
	Short: "Describe members in a pool",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		ctx := context.Background()
		poolID, _ := cmd.Flags().GetString("pool-id")
		memberID, _ := cmd.Flags().GetString("member-id")

		if memberID != "" {
			result, err := client.GetMember(ctx, poolID, memberID)
			if err != nil {
				exitWithError("Failed to get member", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			m := result.Member
			fmt.Printf("ID:      %s\n", m.ID)
			fmt.Printf("Name:    %s\n", m.Name)
			fmt.Printf("Address: %s\n", m.Address)
			fmt.Printf("Port:    %d\n", m.ProtocolPort)
			fmt.Printf("Weight:  %d\n", m.Weight)
			fmt.Printf("Subnet:  %s\n", m.SubnetID)
			fmt.Printf("Status:  %s\n", m.OperatingStatus)
		} else {
			result, err := client.ListMembers(ctx, poolID)
			if err != nil {
				exitWithError("Failed to list members", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tADDRESS\tPORT\tWEIGHT\tSTATUS")
			for _, m := range result.Members {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%s\n",
					m.ID, m.Name, m.Address, m.ProtocolPort, m.Weight, m.OperatingStatus)
			}
			w.Flush()
		}
	},
}

var lbCreateMemberCmd = &cobra.Command{
	Use:   "create-member",
	Short: "Create a new member",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		poolID, _ := cmd.Flags().GetString("pool-id")
		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetInt("port")
		weight, _ := cmd.Flags().GetInt("weight")
		subnetID, _ := cmd.Flags().GetString("subnet-id")

		input := &loadbalancer.CreateMemberInput{
			Address:      address,
			ProtocolPort: port,
			Weight:       weight,
			SubnetID:     subnetID,
		}

		result, err := client.CreateMember(context.Background(), poolID, input)
		if err != nil {
			exitWithError("Failed to create member", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Member created: %s\n", result.Member.ID)
		fmt.Printf("Address: %s:%d\n", result.Member.Address, result.Member.ProtocolPort)
	},
}

var lbDeleteMemberCmd = &cobra.Command{
	Use:   "delete-member",
	Short: "Delete a member",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		poolID, _ := cmd.Flags().GetString("pool-id")
		memberID, _ := cmd.Flags().GetString("member-id")

		if err := client.DeleteMember(context.Background(), poolID, memberID); err != nil {
			exitWithError("Failed to delete member", err)
		}
		fmt.Printf("Member %s deleted\n", memberID)
	},
}
