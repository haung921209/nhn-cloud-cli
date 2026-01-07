package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/colocationgw"
	"github.com/spf13/cobra"
)

var colocationgwCmd = &cobra.Command{
	Use:     "colocation-gateway",
	Aliases: []string{"colocation-gw", "cologw"},
	Short:   "Manage Colocation Gateways",
}

func init() {
	rootCmd.AddCommand(colocationgwCmd)

	colocationgwCmd.AddCommand(colocationgwListCmd)
	colocationgwCmd.AddCommand(colocationgwGetCmd)
}

func newColocationGWClient() *colocationgw.Client {
	return colocationgw.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

var colocationgwListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all colocation gateways",
	Run: func(cmd *cobra.Command, args []string) {
		client := newColocationGWClient()
		result, err := client.List(context.Background())
		if err != nil {
			exitWithError("Failed to list colocation gateways", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tROUTER_ID\tVLAN_ID")
		for _, gw := range result.ColocationGateways {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
				gw.ID, gw.Name, gw.Status, gw.RouterID, gw.VLANID)
		}
		w.Flush()
	},
}

var colocationgwGetCmd = &cobra.Command{
	Use:   "get [gateway-id]",
	Short: "Get colocation gateway details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newColocationGWClient()
		result, err := client.Get(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get colocation gateway", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		gw := result.ColocationGateway
		fmt.Printf("ID:                %s\n", gw.ID)
		fmt.Printf("Name:              %s\n", gw.Name)
		fmt.Printf("Description:       %s\n", gw.Description)
		fmt.Printf("Status:            %s\n", gw.Status)
		fmt.Printf("Tenant ID:         %s\n", gw.TenantID)
		fmt.Printf("Router ID:         %s\n", gw.RouterID)
		fmt.Printf("Subnet ID:         %s\n", gw.SubnetID)
		fmt.Printf("Network ID:        %s\n", gw.NetworkID)
		fmt.Printf("Local IP Address:  %s\n", gw.LocalIPAddress)
		fmt.Printf("Remote IP Address: %s\n", gw.RemoteIPAddress)
		fmt.Printf("VLAN ID:           %d\n", gw.VLANID)
		fmt.Printf("Connection Type:   %s\n", gw.ConnectionType)
		fmt.Printf("Created At:        %s\n", gw.CreatedAt)
		fmt.Printf("Updated At:        %s\n", gw.UpdatedAt)
	},
}
