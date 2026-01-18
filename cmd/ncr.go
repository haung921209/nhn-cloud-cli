package cmd

import (
	"fmt"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/container/ncr"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/spf13/cobra"
)

var ncrCmd = &cobra.Command{
	Use:     "ncr",
	Aliases: []string{"registry"},
	Short:   "Manage NHN Container Registry (NCR)",
	Long:    `Manage container registries, images, and tags.`,
}

func init() {
	rootCmd.AddCommand(ncrCmd)
}

func getNCRClient() *ncr.Client {
	creds := credentials.NewStatic(getAccessKey(), getSecretKey())
	return ncr.NewClient(getRegion(), getNCRAppKey(), creds, nil, debug)
}

// formatSize is used by ncr_repositories.go
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
