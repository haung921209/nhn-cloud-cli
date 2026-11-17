package cmd

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/loadbalancer"
	"github.com/spf13/cobra"
)

// Listener options the SDK input types do not carry. Each entry maps a flag
// to the listener field of the Load Balancer API.
type listenerOpt struct {
	flag, field, kind, usage string
}

var listenerExtraOpts = []listenerOpt{
	{"keepalive-timeout", "keepalive_timeout", "int", "Keepalive timeout in seconds (API default 300)"},
	{"keepalive-enable", "keepalive_enable", "bool", "Enable keepalive"},
	{"proxy-protocol", "proxy_protocol", "bool", "Send the PROXY protocol header to members"},
	{"x-forwarded-for", "enable_x_forwarded_for", "bool", "Add X-Forwarded-For (HTTP listeners)"},
	{"x-forwarded-proto", "enable_x_forwarded_proto", "bool", "Add X-Forwarded-Proto (HTTP listeners)"},
	{"x-forwarded-port", "enable_x_forwarded_port", "bool", "Add X-Forwarded-Port (HTTP listeners)"},
	{"block-invalid-http-request", "block_invalid_http_request", "bool", "Block malformed HTTP requests"},
	{"tls-version", "tls_version", "string", "Minimum TLS version (TERMINATED_HTTPS)"},
	{"ssl-policy-id", "ssl_policy_id", "string", "SSL policy ID (TERMINATED_HTTPS)"},
	{"tls-container-ref", "default_tls_container_ref", "string", "Default TLS certificate container ref (TERMINATED_HTTPS)"},
	{"sni-container-refs", "sni_container_refs", "strings", "SNI certificate container refs (TERMINATED_HTTPS)"},
	{"description", "description", "string", "Description"},
	{"admin-state-up", "admin_state_up", "bool", "Administrative state"},
}

// createOnlyListenerOpts are accepted by create but not by update.
var createOnlyListenerOpts = []listenerOpt{
	{"protocol-version", "protocol_version", "string", "Protocol version (HTTP listeners)"},
}

func addListenerOptFlags(cmd *cobra.Command, opts []listenerOpt) {
	for _, o := range opts {
		switch o.kind {
		case "int":
			cmd.Flags().Int(o.flag, 0, o.usage)
		case "bool":
			cmd.Flags().Bool(o.flag, false, o.usage)
		case "strings":
			cmd.Flags().StringSlice(o.flag, nil, o.usage)
		default:
			cmd.Flags().String(o.flag, "", o.usage)
		}
	}
}

// changedListenerFields returns API fields for the flags the user set.
func changedListenerFields(cmd *cobra.Command, opts []listenerOpt) map[string]interface{} {
	m := map[string]interface{}{}
	for _, o := range opts {
		if !cmd.Flags().Changed(o.flag) {
			continue
		}
		switch o.kind {
		case "int":
			m[o.field], _ = cmd.Flags().GetInt(o.flag)
		case "bool":
			m[o.field], _ = cmd.Flags().GetBool(o.flag)
		case "strings":
			m[o.field], _ = cmd.Flags().GetStringSlice(o.flag)
		default:
			m[o.field], _ = cmd.Flags().GetString(o.flag)
		}
	}
	return m
}

func init() {
	addListenerOptFlags(lbCreateListenerCmd, listenerExtraOpts)
	addListenerOptFlags(lbCreateListenerCmd, createOnlyListenerOpts)

	lbDescribeListenersCmd.Flags().Bool("all-fields", false, "Print the full API object as JSON (keepalive, proxy protocol, X-Forwarded-*, TLS settings)")

	loadbalancerCmd.AddCommand(lbUpdateListenerCmd)
	lbUpdateListenerCmd.Flags().String("listener-id", "", requiredHint("Listener ID", hintListenerID))
	lbUpdateListenerCmd.Flags().String("name", "", "Listener name")
	lbUpdateListenerCmd.Flags().String("pool-id", "", "Default pool ID ("+hintPoolID+")")
	lbUpdateListenerCmd.Flags().Int("connection-limit", 0, "Connection limit")
	lbUpdateListenerCmd.MarkFlagRequired("listener-id")
	addListenerOptFlags(lbUpdateListenerCmd, listenerExtraOpts)
}

// createListenerWithOptions posts the SDK create input plus the extra fields
// and prints the API's listener object.
func createListenerWithOptions(input *loadbalancer.CreateListenerInput, extra map[string]interface{}) {
	listener := map[string]interface{}{}
	if err := remarshal(input, &listener); err != nil {
		exitWithError("Failed to build request", err)
	}
	for k, v := range extra {
		listener[k] = v
	}
	var res loadbalancer.GetListenerOutput
	raw := apiCall(http.MethodPost, "/v2.0/lbaas/listeners", map[string]interface{}{"listener": listener}, &res)
	emit(raw, func() {
		fmt.Printf("Listener created: %s\n", res.Listener.ID)
		fmt.Printf("Name: %s\n", res.Listener.Name)
	})
}

func describeListenerAllFields(id string) {
	path := "/v2.0/lbaas/listeners"
	if id != "" {
		path += "/" + url.PathEscape(id)
	}
	raw := apiCall(http.MethodGet, path, nil, nil)
	if output == "yaml" || query != "" {
		emit(raw, nil)
		return
	}
	printJSON(raw)
}

var lbUpdateListenerCmd = &cobra.Command{
	Use:   "update-listener",
	Short: "Update a listener (keepalive, proxy protocol, connection limit, TLS, ...)",
	Example: `  nhncloud loadbalancer update-listener --listener-id <id> --keepalive-timeout 3600
  nhncloud loadbalancer update-listener --listener-id <id> --proxy-protocol=true`,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("listener-id")
		fields := changedListenerFields(cmd, listenerExtraOpts)
		if cmd.Flags().Changed("name") {
			fields["name"], _ = cmd.Flags().GetString("name")
		}
		if cmd.Flags().Changed("pool-id") {
			fields["default_pool_id"], _ = cmd.Flags().GetString("pool-id")
		}
		if cmd.Flags().Changed("connection-limit") {
			fields["connection_limit"], _ = cmd.Flags().GetInt("connection-limit")
		}
		if len(fields) == 0 {
			exitWithError("nothing to update: pass at least one option (see --help)", nil)
		}
		var res loadbalancer.GetListenerOutput
		raw := apiCall(http.MethodPut, "/v2.0/lbaas/listeners/"+url.PathEscape(id), map[string]interface{}{"listener": fields}, &res)
		emit(raw, func() { fmt.Printf("Listener updated: %s\n", res.Listener.ID) })
	},
}
