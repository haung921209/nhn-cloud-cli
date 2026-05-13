// helpers/nks-control — NKS cluster lifecycle helper used by the
// nks-rds-p99 scenario.
//
// Why this exists: the default-build nhn-cloud-cli only ships node-group
// commands; cluster create/delete/get-kubeconfig live in cmd/nks_clusters.go
// behind the `cli_full` build tag and that file currently doesn't compile
// against the latest SDK because of unrelated drift in other cli_full
// commands. Rather than ungate cli_full just for this scenario, we use
// the SDK directly here — same pattern as scenarios/.../helpers/attach-fip.
//
// Auth env (all required):
//
//	NHNCLOUD_USERNAME, NHNCLOUD_PASSWORD, NHNCLOUD_TENANT_ID
//	NHNCLOUD_REGION  (default kr1)
//
// Subcommands:
//
//	list-templates                     # tab-separated id\tname\tcoe\tserver_type
//	list-versions                      # one supported_k8s key per line
//	create   --name X --template-id T  --node-count N --keypair K --subnet-id S [--node-flavor-id F]
//	get      --cluster-id ID            # JSON cluster object on stdout
//	wait     --cluster-id ID --for-state ACTIVE|DELETED [--timeout 60m]
//	kubeconfig --cluster-id ID [--out path]
//	delete   --cluster-id ID            # submit delete; does NOT block (use `wait --for-state DELETED`)
//
// Stdout for `create` (machine-parseable, one KEY=VALUE per line):
//
//	CLUSTER_ID=<uuid>
//	CLUSTER_NAME=<name>
//
// Exit codes: 0 success, 1 generic failure, 2 bad arguments, 3 timeout.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/nks"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/image"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/vpc"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	sub := os.Args[1]
	args := os.Args[2:]

	creds, region, err := loadCreds()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(2)
	}
	debug := os.Getenv("NKS_CONTROL_DEBUG") != ""
	client := nks.NewClient(region, creds, nil, debug)

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Minute)
	defer cancel()

	switch sub {
	case "list-templates":
		os.Exit(runListTemplates(ctx, client))
	case "list-versions":
		os.Exit(runListVersions(ctx, client))
	case "create":
		os.Exit(runCreate(ctx, client, args))
	case "get":
		os.Exit(runGet(ctx, client, args))
	case "wait":
		os.Exit(runWait(ctx, client, args))
	case "kubeconfig":
		os.Exit(runKubeconfig(ctx, client, args))
	case "delete":
		os.Exit(runDelete(ctx, client, args))
	case "resolve-vpc":
		// uses VPC client, not NKS — handled separately
		os.Exit(runResolveVPC(ctx, creds, region, args))
	case "list-node-images":
		os.Exit(runListNodeImages(ctx, creds, region, args))
	case "-h", "--help", "help":
		usage()
		return
	default:
		fmt.Fprintf(os.Stderr, "ERROR: unknown subcommand %q\n", sub)
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: nks-control <subcommand> [flags]")
	fmt.Fprintln(os.Stderr, "  list-templates | list-versions")
	fmt.Fprintln(os.Stderr, "  create    --name N --template-id T --node-count K --keypair P --subnet-id S [--node-flavor-id F]")
	fmt.Fprintln(os.Stderr, "  get       --cluster-id ID")
	fmt.Fprintln(os.Stderr, "  wait      --cluster-id ID --for-state ACTIVE|DELETED [--timeout 60m]")
	fmt.Fprintln(os.Stderr, "  kubeconfig --cluster-id ID [--out path]")
	fmt.Fprintln(os.Stderr, "  delete    --cluster-id ID")
	fmt.Fprintln(os.Stderr, "  resolve-vpc --subnet-id ID   # prints VPC_ID=<uuid> (NKS_NETWORK_ID for create)")
	fmt.Fprintln(os.Stderr, "  list-node-images [--name-glob=*]   # NKS-filtered Glance image list (id\\tname)")
}

func loadCreds() (credentials.IdentityCredentials, string, error) {
	u := os.Getenv("NHNCLOUD_USERNAME")
	p := os.Getenv("NHNCLOUD_PASSWORD")
	t := os.Getenv("NHNCLOUD_TENANT_ID")
	if u == "" || p == "" || t == "" {
		return nil, "", errors.New("NHNCLOUD_USERNAME / NHNCLOUD_PASSWORD / NHNCLOUD_TENANT_ID must all be set")
	}
	region := os.Getenv("NHNCLOUD_REGION")
	if region == "" {
		region = "kr1"
	}
	return credentials.NewStaticIdentity(u, p, t), region, nil
}

func runListTemplates(ctx context.Context, c *nks.Client) int {
	out, err := c.ListClusterTemplates(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ListClusterTemplates:", err)
		return 1
	}
	for _, t := range out.ClusterTemplates {
		fmt.Printf("%s\t%s\t%s\t%s\n", t.ID, t.Name, t.COE, t.ServerType)
	}
	return 0
}

func runListVersions(ctx context.Context, c *nks.Client) int {
	out, err := c.GetSupportedVersions(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: GetSupportedVersions:", err)
		return 1
	}
	for k, ok := range out.SupportedK8s {
		if ok {
			fmt.Println(k)
		}
	}
	return 0
}

// repeatable flag: --label k=v  (collected into labels map)
type labelFlags []string

func (l *labelFlags) String() string     { return strings.Join(*l, ",") }
func (l *labelFlags) Set(v string) error { *l = append(*l, v); return nil }

func runCreate(ctx context.Context, c *nks.Client, args []string) int {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	name := fs.String("name", "", "cluster name (required)")
	// NHN's API requires the literal string "iaas_console" — the OpenStack-style
	// template UUID listing endpoint is not exposed (404). Default to that;
	// allow override only if NHN adds true templates later.
	tplID := fs.String("template-id", "iaas_console", `cluster_template_id (NHN requires literal "iaas_console")`)
	nodeCount := fs.Int("node-count", 1, "worker node count")
	nodeFlavor := fs.String("node-flavor-id", "", "worker flavor_id (root field, required by spec)")
	keypair := fs.String("keypair", "", "keypair name (required)")
	subnetID := fs.String("subnet-id", "", "fixed_subnet UUID (required)")
	networkID := fs.String("network-id", "", "fixed_network UUID (required by spec)")
	var labels labelFlags
	fs.Var(&labels, "label", `repeatable: --label key=value (NHN requires availability_zone, kube_tag, node_image, boot_volume_type, boot_volume_size, cert_manager_api, ca_enable, master_lb_floating_ip_enabled at minimum)`)
	_ = fs.Parse(args)

	if *name == "" || *keypair == "" || *subnetID == "" || *networkID == "" || *nodeFlavor == "" {
		fmt.Fprintln(os.Stderr, "ERROR: --name, --keypair, --subnet-id, --network-id, --node-flavor-id are all required")
		return 2
	}

	labelMap := make(map[string]string, len(labels))
	for _, kv := range labels {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			fmt.Fprintf(os.Stderr, "ERROR: bad --label %q (expected key=value)\n", kv)
			return 2
		}
		labelMap[parts[0]] = parts[1]
	}

	in := &nks.CreateClusterInput{
		Name:              *name,
		ClusterTemplateID: *tplID,
		NodeCount:         *nodeCount,
		KeyPair:           *keypair,
		SubnetID:          *subnetID,
		NetworkID:         *networkID,
		FlavorID:          *nodeFlavor,
		Labels:            labelMap,
	}

	out, err := c.CreateCluster(ctx, in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: CreateCluster:", err)
		return 1
	}
	fmt.Printf("CLUSTER_ID=%s\n", out.Cluster.ID)
	fmt.Printf("CLUSTER_NAME=%s\n", out.Cluster.Name)
	return 0
}

func runGet(ctx context.Context, c *nks.Client, args []string) int {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	id := fs.String("cluster-id", "", "cluster UUID (required)")
	_ = fs.Parse(args)
	if *id == "" {
		fmt.Fprintln(os.Stderr, "ERROR: --cluster-id is required")
		return 2
	}
	out, err := c.GetCluster(ctx, *id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: GetCluster:", err)
		return 1
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out.Cluster); err != nil {
		return 1
	}
	return 0
}

// runWait polls GetCluster every 15s until the requested terminal state.
//
// "ACTIVE": status field == "CREATE_COMPLETE" or "UPDATE_COMPLETE" (NHN
// reports both for healthy clusters depending on history; either is fine).
// Anything containing "FAILED" or "ERROR" terminates non-zero.
//
// "DELETED": GetCluster returns a 404-shaped error, OR status contains
// "DELETE_COMPLETE". 404 detection is best-effort string match — if the
// SDK ever wraps in a typed not-found error we should switch to errors.Is.
func runWait(ctx context.Context, c *nks.Client, args []string) int {
	fs := flag.NewFlagSet("wait", flag.ExitOnError)
	id := fs.String("cluster-id", "", "cluster UUID (required)")
	state := fs.String("for-state", "ACTIVE", "ACTIVE | DELETED")
	timeout := fs.Duration("timeout", 60*time.Minute, "max wait")
	interval := fs.Duration("interval", 15*time.Second, "poll interval")
	_ = fs.Parse(args)
	if *id == "" {
		fmt.Fprintln(os.Stderr, "ERROR: --cluster-id is required")
		return 2
	}

	deadline := time.Now().Add(*timeout)
	last := ""
	for {
		if time.Now().After(deadline) {
			fmt.Fprintf(os.Stderr, "\nERROR: timeout waiting for state=%s (last status=%q)\n", *state, last)
			return 3
		}
		out, err := c.GetCluster(ctx, *id)
		if err != nil {
			lower := strings.ToLower(err.Error())
			if *state == "DELETED" && (strings.Contains(lower, "404") || strings.Contains(lower, "not found")) {
				fmt.Println("\n[wait] cluster fully deleted (404)")
				return 0
			}
			fmt.Fprintf(os.Stderr, "\n[wait] GetCluster err: %v — retrying\n", err)
			time.Sleep(*interval)
			continue
		}
		st := out.Cluster.Status
		if st != last {
			fmt.Fprintf(os.Stderr, "[wait] status=%s\n", st)
			last = st
		} else {
			fmt.Fprint(os.Stderr, ".")
		}
		switch *state {
		case "ACTIVE":
			if strings.Contains(st, "COMPLETE") && !strings.Contains(st, "DELETE") {
				fmt.Fprintln(os.Stderr)
				return 0
			}
		case "DELETED":
			if strings.Contains(st, "DELETE_COMPLETE") {
				fmt.Fprintln(os.Stderr)
				return 0
			}
		}
		if strings.Contains(st, "FAILED") || strings.Contains(st, "ERROR") {
			fmt.Fprintf(os.Stderr, "\nERROR: cluster reached terminal-failure status=%s reason=%q\n", st, out.Cluster.StatusReason)
			return 1
		}
		time.Sleep(*interval)
	}
}

func runKubeconfig(ctx context.Context, c *nks.Client, args []string) int {
	fs := flag.NewFlagSet("kubeconfig", flag.ExitOnError)
	id := fs.String("cluster-id", "", "cluster UUID (required)")
	out := fs.String("out", "", "write kubeconfig to file (default: stdout)")
	_ = fs.Parse(args)
	if *id == "" {
		fmt.Fprintln(os.Stderr, "ERROR: --cluster-id is required")
		return 2
	}
	res, err := c.GetKubeconfig(ctx, *id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: GetKubeconfig:", err)
		return 1
	}
	if *out != "" {
		// 0600: kubeconfig embeds bearer-equivalent credentials.
		if err := os.WriteFile(*out, []byte(res.Kubeconfig), 0600); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR: write kubeconfig:", err)
			return 1
		}
		fmt.Fprintf(os.Stderr, "[kubeconfig] wrote %d bytes to %s\n", len(res.Kubeconfig), *out)
		return 0
	}
	fmt.Print(res.Kubeconfig)
	return 0
}

// runListNodeImages calls the Glance compute images endpoint with NHN's
// "NKS-only" query-string filter (`nhncloud_allow_nks_cpu_flavor=true&
// visibility=public`, per docs/api-specs/container/nks.md "베이스 이미지 UUID").
// Output: tab-separated id\tname\tstatus, one row per image.
//
// Optional `--name-glob` is a substring match on Name; if non-empty only
// matching rows are emitted (useful for "give me Ubuntu rows only").
func runListNodeImages(ctx context.Context, creds credentials.IdentityCredentials, region string, args []string) int {
	fs := flag.NewFlagSet("list-node-images", flag.ExitOnError)
	nameGlob := fs.String("name-glob", "", "substring filter on Name (case-sensitive)")
	_ = fs.Parse(args)

	c := image.NewClient(region, creds, nil, false)
	// Visibility goes through the typed field; the NHN-only flag rides on
	// ExtraParams (added in SDK v0.1.34).
	out, err := c.ListImages(ctx, &image.ListImagesInput{
		Visibility: "public",
		ExtraParams: map[string]string{
			"nhncloud_allow_nks_cpu_flavor": "true",
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ListImages (NKS-filtered):", err)
		return 1
	}
	for _, img := range out.Images {
		if *nameGlob != "" && !strings.Contains(img.Name, *nameGlob) {
			continue
		}
		fmt.Printf("%s\t%s\t%s\n", img.ID, img.Name, img.Status)
	}
	return 0
}

// runResolveVPC asks the VPC service for the parent network UUID of a subnet.
// NHN's API populates Subnet.VPCID (json:"vpc_id") — added to the SDK in
// v0.1.32. Output: VPC_ID=<uuid> on stdout, suitable for `eval $(... resolve-vpc)`
// or for the preflight step to grep out and write to state.env.
func runResolveVPC(ctx context.Context, creds credentials.IdentityCredentials, region string, args []string) int {
	fs := flag.NewFlagSet("resolve-vpc", flag.ExitOnError)
	subnetID := fs.String("subnet-id", "", "subnet UUID (required)")
	_ = fs.Parse(args)
	if *subnetID == "" {
		fmt.Fprintln(os.Stderr, "ERROR: --subnet-id is required")
		return 2
	}
	c := vpc.NewClient(region, creds, nil, false)
	out, err := c.GetSubnet(ctx, *subnetID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: GetSubnet:", err)
		return 1
	}
	vpcID := out.Subnet.VPCID
	if vpcID == "" {
		// Fall back to the OpenStack-style network_id if NHN ever stops
		// returning vpc_id (none of the live observations show that, but
		// safer than empty output).
		vpcID = out.Subnet.NetworkID
	}
	if vpcID == "" {
		fmt.Fprintln(os.Stderr, "ERROR: response had neither vpc_id nor network_id populated")
		return 1
	}
	fmt.Printf("VPC_ID=%s\n", vpcID)
	return 0
}

func runDelete(ctx context.Context, c *nks.Client, args []string) int {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	id := fs.String("cluster-id", "", "cluster UUID (required)")
	_ = fs.Parse(args)
	if *id == "" {
		fmt.Fprintln(os.Stderr, "ERROR: --cluster-id is required")
		return 2
	}
	if err := c.DeleteCluster(ctx, *id); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: DeleteCluster:", err)
		return 1
	}
	fmt.Println("DELETE_SUBMITTED=true")
	return 0
}
