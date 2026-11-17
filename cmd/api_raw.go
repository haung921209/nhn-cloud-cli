package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(apiRawCmd)
	apiRawCmd.Flags().StringP("method", "X", "GET", "HTTP method")
	apiRawCmd.Flags().String("service", "network", "Service type in the IaaS service catalog (network, compute, image, volumev2, ...)")
	apiRawCmd.Flags().String("body", "", "JSON request body, or @file to read it from a file")
}

var apiRawCmd = &cobra.Command{
	Use:   "api <path>",
	Short: "Call an IaaS API path directly (for APIs without a dedicated command)",
	Long: `Call an NHN Cloud IaaS API with the configured API credentials and print the JSON response.

The endpoint comes from the IaaS token's service catalog for the current region.
Use this for APIs that have no dedicated command yet.

Examples:
  nhncloud api /v2.0/routingtables
  nhncloud api --service compute /servers/detail   (the compute endpoint already ends with /v2/<tenant-id>)
  nhncloud api -X PUT /v2.0/routingtables/<id>/set_as_default`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		method, _ := cmd.Flags().GetString("method")
		service, _ := cmd.Flags().GetString("service")
		bodyArg, _ := cmd.Flags().GetString("body")

		var body interface{}
		if bodyArg != "" {
			data := []byte(bodyArg)
			if strings.HasPrefix(bodyArg, "@") {
				var err error
				if data, err = os.ReadFile(bodyArg[1:]); err != nil {
					exitWithError("Failed to read body file", err)
				}
			}
			if err := json.Unmarshal(data, &body); err != nil {
				exitWithError("Body is not valid JSON", err)
			}
		}

		path := args[0]
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		var out interface{}
		if err := iaasRequest(context.Background(), service, strings.ToUpper(method), path, body, &out); err != nil {
			exitWithError("API call failed", err)
		}
		if out == nil {
			fmt.Println("OK")
			return
		}
		emit(out, func() { printJSON(out) })
	},
}
