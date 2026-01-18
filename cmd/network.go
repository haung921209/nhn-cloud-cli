package cmd

import (
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/spf13/cobra"
)

var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "Manage VPCs and subnets (Legacy - use describe-vpcs, etc.)",
}

var networkSecurityGroupCmd = &cobra.Command{
	Use:     "security-group",
	Aliases: []string{"sg"},
	Short:   "Manage network security groups (Legacy - use describe-security-groups, etc.)",
}

var floatingIPCmd = &cobra.Command{
	Use:     "floating-ip",
	Aliases: []string{"fip"},
	Short:   "Manage floating IPs (Legacy - use describe-floating-ips, etc.)",
}

var portCmd = &cobra.Command{
	Use:   "port",
	Short: "Manage network ports",
}

func init() {
	// Root commands for network
	// We are keeping these as parent commands because new AWS commands attach to them?
	// Wait, AWS commands are: describe-vpcs.
	// If I attach describe-vpcs to vpcCmd, it becomes "nhncloud vpc describe-vpcs".
	// That is NOT AWS style. AWS style is "nhncloud describe-vpcs" (flat).
	// OR "nhncloud network describe-vpcs"?
	// "nhncloud compute describe-instances" was attached to "computeCmd".
	// "nhncloud describe-instances" logic: RDS commands are "nhncloud rds-mysql describe-db-instances".
	// So "nhncloud network describe-vpcs" is probably the goal.
	// In the previous steps, I attached to vpcCmd?
	// Let's check: in network_vpc.go I did: `vpcCmd.AddCommand(describeVPCsCmd)`.
	// THIS MEANS "nhncloud vpc describe-vpcs". This is nested, not flat under network.
	// In Compute, I attached to `computeCmd`. `nhncloud compute describe-instances`. That is flat under compute.

	// Issue: "vpc" IS the service command in the old model? No "vpc" was a subcommand of ROOT?
	// In original code: `rootCmd.AddCommand(vpcCmd)`.
	// So `nhncloud vpc list`.
	// If I want `nhncloud network describe-vpcs`, I need a `networkCmd`.
	// But there was no `networkCmd` in original code! Each resource was a top level command?
	// `rootCmd.AddCommand(vpcCmd)` -> `nhncloud vpc ...`
	// `rootCmd.AddCommand(networkSecurityGroupCmd)` -> `nhncloud security-group ...`
	// `rootCmd.AddCommand(floatingIPCmd)` -> `nhncloud floating-ip ...`

	// I should probably CREATE a `networkCmd` to group them if I want `nhncloud network describe-vpcs`?
	// OR keep them top level?
	// "Compute & Network" implies they are services.
	// `nhncloud compute ...` exists.
	// `nhncloud network ...` does NOT exist.
	// RDS is `nhncloud rds-mysql ...`.

	// User request: "AWS CLI Alignment".
	// AWS is `aws ec2 describe-instances`. `aws ec2` is the service.
	// `aws ec2 describe-vpcs`.
	// So `vpcs` belong to `ec2` (or `network` service).
	// If I want `nhncloud compute describe-vpcs`, I should attach to `computeCmd`.
	// But VPC is network.
	// Maybe I should introduce `networkCmd`?
	// If I do `nhncloud network describe-vpcs`.

	// However, I already attached them to `vpcCmd` in `network_vpc.go`.
	// If `vpcCmd` is added to `rootCmd`, then it is `nhncloud vpc describe-vpcs`.
	// Just `describe-vpcs` as a subcommand of `vpc` is weird. `vpc describe-vpcs`. Redundant.

	// I should attach `describeVPCsCmd` to `networkCmd` or `computeCmd` (since NHN might group them).
	// Or `rootCmd`? AWS CLI has one root? No `aws` is root.

	// PROPOSAL:
	// Create `networkCmd` (alias `vpc`?).
	// `nhncloud network describe-vpcs`.
	// `nhncloud network create-vpc`.
	// `nhncloud network describe-subnets`.
	// `nhncloud network describe-security-groups`.

	// To do this:
	// 1. Define `networkCmd` in `cmd/network.go`.
	// 2. Attach `describeVPCsCmd` etc. to `networkCmd` instead of `vpcCmd`.
	// 3. Update `network_vpc.go` etc to use `networkCmd`.
	// 4. Update imports.

	// Wait, I already wrote `network_vpc.go` using `vpcCmd`.
	// I need to change `vpcCmd` to `networkCmd` in those files.
	// Or I can define `vpcCmd` as `networkCmd`?
	// var networkCmd = &cobra.Command{ Use: "network", ... }
	// And usage in `network_vpc.go` -> `networkCmd.AddCommand(...)`.

	// I will update `cmd/network.go` to define `networkCmd`.
	// And I will update `cmd/network_vpc.go`, `cmd/network_security_groups.go`, `cmd/network_floating_ip.go` to use `networkCmd`.

	rootCmd.AddCommand(networkCmd)

	// Port?
	networkCmd.AddCommand(portCmd)
}

// Define networkCmd
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage Network resources (VPC, Subnets, SG, FIP)",
}

// Helper
func getIdentityCreds() credentials.IdentityCredentials {
	return credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
}
