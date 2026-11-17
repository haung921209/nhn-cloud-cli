package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

// Secret flags accept "-" to read the value from stdin (or a hidden prompt on
// a terminal), so passwords stay out of ps output and shell history.
var secretFlagNames = map[string]bool{
	"password": true, "db-password": true, "master-user-password": true, "secret": true,
}

const secretStdinHint = " (use - to read it from stdin)"

// markSecretFlags adds the stdin hint to every secret flag's help text.
func markSecretFlags(c *cobra.Command) {
	mark := func(f *pflag.Flag) {
		if secretFlagNames[f.Name] && !strings.HasSuffix(f.Usage, secretStdinHint) {
			f.Usage += secretStdinHint
		}
	}
	c.LocalNonPersistentFlags().VisitAll(mark)
	c.PersistentFlags().VisitAll(mark)
	for _, sub := range c.Commands() {
		markSecretFlags(sub)
	}
}

// readSecretFlags replaces a secret flag set to "-" with one line from in.
// Only one flag can read stdin per command.
func readSecretFlags(cmd *cobra.Command, in io.Reader) error {
	var dash []*pflag.Flag
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if secretFlagNames[f.Name] && f.Changed && f.Value.String() == "-" {
			dash = append(dash, f)
		}
	})
	if len(dash) == 0 {
		return nil
	}
	if len(dash) > 1 {
		return fmt.Errorf("only one secret flag can read stdin (--%s and --%s both use -)", dash[0].Name, dash[1].Name)
	}
	f := dash[0]
	var value string
	if file, ok := in.(*os.File); ok && term.IsTerminal(int(file.Fd())) {
		fmt.Fprintf(os.Stderr, "%s: ", f.Name)
		b, err := term.ReadPassword(int(file.Fd()))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return fmt.Errorf("read --%s: %w", f.Name, err)
		}
		value = string(b)
	} else {
		line, err := bufio.NewReader(in).ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("read --%s from stdin: %w", f.Name, err)
		}
		value = strings.TrimRight(line, "\r\n")
	}
	if value == "" {
		return fmt.Errorf("--%s: stdin was empty", f.Name)
	}
	return f.Value.Set(value)
}

func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return readSecretFlags(cmd, os.Stdin)
	}
}
