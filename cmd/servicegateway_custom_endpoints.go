package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// User-defined service endpoints, their allowed projects and endpoint usage
// (Network > Service Gateway API v2: /v2.0/gateways/myserviceendpoints,
// /v2.0/gateways/serviceendpointallowprojects, /v2.0/gateways/serviceendpointusages).

type myServiceEndpoint struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	ServiceName  string `json:"service_name"`
	PortID       string `json:"port_id"`
	EndpointType string `json:"endpoint_type"`
	MaxCount     int    `json:"max_count"`
	CurrentCount int    `json:"current_count"`
	Description  string `json:"description"`
}

type allowProject struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	ServiceEndpointID string `json:"service_endpoint_id"`
	TargetTenantID    string `json:"target_tenant_id"`
	Description       string `json:"description"`
}

type endpointUsage struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	ServiceEndpointID string `json:"service_endpoint_id"`
	NetworkID         string `json:"network_id"`
	SubnetID          string `json:"subnet_id"`
	FixedIP           string `json:"fixed_ip"`
	Status            string `json:"status"`
}

const hintMyEndpointID = "find it with: nhncloud service-gateway describe-my-service-endpoints"

func init() {
	serviceGatewayCmd.AddCommand(mseDescribeCmd, mseCreateCmd, mseUpdateCmd, mseDeleteCmd, mseGenerateNameCmd,
		sapDescribeCmd, sapAddCmd, sapUpdateCmd, sapRemoveCmd, seUsageCmd)

	mseDescribeCmd.Flags().String("endpoint-id", "", "User-defined service endpoint ID")

	mseCreateCmd.Flags().String("name", "", "Endpoint name (required)")
	mseCreateCmd.Flags().String("port-id", "", requiredHint("Port that receives the traffic, e.g. a load balancer VIP port", hintPortID))
	mseCreateCmd.Flags().String("display-name", "", "Name shown to projects that use the endpoint")
	mseCreateCmd.Flags().Int("max-count", 0, "Maximum number of service gateways that can use the endpoint")
	mseCreateCmd.Flags().String("description", "", "Description")
	mseCreateCmd.MarkFlagRequired("name")
	mseCreateCmd.MarkFlagRequired("port-id")

	mseUpdateCmd.Flags().String("endpoint-id", "", requiredHint("User-defined service endpoint ID", hintMyEndpointID))
	mseUpdateCmd.Flags().String("name", "", "Endpoint name")
	mseUpdateCmd.Flags().String("display-name", "", "Display name")
	mseUpdateCmd.Flags().Int("max-count", 0, "Maximum number of service gateways")
	mseUpdateCmd.Flags().String("description", "", "Description")
	mseUpdateCmd.MarkFlagRequired("endpoint-id")

	for _, c := range []*cobra.Command{mseDeleteCmd, mseGenerateNameCmd} {
		c.Flags().String("endpoint-id", "", requiredHint("User-defined service endpoint ID", hintMyEndpointID))
		c.MarkFlagRequired("endpoint-id")
	}

	sapDescribeCmd.Flags().String("endpoint-id", "", "Only projects allowed on this endpoint ("+hintMyEndpointID+")")
	sapDescribeCmd.Flags().String("target-project-id", "", "Only entries for this project ID (tenant ID)")

	sapAddCmd.Flags().String("endpoint-id", "", requiredHint("User-defined service endpoint ID", hintMyEndpointID))
	sapAddCmd.Flags().String("target-project-id", "", requiredHint("Tenant ID of the project to allow", "find it with: nhncloud auth whoami, run with the other project's credentials"))
	sapAddCmd.Flags().String("name", "", "Name")
	sapAddCmd.Flags().String("description", "", "Description")
	sapAddCmd.MarkFlagRequired("endpoint-id")
	sapAddCmd.MarkFlagRequired("target-project-id")

	sapUpdateCmd.Flags().String("allow-project-id", "", requiredHint("Allowed project entry ID", hintAllowProjectID))
	sapUpdateCmd.Flags().String("name", "", "Name")
	sapUpdateCmd.Flags().String("description", "", "Description")
	sapUpdateCmd.MarkFlagRequired("allow-project-id")

	sapRemoveCmd.Flags().String("allow-project-id", "", requiredHint("Allowed project entry ID", hintAllowProjectID))
	sapRemoveCmd.MarkFlagRequired("allow-project-id")

	seUsageCmd.Flags().String("vpc-id", "", "Only usages in this VPC")
	seUsageCmd.Flags().String("subnet-id", "", "Only usages in this subnet")
}

// optionalFields collects string/int flags the user set into API fields.
func optionalFields(cmd *cobra.Command, flagToField map[string]string) map[string]interface{} {
	m := map[string]interface{}{}
	for flag, field := range flagToField {
		if !cmd.Flags().Changed(flag) {
			continue
		}
		if f := cmd.Flags().Lookup(flag); f.Value.Type() == "int" {
			m[field], _ = cmd.Flags().GetInt(flag)
		} else {
			m[field], _ = cmd.Flags().GetString(flag)
		}
	}
	return m
}

var mseFields = map[string]string{"name": "name", "display-name": "display_name", "max-count": "max_count", "description": "description"}

func printMyServiceEndpoints(eps []myServiceEndpoint) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSERVICE_NAME\tPORT_ID\tUSED/MAX")
	for _, e := range eps {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d/%d\n", e.ID, e.Name, e.ServiceName, e.PortID, e.CurrentCount, e.MaxCount)
	}
	w.Flush()
}

var mseDescribeCmd = &cobra.Command{
	Use:   "describe-my-service-endpoints",
	Short: "Describe user-defined service endpoints of this project",
	Run: func(cmd *cobra.Command, args []string) {
		if id, _ := cmd.Flags().GetString("endpoint-id"); id != "" {
			var res struct {
				Endpoint myServiceEndpoint `json:"myserviceendpoint"`
			}
			raw := apiCall(http.MethodGet, "/v2.0/gateways/myserviceendpoints/"+url.PathEscape(id), nil, &res)
			emit(raw, func() { printMyServiceEndpoints([]myServiceEndpoint{res.Endpoint}) })
			return
		}
		var res struct {
			Endpoints []myServiceEndpoint `json:"myserviceendpoints"`
		}
		raw := apiCall(http.MethodGet, "/v2.0/gateways/myserviceendpoints", nil, &res)
		emit(raw, func() { printMyServiceEndpoints(res.Endpoints) })
	},
}

var mseCreateCmd = &cobra.Command{
	Use:   "create-my-service-endpoint",
	Short: "Create a user-defined service endpoint that other projects reach through a service gateway",
	Example: `  nhncloud service-gateway create-my-service-endpoint --name my-svc --port-id <lb-vip-port-id>
  nhncloud service-gateway add-endpoint-allowed-project --endpoint-id <id> --target-project-id <tenant-id>`,
	Run: func(cmd *cobra.Command, args []string) {
		fields := optionalFields(cmd, mseFields)
		fields["port_id"], _ = cmd.Flags().GetString("port-id")
		var res struct {
			Endpoint myServiceEndpoint `json:"myserviceendpoint"`
		}
		raw := apiCall(http.MethodPost, "/v2.0/gateways/myserviceendpoints", map[string]interface{}{"myserviceendpoint": fields}, &res)
		emit(raw, func() {
			fmt.Printf("Service endpoint created: %s\n", res.Endpoint.ID)
			fmt.Printf("Service name: %s\n", res.Endpoint.ServiceName)
		})
	},
}

var mseUpdateCmd = &cobra.Command{
	Use:   "update-my-service-endpoint",
	Short: "Update a user-defined service endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("endpoint-id")
		fields := optionalFields(cmd, mseFields)
		if len(fields) == 0 {
			exitWithError("nothing to update: pass --name, --display-name, --max-count or --description", nil)
		}
		var res struct {
			Endpoint myServiceEndpoint `json:"myserviceendpoint"`
		}
		raw := apiCall(http.MethodPut, "/v2.0/gateways/myserviceendpoints/"+url.PathEscape(id), map[string]interface{}{"myserviceendpoint": fields}, &res)
		emit(raw, func() { fmt.Printf("Service endpoint updated: %s\n", res.Endpoint.ID) })
	},
}

var mseDeleteCmd = &cobra.Command{
	Use:   "delete-my-service-endpoint",
	Short: "Delete a user-defined service endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("endpoint-id")
		apiCall(http.MethodDelete, "/v2.0/gateways/myserviceendpoints/"+url.PathEscape(id), nil, nil)
		fmt.Printf("Service endpoint %s deleted\n", id)
	},
}

var mseGenerateNameCmd = &cobra.Command{
	Use:   "regenerate-service-endpoint-name",
	Short: "Issue a new service name for a user-defined service endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("endpoint-id")
		raw := apiCall(http.MethodPut, "/v2.0/gateways/serviceendpoints/"+url.PathEscape(id)+"/generate_service_name", nil, nil)
		emit(raw, func() { fmt.Printf("Service name regenerated for %s\n", id) })
	},
}

func printAllowProjects(items []allowProject) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSERVICE_ENDPOINT_ID\tTARGET_PROJECT_ID")
	for _, a := range items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", a.ID, a.Name, a.ServiceEndpointID, a.TargetTenantID)
	}
	w.Flush()
}

var sapDescribeCmd = &cobra.Command{
	Use:   "describe-endpoint-allowed-projects",
	Short: "Describe projects allowed to use user-defined service endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		q := url.Values{}
		if v, _ := cmd.Flags().GetString("endpoint-id"); v != "" {
			q.Set("service_endpoint_id", v)
		}
		if v, _ := cmd.Flags().GetString("target-project-id"); v != "" {
			q.Set("target_tenant_id", v)
		}
		path := "/v2.0/gateways/serviceendpointallowprojects"
		if len(q) > 0 {
			path += "?" + q.Encode()
		}
		var res struct {
			Items []allowProject `json:"serviceendpointallowprojects"`
		}
		raw := apiCall(http.MethodGet, path, nil, &res)
		emit(raw, func() { printAllowProjects(res.Items) })
	},
}

var sapAddCmd = &cobra.Command{
	Use:   "add-endpoint-allowed-project",
	Short: "Allow another project to use a user-defined service endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		fields := optionalFields(cmd, map[string]string{"name": "name", "description": "description"})
		fields["service_endpoint_id"], _ = cmd.Flags().GetString("endpoint-id")
		fields["target_tenant_id"], _ = cmd.Flags().GetString("target-project-id")
		var res struct {
			Item allowProject `json:"serviceendpointallowproject"`
		}
		raw := apiCall(http.MethodPost, "/v2.0/gateways/serviceendpointallowprojects", map[string]interface{}{"serviceendpointallowproject": fields}, &res)
		emit(raw, func() { fmt.Printf("Allowed project added: %s\n", res.Item.ID) })
	},
}

var sapUpdateCmd = &cobra.Command{
	Use:   "update-endpoint-allowed-project",
	Short: "Rename or describe an allowed project entry",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("allow-project-id")
		fields := optionalFields(cmd, map[string]string{"name": "name", "description": "description"})
		if len(fields) == 0 {
			exitWithError("nothing to update: pass --name and/or --description", nil)
		}
		var res struct {
			Item allowProject `json:"serviceendpointallowproject"`
		}
		raw := apiCall(http.MethodPut, "/v2.0/gateways/serviceendpointallowprojects/"+url.PathEscape(id), map[string]interface{}{"serviceendpointallowproject": fields}, &res)
		emit(raw, func() { fmt.Printf("Allowed project updated: %s\n", res.Item.ID) })
	},
}

var sapRemoveCmd = &cobra.Command{
	Use:   "remove-endpoint-allowed-project",
	Short: "Stop allowing a project to use a user-defined service endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("allow-project-id")
		apiCall(http.MethodDelete, "/v2.0/gateways/serviceendpointallowprojects/"+url.PathEscape(id), nil, nil)
		fmt.Printf("Allowed project %s removed\n", id)
	},
}

var seUsageCmd = &cobra.Command{
	Use:   "describe-service-endpoint-usages",
	Short: "Describe service gateways that use this project's user-defined endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		q := url.Values{}
		if v, _ := cmd.Flags().GetString("vpc-id"); v != "" {
			q.Set("network_id", v)
		}
		if v, _ := cmd.Flags().GetString("subnet-id"); v != "" {
			q.Set("subnet_id", v)
		}
		path := "/v2.0/gateways/serviceendpointusages"
		if len(q) > 0 {
			path += "?" + q.Encode()
		}
		var res struct {
			Items []endpointUsage `json:"serviceendpointusages"`
		}
		raw := apiCall(http.MethodGet, path, nil, &res)
		emit(raw, func() {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tSERVICE_ENDPOINT_ID\tVPC_ID\tSUBNET_ID\tFIXED_IP\tSTATUS")
			for _, u := range res.Items {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", u.ID, u.Name, u.ServiceEndpointID, u.NetworkID, u.SubnetID, u.FixedIP, u.Status)
			}
			w.Flush()
		})
	},
}
