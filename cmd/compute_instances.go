package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/compute"
	"github.com/spf13/cobra"
)

func init() {
	computeCmd.AddCommand(computeDescribeInstancesCmd)
	computeCmd.AddCommand(computeCreateInstanceCmd)
	computeCmd.AddCommand(computeDeleteInstanceCmd)
	computeCmd.AddCommand(computeStartInstancesCmd)
	computeCmd.AddCommand(computeStopInstancesCmd)
	computeCmd.AddCommand(computeRebootInstancesCmd)

	computeDescribeInstancesCmd.Flags().String("instance-id", "", "ID of the instance to describe")

	computeCreateInstanceCmd.Flags().String("name", "", "Instance name (required)")
	computeCreateInstanceCmd.Flags().String("image-id", "", "Image ID (required)")
	computeCreateInstanceCmd.Flags().String("flavor-id", "", "Flavor ID (required)")
	computeCreateInstanceCmd.Flags().String("subnet-id", "", "Network/Subnet ID (required)")
	computeCreateInstanceCmd.Flags().String("key-name", "", "SSH keypair name")
	computeCreateInstanceCmd.Flags().String("security-group-ids", "", "Security group name (comma separated)") // AWS uses IDs, NHN uses names typically? Old cmd used name.
	computeCreateInstanceCmd.Flags().String("availability-zone", "", "Availability zone")
	computeCreateInstanceCmd.Flags().Int("block-device-mapping-v2-boot-volume-size", 20, "Boot volume size in GB")
	computeCreateInstanceCmd.Flags().Bool("wait", false, "Wait for instance to be available") // Future
	computeCreateInstanceCmd.MarkFlagRequired("name")
	computeCreateInstanceCmd.MarkFlagRequired("image-id")
	computeCreateInstanceCmd.MarkFlagRequired("flavor-id")
	computeCreateInstanceCmd.MarkFlagRequired("subnet-id")

	computeDeleteInstanceCmd.Flags().String("instance-id", "", "Instance ID (required)")
	computeDeleteInstanceCmd.MarkFlagRequired("instance-id")

	computeStartInstancesCmd.Flags().String("instance-id", "", "Instance ID (required)")
	computeStartInstancesCmd.MarkFlagRequired("instance-id")

	computeStopInstancesCmd.Flags().String("instance-id", "", "Instance ID (required)")
	computeStopInstancesCmd.MarkFlagRequired("instance-id")

	computeRebootInstancesCmd.Flags().String("instance-id", "", "Instance ID (required)")
	computeRebootInstancesCmd.Flags().Bool("hard", false, "Hard reboot")
	computeRebootInstancesCmd.MarkFlagRequired("instance-id")
}

var computeDescribeInstancesCmd = &cobra.Command{
	Use:   "describe-instances",
	Short: "Describe compute instances",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()
		instanceID, _ := cmd.Flags().GetString("instance-id")

		if instanceID != "" {
			// Get Single
			result, err := client.GetServer(ctx, instanceID)
			if err != nil {
				exitWithError("Failed to get instance", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			s := result.Server
			fmt.Printf("ID:                %s\n", s.ID)
			fmt.Printf("Name:              %s\n", s.Name)
			fmt.Printf("Status:            %s\n", s.Status)
			fmt.Printf("Key Name:          %s\n", s.KeyName)
			fmt.Printf("Availability Zone: %s\n", s.AvailabilityZone)
			fmt.Printf("Created:           %s\n", s.Created)
			if len(s.Addresses) > 0 {
				fmt.Println("\nAddresses:")
				for network, addrs := range s.Addresses {
					for _, addr := range addrs {
						fmt.Printf("  %s: %s (%s)\n", network, addr.Addr, addr.Type)
					}
				}
			}
		} else {
			// List All
			result, err := client.ListServers(ctx)
			if err != nil {
				exitWithError("Failed to list instances", err)
			}
			if output == "json" {
				printJSON(result)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tNAME\tSTATUS\tKEY\tAZ\tCREATED")
			for _, s := range result.Servers {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
					s.ID, s.Name, s.Status, s.KeyName, s.AvailabilityZone, s.Created)
			}
			w.Flush()
		}
	},
}

var computeCreateInstanceCmd = &cobra.Command{
	Use:   "create-instance",
	Short: "Create a new compute instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		imageID, _ := cmd.Flags().GetString("image-id")
		flavorID, _ := cmd.Flags().GetString("flavor-id")
		subnetID, _ := cmd.Flags().GetString("subnet-id")
		keyName, _ := cmd.Flags().GetString("key-name")
		az, _ := cmd.Flags().GetString("availability-zone")
		volumeSize, _ := cmd.Flags().GetInt("block-device-mapping-v2-boot-volume-size")
		// sgIDs, _ := cmd.Flags().GetString("security-group-ids") // TODO: Implement SG support properly (old cmd used name)

		input := &compute.CreateServerInput{
			Name: name,
			// ImageRef set conditionally below
			FlavorRef:        flavorID,
			KeyName:          keyName,
			AvailabilityZone: az,
			Networks: []compute.ServerNetwork{
				{UUID: subnetID},
			},
		}

		if volumeSize > 0 {
			input.BlockDeviceMapping = []compute.BlockDeviceMapping{
				{
					BootIndex:           0,
					UUID:                imageID,
					SourceType:          "image",
					DestinationType:     "volume",
					VolumeSize:          volumeSize,
					DeleteOnTermination: true,
				},
			}
		} else {
			input.ImageRef = imageID
		}

		result, err := client.CreateServer(ctx, input)
		if err != nil {
			exitWithError("Failed to create instance", err)
		}

		if output == "json" {
			printJSON(result)
			return
		}

		fmt.Printf("Instance created successfully: %s (%s)\n", result.Server.ID, result.Server.Name)
	},
}

var computeDeleteInstanceCmd = &cobra.Command{
	Use:   "delete-instance",
	Short: "Delete a compute instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()
		instanceID, _ := cmd.Flags().GetString("instance-id")

		if err := client.DeleteServer(ctx, instanceID); err != nil {
			exitWithError("Failed to delete instance", err)
		}
		fmt.Printf("Instance %s deleted successfully\n", instanceID)
	},
}

var computeStartInstancesCmd = &cobra.Command{
	Use:   "start-instances",
	Short: "Start compute instances",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()
		instanceID, _ := cmd.Flags().GetString("instance-id")

		if err := client.StartServer(ctx, instanceID); err != nil {
			exitWithError("Failed to start instance", err)
		}
		fmt.Printf("Instance %s started\n", instanceID)
	},
}

var computeStopInstancesCmd = &cobra.Command{
	Use:   "stop-instances",
	Short: "Stop compute instances",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()
		instanceID, _ := cmd.Flags().GetString("instance-id")

		if err := client.StopServer(ctx, instanceID); err != nil {
			exitWithError("Failed to stop instance", err)
		}
		fmt.Printf("Instance %s stopped\n", instanceID)
	},
}

var computeRebootInstancesCmd = &cobra.Command{
	Use:   "reboot-instances",
	Short: "Reboot compute instances",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()
		instanceID, _ := cmd.Flags().GetString("instance-id")
		hard, _ := cmd.Flags().GetBool("hard")

		if err := client.RebootServer(ctx, instanceID, hard); err != nil {
			exitWithError("Failed to reboot instance", err)
		}
		fmt.Printf("Instance %s rebooted\n", instanceID)
	},
}
