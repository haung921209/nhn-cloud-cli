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

		recursive, _ := cmd.Flags().GetBool("recursive")
		prefix := path.Object

		input := &object.ListObjectsInput{
			Prefix: prefix,
		}

		if !recursive {
			input.Delimiter = "/"
		}

		output, err := client.ListObjects(ctx, path.Container, input)
		if err != nil {
			exitWithError("Failed to list objects", err)
		}

		// Print Common Prefixes (Virtual Directories)
		for _, p := range output.CommonPrefixes {
			fmt.Printf("                           PRE %s\n", p)
		}

		// Print Objects
		for _, o := range output.Objects {
			// If pseudo-directory itself is listed, skip it
			if o.Name == prefix && strings.HasSuffix(prefix, "/") {
				continue
			}
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
		recursive, _ := cmd.Flags().GetBool("recursive")

		if !srcPath.IsRemote && destPath.IsRemote {
			// Local -> OBS (Upload)
			if err := uploadToOBS(ctx, client, srcPath, destPath, segSize, recursive); err != nil {
				exitWithError("Upload failed", err)
			}
		} else if srcPath.IsRemote && !destPath.IsRemote {
			// OBS -> Local (Download)
			if err := downloadFromOBS(ctx, client, srcPath, destPath, recursive); err != nil {
				exitWithError("Download failed", err)
			}
		} else if srcPath.IsRemote && destPath.IsRemote {
			// OBS -> OBS (Copy)
			if err := copyInOBS(ctx, client, srcPath, destPath, recursive); err != nil {
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

	// Flags
	obsCpCmd.Flags().Int64("segment-size", 1024*1024*1024, "Segment size in bytes for multipart upload (default 1GB)")
	obsCpCmd.Flags().BoolP("recursive", "r", false, "Command is performed on all files or objects under the specified directory or prefix")

	obsLsCmd.Flags().BoolP("recursive", "r", false, "Command is performed on all files or objects under the specified directory or prefix")
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

func uploadToOBS(ctx context.Context, client *object.Client, src *OBSPath, dest *OBSPath, segmentSize int64, recursive bool) error {
	fi, err := os.Stat(src.RawPath)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		if !recursive {
			return fmt.Errorf("source is a directory, use --recursive to upload")
		}
		return uploadDirectory(ctx, client, src.RawPath, dest.Container, dest.Object, segmentSize)
	}

	return uploadFile(ctx, client, src.RawPath, dest.Container, dest.Object, fi.Size(), segmentSize)
}

func uploadFile(ctx context.Context, client *object.Client, srcPath, container, objectName string, size int64, segmentSize int64) error {
	// Adjust object name if destination implies directory
	if objectName == "" || strings.HasSuffix(objectName, "/") {
		objectName = filepath.Join(objectName, filepath.Base(srcPath))
	}

	// Threshold: 5GB
	const multipartThreshold = 5 * 1024 * 1024 * 1024 // 5GB

	if size > multipartThreshold {
		f, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer f.Close()
		return uploadMultipartSLO(ctx, client, f, size, container, objectName, segmentSize)
	}

	// Simple Upload
	fmt.Printf("Uploading %s to obs://%s/%s (Size: %d bytes)...\n", srcPath, container, objectName, size)
	f, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()

	input := &object.PutObjectInput{
		Container:   container,
		ObjectName:  objectName,
		Body:        f,
		ContentType: "application/octet-stream", // TODO: Detect mime type
	}

	_, err = client.PutObject(ctx, input)
	if err == nil {
		fmt.Printf("Upload complete: %s\n", srcPath)
	}
	return err
}

func uploadDirectory(ctx context.Context, client *object.Client, localDir, container, prefix string, segmentSize int64) error {
	var files []string
	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("Uploading directory %s (%d files) to obs://%s/%s\n", localDir, len(files), container, prefix)

	// Simple sequential upload for now, can be parallelized later
	for _, path := range files {
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		// Construct object key
		objName := filepath.Join(prefix, relPath)

		fi, err := os.Stat(path)
		if err != nil {
			return err
		}

		if err := uploadFile(ctx, client, path, container, objName, fi.Size(), segmentSize); err != nil {
			return err
		}
	}
	return nil
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

func downloadFromOBS(ctx context.Context, client *object.Client, src *OBSPath, dest *OBSPath, recursive bool) error {
	destPath := dest.RawPath
	if destPath == "" {
		destPath = "."
	}

	if recursive {
		return downloadDirectory(ctx, client, src.Container, src.Object, destPath)
	}

	// Single file download
	if destPath == "." || strings.HasSuffix(destPath, "/") {
		destPath = filepath.Join(destPath, filepath.Base(src.Object))
	}
	return downloadFile(ctx, client, src.Container, src.Object, destPath)
}

func downloadFile(ctx context.Context, client *object.Client, container, objectName, localPath string) error {
	// Ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}

	// Open file for writing
	f, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("create local file %s: %w", localPath, err)
	}
	defer f.Close()

	fmt.Printf("Downloading obs://%s/%s to %s...\n", container, objectName, localPath)

	// GetObject
	out, err := client.GetObject(ctx, container, objectName)
	if err != nil {
		os.Remove(localPath) // cleanup on error
		return err
	}
	defer out.Body.Close()

	// Stream to file
	written, err := io.Copy(f, out.Body)
	if err != nil {
		return fmt.Errorf("write stream: %w", err)
	}

	fmt.Printf("Download complete: %s (%d bytes)\n", localPath, written)
	return nil
}

func downloadDirectory(ctx context.Context, client *object.Client, container, prefix, localDir string) error {
	// List all objects recursively
	input := &object.ListObjectsInput{
		Prefix: prefix,
	}
	output, err := client.ListObjects(ctx, container, input)
	if err != nil {
		return err
	}

	fmt.Printf("Downloading directory obs://%s/%s (%d objects) to %s\n", container, prefix, len(output.Objects), localDir)

	for _, o := range output.Objects {
		// Calculate relative path
		// e.g. prefix="data/", obj="data/conf/file.txt" -> "conf/file.txt"
		relPath := strings.TrimPrefix(o.Name, prefix)
		// Removing leading slash if any
		relPath = strings.TrimPrefix(relPath, "/")

		if relPath == "" {
			continue // Skip the directory itself marker if exists
		}

		localPath := filepath.Join(localDir, relPath)
		if err := downloadFile(ctx, client, container, o.Name, localPath); err != nil {
			return err
		}
	}
	return nil
}

func copyInOBS(ctx context.Context, client *object.Client, src *OBSPath, dest *OBSPath, recursive bool) error {
	if recursive {
		return copyDirectory(ctx, client, src.Container, src.Object, dest.Container, dest.Object)
	}

	srcObj := src.Object
	destObj := dest.Object

	if destObj == "" || strings.HasSuffix(destObj, "/") {
		destObj = filepath.Join(destObj, filepath.Base(srcObj))
	}

	return copyFile(ctx, client, src.Container, srcObj, dest.Container, destObj)
}

func copyFile(ctx context.Context, client *object.Client, srcContainer, srcObj, destContainer, destObj string) error {
	fmt.Printf("Copying obs://%s/%s to obs://%s/%s...\n", srcContainer, srcObj, destContainer, destObj)

	input := &object.CopyObjectInput{
		SourceContainer:       srcContainer,
		SourceObjectName:      srcObj,
		DestinationContainer:  destContainer,
		DestinationObjectName: destObj,
	}

	if err := client.CopyObject(ctx, input); err != nil {
		return err
	}

	fmt.Printf("Copy complete: %s\n", destObj)
	return nil
}

func copyDirectory(ctx context.Context, client *object.Client, srcContainer, srcPrefix, destContainer, destPrefix string) error {
	// List source objects
	input := &object.ListObjectsInput{
		Prefix: srcPrefix,
	}
	output, err := client.ListObjects(ctx, srcContainer, input)
	if err != nil {
		return err
	}

	fmt.Printf("Copying directory obs://%s/%s (%d objects) to obs://%s/%s\n", srcContainer, srcPrefix, len(output.Objects), destContainer, destPrefix)

	for _, o := range output.Objects {
		relPath := strings.TrimPrefix(o.Name, srcPrefix)
		relPath = strings.TrimPrefix(relPath, "/")

		if relPath == "" {
			continue
		}

		newObjName := filepath.Join(destPrefix, relPath)
		if err := copyFile(ctx, client, srcContainer, o.Name, destContainer, newObjName); err != nil {
			return err
		}
	}
	return nil
}
