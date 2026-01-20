package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/object"
	"github.com/spf13/cobra"
)

var objectStorageCmd = &cobra.Command{
	Use:     "object-storage",
	Aliases: []string{"obs"},
	Short:   "Manage Object Storage",
}

var obsLsCmd = &cobra.Command{
	Use:   "ls [obs://container]",
	Short: "List containers or objects",
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		if len(args) == 0 {
			// List Containers
			output, err := client.ListContainers(ctx, nil)
			if err != nil {
				exitWithError("Failed to list containers", err)
			}
			for _, c := range output.Containers {
				fmt.Printf("%s\t%d bytes\t%d objects\n", c.Name, c.Bytes, c.Count)
			}
			return
		}

		// List Objects in Container
		path, err := parseOBSPath(args[0])
		if err != nil {
			exitWithError("Invalid path", err)
		}
		if !path.IsRemote {
			exitWithError("Argument must be obs:// path", nil)
		}

		input := &object.ListObjectsInput{
			Prefix: path.Object,
		}
		output, err := client.ListObjects(ctx, path.Container, input)
		if err != nil {
			exitWithError("Failed to list objects", err)
		}
		for _, o := range output.Objects {
			fmt.Printf("%s\t%d\t%s\n", o.Name, o.Bytes, o.LastModified)
		}
	},
}

var obsCpCmd = &cobra.Command{
	Use:   "cp <source> <destination>",
	Short: "Copy files to/from Object Storage",
	Long: `Copy files between local filesystem and Object Storage, or between Object Storage containers.
Uses 'obs://container/object' syntax for remote paths.
Large files (>5GB) are automatically uploaded as Static Large Objects (SLO).`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getObjectStorageClient()
		ctx := context.Background()

		srcPath, err := parseOBSPath(args[0])
		if err != nil {
			exitWithError("Invalid source path", err)
		}

		destPath, err := parseOBSPath(args[1])
		if err != nil {
			exitWithError("Invalid destination path", err)
		}

		segSize, _ := cmd.Flags().GetInt64("segment-size")

		if !srcPath.IsRemote && destPath.IsRemote {
			// Local -> OBS (Upload)
			if err := uploadToOBS(ctx, client, srcPath, destPath, segSize); err != nil {
				exitWithError("Upload failed", err)
			}
		} else if srcPath.IsRemote && !destPath.IsRemote {
			// OBS -> Local (Download)
			if err := downloadFromOBS(ctx, client, srcPath, destPath); err != nil {
				exitWithError("Download failed", err)
			}
		} else if srcPath.IsRemote && destPath.IsRemote {
			// OBS -> OBS (Copy)
			if err := copyInOBS(ctx, client, srcPath, destPath); err != nil {
				exitWithError("Copy failed", err)
			}
		} else {
			// Local -> Local (Not supported/Out of scope but can fallback to cp)
			exitWithError("Local to Local copy is not supported by this tool", fmt.Errorf("use standard cp command"))
		}
	},
}

func init() {
	rootCmd.AddCommand(objectStorageCmd)
	objectStorageCmd.AddCommand(obsCpCmd)
	objectStorageCmd.AddCommand(obsLsCmd)

	// Future flags: --recursive
	obsCpCmd.Flags().Int64("segment-size", 1024*1024*1024, "Segment size in bytes for multipart upload (default 1GB)")
}

func getObjectStorageClient() *object.Client {
	cfg := LoadConfig()
	creds := getIdentityCredentials()
	return object.NewClient(cfg.Region, creds, nil, false) // debug=false for now
}

func getIdentityCredentials() credentials.IdentityCredentials {
	cfg := LoadConfig()
	tenantID := cfg.TenantID
	if cfg.OBSTenantID != "" {
		tenantID = cfg.OBSTenantID
	}
	return credentials.NewStaticIdentity(cfg.Username, cfg.APIPassword, tenantID)
}

func uploadToOBS(ctx context.Context, client *object.Client, src *OBSPath, dest *OBSPath, segmentSize int64) error {
	f, err := os.Open(src.RawPath)
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	objectName := dest.Object
	if objectName == "" || strings.HasSuffix(objectName, "/") {
		objectName = filepath.Join(objectName, filepath.Base(src.RawPath))
	}

	// Threshold: 5GB
	const multipartThreshold = 5 * 1024 * 1024 * 1024 // 5GB

	if info.Size() > multipartThreshold {
		return uploadMultipartSLO(ctx, client, f, info.Size(), dest.Container, objectName, segmentSize)
	}

	// Simple Upload
	fmt.Printf("Uploading %s to obs://%s/%s (Size: %d bytes)...\n", src.RawPath, dest.Container, objectName, info.Size())
	input := &object.PutObjectInput{
		Container:   dest.Container,
		ObjectName:  objectName,
		Body:        f,
		ContentType: "application/octet-stream", // TODO: Detect mime type
	}

	_, err = client.PutObject(ctx, input)
	if err == nil {
		fmt.Printf("Upload complete.\n")
	}
	return err
}

func uploadMultipartSLO(ctx context.Context, client *object.Client, f *os.File, fileSize int64, container, objectName string, segmentSize int64) error {
	fmt.Printf("Large file detected (%d bytes). Using SLO Multipart Upload (Segment Size: %d bytes)...\n", fileSize, segmentSize)

	segmentContainer := container + "_segments"
	// Ensure segment container exists
	err := client.CreateContainer(ctx, &object.CreateContainerInput{Name: segmentContainer})
	if err != nil {
		// Ignore if already exists (409/202) -> SDK ensure logic might be needed or just try
		// CreateContainer returns error on non-201/202 headers usually.
		fmt.Printf("Note: Segment container creation attempt: %v\n", err)
	}

	totalSegments := (fileSize + segmentSize - 1) / segmentSize
	var segments []object.SLOSegment

	for i := int64(0); i < totalSegments; i++ {
		offset := i * segmentSize
		remaining := fileSize - offset
		if remaining > segmentSize {
			remaining = segmentSize
		}

		// Seek to offset (redundant if sequential, but safe)
		_, err := f.Seek(offset, 0)
		if err != nil {
			return fmt.Errorf("seek segment %d: %w", i, err)
		}

		// Limit reader
		partReader := io.LimitReader(f, remaining)

		fmt.Printf("Uploading segment %d/%d (%d bytes) to %s/%s...\n", i+1, totalSegments, remaining, segmentContainer, objectName)

		input := &object.UploadSegmentInput{
			Container:    segmentContainer,
			ObjectName:   objectName,
			SegmentIndex: int(i + 1),
			Body:         partReader,
			ContentType:  "application/octet-stream",
		}

		out, err := client.UploadSegment(ctx, input)
		if err != nil {
			return fmt.Errorf("upload segment %d failed: %w", i+1, err)
		}

		// SDK UploadSegment constructs path as `container/objectName/001`.
		// SLO Segment path should match that.
		// `/%s/%s/%03d` -> `container/objectName/001`
		// The `Path` in SLOSegment must be `/{segment-container}/{object-name}/{index}`.
		segmentPath := fmt.Sprintf("/%s/%s/%03d", segmentContainer, objectName, i+1)

		segments = append(segments, object.SLOSegment{
			Path:      segmentPath,
			ETag:      out.ETag,
			SizeBytes: remaining,
		})
	}

	fmt.Printf("All segments uploaded. Creating SLO manifest...\n")

	manifestInput := &object.CreateSLOManifestInput{
		Container:   container,
		ObjectName:  objectName,
		Segments:    segments,
		ContentType: "application/octet-stream", // or original mime type
	}

	if err := client.CreateSLOManifest(ctx, manifestInput); err != nil {
		return fmt.Errorf("create manifest failed: %w", err)
	}

	fmt.Printf("SLO Upload complete.\n")
	return nil
}

func downloadFromOBS(ctx context.Context, client *object.Client, src *OBSPath, dest *OBSPath) error {
	srcObject := src.Object
	destPath := dest.RawPath

	if destPath == "" || destPath == "." {
		destPath = filepath.Base(srcObject)
	}

	// Open file for writing
	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create local file %s: %w", destPath, err)
	}
	defer f.Close()

	fmt.Printf("Downloading obs://%s/%s to %s...\n", src.Container, srcObject, destPath)

	// GetObject
	out, err := client.GetObject(ctx, src.Container, srcObject)
	if err != nil {
		os.Remove(destPath) // cleanup on error
		return err
	}
	defer out.Body.Close()

	// Stream to file
	written, err := io.Copy(f, out.Body)
	if err != nil {
		return fmt.Errorf("write stream: %w", err)
	}

	fmt.Printf("Download complete (%d bytes).\n", written)
	return nil
}

func copyInOBS(ctx context.Context, client *object.Client, src *OBSPath, dest *OBSPath) error {
	srcObj := src.Object
	destObj := dest.Object

	if destObj == "" || strings.HasSuffix(destObj, "/") {
		destObj = filepath.Join(destObj, filepath.Base(srcObj))
	}

	fmt.Printf("Copying obs://%s/%s to obs://%s/%s...\n", src.Container, srcObj, dest.Container, destObj)

	input := &object.CopyObjectInput{
		SourceContainer:       src.Container,
		SourceObjectName:      srcObj,
		DestinationContainer:  dest.Container,
		DestinationObjectName: destObj,
	}

	if err := client.CopyObject(ctx, input); err != nil {
		return err
	}

	fmt.Printf("Copy complete.\n")
	return nil
}
