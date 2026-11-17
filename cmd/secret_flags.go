package cmd

import (
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

// secretAnnotation marks one flag as secret (catalog commands), unlike
// secretFlagNames which applies to every flag of that name.
const secretAnnotation = "nhncloud_secret"

func isSecretFlag(f *pflag.Flag) bool {
	if a := f.Annotations[secretAnnotation]; len(a) > 0 {
		return a[0] == "true"
	}
	return secretFlagNames[f.Name]
}

// markSecretFlags adds the stdin hint to every secret flag's help text.
func markSecretFlags(c *cobra.Command) {
	mark := func(f *pflag.Flag) {
		if isSecretFlag(f) && !strings.HasSuffix(f.Usage, secretStdinHint) {
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
		if isSecretFlag(f) && f.Changed && f.Value.String() == "-" {
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
	} else if len(f.Annotations[secretAnnotation]) > 0 && stdinUsers(cmd) == 1 {
		// A catalog secret is the only stdin user: take all of it (multi-line
		// PEM keys). Hand-written commands keep reading one line, since some
		// read the rest of stdin themselves.
		b, err := io.ReadAll(in)
		if err != nil {
			return fmt.Errorf("read --%s from stdin: %w", f.Name, err)
		}
		value = strings.TrimRight(string(b), "\r\n")
	} else {
		line, err := readLine(in)
		if err != nil {
			return fmt.Errorf("read --%s from stdin: %w", f.Name, err)
		}
		value = strings.TrimRight(line, "\r")
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

// readLine reads up to and including the first newline one byte at a time,
// so input after the secret stays in stdin for the command (a buffered reader
// would swallow it).
func readLine(in io.Reader) (string, error) {
	var b strings.Builder
	buf := make([]byte, 1)
	for {
		n, err := in.Read(buf)
		if n == 1 {
			if buf[0] == '\n' {
				return b.String(), nil
			}
			b.WriteByte(buf[0])
		}
		if err == io.EOF {
			return b.String(), nil
		}
		if err != nil {
			return "", err
		}
	}
}

// stdinUsers counts the flags set to "-" (stdin), e.g. a secret and
// --data-file.
func stdinUsers(cmd *cobra.Command) int {
	n := 0
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed && f.Value.String() == "-" {
			n++
		}
	})
	return n
}
