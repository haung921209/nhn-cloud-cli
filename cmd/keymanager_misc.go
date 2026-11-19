package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	keymanagerCmd.AddCommand(kmGetClientInfoCmd)
	addClientMACFlag(kmGetClientInfoCmd)
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

		fmt.Printf("Client IP:          %s\n", result.Body.ClientIP)
		fmt.Printf("Client MAC Header:  %s\n", result.Body.ClientMACHeader)
		fmt.Printf("Client Certificate: %v\n", result.Body.ClientSentCertificate)
	},
}
