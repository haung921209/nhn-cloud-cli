package main

import (
	"os"

	"github.com/haung921209/nhn-cloud-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
