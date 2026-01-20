package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	ncrCmd.AddCommand(ncrLoginCmd)
}

var ncrLoginCmd = &cobra.Command{
	Use:   "login [registry-name]",
	Short: "Log in to a container registry via Docker",
	Long: `Authenticate the Docker client with your NHN Cloud Container Registry.
This command automatically retrieves your User Access Key and Secret Key from the CLI configuration
and executes 'docker login' for the specified registry.

If no registry name is provided, it will attempt to find a single available registry.
If you know the Registry URI directly (e.g., myreg.kr1.ncr.nhncloud.com), you can provide that as the argument.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getNCRClient()
		ctx := context.Background()

		var registryURI string

		if len(args) > 0 {
			input := args[0]
			// heuristic: if it contains dots, treat as URI
			if strings.Contains(input, ".") {
				registryURI = input
			} else {
				// Treat as registry name, look it up
				// We need to list registries to find the one with this name
				// Since SDK might not support GetRegistryByName directly, we list and filter.
				// Optimization: If GetRegistry(id) works with name? unlikely.
				// Let's list.
				regs, err := client.ListRegistries(ctx)
				if err != nil {
					exitWithError("Failed to list registries", err)
				}
				found := false
				for _, r := range regs.Registries {
					if r.Name == input {
						registryURI = r.URI
						found = true
						break
					}
				}
				if !found {
					exitWithError(fmt.Sprintf("Registry '%s' not found", input), nil)
				}
			}
		} else {
			// No argument, try to auto-detect
			regs, err := client.ListRegistries(ctx)
			if err != nil {
				exitWithError("Failed to list registries", err)
			}
			if len(regs.Registries) == 0 {
				exitWithError("No registries found. Create one first.", nil)
			} else if len(regs.Registries) == 1 {
				registryURI = regs.Registries[0].URI
				fmt.Printf("Found single registry: %s (%s)\n", regs.Registries[0].Name, registryURI)
			} else {
				fmt.Println("Multiple registries found:")
				for _, r := range regs.Registries {
					fmt.Printf(" - %s (%s)\n", r.Name, r.URI)
				}
				exitWithError("Please specify a registry name or URI", nil)
			}
		}

		// Check credentials
		accessKey := getAccessKey()
		secretKey := getSecretKey()

		if accessKey == "" || secretKey == "" {
			exitWithError("User Access Key or Secret Key is missing. Configure them via environment variables or flags.", nil)
		}

		fmt.Printf("Logging in to %s...\n", registryURI)

		// Prepare docker login command
		// docker login <URI> -u <AccessKey> --password-stdin
		dockerCmd := exec.Command("docker", "login", registryURI, "-u", accessKey, "--password-stdin")

		// Pipe secret key to stdin
		dockerCmd.Stdin = strings.NewReader(secretKey)
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr

		if err := dockerCmd.Run(); err != nil {
			exitWithError("Docker login failed", err)
		}

		fmt.Println("Login Succeeded! You can now use 'docker push/pull'.")
	},
}
