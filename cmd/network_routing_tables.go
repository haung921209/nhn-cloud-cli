package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// Routing tables and routes (Network > VPC API v2: /v2.0/routingtables,
// /v2.0/routes, /v2.0/vpcsubnets/{id}/attach_routingtable).

const (
	hintVPCID          = "find it with: nhncloud network describe-vpcs"
	hintSubnetID       = "find it with: nhncloud network describe-subnets"
	hintRoutingTableID = "find it with: nhncloud network describe-routing-tables --vpc-id <vpc-id>"
	hintIGWID          = "find it with: nhncloud internet-gateway describe-internet-gateways"
	hintExtNetID       = "find it with: nhncloud internet-gateway describe-external-networks"
	hintRouteID        = "find it with: nhncloud network describe-routes --routing-table-id <id>"
	hintPortID         = "find it with: nhncloud network describe-ports --device-id <instance-or-lb-id>"
	hintLBID           = "find it with: nhncloud loadbalancer describe-load-balancers"
	hintListenerID     = "find it with: nhncloud loadbalancer describe-listeners"
	hintPoolID         = "find it with: nhncloud loadbalancer describe-pools"
	hintSvcEndpointID  = "find it with: nhncloud service-gateway describe-service-endpoints (or describe-my-service-endpoints)"
	hintSvcGatewayID   = "find it with: nhncloud service-gateway describe-service-gateways"
	hintAllowProjectID = "find it with: nhncloud service-gateway describe-endpoint-allowed-projects"
	hintFloatingIPID   = "find it with: nhncloud network describe-floating-ips"
)

func requiredHint(desc, hint string) string { return desc + " (required; " + hint + ")" }

type idName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UnmarshalJSON accepts both {"id":..,"name":..} (list with detail=true) and
// a bare ID string (single routing table response).
func (n *idName) UnmarshalJSON(b []byte) error {
	var id string
	if err := json.Unmarshal(b, &id); err == nil {
		*n = idName{ID: id}
		return nil
	}
	type plain idName
	return json.Unmarshal(b, (*plain)(n))
}

type routingTable struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	DefaultTable bool     `json:"default_table"`
	Distributed  bool     `json:"distributed"`
	GatewayID    string   `json:"gateway_id"`
	GatewayName  string   `json:"gateway_name"`
	State        string   `json:"state"`
	VPCs         []idName `json:"vpcs"`
	Subnets      []idName `json:"subnets"`
	Routes       []route  `json:"routes"`
	CreateTime   string   `json:"create_time"`
}

type route struct {
	ID             string `json:"id"`
	CIDR           string `json:"cidr"`
	Gateway        string `json:"gateway"`
	GatewayID      string `json:"gateway_id"`
	RoutingTableID string `json:"routingtable_id"`
	Description    string `json:"description"`
}

func (t routingTable) vpcID() string {
	if len(t.VPCs) > 0 {
		return t.VPCs[0].ID
	}
	return ""
}

// filterRoutingTablesByVPC keeps tables attached to vpcID ("" keeps all).
func filterRoutingTablesByVPC(tables []routingTable, vpcID string) []routingTable {
	if vpcID == "" {
		return tables
	}
	var out []routingTable
	for _, t := range tables {
		for _, v := range t.VPCs {
			if v.ID == vpcID {
				out = append(out, t)
				break
			}
		}
	}
	return out
}

// emit prints raw (the API response) for json/yaml/--query output and calls
// table otherwise.
func emit(raw interface{}, table func()) {
	if output == "json" || output == "yaml" || query != "" {
		if err := printOutput(raw); err != nil {
			exitWithError("Failed to print output", err)
		}
		return
	}
	table()
}

// apiCall runs one network API request and decodes the JSON response into
// both a generic value (for json/yaml output) and typed (for tables).
func apiCall(method, path string, body interface{}, typed interface{}) interface{} {
	var raw interface{}
	if err := networkRequest(context.Background(), method, path, body, &raw); err != nil {
		exitWithError(fmt.Sprintf("%s %s failed", method, path), err)
	}
	if typed != nil && raw != nil {
		if err := remarshal(raw, typed); err != nil {
			exitWithError("Failed to decode response", err)
		}
	}
	return raw
}

func init() {
	networkCmd.AddCommand(rtDescribeCmd, rtCreateCmd, rtUpdateCmd, rtDeleteCmd,
		rtAttachGatewayCmd, rtDetachGatewayCmd, rtSetDefaultCmd,
		rtAttachSubnetCmd, rtDetachSubnetCmd,
		routeDescribeCmd, routeCreateCmd, routeDeleteCmd)

	rtDescribeCmd.Flags().String("routing-table-id", "", "Routing table ID (shows routes and attached subnets)")
	rtDescribeCmd.Flags().String("vpc-id", "", "Only routing tables of this VPC ("+hintVPCID+")")

	rtCreateCmd.Flags().String("name", "", "Routing table name (required)")
	rtCreateCmd.Flags().String("vpc-id", "", requiredHint("VPC ID", hintVPCID))
	rtCreateCmd.Flags().Bool("distributed", true, "Distributed routing (false = centralized)")
	rtCreateCmd.MarkFlagRequired("name")
	rtCreateCmd.MarkFlagRequired("vpc-id")

	rtUpdateCmd.Flags().String("routing-table-id", "", requiredHint("Routing table ID", hintRoutingTableID))
	rtUpdateCmd.Flags().String("name", "", "New name")
	rtUpdateCmd.Flags().Bool("distributed", true, "Distributed routing (false = centralized)")
	rtUpdateCmd.MarkFlagRequired("routing-table-id")

	for _, c := range []*cobra.Command{rtDeleteCmd, rtDetachGatewayCmd, rtSetDefaultCmd} {
		c.Flags().String("routing-table-id", "", requiredHint("Routing table ID", hintRoutingTableID))
		c.MarkFlagRequired("routing-table-id")
	}

	rtAttachGatewayCmd.Flags().String("routing-table-id", "", requiredHint("Routing table ID", hintRoutingTableID))
	rtAttachGatewayCmd.Flags().String("gateway-id", "", requiredHint("Internet gateway ID", hintIGWID))
	rtAttachGatewayCmd.MarkFlagRequired("routing-table-id")
	rtAttachGatewayCmd.MarkFlagRequired("gateway-id")

	rtAttachSubnetCmd.Flags().String("subnet-id", "", requiredHint("Subnet ID", hintSubnetID))
	rtAttachSubnetCmd.Flags().String("routing-table-id", "", requiredHint("Routing table ID", hintRoutingTableID))
	rtAttachSubnetCmd.MarkFlagRequired("subnet-id")
	rtAttachSubnetCmd.MarkFlagRequired("routing-table-id")

	rtDetachSubnetCmd.Flags().String("subnet-id", "", requiredHint("Subnet ID", hintSubnetID))
	rtDetachSubnetCmd.MarkFlagRequired("subnet-id")

	routeDescribeCmd.Flags().String("routing-table-id", "", "Only routes of this routing table ("+hintRoutingTableID+")")
	routeDescribeCmd.Flags().String("route-id", "", "Route ID")

	routeCreateCmd.Flags().String("routing-table-id", "", requiredHint("Routing table ID", hintRoutingTableID))
	routeCreateCmd.Flags().String("cidr", "", "Destination CIDR (required, e.g. 10.10.0.0/24)")
	routeCreateCmd.Flags().String("gateway", "", "Next-hop IP address (required)")
	routeCreateCmd.Flags().String("description", "", "Description")
	routeCreateCmd.MarkFlagRequired("routing-table-id")
	routeCreateCmd.MarkFlagRequired("cidr")
	routeCreateCmd.MarkFlagRequired("gateway")

	routeDeleteCmd.Flags().String("route-id", "", requiredHint("Route ID", hintRouteID))
	routeDeleteCmd.MarkFlagRequired("route-id")
}

func printRoutingTableDetail(t routingTable) {
	fmt.Printf("ID:           %s\n", t.ID)
	fmt.Printf("Name:         %s\n", t.Name)
	fmt.Printf("VPC:          %s\n", t.vpcID())
	fmt.Printf("Default:      %t\n", t.DefaultTable)
	fmt.Printf("Distributed:  %t\n", t.Distributed)
	fmt.Printf("Gateway:      %s %s\n", t.GatewayID, t.GatewayName)
	fmt.Printf("State:        %s\n", t.State)
	if len(t.Subnets) > 0 {
		fmt.Println("Subnets:")
		for _, s := range t.Subnets {
			fmt.Printf("  %s  %s\n", s.ID, s.Name)
		}
	}
	if len(t.Routes) > 0 {
		fmt.Println("Routes:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  ID\tCIDR\tGATEWAY")
		for _, r := range t.Routes {
			fmt.Fprintf(w, "  %s\t%s\t%s\n", r.ID, r.CIDR, r.Gateway)
		}
		w.Flush()
	}
}

var rtDescribeCmd = &cobra.Command{
	Use:   "describe-routing-tables",
	Short: "Describe routing tables (find the default routing table of a VPC)",
	Example: `  nhncloud network describe-routing-tables --vpc-id <vpc-id>
  nhncloud network describe-routing-tables --routing-table-id <id>`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("routing-table-id")
		vpcID, _ := cmd.Flags().GetString("vpc-id")

		if id != "" {
			var res struct {
				RoutingTable routingTable `json:"routingtable"`
			}
			raw := apiCall(http.MethodGet, "/v2.0/routingtables/"+url.PathEscape(id), nil, &res)
			emit(raw, func() { printRoutingTableDetail(res.RoutingTable) })
			return
		}

		// detail=true adds the VPC and subnets of each table; the list API has
		// no VPC filter, so filter here.
		var res struct {
			RoutingTables []routingTable `json:"routingtables"`
		}
		raw := apiCall(http.MethodGet, "/v2.0/routingtables?detail=true", nil, &res)
		tables := filterRoutingTablesByVPC(res.RoutingTables, vpcID)
		if vpcID != "" {
			raw = map[string]interface{}{"routingtables": tables}
		}
		emit(raw, func() {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tVPC_ID\tDEFAULT\tDISTRIBUTED\tGATEWAY_ID\tSUBNETS\tSTATE")
			for _, t := range tables {
				fmt.Fprintf(w, "%s\t%s\t%s\t%t\t%t\t%s\t%d\t%s\n",
					t.ID, t.Name, t.vpcID(), t.DefaultTable, t.Distributed, t.GatewayID, len(t.Subnets), t.State)
			}
			w.Flush()
		})
	},
}

var rtCreateCmd = &cobra.Command{
	Use:   "create-routing-table",
	Short: "Create a routing table in a VPC",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		vpcID, _ := cmd.Flags().GetString("vpc-id")
		distributed, _ := cmd.Flags().GetBool("distributed")
		body := map[string]interface{}{"routingtable": map[string]interface{}{
			"name": name, "vpc_id": vpcID, "distributed": distributed,
		}}
		var res struct {
			RoutingTable routingTable `json:"routingtable"`
		}
		raw := apiCall(http.MethodPost, "/v2.0/routingtables", body, &res)
		emit(raw, func() { fmt.Printf("Routing table created: %s\n", res.RoutingTable.ID) })
	},
}

var rtUpdateCmd = &cobra.Command{
	Use:   "update-routing-table",
	Short: "Rename a routing table or change its routing mode",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("routing-table-id")
		rt := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			rt["name"], _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("distributed") {
			rt["distributed"], _ = cmd.Flags().GetBool("distributed")
		}
		if len(rt) == 0 {
			exitWithError("nothing to update: pass --name and/or --distributed", nil)
		}
		var res struct {
			RoutingTable routingTable `json:"routingtable"`
		}
		raw := apiCall(http.MethodPut, "/v2.0/routingtables/"+url.PathEscape(id), map[string]interface{}{"routingtable": rt}, &res)
		emit(raw, func() { fmt.Printf("Routing table updated: %s\n", res.RoutingTable.ID) })
	},
}

var rtDeleteCmd = &cobra.Command{
	Use:   "delete-routing-table",
	Short: "Delete a routing table",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("routing-table-id")
		apiCall(http.MethodDelete, "/v2.0/routingtables/"+url.PathEscape(id), nil, nil)
		fmt.Printf("Routing table %s deleted\n", id)
	},
}

// rtAction runs a PUT action on a routing table and prints the result.
func rtAction(id, action string, body interface{}, done string) {
	var res struct {
		RoutingTable routingTable `json:"routingtable"`
	}
	raw := apiCall(http.MethodPut, "/v2.0/routingtables/"+url.PathEscape(id)+"/"+action, body, &res)
	emit(raw, func() { fmt.Printf("%s: %s\n", done, id) })
}

var rtAttachGatewayCmd = &cobra.Command{
	Use:   "attach-routing-table-gateway",
	Short: "Attach an internet gateway to a routing table",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("routing-table-id")
		gw, _ := cmd.Flags().GetString("gateway-id")
		rtAction(id, "attach_gateway", map[string]string{"gateway_id": gw}, "Gateway attached to routing table")
	},
}

var rtDetachGatewayCmd = &cobra.Command{
	Use:   "detach-routing-table-gateway",
	Short: "Detach the internet gateway from a routing table",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("routing-table-id")
		rtAction(id, "detach_gateway", nil, "Gateway detached from routing table")
	},
}

var rtSetDefaultCmd = &cobra.Command{
	Use:   "set-default-routing-table",
	Short: "Make a routing table the default of its VPC",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("routing-table-id")
		rtAction(id, "set_as_default", nil, "Default routing table set")
	},
}

func subnetAction(subnetID, action string, body interface{}, done string) {
	var raw interface{}
	raw = apiCall(http.MethodPut, "/v2.0/vpcsubnets/"+url.PathEscape(subnetID)+"/"+action, body, nil)
	emit(raw, func() { fmt.Printf("%s: %s\n", done, subnetID) })
}

var rtAttachSubnetCmd = &cobra.Command{
	Use:   "attach-subnet-routing-table",
	Short: "Attach a subnet to a routing table",
	Run: func(cmd *cobra.Command, args []string) {
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		rtID, _ := cmd.Flags().GetString("routing-table-id")
		subnetAction(subnetID, "attach_routingtable", map[string]string{"routingtable_id": rtID}, "Subnet attached to routing table")
	},
}

var rtDetachSubnetCmd = &cobra.Command{
	Use:   "detach-subnet-routing-table",
	Short: "Detach a subnet from its routing table (it falls back to the VPC default)",
	Run: func(cmd *cobra.Command, args []string) {
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		subnetAction(subnetID, "detach_routingtable", nil, "Subnet detached from routing table")
	},
}

var routeDescribeCmd = &cobra.Command{
	Use:   "describe-routes",
	Short: "Describe routes of routing tables",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("route-id")
		rtID, _ := cmd.Flags().GetString("routing-table-id")
		if id != "" {
			var res struct {
				Route route `json:"route"`
			}
			raw := apiCall(http.MethodGet, "/v2.0/routes/"+url.PathEscape(id), nil, &res)
			emit(raw, func() {
				r := res.Route
				fmt.Printf("ID:            %s\nCIDR:          %s\nGateway:       %s\nRouting Table: %s\nDescription:   %s\n",
					r.ID, r.CIDR, r.Gateway, r.RoutingTableID, r.Description)
			})
			return
		}
		path := "/v2.0/routes"
		if rtID != "" {
			path += "?routingtable_id=" + url.QueryEscape(rtID)
		}
		var res struct {
			Routes []route `json:"routes"`
		}
		raw := apiCall(http.MethodGet, path, nil, &res)
		emit(raw, func() {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tCIDR\tGATEWAY\tROUTING_TABLE_ID\tDESCRIPTION")
			for _, r := range res.Routes {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", r.ID, r.CIDR, r.Gateway, r.RoutingTableID, r.Description)
			}
			w.Flush()
		})
	},
}

var routeCreateCmd = &cobra.Command{
	Use:   "create-route",
	Short: "Add a static route to a routing table",
	Run: func(cmd *cobra.Command, args []string) {
		rtID, _ := cmd.Flags().GetString("routing-table-id")
		cidr, _ := cmd.Flags().GetString("cidr")
		gw, _ := cmd.Flags().GetString("gateway")
		desc, _ := cmd.Flags().GetString("description")
		// The API marks description as required; send an empty one when unset.
		body := map[string]interface{}{"route": map[string]string{
			"routingtable_id": rtID, "cidr": cidr, "gateway": gw, "description": desc,
		}}
		var res struct {
			Route route `json:"route"`
		}
		raw := apiCall(http.MethodPost, "/v2.0/routes", body, &res)
		emit(raw, func() { fmt.Printf("Route created: %s\n", res.Route.ID) })
	},
}

var routeDeleteCmd = &cobra.Command{
	Use:   "delete-route",
	Short: "Delete a static route",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("route-id")
		apiCall(http.MethodDelete, "/v2.0/routes/"+url.PathEscape(id), nil, nil)
		fmt.Printf("Route %s deleted\n", id)
	},
}
