package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/floatingip"
	"github.com/spf13/cobra"
)

func init() {
	networkCmd.AddCommand(networkDescribeFloatingIPsCmd)
	networkCmd.AddCommand(networkAllocateFloatingIPCmd)
	networkCmd.AddCommand(networkReleaseFloatingIPCmd)
	networkCmd.AddCommand(networkAssociateFloatingIPCmd)
	networkCmd.AddCommand(networkDisassociateFloatingIPCmd)

	networkDescribeFloatingIPsCmd.Flags().String("floating-ip-id", "", "Floating IP ID")

	networkAllocateFloatingIPCmd.Flags().String("network-id", "", requiredHint("External network ID", hintExtNetID))
	networkAllocateFloatingIPCmd.MarkFlagRequired("network-id")

	networkReleaseFloatingIPCmd.Flags().String("floating-ip-id", "", requiredHint("Floating IP ID", hintFloatingIPID))
	networkReleaseFloatingIPCmd.MarkFlagRequired("floating-ip-id")

	networkAssociateFloatingIPCmd.Flags().String("floating-ip-id", "", requiredHint("Floating IP ID", hintFloatingIPID))
	networkAssociateFloatingIPCmd.Flags().String("port-id", "", "Port ID to associate (this or --instance-id is required; "+hintPortID+")")
	networkAssociateFloatingIPCmd.Flags().String("instance-id", "", "Associate with the instance's port (the instance must have exactly one port)")
	networkAssociateFloatingIPCmd.MarkFlagRequired("floating-ip-id")
	networkAssociateFloatingIPCmd.MarkFlagsOneRequired("port-id", "instance-id")
	networkAssociateFloatingIPCmd.MarkFlagsMutuallyExclusive("port-id", "instance-id")

	networkDisassociateFloatingIPCmd.Flags().String("floating-ip-id", "", requiredHint("Floating IP ID", hintFloatingIPID))
	networkDisassociateFloatingIPCmd.MarkFlagRequired("floating-ip-id")
}

var networkDescribeFloatingIPsCmd = &cobra.Command{
	Use:   "describe-floating-ips",
	Short: "Describe floating IPs",
	Run: func(cmd *cobra.Command, args []string) {
		client := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("floating-ip-id")

		if id != "" {
			result, err := client.GetFloatingIP(ctx, id)
			if err != nil {
				exitWithError("Failed to get floating IP", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			fip := result.FloatingIP
			fmt.Printf("ID:          %s\n", fip.ID)
			fmt.Printf("Floating IP: %s\n", fip.FloatingIPAddress)
			fmt.Printf("Fixed IP:    %s\n", fip.FixedIPAddress)
			fmt.Printf("Status:      %s\n", fip.Status)
			portID := ""
			if fip.PortID != nil {
				portID = *fip.PortID
			}
			fmt.Printf("Port ID:     %s\n", portID)
		} else {
			result, err := client.ListFloatingIPs(ctx)
			if err != nil {
				exitWithError("Failed to list floating IPs", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tFLOATING_IP\tFIXED_IP\tSTATUS\tPORT_ID")
			for _, fip := range result.FloatingIPs {
				portID := ""
				if fip.PortID != nil {
					portID = *fip.PortID
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					fip.ID, fip.FloatingIPAddress, fip.FixedIPAddress, fip.Status, portID)
			}
			w.Flush()
		}
	},
}

var networkAllocateFloatingIPCmd = &cobra.Command{
	Use:   "allocate-floating-ip",
	Short: "Allocate a new floating IP",
	Run: func(cmd *cobra.Command, args []string) {
		client := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		networkID, _ := cmd.Flags().GetString("network-id")

		input := &floatingip.CreateFloatingIPInput{
			FloatingNetworkID: networkID,
		}

		result, err := client.CreateFloatingIP(ctx, input)
		if err != nil {
			exitWithError("Failed to allocate floating IP", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Floating IP allocated: %s\n", result.FloatingIP.FloatingIPAddress)
		fmt.Printf("ID: %s\n", result.FloatingIP.ID)
	},
}

var networkReleaseFloatingIPCmd = &cobra.Command{
	Use:   "release-floating-ip",
	Short: "Release a floating IP",
	Run: func(cmd *cobra.Command, args []string) {
		client := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()
		id, _ := cmd.Flags().GetString("floating-ip-id")

		if err := client.DeleteFloatingIP(ctx, id); err != nil {
			exitWithError("Failed to release floating IP", err)
		}

		fmt.Printf("Floating IP %s released\n", id)
	},
}

var networkAssociateFloatingIPCmd = &cobra.Command{
	Use:   "associate-floating-ip",
	Short: "Associate floating IP with a port",
	Run: func(cmd *cobra.Command, args []string) {
		client := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		id, _ := cmd.Flags().GetString("floating-ip-id")
		portID, _ := cmd.Flags().GetString("port-id")
		if instanceID, _ := cmd.Flags().GetString("instance-id"); instanceID != "" {
			portID = instancePortID(instanceID)
		}
		if portID == "" {
			exitWithError("--port-id or --instance-id must not be empty", nil)
		}

		input := &floatingip.UpdateFloatingIPInput{
			PortID: &portID,
		}

		result, err := client.UpdateFloatingIP(ctx, id, input)
		if err != nil {
			exitWithError("Failed to associate floating IP", err)
		}

		if err := printWriteResult(result, "Floating IP %s associated with port %s\n",
			result.FloatingIP.FloatingIPAddress, portID); err != nil {
			exitWithError("Failed to print output", err)
		}
	},
}

var networkDisassociateFloatingIPCmd = &cobra.Command{
	Use:   "disassociate-floating-ip",
	Short: "Disassociate floating IP from port",
	Run: func(cmd *cobra.Command, args []string) {
		client := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		ctx := context.Background()

		id, _ := cmd.Flags().GetString("floating-ip-id")

		input := &floatingip.UpdateFloatingIPInput{
			PortID: nil,
		}

		result, err := client.UpdateFloatingIP(ctx, id, input)
		if err != nil {
			exitWithError("Failed to disassociate floating IP", err)
		}

		if err := printWriteResult(result, "Floating IP %s disassociated\n", id); err != nil {
			exitWithError("Failed to print output", err)
		}
	},
}
