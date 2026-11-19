package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// isSecretName matches profile keys and env names that hold credentials.
func isSecretName(name string) bool {
	n := strings.ToLower(name)
	for _, s := range []string{"key", "secret", "password", "token"} {
		if strings.Contains(n, s) {
			return true
		}
	}
	return false
}

// knownSecrets are the credential values this run could send: flags,
// NHN_CLOUD_* env vars and the profile's key, secret and password entries.
func knownSecrets() []string {
	secrets := []string{appKey, password, secretKeyFlag}
	for _, kv := range os.Environ() {
		name, v, _ := strings.Cut(kv, "=")
		if strings.HasPrefix(name, "NHN_CLOUD_") && isSecretName(name) {
			secrets = append(secrets, v)
		}
	}
	cfg := LoadConfig()
	secrets = append(secrets, cfg.AccessKeyID, cfg.SecretAccessKey, cfg.AppKey, cfg.APIPassword,
		cfg.RDSAppKey, cfg.RDSPostgreSQLAppKey, cfg.RDSMariaDBAppKey, cfg.NCRAppKey)
	for k, v := range cfg.Extra {
		if isSecretName(k) {
			secrets = append(secrets, v)
		}
	}
	return secrets
}

// errorText is an error message with every known credential value hidden.
// Servers echo some of them back ("Invalid appKey: <value>").
func errorText(s string) string {
	return redact(s, knownSecrets())
}

// rejectUnknownSubcommands makes a group command (one with subcommands and
// no action) fail on an unknown subcommand instead of printing its help and
// exiting 0, so scripts see the mistake.
func rejectUnknownSubcommands(c *cobra.Command) {
	for _, sub := range c.Commands() {
		rejectUnknownSubcommands(sub)
	}
	if !c.HasSubCommands() || c.Run != nil || c.RunE != nil || c == rootCmd {
		return
	}
	c.Args = cobra.ArbitraryArgs
	c.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
		}
		return cmd.Help()
	}
}
