package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	keymanagerCmd.AddCommand(kmGetClientInfoCmd)
}

var kmGetClientInfoCmd = &cobra.Command{
	Use:   "get-client-info",
	Short: "Get client information",
	Run: func(cmd *cobra.Command, args []string) {
		client := newKeyManagerClient()
		ctx := context.Background()

		result, err := client.GetClientInfo(ctx)
		if err != nil {
			exitWithError("Failed to get client info", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("App Key:     %s\n", result.Body.AppKey)
		fmt.Printf("IP Address:  %s\n", result.Body.IPAddress)
		fmt.Printf("MAC Address: %s\n", result.Body.MACAddress)
	},
}
