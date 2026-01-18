package cmd

import (
	"context"
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/nas"
	"github.com/spf13/cobra"
)

func init() {
	nasCmd.AddCommand(nasCreateInterfaceCmd)
	nasCmd.AddCommand(nasDeleteInterfaceCmd)

	nasCreateInterfaceCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasCreateInterfaceCmd.Flags().String("subnet-id", "", "Subnet ID (required)")
	nasCreateInterfaceCmd.MarkFlagRequired("volume-id")
	nasCreateInterfaceCmd.MarkFlagRequired("subnet-id")

	nasDeleteInterfaceCmd.Flags().String("volume-id", "", "Volume ID (required)")
	nasDeleteInterfaceCmd.Flags().String("interface-id", "", "Interface ID (required)")
	nasDeleteInterfaceCmd.MarkFlagRequired("volume-id")
	nasDeleteInterfaceCmd.MarkFlagRequired("interface-id")
}

var nasCreateInterfaceCmd = &cobra.Command{
	Use:   "create-interface",
	Short: "Create a new interface for a volume",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		ctx := context.Background()
		volID, _ := cmd.Flags().GetString("volume-id")
		subnetID, _ := cmd.Flags().GetString("subnet-id")

		input := &nas.CreateInterfaceInput{
			SubnetID: subnetID,
		}

		result, err := client.CreateInterface(ctx, volID, input)
		if err != nil {
			exitWithError("Failed to create interface", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Interface created: %s\n", result.Interface.ID)
		fmt.Printf("Path: %s\n", result.Interface.Path)
	},
}

var nasDeleteInterfaceCmd = &cobra.Command{
	Use:   "delete-interface",
	Short: "Delete an interface from a volume",
	Run: func(cmd *cobra.Command, args []string) {
		client := newNASClient()
		ctx := context.Background()
		volID, _ := cmd.Flags().GetString("volume-id")
		ifaceID, _ := cmd.Flags().GetString("interface-id")

		if err := client.DeleteInterface(ctx, volID, ifaceID); err != nil {
			exitWithError("Failed to delete interface", err)
		}

		fmt.Printf("Interface %s deleted\n", ifaceID)
	},
}
