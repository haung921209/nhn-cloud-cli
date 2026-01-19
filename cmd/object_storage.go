package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/object"
	"github.com/spf13/cobra"
)

var objectStorageCmd = &cobra.Command{
	Use:     "object-storage",
	Aliases: []string{"os", "swift"},
	Short:   "Manage Object Storage containers and objects",
	Long:    `Manage object storage containers and objects (Swift-compatible).`,
}

func init() {
	rootCmd.AddCommand(objectStorageCmd)
	objectStorageCmd.AddCommand(osAccountInfoCmd)
}

func getObjectStorageClient() *object.Client {
	// Object Storage often uses a separate Tenant ID from Compute
	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getObjectStorageTenantID())
	return object.NewClient(getRegion(), creds, nil, debug)
}

// Helpers retained for shared use
func detectContentType(filePath string) string {
	ext := strings.ToLower(filePath)
	switch {
	case strings.HasSuffix(ext, ".html"), strings.HasSuffix(ext, ".htm"):
		return "text/html"
	case strings.HasSuffix(ext, ".css"):
		return "text/css"
	case strings.HasSuffix(ext, ".js"):
		return "application/javascript"
	case strings.HasSuffix(ext, ".json"):
		return "application/json"
	case strings.HasSuffix(ext, ".xml"):
		return "application/xml"
	case strings.HasSuffix(ext, ".txt"):
		return "text/plain"
	case strings.HasSuffix(ext, ".pdf"):
		return "application/pdf"
	case strings.HasSuffix(ext, ".png"):
		return "image/png"
	case strings.HasSuffix(ext, ".jpg"), strings.HasSuffix(ext, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(ext, ".gif"):
		return "image/gif"
	case strings.HasSuffix(ext, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(ext, ".zip"):
		return "application/zip"
	case strings.HasSuffix(ext, ".tar"):
		return "application/x-tar"
	case strings.HasSuffix(ext, ".gz"):
		return "application/gzip"
	default:
		return "application/octet-stream"
	}
}

func parseOBSTime(s string) time.Time {
	layouts := []string{
		"2006-01-02T15:04:05.000000",
		"2006-01-02T15:04:05",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// Retain account-info as it's useful
var osAccountInfoCmd = &cobra.Command{
	Use:   "account-info",
	Short: "Show storage account information",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		info, err := client.GetAccountInfo(context.Background())
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Container Count: %d\n", info.ContainerCount)
		fmt.Printf("Object Count: %d\n", info.ObjectCount)
		fmt.Printf("Bytes Used: %d\n", info.BytesUsed)
	},
}
