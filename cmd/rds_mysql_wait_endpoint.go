package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// ============================================================================
// wait-db-instance: poll GetInstance until dbInstanceStatus matches.
//
// Polls GetInstance every --interval (default 15s) until dbInstanceStatus
// matches --for-state, or --timeout elapses.
//
//	exit 0 — desired state reached
//	exit 1 — timeout
//	exit 2 — terminal-error state (any state containing "FAIL")
//
// Ref: docs/api-specs/database/rds-mysql-v4.0.md#db-인스턴스-목록-보기
//
// Valid --for-state values per spec (response field dbInstanceStatus):
//
//	AVAILABLE, BEFORE_CREATE, STORAGE_FULL, FAIL_TO_CREATE,
//	FAIL_TO_CONNECT, REPLICATION_STOP, FAILOVER, SHUTDOWN, DELETED
//
// (states are passed through as opaque strings; the spec is authoritative)
// ============================================================================

var waitDBInstanceCmd = &cobra.Command{
	Use:   "wait-db-instance",
	Short: "Wait for a MySQL DB instance to reach a target dbInstanceStatus",
	Long: `Polls GetInstance until dbInstanceStatus matches --for-state, or --timeout fires.

Exit codes:
  0  desired state reached
  1  timeout elapsed before reaching desired state
  2  instance entered a terminal-error state (any value containing "FAIL")

Valid --for-state values (per v4.0 spec, response.dbInstance.dbInstanceStatus):
  AVAILABLE, BEFORE_CREATE, STORAGE_FULL, FAIL_TO_CREATE,
  FAIL_TO_CONNECT, REPLICATION_STOP, FAILOVER, SHUTDOWN, DELETED

Example:
  nhncloud rds-mysql wait-db-instance \
    --db-instance-identifier mydb --for-state AVAILABLE --timeout 30m

Ref: docs/api-specs/database/rds-mysql-v4.0.md#db-인스턴스-목록-보기`,
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, err := getResolvedInstanceID(cmd, newMySQLClient())
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		forState, _ := cmd.Flags().GetString("for-state")
		if forState == "" {
			exitWithError("--for-state is required", nil)
		}
		timeoutStr, _ := cmd.Flags().GetString("timeout")
		intervalStr, _ := cmd.Flags().GetString("interval")

		timeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			exitWithError(fmt.Sprintf("invalid --timeout %q", timeoutStr), err)
		}
		interval, err := time.ParseDuration(intervalStr)
		if err != nil {
			exitWithError(fmt.Sprintf("invalid --interval %q", intervalStr), err)
		}

		client := newMySQLClient()
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		want := strings.ToUpper(strings.TrimSpace(forState))
		deadline := time.Now().Add(timeout)

		for {
			resp, err := client.GetInstance(ctx, dbInstanceID)
			if err != nil {
				// Network/transient errors: keep polling unless deadline passed.
				if time.Now().After(deadline) {
					fmt.Fprintf(os.Stderr, "wait timed out after %s while polling: %v\n", timeout, err)
					os.Exit(1)
				}
			} else {
				got := strings.ToUpper(string(resp.DBInstanceStatus))
				if got == want {
					fmt.Printf("instance %s reached state %s\n", dbInstanceID, got)
					return // exit 0
				}
				// Terminal error: any state containing "FAIL"
				if strings.Contains(got, "FAIL") {
					fmt.Fprintf(os.Stderr, "instance %s entered terminal state %s (wanted %s)\n",
						dbInstanceID, got, want)
					os.Exit(2)
				}
				fmt.Printf("instance %s state=%s (wanted %s) — sleeping %s...\n",
					dbInstanceID, got, want, interval)
			}

			if time.Now().Add(interval).After(deadline) {
				fmt.Fprintf(os.Stderr, "wait timed out after %s (last poll)\n", timeout)
				os.Exit(1)
			}
			select {
			case <-ctx.Done():
				fmt.Fprintf(os.Stderr, "wait context cancelled: %v\n", ctx.Err())
				os.Exit(1)
			case <-time.After(interval):
			}
		}
	},
}

// ============================================================================
// show-db-endpoint: thin wrapper that prints "<host>:<port>\n" for shell sub.
//
// Resolves host using the same fallback chain as `connect`:
//
//	Network.DomainName → FloatingIP → PublicIP → IPAddress → GetNetworkInfo EXTERNAL
//
// Ref: docs/api-specs/database/rds-mysql-v4.0.md#db-인스턴스-상세-보기
// Ref: docs/api-specs/database/rds-mysql-v4.0.md#네트워크-정보-보기
// ============================================================================

var showDBEndpointCmd = &cobra.Command{
	Use:   "show-db-endpoint",
	Short: "Print '<host>:<port>' for a MySQL DB instance (for shell substitution)",
	Long: `Prints "<host>:<port>\n" to stdout for a MySQL DB instance.

Designed for shell substitution in scenario runners:

  ENDPOINT=$(nhncloud rds-mysql show-db-endpoint --db-instance-identifier $ID)
  HOST="${ENDPOINT%:*}"
  PORT="${ENDPOINT##*:}"

Ref: docs/api-specs/database/rds-mysql-v4.0.md#db-인스턴스-상세-보기`,
	Run: func(cmd *cobra.Command, args []string) {
		dbInstanceID, err := getResolvedInstanceID(cmd, newMySQLClient())
		if err != nil {
			exitWithError("failed to resolve instance ID", err)
		}

		client := newMySQLClient()
		inst, err := client.GetInstance(context.Background(), dbInstanceID)
		if err != nil {
			exitWithError("failed to get instance details", err)
		}

		host := ""
		if inst.Network != nil {
			host = inst.Network.DomainName
			if host == "" {
				host = inst.Network.FloatingIP
			}
			if host == "" {
				host = inst.Network.PublicIP
			}
			if host == "" {
				host = inst.Network.IPAddress
			}
		}

		// Fallback to GetNetworkInfo (same chain as `connect`).
		if host == "" {
			netInfo, err := client.GetNetworkInfo(context.Background(), dbInstanceID)
			if err == nil && netInfo != nil {
				for _, ep := range netInfo.EndPoints {
					if ep.EndPointType == "EXTERNAL" {
						if ep.Domain != "" {
							host = ep.Domain
						} else if ep.IPAddress != "" {
							host = ep.IPAddress
						}
						break
					}
					if host == "" {
						if ep.Domain != "" {
							host = ep.Domain
						} else if ep.IPAddress != "" {
							host = ep.IPAddress
						}
					}
				}
			}
		}

		if host == "" {
			exitWithError(fmt.Sprintf(
				"unable to determine endpoint host for %q — instance may not have public access enabled",
				inst.DBInstanceName), nil)
		}

		fmt.Printf("%s:%d\n", host, inst.DBPort)
	},
}

// ============================================================================
// init: register subcommands + flags
// ============================================================================

func init() {
	rdsMySQLCmd.AddCommand(waitDBInstanceCmd)
	rdsMySQLCmd.AddCommand(showDBEndpointCmd)

	// wait-db-instance flags
	waitDBInstanceCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
	waitDBInstanceCmd.Flags().String("for-state", "", "Target dbInstanceStatus (e.g. AVAILABLE) (required)")
	waitDBInstanceCmd.Flags().String("timeout", "30m", "Max time to wait (Go duration, e.g. 30m, 1h)")
	waitDBInstanceCmd.Flags().String("interval", "15s", "Polling interval (Go duration)")

	// show-db-endpoint flags
	showDBEndpointCmd.Flags().String("db-instance-identifier", "", "DB instance identifier (required)")
}
