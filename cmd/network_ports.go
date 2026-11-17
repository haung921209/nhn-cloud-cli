package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	networkCmd.AddCommand(networkDescribePortsCmd)
	networkDescribePortsCmd.Flags().String("port-id", "", "Port ID")
	networkDescribePortsCmd.Flags().String("device-id", "", "Only ports of this device (instance ID, load balancer ID, ...)")
	networkDescribePortsCmd.Flags().String("vpc-id", "", "Only ports in this VPC ("+hintVPCID+")")
	networkDescribePortsCmd.Flags().String("fixed-ip", "", "Only the port holding this IP address")
	networkDescribePortsCmd.Flags().String("subnet-id", "", "Only ports in this subnet ("+hintSubnetID+")")
}

type netPort struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	NetworkID   string `json:"network_id"`
	DeviceID    string `json:"device_id"`
	DeviceOwner string `json:"device_owner"`
	Status      string `json:"status"`
	FixedIPs    []struct {
		SubnetID  string `json:"subnet_id"`
		IPAddress string `json:"ip_address"`
	} `json:"fixed_ips"`
}

var networkDescribePortsCmd = &cobra.Command{
	Use:   "describe-ports",
	Short: "Describe network ports (find the port ID of an instance or load balancer)",
	Example: `  nhncloud network describe-ports --device-id <instance-id>
  nhncloud network describe-ports --fixed-ip 192.168.0.10`,
	Run: func(cmd *cobra.Command, args []string) {
		if id, _ := cmd.Flags().GetString("port-id"); id != "" {
			// A single port has no useful table form; emit falls back to JSON.
			emit(apiCall(http.MethodGet, "/v2.0/ports/"+url.PathEscape(id), nil, nil), nil)
			return
		}
		q := url.Values{}
		if v, _ := cmd.Flags().GetString("device-id"); v != "" {
			q.Set("device_id", v)
		}
		if v, _ := cmd.Flags().GetString("vpc-id"); v != "" {
			q.Set("network_id", v)
		}
		if v, _ := cmd.Flags().GetString("fixed-ip"); v != "" {
			q.Add("fixed_ips", "ip_address="+v)
		}
		if v, _ := cmd.Flags().GetString("subnet-id"); v != "" {
			q.Add("fixed_ips", "subnet_id="+v)
		}
		path := "/v2.0/ports"
		if len(q) > 0 {
			path += "?" + q.Encode()
		}
		var res struct {
			Ports []netPort `json:"ports"`
		}
		raw := apiCall(http.MethodGet, path, nil, &res)
		emit(raw, func() {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tFIXED_IPS\tDEVICE_OWNER\tDEVICE_ID\tSTATUS")
			for _, p := range res.Ports {
				var ips []string
				for _, ip := range p.FixedIPs {
					ips = append(ips, ip.IPAddress)
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					p.ID, p.Name, strings.Join(ips, ","), p.DeviceOwner, p.DeviceID, p.Status)
			}
			w.Flush()
		})
	},
}
