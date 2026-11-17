package cmd

import (
	"os"
	"os/exec"
	"os/signal"

	"github.com/haung921209/nhn-cloud-cli/internal/cert"
)

// runWithSecret runs a database client with the password passed outside argv
// (see cert.ClientSecret) and removes any temporary secret file before
// returning, also when the client fails. While the client runs, the CLI
// ignores Ctrl-C so the client handles it and the cleanup still happens.
func runWithSecret(c *exec.Cmd, serviceType, password string) error {
	args, env, cleanup, err := cert.ClientSecret(serviceType, password)
	if err != nil {
		return err
	}
	defer cleanup()
	if len(args) > 0 {
		c.Args = append(append([]string{c.Args[0]}, args...), c.Args[1:]...)
	}
	if len(env) > 0 {
		c.Env = append(os.Environ(), env...)
	}
	signal.Ignore(os.Interrupt)
	defer signal.Reset(os.Interrupt)
	return c.Run()
}
