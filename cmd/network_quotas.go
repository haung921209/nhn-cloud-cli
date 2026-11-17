package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// NHN Cloud does not expose a quota API for network resources (the
// OpenStack quota paths return 404), so this command counts what the project
// uses and shows it next to the default limits from the public resource
// policy. A project whose limits were raised on request has higher limits.
// ponytail: limits are hard-coded from the public policy; switch to an API
// when one is published.

type quotaRow struct {
	Resource string `json:"resource"`
	Scope    string `json:"scope"`
	Used     string `json:"used"`
	Limit    int    `json:"default_limit"`
}

type quotaInputs struct {
	tenantID        string
	vpcs            []map[string]interface{}
	subnets         []map[string]interface{}
	routingTables   []routingTable
	projectCounts   map[string]int    // resource -> count, for project-wide resources
	projectFailures map[string]string // resource -> error, when the list call failed
}

// projectQuotas lists project-wide resources: name, list path, response key, default limit.
var projectQuotas = []struct {
	name, path, key string
	limit           int
}{
	{"internet gateway", "/v2.0/internetgateways", "internetgateways", 3},
	{"NAT gateway", "/v2.0/natgateways", "natgateways", 3},
	{"service gateway", "/v2.0/gateways/servicegateways", "servicegateways", 10},
	{"user-defined service endpoint", "/v2.0/gateways/myserviceendpoints", "myserviceendpoints", 5},
	{"floating IP", "/v2.0/floatingips", "floatingips", 50},
	{"load balancer", "/v2.0/lbaas/loadbalancers", "loadbalancers", 10},
}

const (
	vpcDefaultLimit         = 3
	subnetPerVPCLimit       = 10
	routingTablePerVPCLimit = 10
)

// buildQuotaRows turns listed resources into quota rows. Only VPCs owned by
// the tenant count (the shared public network is listed for everyone).
func buildQuotaRows(in quotaInputs) []quotaRow {
	var own []string
	names := map[string]string{}
	for _, v := range in.vpcs {
		id, _ := v["id"].(string)
		tenant, _ := v["tenant_id"].(string)
		external, _ := v["router:external"].(bool)
		if external || (in.tenantID != "" && tenant != in.tenantID) {
			continue
		}
		own = append(own, id)
		names[id], _ = v["name"].(string)
	}
	sort.Strings(own)

	used := func(name string, n int) string {
		if e, ok := in.projectFailures[name]; ok {
			return "unknown (" + e + ")"
		}
		return strconv.Itoa(n)
	}
	rows := []quotaRow{{"VPC", "project", used("VPC", len(own)), vpcDefaultLimit}}

	subnets := map[string]int{}
	for _, s := range in.subnets {
		if id, _ := s["vpc_id"].(string); id != "" {
			subnets[id]++
		}
	}
	tables := map[string]int{}
	for _, t := range in.routingTables {
		for _, v := range t.VPCs {
			tables[v.ID]++
		}
	}
	for _, id := range own {
		scope := "VPC " + names[id] + " (" + id + ")"
		rows = append(rows,
			quotaRow{"subnet", scope, used("subnet", subnets[id]), subnetPerVPCLimit},
			quotaRow{"routing table", scope, used("routing table", tables[id]), routingTablePerVPCLimit})
	}

	for _, q := range projectQuotas {
		rows = append(rows, quotaRow{q.name, "project", used(q.name, in.projectCounts[q.name]), q.limit})
	}
	return rows
}

func init() {
	networkCmd.AddCommand(networkDescribeQuotasCmd)
}

var networkDescribeQuotasCmd = &cobra.Command{
	Use:   "describe-quotas",
	Short: "Show network resource usage against the default project limits",
	Long: `Show how many network resources the project uses next to the default
limits of the NHN Cloud resource policy (per region).

NHN Cloud has no public quota API for these resources, so the limits are the
published defaults. If your project's limits were raised, the real limit is
higher. Creating beyond the limit fails with 409 OverQuota.`,
	Run: func(cmd *cobra.Command, args []string) {
		in := quotaInputs{
			tenantID:        getTenantID(),
			projectCounts:   map[string]int{},
			projectFailures: map[string]string{},
		}
		// A failed list marks its rows unknown instead of stopping the command.
		listInto := func(name, path string, out interface{}) {
			if err := networkRequest(context.Background(), http.MethodGet, path, nil, out); err != nil {
				in.projectFailures[name] = "list failed"
				if debug {
					fmt.Fprintf(os.Stderr, "[DEBUG] %s: %v\n", name, err)
				}
			}
		}
		var vpcs struct {
			VPCs []map[string]interface{} `json:"vpcs"`
		}
		listInto("VPC", "/v2.0/vpcs", &vpcs)
		var subnets struct {
			Subnets []map[string]interface{} `json:"vpcsubnets"`
		}
		listInto("subnet", "/v2.0/vpcsubnets", &subnets)
		var tables struct {
			RoutingTables []routingTable `json:"routingtables"`
		}
		listInto("routing table", "/v2.0/routingtables?detail=true", &tables)
		in.vpcs, in.subnets, in.routingTables = vpcs.VPCs, subnets.Subnets, tables.RoutingTables

		for _, q := range projectQuotas {
			var res map[string]json.RawMessage
			var items []json.RawMessage
			err := networkRequest(context.Background(), http.MethodGet, q.path, nil, &res)
			if err == nil {
				err = json.Unmarshal(res[q.key], &items)
			}
			if err != nil {
				in.projectFailures[q.name] = "list failed"
				if debug {
					fmt.Fprintf(os.Stderr, "[DEBUG] %s: %v\n", q.name, err)
				}
				continue
			}
			in.projectCounts[q.name] = len(items)
		}

		rows := buildQuotaRows(in)
		emit(map[string]interface{}{"quotas": rows}, func() {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "RESOURCE\tSCOPE\tUSED\tDEFAULT_LIMIT")
			for _, r := range rows {
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", r.Resource, r.Scope, r.Used, r.Limit)
			}
			w.Flush()
			fmt.Println("\nDefault limits from the public resource policy; raised project limits are not visible here.")
		})
	},
}
