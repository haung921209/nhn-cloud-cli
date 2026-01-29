package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/haung921209/nhn-cloud-cli/internal/sshkeys"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/floatingip"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/securitygroup"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/vpc"
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
Automatically handles:
1. Floating IP: Associates an available one or allocates a new one if missing.
2. Security Group: Ensures a security group allowing SSH (port 22) is attached.
3. SSH Key: Resolves the private key from ~/.ssh/ or managed keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		computeClient := getComputeClient()
		ctx := context.Background()

		instanceID, _ := cmd.Flags().GetString("instance-id")
		username, _ := cmd.Flags().GetString("username")
		identityFile, _ := cmd.Flags().GetString("identity-file")

		// 1. Get Instance Details
		serverOutput, err := computeClient.GetServer(ctx, instanceID)
		if err != nil {
			exitWithError("Failed to get instance details", err)
		}
		server := serverOutput.Server

		// 1.1 Auto-detect Username from Metadata
		if !cmd.Flags().Changed("username") {
			if loginUser, ok := server.Metadata["login_username"]; ok && loginUser != "" {
				username = loginUser
				fmt.Printf("Auto-detected username: %s\n", username)
			}
		}

		// 2. Check & Setup Network (Public IP)
		// We need to find the port associated with this instance to do network operations
		fipClient := floatingip.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		sgClient := securitygroup.NewClient(getRegion(), getIdentityCreds(), nil, debug)
		vpcClient := vpc.NewClient(getRegion(), getIdentityCreds(), nil, debug)

		// Check for ports using our extension method on FIP client
		portsOutput, err := fipClient.ListPorts(ctx, &floatingip.ListPortsOptions{DeviceID: instanceID})
		if err != nil {
			fmt.Printf("Warning: Failed to list ports for instance: %v. Automations might fail.\n", err)
		}

		var targetPort *floatingip.Port
		if portsOutput != nil && len(portsOutput.Ports) > 0 {
			targetPort = &portsOutput.Ports[0]
		}

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

		// Auto-assign Floating IP if missing
		if publicIP == "" {
			if targetPort == nil {
				exitWithError("Instance has no public IP and failed to find instance port to attach one.", nil)
			}
			fmt.Println("Instance has no floating IP. Attempting to assign one...")

			// Check for available floating IP
			fips, err := fipClient.ListFloatingIPs(ctx)
			if err != nil {
				exitWithError("Failed to list floating IPs", err)
			}

			var fipID string
			for _, fip := range fips.FloatingIPs {
				if fip.Status == "DOWN" && fip.PortID == nil {
					fipID = fip.ID
					publicIP = fip.FloatingIPAddress
					fmt.Printf("Found available floating IP: %s\n", publicIP)
					break
				}
			}

			// Allocate new if none available
			if fipID == "" {
				fmt.Println("No available floating IP found. Allocating a new one...")
				// Find external network
				vpcs, err := vpcClient.ListVPCs(ctx)
				if err != nil {
					exitWithError("Failed to list networks", err)
				}
				var extNetID string
				for _, v := range vpcs.VPCs {
					if v.RouterExternal {
						extNetID = v.ID
						break
					}
				}
				if extNetID == "" {
					exitWithError("Failed to find an external network to allocate floating IP.", nil)
				}

				newFip, err := fipClient.CreateFloatingIP(ctx, &floatingip.CreateFloatingIPInput{
					FloatingNetworkID: extNetID,
				})
				if err != nil {
					exitWithError("Failed to create floating IP", err)
				}
				fipID = newFip.FloatingIP.ID
				publicIP = newFip.FloatingIP.FloatingIPAddress
				fmt.Printf("Allocated new floating IP: %s\n", publicIP)
			}

			// Associate
			portID := targetPort.ID
			if _, err := fipClient.UpdateFloatingIP(ctx, fipID, &floatingip.UpdateFloatingIPInput{PortID: &portID}); err != nil {
				exitWithError("Failed to associate floating IP", err)
			}
			fmt.Printf("Associated floating IP %s to instance.\n", publicIP)
		}

		// 3. Check & Setup Security Groups
		if targetPort != nil {
			// Check if SSH is allowed
			sshAllowed := false
			currentSGs := targetPort.SecurityGroups

			for _, sgID := range currentSGs {
				sg, err := sgClient.GetSecurityGroup(ctx, sgID)
				if err != nil {
					continue
				}
				for _, rule := range sg.SecurityGroup.Rules {
					// Check for Allow SSH (TCP 22) from everywhere or at least active
					if rule.Direction == "ingress" && rule.Protocol != nil && *rule.Protocol == "tcp" &&
						rule.PortRangeMin != nil && *rule.PortRangeMin <= 22 && *rule.PortRangeMax >= 22 {
						sshAllowed = true
						break
					}
				}
				if sshAllowed {
					break
				}
			}

			if !sshAllowed {
				fmt.Println("SSH access (port 22) seems to be blocked. configuring security group...")

				// Find or Create "default-ssh"
				sgs, err := sgClient.ListSecurityGroups(ctx)
				var sshSGID string
				if err == nil {
					for _, sg := range sgs.SecurityGroups {
						if sg.Name == "default-ssh" {
							sshSGID = sg.ID
							break
						}
					}
				}

				if sshSGID == "" {
					fmt.Println("Creating 'default-ssh' security group...")
					newSG, err := sgClient.CreateSecurityGroup(ctx, &securitygroup.CreateSecurityGroupInput{
						Name:        "default-ssh",
						Description: "Auto-created by CLI for SSH access",
					})
					if err != nil {
						fmt.Printf("Warning: Failed to create security group: %v\n", err)
					} else {
						sshSGID = newSG.SecurityGroup.ID

						// Detect Public IP for Rule
						userIP := getPublicIP()
						remotePrefix := "0.0.0.0/0"
						if userIP != "" {
							remotePrefix = userIP + "/32"
							fmt.Printf("Authorizing SSH access for your detected IP: %s\n", userIP)
						} else {
							fmt.Println("Warning: Failed to detect your public IP. Allowing all IPs (0.0.0.0/0).")
						}

						// Add Rule
						portVal := 22
						_, err := sgClient.CreateRule(ctx, &securitygroup.CreateRuleInput{
							SecurityGroupID: sshSGID,
							Direction:       "ingress",
							EtherType:       "IPv4",
							Protocol:        "tcp",
							PortRangeMin:    &portVal,
							PortRangeMax:    &portVal,
							RemoteIPPrefix:  remotePrefix,
						})
						if err != nil {
							fmt.Printf("Warning: Failed to create SSH rule: %v\n", err)
						}
					}
				}

				if sshSGID != "" {
					// Attach to port
					newSGList := append(currentSGs, sshSGID)
					_, err := fipClient.UpdatePort(ctx, targetPort.ID, &floatingip.UpdatePortInput{SecurityGroups: &newSGList})
					if err != nil {
						fmt.Printf("Warning: Failed to attach security group to port: %v\n", err)
					} else {
						fmt.Println("Attached 'default-ssh' security group to instance.")
					}
				}
			}
		}

		// 4. Resolve Private Key
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

		// 5. Construct SSH Command
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

func getPublicIP() string {
	client := http.Client{
		Timeout: 2 * time.Second, // Short timeout
	}
	resp, err := client.Get("https://checkip.amazonaws.com")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(body))
}
