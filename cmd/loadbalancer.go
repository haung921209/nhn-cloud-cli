package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/loadbalancer"
	"github.com/spf13/cobra"
)

var loadbalancerCmd = &cobra.Command{
	Use:     "loadbalancer",
	Aliases: []string{"lb"},
	Short:   "Manage Load Balancers",
}

var lbListenerCmd = &cobra.Command{
	Use:   "listener",
	Short: "Manage listeners",
}

var lbPoolCmd = &cobra.Command{
	Use:   "pool",
	Short: "Manage pools",
}

var lbMemberCmd = &cobra.Command{
	Use:   "member",
	Short: "Manage pool members",
}

var lbHealthMonitorCmd = &cobra.Command{
	Use:     "healthmonitor",
	Aliases: []string{"hm"},
	Short:   "Manage health monitors",
}

func init() {
	rootCmd.AddCommand(loadbalancerCmd)

	loadbalancerCmd.AddCommand(lbListCmd)
	loadbalancerCmd.AddCommand(lbGetCmd)
	loadbalancerCmd.AddCommand(lbCreateCmd)
	loadbalancerCmd.AddCommand(lbDeleteCmd)
	loadbalancerCmd.AddCommand(lbListenerCmd)
	loadbalancerCmd.AddCommand(lbPoolCmd)
	loadbalancerCmd.AddCommand(lbHealthMonitorCmd)

	lbCreateCmd.Flags().String("name", "", "Load balancer name (required)")
	lbCreateCmd.Flags().String("description", "", "Description")
	lbCreateCmd.Flags().String("subnet-id", "", "VIP subnet ID (required)")
	lbCreateCmd.Flags().String("vip-address", "", "VIP address (optional)")
	lbCreateCmd.Flags().String("provider", "", "Provider (optional)")
	lbCreateCmd.MarkFlagRequired("name")
	lbCreateCmd.MarkFlagRequired("subnet-id")

	lbListenerCmd.AddCommand(lbListenerListCmd)
	lbListenerCmd.AddCommand(lbListenerGetCmd)
	lbListenerCmd.AddCommand(lbListenerCreateCmd)
	lbListenerCmd.AddCommand(lbListenerDeleteCmd)

	lbListenerCreateCmd.Flags().String("name", "", "Listener name (required)")
	lbListenerCreateCmd.Flags().String("lb-id", "", "Load balancer ID (required)")
	lbListenerCreateCmd.Flags().String("protocol", "TCP", "Protocol (TCP/HTTP/HTTPS/TERMINATED_HTTPS)")
	lbListenerCreateCmd.Flags().Int("port", 80, "Protocol port (required)")
	lbListenerCreateCmd.Flags().String("pool-id", "", "Default pool ID")
	lbListenerCreateCmd.Flags().Int("connection-limit", -1, "Connection limit")
	lbListenerCreateCmd.MarkFlagRequired("name")
	lbListenerCreateCmd.MarkFlagRequired("lb-id")
	lbListenerCreateCmd.MarkFlagRequired("port")

	lbPoolCmd.AddCommand(lbPoolListCmd)
	lbPoolCmd.AddCommand(lbPoolGetCmd)
	lbPoolCmd.AddCommand(lbPoolCreateCmd)
	lbPoolCmd.AddCommand(lbPoolDeleteCmd)
	lbPoolCmd.AddCommand(lbMemberCmd)

	lbPoolCreateCmd.Flags().String("name", "", "Pool name (required)")
	lbPoolCreateCmd.Flags().String("protocol", "TCP", "Protocol (TCP/HTTP/HTTPS/PROXY)")
	lbPoolCreateCmd.Flags().String("algorithm", "ROUND_ROBIN", "LB algorithm (ROUND_ROBIN/LEAST_CONNECTIONS/SOURCE_IP)")
	lbPoolCreateCmd.Flags().String("lb-id", "", "Load balancer ID")
	lbPoolCreateCmd.Flags().String("listener-id", "", "Listener ID")
	lbPoolCreateCmd.MarkFlagRequired("name")

	lbMemberCmd.AddCommand(lbMemberListCmd)
	lbMemberCmd.AddCommand(lbMemberGetCmd)
	lbMemberCmd.AddCommand(lbMemberCreateCmd)
	lbMemberCmd.AddCommand(lbMemberDeleteCmd)

	lbMemberCreateCmd.Flags().String("pool-id", "", "Pool ID (required)")
	lbMemberCreateCmd.Flags().String("address", "", "Member IP address (required)")
	lbMemberCreateCmd.Flags().Int("port", 80, "Member port (required)")
	lbMemberCreateCmd.Flags().Int("weight", 1, "Member weight")
	lbMemberCreateCmd.Flags().String("subnet-id", "", "Subnet ID")
	lbMemberCreateCmd.MarkFlagRequired("pool-id")
	lbMemberCreateCmd.MarkFlagRequired("address")
	lbMemberCreateCmd.MarkFlagRequired("port")

	lbHealthMonitorCmd.AddCommand(lbHMListCmd)
	lbHealthMonitorCmd.AddCommand(lbHMGetCmd)
	lbHealthMonitorCmd.AddCommand(lbHMCreateCmd)
	lbHealthMonitorCmd.AddCommand(lbHMDeleteCmd)

	lbHMCreateCmd.Flags().String("pool-id", "", "Pool ID (required)")
	lbHMCreateCmd.Flags().String("type", "TCP", "Monitor type (TCP/HTTP/HTTPS/PING)")
	lbHMCreateCmd.Flags().Int("delay", 5, "Delay between checks (seconds)")
	lbHMCreateCmd.Flags().Int("timeout", 5, "Check timeout (seconds)")
	lbHMCreateCmd.Flags().Int("max-retries", 3, "Max retries before marking DOWN")
	lbHMCreateCmd.Flags().String("http-method", "GET", "HTTP method (for HTTP/HTTPS)")
	lbHMCreateCmd.Flags().String("url-path", "/", "URL path (for HTTP/HTTPS)")
	lbHMCreateCmd.Flags().String("expected-codes", "200", "Expected HTTP codes")
	lbHMCreateCmd.MarkFlagRequired("pool-id")
}

func newLBClient() *loadbalancer.Client {
	return loadbalancer.NewClient(getRegion(), getIdentityCreds(), nil, debug)
}

var lbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all load balancers",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		result, err := client.ListLoadBalancers(context.Background())
		if err != nil {
			exitWithError("Failed to list load balancers", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tVIP_ADDRESS\tSTATUS\tPROVISIONING")
		for _, lb := range result.LoadBalancers {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				lb.ID, lb.Name, lb.VIPAddress, lb.OperatingStatus, lb.ProvisioningStatus)
		}
		w.Flush()
	},
}

var lbGetCmd = &cobra.Command{
	Use:   "get [lb-id]",
	Short: "Get load balancer details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		result, err := client.GetLoadBalancer(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get load balancer", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		lb := result.LoadBalancer
		fmt.Printf("ID:           %s\n", lb.ID)
		fmt.Printf("Name:         %s\n", lb.Name)
		fmt.Printf("Description:  %s\n", lb.Description)
		fmt.Printf("VIP Address:  %s\n", lb.VIPAddress)
		fmt.Printf("VIP Subnet:   %s\n", lb.VIPSubnetID)
		fmt.Printf("Status:       %s\n", lb.OperatingStatus)
		fmt.Printf("Provisioning: %s\n", lb.ProvisioningStatus)
		fmt.Printf("Provider:     %s\n", lb.Provider)
		fmt.Printf("Created:      %s\n", lb.CreatedAt)
	},
}

var lbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new load balancer",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		vipAddress, _ := cmd.Flags().GetString("vip-address")
		provider, _ := cmd.Flags().GetString("provider")

		input := &loadbalancer.CreateLoadBalancerInput{
			Name:        name,
			Description: description,
			VIPSubnetID: subnetID,
			VIPAddress:  vipAddress,
			Provider:    provider,
		}

		result, err := client.CreateLoadBalancer(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create load balancer", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Load balancer created: %s\n", result.LoadBalancer.ID)
		fmt.Printf("Name: %s\n", result.LoadBalancer.Name)
		fmt.Printf("VIP:  %s\n", result.LoadBalancer.VIPAddress)
	},
}

var lbDeleteCmd = &cobra.Command{
	Use:   "delete [lb-id]",
	Short: "Delete a load balancer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		if err := client.DeleteLoadBalancer(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete load balancer", err)
		}
		fmt.Printf("Load balancer %s deleted\n", args[0])
	},
}

var lbListenerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all listeners",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		result, err := client.ListListeners(context.Background())
		if err != nil {
			exitWithError("Failed to list listeners", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tPROTOCOL\tPORT\tLB_ID\tSTATUS")
		for _, l := range result.Listeners {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%s\n",
				l.ID, l.Name, l.Protocol, l.ProtocolPort, l.LoadBalancerID, l.OperatingStatus)
		}
		w.Flush()
	},
}

var lbListenerGetCmd = &cobra.Command{
	Use:   "get [listener-id]",
	Short: "Get listener details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		result, err := client.GetListener(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get listener", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		l := result.Listener
		fmt.Printf("ID:              %s\n", l.ID)
		fmt.Printf("Name:            %s\n", l.Name)
		fmt.Printf("Protocol:        %s\n", l.Protocol)
		fmt.Printf("Port:            %d\n", l.ProtocolPort)
		fmt.Printf("Load Balancer:   %s\n", l.LoadBalancerID)
		fmt.Printf("Default Pool:    %s\n", l.DefaultPoolID)
		fmt.Printf("Connection Limit:%d\n", l.ConnectionLimit)
		fmt.Printf("Status:          %s\n", l.OperatingStatus)
	},
}

var lbListenerCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new listener",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		name, _ := cmd.Flags().GetString("name")
		lbID, _ := cmd.Flags().GetString("lb-id")
		protocol, _ := cmd.Flags().GetString("protocol")
		port, _ := cmd.Flags().GetInt("port")
		poolID, _ := cmd.Flags().GetString("pool-id")
		connLimit, _ := cmd.Flags().GetInt("connection-limit")

		input := &loadbalancer.CreateListenerInput{
			Name:            name,
			LoadBalancerID:  lbID,
			Protocol:        protocol,
			ProtocolPort:    port,
			DefaultPoolID:   poolID,
			ConnectionLimit: connLimit,
		}

		result, err := client.CreateListener(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create listener", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Listener created: %s\n", result.Listener.ID)
		fmt.Printf("Name: %s\n", result.Listener.Name)
	},
}

var lbListenerDeleteCmd = &cobra.Command{
	Use:   "delete [listener-id]",
	Short: "Delete a listener",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		if err := client.DeleteListener(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete listener", err)
		}
		fmt.Printf("Listener %s deleted\n", args[0])
	},
}

var lbPoolListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pools",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		result, err := client.ListPools(context.Background())
		if err != nil {
			exitWithError("Failed to list pools", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tPROTOCOL\tALGORITHM\tSTATUS")
		for _, p := range result.Pools {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				p.ID, p.Name, p.Protocol, p.LBAlgorithm, p.OperatingStatus)
		}
		w.Flush()
	},
}

var lbPoolGetCmd = &cobra.Command{
	Use:   "get [pool-id]",
	Short: "Get pool details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		result, err := client.GetPool(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get pool", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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
	},
}

var lbPoolCreateCmd = &cobra.Command{
	Use:   "create",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Pool created: %s\n", result.Pool.ID)
		fmt.Printf("Name: %s\n", result.Pool.Name)
	},
}

var lbPoolDeleteCmd = &cobra.Command{
	Use:   "delete [pool-id]",
	Short: "Delete a pool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		if err := client.DeletePool(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete pool", err)
		}
		fmt.Printf("Pool %s deleted\n", args[0])
	},
}

var lbMemberListCmd = &cobra.Command{
	Use:   "list [pool-id]",
	Short: "List members in a pool",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		result, err := client.ListMembers(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to list members", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tADDRESS\tPORT\tWEIGHT\tSTATUS")
		for _, m := range result.Members {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%s\n",
				m.ID, m.Name, m.Address, m.ProtocolPort, m.Weight, m.OperatingStatus)
		}
		w.Flush()
	},
}

var lbMemberGetCmd = &cobra.Command{
	Use:   "get [pool-id] [member-id]",
	Short: "Get member details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		result, err := client.GetMember(context.Background(), args[0], args[1])
		if err != nil {
			exitWithError("Failed to get member", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
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
	},
}

var lbMemberCreateCmd = &cobra.Command{
	Use:   "create",
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
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Member created: %s\n", result.Member.ID)
		fmt.Printf("Address: %s:%d\n", result.Member.Address, result.Member.ProtocolPort)
	},
}

var lbMemberDeleteCmd = &cobra.Command{
	Use:   "delete [pool-id] [member-id]",
	Short: "Delete a member",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		if err := client.DeleteMember(context.Background(), args[0], args[1]); err != nil {
			exitWithError("Failed to delete member", err)
		}
		fmt.Printf("Member %s deleted\n", args[1])
	},
}

var lbHMListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all health monitors",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		result, err := client.ListHealthMonitors(context.Background())
		if err != nil {
			exitWithError("Failed to list health monitors", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tPOOL_ID\tSTATUS")
		for _, h := range result.HealthMonitors {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				h.ID, h.Name, h.Type, h.PoolID, h.OperatingStatus)
		}
		w.Flush()
	},
}

var lbHMGetCmd = &cobra.Command{
	Use:   "get [healthmonitor-id]",
	Short: "Get health monitor details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		result, err := client.GetHealthMonitor(context.Background(), args[0])
		if err != nil {
			exitWithError("Failed to get health monitor", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		h := result.HealthMonitor
		fmt.Printf("ID:         %s\n", h.ID)
		fmt.Printf("Name:       %s\n", h.Name)
		fmt.Printf("Type:       %s\n", h.Type)
		fmt.Printf("Pool ID:    %s\n", h.PoolID)
		fmt.Printf("Delay:      %d\n", h.Delay)
		fmt.Printf("Timeout:    %d\n", h.Timeout)
		fmt.Printf("Max Retries:%d\n", h.MaxRetries)
		fmt.Printf("URL Path:   %s\n", h.URLPath)
		fmt.Printf("Status:     %s\n", h.OperatingStatus)
	},
}

var lbHMCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new health monitor",
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		poolID, _ := cmd.Flags().GetString("pool-id")
		monitorType, _ := cmd.Flags().GetString("type")
		delay, _ := cmd.Flags().GetInt("delay")
		timeout, _ := cmd.Flags().GetInt("timeout")
		maxRetries, _ := cmd.Flags().GetInt("max-retries")
		httpMethod, _ := cmd.Flags().GetString("http-method")
		urlPath, _ := cmd.Flags().GetString("url-path")
		expectedCodes, _ := cmd.Flags().GetString("expected-codes")

		input := &loadbalancer.CreateHealthMonitorInput{
			PoolID:        poolID,
			Type:          monitorType,
			Delay:         delay,
			Timeout:       timeout,
			MaxRetries:    maxRetries,
			HTTPMethod:    httpMethod,
			URLPath:       urlPath,
			ExpectedCodes: expectedCodes,
		}

		result, err := client.CreateHealthMonitor(context.Background(), input)
		if err != nil {
			exitWithError("Failed to create health monitor", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Health monitor created: %s\n", result.HealthMonitor.ID)
		fmt.Printf("Type: %s\n", result.HealthMonitor.Type)
	},
}

var lbHMDeleteCmd = &cobra.Command{
	Use:   "delete [healthmonitor-id]",
	Short: "Delete a health monitor",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newLBClient()
		if err := client.DeleteHealthMonitor(context.Background(), args[0]); err != nil {
			exitWithError("Failed to delete health monitor", err)
		}
		fmt.Printf("Health monitor %s deleted\n", args[0])
	},
}
