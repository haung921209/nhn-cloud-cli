package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/database/postgresql"
	"github.com/spf13/cobra"
)

// ============================================================================
// HBA Rules Commands (PostgreSQL-specific)
// ============================================================================

var describeHBARulesCmd = &cobra.Command{
	Use:   "describe-hba-rules",
	Short: "Describe PostgreSQL HBA rules (pg_hba.conf)",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ListHBARules(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to list HBA rules", err)
		}

		if output == "json" {
			postgresqlPrintJSON(result)
		} else {
			postgresqlPrintHBARuleList(result)
		}
	},
}

var createHBARuleCmd = &cobra.Command{
	Use:   "create-hba-rule",
	Short: "Create a PostgreSQL HBA rule",
	Long: `Create a pg_hba.conf access control rule.

Example:
  nhncloud rds-postgresql create-hba-rule \
    --db-instance-identifier my-pg \
    --address 10.0.0.0/16 \
    --auth-method SCRAM_SHA_256 \
    --database-apply-type ENTIRE \
    --db-user-apply-type ENTIRE`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		address, _ := cmd.Flags().GetString("address")
		authMethod, _ := cmd.Flags().GetString("auth-method")
		dbApplyType, _ := cmd.Flags().GetString("database-apply-type")
		userApplyType, _ := cmd.Flags().GetString("db-user-apply-type")
		connType, _ := cmd.Flags().GetString("connection-type")

		if address == "" {
			exitWithError("--address is required (e.g., 0.0.0.0/0)", nil)
		}
		if authMethod == "" {
			authMethod = "SCRAM_SHA_256"
		}
		if dbApplyType == "" {
			dbApplyType = "ENTIRE"
		}
		if userApplyType == "" {
			userApplyType = "ENTIRE"
		}

		req := &postgresql.CreateHBARuleRequest{
			Address:           address,
			AuthMethod:        authMethod,
			DatabaseApplyType: dbApplyType,
			DBUserApplyType:   userApplyType,
		}
		if connType != "" {
			req.ConnectionType = connType
		}

		result, err := client.CreateHBARule(context.Background(), instanceID, req)
		if err != nil {
			exitWithError("failed to create HBA rule", err)
		}

		fmt.Printf("HBA rule created.\n")
		fmt.Printf("Rule ID: %s\n", result.HBARuleID)
	},
}

var deleteHBARuleCmd = &cobra.Command{
	Use:   "delete-hba-rule",
	Short: "Delete a PostgreSQL HBA rule",
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		ruleID, _ := cmd.Flags().GetString("hba-rule-id")
		if ruleID == "" {
			exitWithError("--hba-rule-id is required", nil)
		}

		_, err = client.DeleteHBARule(context.Background(), instanceID, ruleID)
		if err != nil {
			exitWithError("failed to delete HBA rule", err)
		}

		fmt.Printf("HBA rule deleted successfully.\n")
	},
}

var applyHBARulesCmd = &cobra.Command{
	Use:   "apply-hba-rules",
	Short: "Apply HBA rules to PostgreSQL instance",
	Long:  `Apply pending HBA rule changes to the running PostgreSQL instance.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := newPostgreSQLClient()
		instanceID, err := getResolvedPostgreSQLInstanceID(cmd, client)
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		result, err := client.ApplyHBARules(context.Background(), instanceID)
		if err != nil {
			exitWithError("failed to apply HBA rules", err)
		}

		fmt.Printf("HBA rules application initiated.\n")
		fmt.Printf("Job ID: %s\n", result.JobID)
	},
}

// ============================================================================
// Print Functions
// ============================================================================

func postgresqlPrintHBARuleList(result *postgresql.ListHBARulesResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "RULE_ID\tORDER\tADDRESS\tAUTH_METHOD\tAPPLICABLE")
	for _, rule := range result.HBARules {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%v\n",
			rule.HBARuleID,
			rule.Order,
			rule.Address,
			rule.AuthMethod,
			rule.Applicable,
		)
	}
	w.Flush()
}

func init() {
	// HBA Rules commands
	rdsPostgreSQLCmd.AddCommand(describeHBARulesCmd)
	rdsPostgreSQLCmd.AddCommand(createHBARuleCmd)
	rdsPostgreSQLCmd.AddCommand(deleteHBARuleCmd)
	rdsPostgreSQLCmd.AddCommand(applyHBARulesCmd)

	describeHBARulesCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")

	createHBARuleCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	createHBARuleCmd.Flags().String("address", "", "CIDR address (required, e.g., 10.0.0.0/16)")
	createHBARuleCmd.Flags().String("auth-method", "SCRAM_SHA_256", "Authentication method (SCRAM_SHA_256, MD5, TRUST)")
	createHBARuleCmd.Flags().String("database-apply-type", "ENTIRE", "Database apply type (ENTIRE, SELECTED)")
	createHBARuleCmd.Flags().String("db-user-apply-type", "ENTIRE", "User apply type (ENTIRE, USER_CUSTOM)")
	createHBARuleCmd.Flags().String("connection-type", "", "Connection type (HOST, HOSTSSL, HOSTNOSSL)")

	deleteHBARuleCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	deleteHBARuleCmd.Flags().String("hba-rule-id", "", "HBA rule ID (required)")

	applyHBARulesCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
}
