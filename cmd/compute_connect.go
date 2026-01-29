package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/haung921209/nhn-cloud-cli/internal/sshkeys"
)

func init() {
	computeCmd.AddCommand(computeConnectCmd)

	computeConnectCmd.Flags().String("instance-id", "", "ID of the instance to connect to (required)")
	computeConnectCmd.Flags().StringP("username", "l", "centos", "SSH username (default: centos, or auto-detected from metadata)")
	computeConnectCmd.Flags().StringP("identity-file", "i", "", "Identity file (private key) path")
	computeConnectCmd.MarkFlagRequired("instance-id")
}

var computeConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a compute instance via SSH",
	Long: `Connect to a compute instance via SSH.
Automatically detects the instance's public IP and attempts to find the associated SSH private key in ~/.ssh/.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getComputeClient()
		ctx := context.Background()

		instanceID, _ := cmd.Flags().GetString("instance-id")
		username, _ := cmd.Flags().GetString("username")
		identityFile, _ := cmd.Flags().GetString("identity-file")

		// 1. Get Instance Details
		serverOutput, err := client.GetServer(ctx, instanceID)
		if err != nil {
			exitWithError("Failed to get instance details", err)
		}
		server := serverOutput.Server

		// 1.1 Auto-detect Username from Metadata
		// If user didn't specify a flag, and valid metadata exists, use it.
		// Flag has precedence (but here, flag default is "centos", which is tricky.
		// Cobra flags have default values. We should check if the flag was *changed*?
		// Or just check if username == "centos" (default) AND metadata has something else?
		// Better: Set default to "" in flag definition, handle default logic here.
		// But changing flag default might affect help text?
		// Let's check: if username is "centos" (default) and metadata provides "ubuntu", "rocky", etc., switch?
		// Risky if user *wanted* centos.
		// Safe approach: Check `cmd.Flags().Changed("username")`.

		if !cmd.Flags().Changed("username") {
			if loginUser, ok := server.Metadata["login_username"]; ok && loginUser != "" {
				username = loginUser
				fmt.Printf("Auto-detected username: %s\n", username)
			}
		}

		// 2. Extract Public IP
		publicIP := ""
		for _, addrs := range server.Addresses {
			for _, addr := range addrs {
				if addr.Type == "floating" {
					publicIP = addr.Addr
					break
				}
			}
			if publicIP != "" {
				break
			}
		}

		if publicIP == "" {
			fmt.Println("Error: No floating IP found for this instance. Assign a floating IP to connect.")
			os.Exit(1)
		}

		// 3. Resolve Private Key
		keyPath := identityFile
		if keyPath == "" {
			if server.KeyName == "" {
				fmt.Println("Warning: Instance has no Key Pair associated. Trying standard keys...")
			} else {
				// Try to find key in ~/.ssh/
				homeDir, _ := os.UserHomeDir()
				candidates := []string{
					filepath.Join(homeDir, ".ssh", server.KeyName+".pem"),
					filepath.Join(homeDir, ".ssh", server.KeyName),
					filepath.Join(homeDir, ".ssh", "id_rsa"),
				}

				for _, c := range candidates {
					if _, err := os.Stat(c); err == nil {
						keyPath = c
						fmt.Printf("Found Identity File: %s\n", keyPath)
						break
					}
				}

				// If not found in .ssh, try to find in NHN Cloud CLI managed keys
				if keyPath == "" {
					manager := sshkeys.NewManager()
					if keyInfo, err := manager.Get(server.KeyName); err == nil {
						keyPath = keyInfo.Path
						fmt.Printf("Found Identity File (Managed): %s\n", keyPath)
					}
				}

				if keyPath == "" {
					fmt.Printf("Warning: Key Pair '%s' not found locally in ~/.ssh/ or managed keys. You may need to specify -i manually.\n", server.KeyName)
				}
			}
		}

		// 4. Construct SSH Command
		fmt.Printf("Connecting to %s@%s (%s)...\n", username, publicIP, server.Name)

		sshArgs := []string{
			"-o", "StrictHostKeyChecking=no", // Convenience for cloud
			"-o", "UserKnownHostsFile=/dev/null", // Avoid cluttering known_hosts with ephemeral IPs
		}

		if keyPath != "" {
			sshArgs = append(sshArgs, "-i", keyPath)
		}

		sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", username, publicIP))

		sshCmd := exec.Command("ssh", sshArgs...)
		sshCmd.Stdin = os.Stdin
		sshCmd.Stdout = os.Stdout
		sshCmd.Stderr = os.Stderr

		if err := sshCmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				// Propagate exit code
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					os.Exit(status.ExitStatus())
				}
			}
			exitWithError("SSH connection failed", err)
		}
	},
}
