package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/dnsplus"
	"github.com/spf13/cobra"
)

// dnsplusCmd represents the dnsplus command
var dnsplusCmd = &cobra.Command{
	Use:     "dnsplus",
	Aliases: []string{"dns", "dns-plus"},
	Short:   "Manage DNS Plus service",
	Long:    `Manage DNS zones, record sets, GSLB, and health checks in NHN Cloud DNS Plus.`,
}

// ================================
// Zone Commands
// ================================

var dnsplusZoneCmd = &cobra.Command{
	Use:     "zone",
	Aliases: []string{"zones"},
	Short:   "Manage DNS zones",
}

var dnsplusZoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List DNS zones",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListZones(ctx)
		if err != nil {
			return fmt.Errorf("failed to list zones: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.ZoneList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tRECORDS\tCREATED")
		for _, zone := range result.ZoneList {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				zone.ZoneID,
				zone.ZoneName,
				zone.ZoneStatus,
				zone.RecordSetCount,
				zone.CreatedAt.Format("2006-01-02"),
			)
		}
		return w.Flush()
	},
}

var dnsplusZoneCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a DNS zone",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.CreateZoneInput{
			ZoneName:    name,
			Description: description,
		}

		result, err := client.CreateZone(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to create zone: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Zone)
		}

		fmt.Printf("Zone created successfully: %s (%s)\n", result.Zone.ZoneName, result.Zone.ZoneID)
		return nil
	},
}

var dnsplusZoneUpdateCmd = &cobra.Command{
	Use:   "update <zone-id>",
	Short: "Update a DNS zone",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneID := args[0]
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.UpdateZoneInput{
			Description: description,
			ZoneStatus:  status,
		}

		result, err := client.UpdateZone(ctx, zoneID, input)
		if err != nil {
			return fmt.Errorf("failed to update zone: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Zone)
		}

		fmt.Printf("Zone updated successfully: %s\n", result.Zone.ZoneID)
		return nil
	},
}

var dnsplusZoneDeleteCmd = &cobra.Command{
	Use:   "delete <zone-id>...",
	Short: "Delete DNS zones",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeleteZones(ctx, args)
		if err != nil {
			return fmt.Errorf("failed to delete zones: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Zone deletion initiated for %d zone(s)\n", len(args))
		return nil
	},
}

// ================================
// Record Set Commands
// ================================

var dnsplusRecordsetCmd = &cobra.Command{
	Use:     "recordset",
	Aliases: []string{"recordsets", "record", "records"},
	Short:   "Manage DNS record sets",
}

var dnsplusRecordsetListCmd = &cobra.Command{
	Use:   "list <zone-id>",
	Short: "List record sets in a zone",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneID := args[0]
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListRecordSets(ctx, zoneID)
		if err != nil {
			return fmt.Errorf("failed to list record sets: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.RecordSetList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tTTL\tRECORDS")
		for _, rs := range result.RecordSetList {
			records := make([]string, len(rs.RecordList))
			for i, r := range rs.RecordList {
				records[i] = r.RecordContent
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				rs.RecordSetID,
				rs.RecordSetName,
				rs.RecordSetType,
				rs.TTL,
				strings.Join(records, ", "),
			)
		}
		return w.Flush()
	},
}

var dnsplusRecordsetCreateCmd = &cobra.Command{
	Use:   "create <zone-id>",
	Short: "Create a record set",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneID := args[0]
		name, _ := cmd.Flags().GetString("name")
		recordType, _ := cmd.Flags().GetString("type")
		ttl, _ := cmd.Flags().GetInt("ttl")
		records, _ := cmd.Flags().GetStringSlice("record")

		if name == "" || recordType == "" || len(records) == 0 {
			return fmt.Errorf("--name, --type, and --record are required")
		}

		client := newDNSPlusClient()
		ctx := context.Background()

		recordList := make([]dnsplus.Record, len(records))
		for i, r := range records {
			recordList[i] = dnsplus.Record{RecordContent: r}
		}

		input := &dnsplus.CreateRecordSetInput{
			RecordSetName: name,
			RecordSetType: recordType,
			TTL:           ttl,
			RecordList:    recordList,
		}

		result, err := client.CreateRecordSet(ctx, zoneID, input)
		if err != nil {
			return fmt.Errorf("failed to create record set: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.RecordSet)
		}

		fmt.Printf("Record set created successfully: %s (%s)\n", result.RecordSet.RecordSetName, result.RecordSet.RecordSetID)
		return nil
	},
}

var dnsplusRecordsetDeleteCmd = &cobra.Command{
	Use:   "delete <zone-id> <recordset-id>...",
	Short: "Delete record sets",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		zoneID := args[0]
		recordsetIDs := args[1:]

		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeleteRecordSets(ctx, zoneID, recordsetIDs)
		if err != nil {
			return fmt.Errorf("failed to delete record sets: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Record set(s) deleted successfully\n")
		return nil
	},
}

// ================================
// GSLB Commands
// ================================

var dnsplusGslbCmd = &cobra.Command{
	Use:     "gslb",
	Aliases: []string{"gslbs"},
	Short:   "Manage GSLB (Global Server Load Balancing)",
}

var dnsplusGslbListCmd = &cobra.Command{
	Use:   "list",
	Short: "List GSLBs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListGSLBs(ctx)
		if err != nil {
			return fmt.Errorf("failed to list GSLBs: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.GslbList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tDOMAIN\tSTATUS\tROUTING\tPOOLS")
		for _, gslb := range result.GslbList {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\n",
				gslb.GslbID,
				gslb.GslbName,
				gslb.GslbDomain,
				gslb.GslbStatus,
				gslb.RoutingType,
				gslb.PoolCount,
			)
		}
		return w.Flush()
	},
}

var dnsplusGslbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a GSLB",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		routingType, _ := cmd.Flags().GetString("routing-type")
		ttl, _ := cmd.Flags().GetInt("ttl")
		healthCheckID, _ := cmd.Flags().GetString("health-check-id")

		if name == "" || routingType == "" {
			return fmt.Errorf("--name and --routing-type are required")
		}

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.CreateGSLBInput{
			GslbName:      name,
			Description:   description,
			RoutingType:   routingType,
			TTL:           ttl,
			HealthCheckID: healthCheckID,
		}

		result, err := client.CreateGSLB(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to create GSLB: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Gslb)
		}

		fmt.Printf("GSLB created successfully: %s (%s)\n", result.Gslb.GslbName, result.Gslb.GslbID)
		return nil
	},
}

var dnsplusGslbDeleteCmd = &cobra.Command{
	Use:   "delete <gslb-id>...",
	Short: "Delete GSLBs",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeleteGSLBs(ctx, args)
		if err != nil {
			return fmt.Errorf("failed to delete GSLBs: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("GSLB(s) deleted successfully\n")
		return nil
	},
}

// ================================
// Pool Commands
// ================================

var dnsplusPoolCmd = &cobra.Command{
	Use:     "pool",
	Aliases: []string{"pools"},
	Short:   "Manage GSLB pools",
}

var dnsplusPoolListCmd = &cobra.Command{
	Use:   "list <gslb-id>",
	Short: "List pools in a GSLB",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID := args[0]
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListPools(ctx, gslbID)
		if err != nil {
			return fmt.Errorf("failed to list pools: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.PoolList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tPRIORITY\tWEIGHT\tENDPOINTS")
		for _, pool := range result.PoolList {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%d\n",
				pool.PoolID,
				pool.PoolName,
				pool.PoolStatus,
				pool.Priority,
				pool.Weight,
				pool.EndpointCount,
			)
		}
		return w.Flush()
	},
}

var dnsplusPoolCreateCmd = &cobra.Command{
	Use:   "create <gslb-id>",
	Short: "Create a pool",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID := args[0]
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		priority, _ := cmd.Flags().GetInt("priority")
		weight, _ := cmd.Flags().GetInt("weight")
		region, _ := cmd.Flags().GetString("pool-region")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.CreatePoolInput{
			PoolName:    name,
			Description: description,
			Priority:    priority,
			Weight:      weight,
			Region:      region,
		}

		result, err := client.CreatePool(ctx, gslbID, input)
		if err != nil {
			return fmt.Errorf("failed to create pool: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Pool)
		}

		fmt.Printf("Pool created successfully: %s (%s)\n", result.Pool.PoolName, result.Pool.PoolID)
		return nil
	},
}

var dnsplusPoolDeleteCmd = &cobra.Command{
	Use:   "delete <gslb-id> <pool-id>...",
	Short: "Delete pools",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID := args[0]
		poolIDs := args[1:]

		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeletePools(ctx, gslbID, poolIDs)
		if err != nil {
			return fmt.Errorf("failed to delete pools: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Pool(s) deleted successfully\n")
		return nil
	},
}

// ================================
// Endpoint Commands
// ================================

var dnsplusEndpointCmd = &cobra.Command{
	Use:     "endpoint",
	Aliases: []string{"endpoints"},
	Short:   "Manage pool endpoints",
}

var dnsplusEndpointListCmd = &cobra.Command{
	Use:   "list <gslb-id> <pool-id>",
	Short: "List endpoints in a pool",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID := args[0]
		poolID := args[1]

		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListEndpoints(ctx, gslbID, poolID)
		if err != nil {
			return fmt.Errorf("failed to list endpoints: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.EndpointList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tADDRESS\tSTATUS\tWEIGHT\tHEALTH")
		for _, ep := range result.EndpointList {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
				ep.EndpointID,
				ep.EndpointAddress,
				ep.EndpointStatus,
				ep.Weight,
				ep.HealthStatus,
			)
		}
		return w.Flush()
	},
}

var dnsplusEndpointCreateCmd = &cobra.Command{
	Use:   "create <gslb-id> <pool-id>",
	Short: "Create an endpoint",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID := args[0]
		poolID := args[1]
		address, _ := cmd.Flags().GetString("address")
		weight, _ := cmd.Flags().GetInt("weight")
		description, _ := cmd.Flags().GetString("description")

		if address == "" {
			return fmt.Errorf("--address is required")
		}

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.CreateEndpointInput{
			EndpointAddress: address,
			Weight:          weight,
			Description:     description,
		}

		result, err := client.CreateEndpoint(ctx, gslbID, poolID, input)
		if err != nil {
			return fmt.Errorf("failed to create endpoint: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.Endpoint)
		}

		fmt.Printf("Endpoint created successfully: %s (%s)\n", result.Endpoint.EndpointAddress, result.Endpoint.EndpointID)
		return nil
	},
}

var dnsplusEndpointDeleteCmd = &cobra.Command{
	Use:   "delete <gslb-id> <pool-id> <endpoint-id>...",
	Short: "Delete endpoints",
	Args:  cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		gslbID := args[0]
		poolID := args[1]
		endpointIDs := args[2:]

		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeleteEndpoints(ctx, gslbID, poolID, endpointIDs)
		if err != nil {
			return fmt.Errorf("failed to delete endpoints: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Endpoint(s) deleted successfully\n")
		return nil
	},
}

// ================================
// Health Check Commands
// ================================

var dnsplusHealthCheckCmd = &cobra.Command{
	Use:     "health-check",
	Aliases: []string{"health-checks", "healthcheck"},
	Short:   "Manage health checks",
}

var dnsplusHealthCheckListCmd = &cobra.Command{
	Use:   "list",
	Short: "List health checks",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.ListHealthChecks(ctx)
		if err != nil {
			return fmt.Errorf("failed to list health checks: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.HealthCheckList)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tPROTOCOL\tPORT\tINTERVAL\tTIMEOUT")
		for _, hc := range result.HealthCheckList {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%d\n",
				hc.HealthCheckID,
				hc.HealthCheckName,
				hc.Protocol,
				hc.Port,
				hc.Interval,
				hc.Timeout,
			)
		}
		return w.Flush()
	},
}

var dnsplusHealthCheckCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a health check",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		protocol, _ := cmd.Flags().GetString("protocol")
		port, _ := cmd.Flags().GetInt("port")
		path, _ := cmd.Flags().GetString("path")
		host, _ := cmd.Flags().GetString("host")
		interval, _ := cmd.Flags().GetInt("interval")
		timeout, _ := cmd.Flags().GetInt("timeout")
		retries, _ := cmd.Flags().GetInt("retries")
		expectedCodes, _ := cmd.Flags().GetString("expected-codes")

		if name == "" || protocol == "" || port == 0 {
			return fmt.Errorf("--name, --protocol, and --port are required")
		}

		client := newDNSPlusClient()
		ctx := context.Background()

		input := &dnsplus.CreateHealthCheckInput{
			HealthCheckName: name,
			Description:     description,
			Protocol:        protocol,
			Port:            port,
			Path:            path,
			Host:            host,
			Interval:        interval,
			Timeout:         timeout,
			Retries:         retries,
			ExpectedCodes:   expectedCodes,
		}

		result, err := client.CreateHealthCheck(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to create health check: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result.HealthCheck)
		}

		fmt.Printf("Health check created successfully: %s (%s)\n", result.HealthCheck.HealthCheckName, result.HealthCheck.HealthCheckID)
		return nil
	},
}

var dnsplusHealthCheckDeleteCmd = &cobra.Command{
	Use:   "delete <health-check-id>...",
	Short: "Delete health checks",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := newDNSPlusClient()
		ctx := context.Background()

		result, err := client.DeleteHealthChecks(ctx, args)
		if err != nil {
			return fmt.Errorf("failed to delete health checks: %w", err)
		}

		if output == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		fmt.Printf("Health check(s) deleted successfully\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dnsplusCmd)

	// Zone commands
	dnsplusCmd.AddCommand(dnsplusZoneCmd)
	dnsplusZoneCmd.AddCommand(dnsplusZoneListCmd)
	dnsplusZoneCmd.AddCommand(dnsplusZoneCreateCmd)
	dnsplusZoneCmd.AddCommand(dnsplusZoneUpdateCmd)
	dnsplusZoneCmd.AddCommand(dnsplusZoneDeleteCmd)

	dnsplusZoneCreateCmd.Flags().String("name", "", "Zone name (e.g., example.com)")
	dnsplusZoneCreateCmd.Flags().String("description", "", "Zone description")
	dnsplusZoneCreateCmd.MarkFlagRequired("name")

	dnsplusZoneUpdateCmd.Flags().String("description", "", "Zone description")
	dnsplusZoneUpdateCmd.Flags().String("status", "", "Zone status (USE, STOP)")

	// Record set commands
	dnsplusCmd.AddCommand(dnsplusRecordsetCmd)
	dnsplusRecordsetCmd.AddCommand(dnsplusRecordsetListCmd)
	dnsplusRecordsetCmd.AddCommand(dnsplusRecordsetCreateCmd)
	dnsplusRecordsetCmd.AddCommand(dnsplusRecordsetDeleteCmd)

	dnsplusRecordsetCreateCmd.Flags().String("name", "", "Record set name")
	dnsplusRecordsetCreateCmd.Flags().String("type", "", "Record type (A, AAAA, CNAME, MX, TXT, etc.)")
	dnsplusRecordsetCreateCmd.Flags().Int("ttl", 300, "TTL in seconds")
	dnsplusRecordsetCreateCmd.Flags().StringSlice("record", nil, "Record content (can specify multiple)")

	// GSLB commands
	dnsplusCmd.AddCommand(dnsplusGslbCmd)
	dnsplusGslbCmd.AddCommand(dnsplusGslbListCmd)
	dnsplusGslbCmd.AddCommand(dnsplusGslbCreateCmd)
	dnsplusGslbCmd.AddCommand(dnsplusGslbDeleteCmd)

	dnsplusGslbCreateCmd.Flags().String("name", "", "GSLB name")
	dnsplusGslbCreateCmd.Flags().String("description", "", "GSLB description")
	dnsplusGslbCreateCmd.Flags().String("routing-type", "", "Routing type (FAILOVER, RANDOM, GEOLOCATION)")
	dnsplusGslbCreateCmd.Flags().Int("ttl", 300, "TTL in seconds")
	dnsplusGslbCreateCmd.Flags().String("health-check-id", "", "Health check ID")

	// Pool commands
	dnsplusCmd.AddCommand(dnsplusPoolCmd)
	dnsplusPoolCmd.AddCommand(dnsplusPoolListCmd)
	dnsplusPoolCmd.AddCommand(dnsplusPoolCreateCmd)
	dnsplusPoolCmd.AddCommand(dnsplusPoolDeleteCmd)

	dnsplusPoolCreateCmd.Flags().String("name", "", "Pool name")
	dnsplusPoolCreateCmd.Flags().String("description", "", "Pool description")
	dnsplusPoolCreateCmd.Flags().Int("priority", 1, "Pool priority")
	dnsplusPoolCreateCmd.Flags().Int("weight", 1, "Pool weight")
	dnsplusPoolCreateCmd.Flags().String("pool-region", "", "Pool region")

	// Endpoint commands
	dnsplusCmd.AddCommand(dnsplusEndpointCmd)
	dnsplusEndpointCmd.AddCommand(dnsplusEndpointListCmd)
	dnsplusEndpointCmd.AddCommand(dnsplusEndpointCreateCmd)
	dnsplusEndpointCmd.AddCommand(dnsplusEndpointDeleteCmd)

	dnsplusEndpointCreateCmd.Flags().String("address", "", "Endpoint address (IP or domain)")
	dnsplusEndpointCreateCmd.Flags().Int("weight", 1, "Endpoint weight")
	dnsplusEndpointCreateCmd.Flags().String("description", "", "Endpoint description")

	// Health check commands
	dnsplusCmd.AddCommand(dnsplusHealthCheckCmd)
	dnsplusHealthCheckCmd.AddCommand(dnsplusHealthCheckListCmd)
	dnsplusHealthCheckCmd.AddCommand(dnsplusHealthCheckCreateCmd)
	dnsplusHealthCheckCmd.AddCommand(dnsplusHealthCheckDeleteCmd)

	dnsplusHealthCheckCreateCmd.Flags().String("name", "", "Health check name")
	dnsplusHealthCheckCreateCmd.Flags().String("description", "", "Health check description")
	dnsplusHealthCheckCreateCmd.Flags().String("protocol", "", "Protocol (HTTP, HTTPS, TCP)")
	dnsplusHealthCheckCreateCmd.Flags().Int("port", 0, "Port number")
	dnsplusHealthCheckCreateCmd.Flags().String("path", "/", "HTTP path")
	dnsplusHealthCheckCreateCmd.Flags().String("host", "", "HTTP host header")
	dnsplusHealthCheckCreateCmd.Flags().Int("interval", 30, "Check interval in seconds")
	dnsplusHealthCheckCreateCmd.Flags().Int("timeout", 10, "Timeout in seconds")
	dnsplusHealthCheckCreateCmd.Flags().Int("retries", 3, "Number of retries")
	dnsplusHealthCheckCreateCmd.Flags().String("expected-codes", "200", "Expected HTTP status codes")
}

func newDNSPlusClient() *dnsplus.Client {
	return dnsplus.NewClient(getAppKey(), nil, debug)
}
