package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/haung921209/nhn-cloud-cli/internal/sshkeys"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/port"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/network/securitygroup"
)

func init() {
	computeCmd.AddCommand(computeConnectCmd)

	computeConnectCmd.Flags().String("instance-id", "", requiredHint("ID of the instance to connect to", "find it with: nhncloud compute describe-instances"))
	computeConnectCmd.Flags().StringP("ssh-user", "l", "centos", "SSH login user (default: the image's login_username metadata, else centos)")
	computeConnectCmd.Flags().StringP("identity-file", "i", "", "Identity file (private key) path")
	// Kept so existing scripts keep working: --username used to be the SSH
	// user of this command, and --allow-cidr used to open TCP 22.
	computeConnectCmd.Flags().String("username", "", "SSH login user (same as --ssh-user)")
	computeConnectCmd.Flags().MarkDeprecated("username", "use --ssh-user (-l); --username here is the SSH user, not the API user")
	computeConnectCmd.Flags().String("allow-cidr", "", "no longer used")
	computeConnectCmd.Flags().MarkDeprecated("allow-cidr", "connect no longer changes security groups; add a TCP 22 rule with network authorize-security-group-ingress")
	computeConnectCmd.MarkFlagRequired("instance-id")
}

var computeConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a compute instance via SSH",
	Long: `Connect to a compute instance via SSH.
The instance needs a floating IP; connect does not assign one or change
security groups (it warns when no rule seems to allow TCP 22).
The private key is resolved from ~/.ssh/ or managed keys unless -i is given.
The SSH user is --ssh-user (-l); --username is accepted as the same for
existing scripts. Set the API user with NHN_CLOUD_USERNAME or the profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		computeClient := getComputeClient()
		ctx := context.Background()

		instanceID, _ := cmd.Flags().GetString("instance-id")
		sshUser, _ := cmd.Flags().GetString("ssh-user")
		userGiven := cmd.Flags().Changed("ssh-user")
		if u, _ := cmd.Flags().GetString("username"); cmd.Flags().Changed("username") && !userGiven {
			sshUser, userGiven = u, true
		}
		identityFile, _ := cmd.Flags().GetString("identity-file")

		// 1. Get Instance Details
		serverOutput, err := computeClient.GetServer(ctx, instanceID)
		if err != nil {
			exitWithError("Failed to get instance details", err)
		}
		server := serverOutput.Server

		// 1.1 Auto-detect Username from Metadata
		if !userGiven {
			if loginUser, ok := server.Metadata["login_username"]; ok && loginUser != "" {
				sshUser = loginUser
				fmt.Printf("Auto-detected username: %s\n", sshUser)
			}
		}

		// 2. Public IP. connect never takes, allocates or associates a
		// floating IP by itself: an unattached one may belong to someone else.
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
			exitWithError(fmt.Sprintf("Instance %s has no floating IP. Associate one first:\n"+
				"  nhncloud network allocate-floating-ip --network-id <external-network-id>   (or pick your own from describe-floating-ips)\n"+
				"  nhncloud network associate-floating-ip --floating-ip-id <id> --instance-id %s", instanceID, instanceID), nil)
		}

		// 3. SSH rule. Only checked; connect does not change security groups.
		if !sshIngressAllowed(ctx, instanceID) {
			progressf("Warning: no security group on the instance port seems to allow TCP 22. Add a rule, for example:\n" +
				"  nhncloud network authorize-security-group-ingress --group-id <group-id> --protocol tcp --port-min 22 --port-max 22 --remote-ip <your-ip>/32\n")
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
		fmt.Printf("Connecting to %s@%s (%s)...\n", sshUser, publicIP, server.Name)

		sshArgs, err := sshConnectArgs(sshUser, publicIP, keyPath)
		if err != nil {
			exitWithError("Invalid SSH user", err)
		}

		sshCmd := exec.Command("ssh", sshArgs...)
		sshCmd.Stdin = os.Stdin
		sshCmd.Stdout = os.Stdout
		sshCmd.Stderr = os.Stderr

		if err := sshCmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				// Propagate exit code
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					if status.ExitStatus() == 255 {
						progressf("If ssh reported a changed host key because this floating IP was used by another instance, check it and remove the old entry: ssh-keygen -R %s\n", publicIP)
					}
					os.Exit(status.ExitStatus())
				}
			}
			exitWithError("SSH connection failed", err)
		}
	},
}

// sshIngressAllowed reports whether a security group on the instance's
// first port has an ingress TCP rule covering port 22. Lookup errors count as
// allowed, so they never block the connection.
func sshIngressAllowed(ctx context.Context, instanceID string) bool {
	ports, err := port.NewClient(getRegion(), getIdentityCreds(), nil, debug).ListPortsByDevice(ctx, instanceID)
	if err != nil || len(ports.Ports) == 0 {
		return true
	}
	sgClient := securitygroup.NewClient(getRegion(), getIdentityCreds(), nil, debug)
	for _, sgID := range ports.Ports[0].SecurityGroups {
		sg, err := sgClient.GetSecurityGroup(ctx, sgID)
		if err != nil {
			return true
		}
		for _, rule := range sg.SecurityGroup.Rules {
			if rule.Direction == "ingress" && rule.Protocol != nil && *rule.Protocol == "tcp" &&
				rule.PortRangeMin != nil && rule.PortRangeMax != nil && *rule.PortRangeMin <= 22 && *rule.PortRangeMax >= 22 {
				return true
			}
		}
	}
	return false
}

// sshConnectArgs builds the ssh arguments. The host key is checked against
// ~/.ssh/known_hosts and a new host is added on first use (accept-new), so a
// changed key is reported instead of accepted. The user is passed with -l and
// the host after "--", so neither can be read as an ssh option.
func sshConnectArgs(user, host, keyPath string) ([]string, error) {
	if user == "" || strings.HasPrefix(user, "-") || strings.ContainsAny(user, " \t") {
		return nil, fmt.Errorf("%q is not a valid SSH user name", user)
	}
	args := []string{"-o", "StrictHostKeyChecking=accept-new"}
	if keyPath != "" {
		args = append(args, "-i", keyPath)
	}
	return append(args, "-l", user, "--", host), nil
}
