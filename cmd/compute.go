package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/compute"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/spf13/cobra"
)

var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Manage Compute instances (VMs)",
	Long:  `Manage Compute instances including create, delete, start, stop, and more.`,
}

func init() {
	rootCmd.AddCommand(computeCmd)

	computeCmd.AddCommand(computeListCmd)
	computeCmd.AddCommand(computeGetCmd)
	computeCmd.AddCommand(computeCreateCmd)
	computeCmd.AddCommand(computeDeleteCmd)
	computeCmd.AddCommand(computeStartCmd)
	computeCmd.AddCommand(computeStopCmd)
	computeCmd.AddCommand(computeRebootCmd)
	computeCmd.AddCommand(computeFlavorsCmd)
	computeCmd.AddCommand(computeImagesCmd)
	computeCmd.AddCommand(computeKeypairsCmd)
	computeCmd.AddCommand(computeKeypairCreateCmd)
	computeCmd.AddCommand(computeKeypairDeleteCmd)

	computeCreateCmd.Flags().String("name", "", "Instance name (required)")
	computeCreateCmd.Flags().String("image", "", "Image ID (required)")
	computeCreateCmd.Flags().String("flavor", "", "Flavor ID (required)")
	computeCreateCmd.Flags().String("network", "", "Network/Subnet ID (required)")
	computeCreateCmd.Flags().String("key-name", "", "SSH keypair name")
	computeCreateCmd.Flags().String("security-group", "", "Security group name")
	computeCreateCmd.Flags().String("availability-zone", "", "Availability zone")
	computeCreateCmd.Flags().Int("boot-volume-size", 20, "Boot volume size in GB")
	computeCreateCmd.MarkFlagRequired("name")
	computeCreateCmd.MarkFlagRequired("image")
	computeCreateCmd.MarkFlagRequired("flavor")
	computeCreateCmd.MarkFlagRequired("network")

	computeRebootCmd.Flags().Bool("hard", false, "Hard reboot (default: soft)")

	computeKeypairCreateCmd.Flags().String("name", "", "Keypair name (required)")
	computeKeypairCreateCmd.Flags().String("public-key", "", "Public key content (optional)")
	computeKeypairCreateCmd.MarkFlagRequired("name")
}

func getComputeClient() *compute.Client {
	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	return compute.NewClient(getRegion(), creds, nil, debug)
}

var computeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all compute instances",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		result, err := client.ListServers(ctx)
		if err != nil {
			exitWithError("Failed to list instances", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tKEY\tAZ\tCREATED")
		for _, s := range result.Servers {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				s.ID, s.Name, s.Status, s.KeyName, s.AvailabilityZone, s.Created)
		}
		w.Flush()
	},
}

var computeGetCmd = &cobra.Command{
	Use:   "get [instance-id]",
	Short: "Get compute instance details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		result, err := client.GetServer(ctx, args[0])
		if err != nil {
			exitWithError("Failed to get instance", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		s := result.Server
		fmt.Printf("ID:                %s\n", s.ID)
		fmt.Printf("Name:              %s\n", s.Name)
		fmt.Printf("Status:            %s\n", s.Status)
		fmt.Printf("Key Name:          %s\n", s.KeyName)
		fmt.Printf("Availability Zone: %s\n", s.AvailabilityZone)
		fmt.Printf("Created:           %s\n", s.Created)
		fmt.Printf("Updated:           %s\n", s.Updated)

		if len(s.Addresses) > 0 {
			fmt.Println("\nAddresses:")
			for network, addrs := range s.Addresses {
				for _, addr := range addrs {
					fmt.Printf("  %s: %s (%s)\n", network, addr.Addr, addr.Type)
				}
			}
		}

		if len(s.SecurityGroups) > 0 {
			fmt.Println("\nSecurity Groups:")
			for _, sg := range s.SecurityGroups {
				fmt.Printf("  - %s\n", sg.Name)
			}
		}
	},
}

var computeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new compute instance",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		image, _ := cmd.Flags().GetString("image")
		flavor, _ := cmd.Flags().GetString("flavor")
		network, _ := cmd.Flags().GetString("network")
		keyName, _ := cmd.Flags().GetString("key-name")
		sgName, _ := cmd.Flags().GetString("security-group")
		az, _ := cmd.Flags().GetString("availability-zone")
		volumeSize, _ := cmd.Flags().GetInt("boot-volume-size")

		input := &compute.CreateServerInput{
			Name:             name,
			ImageRef:         image,
			FlavorRef:        flavor,
			KeyName:          keyName,
			AvailabilityZone: az,
			Networks: []compute.ServerNetwork{
				{UUID: network},
			},
		}

		if sgName != "" {
			input.SecurityGroups = []compute.SecurityGroup{{Name: sgName}}
		}

		if volumeSize > 0 {
			input.BlockDeviceMapping = []compute.BlockDeviceMapping{
				{
					BootIndex:           0,
					UUID:                image,
					SourceType:          "image",
					DestinationType:     "volume",
					VolumeSize:          volumeSize,
					DeleteOnTermination: true,
				},
			}
		}

		result, err := client.CreateServer(ctx, input)
		if err != nil {
			exitWithError("Failed to create instance", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Instance created successfully!\n")
		fmt.Printf("ID:   %s\n", result.Server.ID)
		fmt.Printf("Name: %s\n", result.Server.Name)
	},
}

var computeDeleteCmd = &cobra.Command{
	Use:   "delete [instance-id]",
	Short: "Delete a compute instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		if err := client.DeleteServer(ctx, args[0]); err != nil {
			exitWithError("Failed to delete instance", err)
		}

		fmt.Printf("Instance %s deleted successfully\n", args[0])
	},
}

var computeStartCmd = &cobra.Command{
	Use:   "start [instance-id]",
	Short: "Start a compute instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		if err := client.StartServer(ctx, args[0]); err != nil {
			exitWithError("Failed to start instance", err)
		}

		fmt.Printf("Instance %s started\n", args[0])
	},
}

var computeStopCmd = &cobra.Command{
	Use:   "stop [instance-id]",
	Short: "Stop a compute instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		if err := client.StopServer(ctx, args[0]); err != nil {
			exitWithError("Failed to stop instance", err)
		}

		fmt.Printf("Instance %s stopped\n", args[0])
	},
}

var computeRebootCmd = &cobra.Command{
	Use:   "reboot [instance-id]",
	Short: "Reboot a compute instance",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		hard, _ := cmd.Flags().GetBool("hard")
		if err := client.RebootServer(ctx, args[0], hard); err != nil {
			exitWithError("Failed to reboot instance", err)
		}

		rebootType := "soft"
		if hard {
			rebootType = "hard"
		}
		fmt.Printf("Instance %s rebooted (%s)\n", args[0], rebootType)
	},
}

var computeFlavorsCmd = &cobra.Command{
	Use:   "flavors",
	Short: "List available compute flavors",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		result, err := client.ListFlavors(ctx)
		if err != nil {
			exitWithError("Failed to list flavors", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tVCPUs\tRAM (MB)\tDISK (GB)")
		for _, f := range result.Flavors {
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\n",
				f.ID, f.Name, f.VCPUs, f.RAM, f.Disk)
		}
		w.Flush()
	},
}

var computeImagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List available compute images",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		result, err := client.ListImages(ctx)
		if err != nil {
			exitWithError("Failed to list images", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tMIN_DISK\tMIN_RAM")
		for _, img := range result.Images {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\n",
				img.ID, img.Name, img.Status, img.MinDisk, img.MinRAM)
		}
		w.Flush()
	},
}

var computeKeypairsCmd = &cobra.Command{
	Use:   "keypairs",
	Short: "List SSH keypairs",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		result, err := client.ListKeyPairs(ctx)
		if err != nil {
			exitWithError("Failed to list keypairs", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tFINGERPRINT")
		for _, kp := range result.KeyPairs {
			fmt.Fprintf(w, "%s\t%s\n", kp.KeyPair.Name, kp.KeyPair.Fingerprint)
		}
		w.Flush()
	},
}

var computeKeypairCreateCmd = &cobra.Command{
	Use:   "keypair-create",
	Short: "Create a new SSH keypair",
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		name, _ := cmd.Flags().GetString("name")
		publicKey, _ := cmd.Flags().GetString("public-key")

		input := &compute.CreateKeyPairInput{
			Name:      name,
			PublicKey: publicKey,
		}

		result, err := client.CreateKeyPair(ctx, input)
		if err != nil {
			exitWithError("Failed to create keypair", err)
		}

		if output == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Keypair created: %s\n", result.KeyPair.Name)
		fmt.Printf("Fingerprint: %s\n", result.KeyPair.Fingerprint)
		if result.KeyPair.PrivateKey != "" {
			fmt.Printf("\nPrivate Key (save this - it won't be shown again):\n%s\n", result.KeyPair.PrivateKey)
		}
	},
}

var computeKeypairDeleteCmd = &cobra.Command{
	Use:   "keypair-delete [name]",
	Short: "Delete an SSH keypair",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		if err := client.DeleteKeyPair(ctx, args[0]); err != nil {
			exitWithError("Failed to delete keypair", err)
		}

		fmt.Printf("Keypair %s deleted\n", args[0])
	},
}
