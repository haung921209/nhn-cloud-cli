package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/storage/object"
	"github.com/spf13/cobra"
)

const obsPrefix = "obs://"

var storageCpCmd = &cobra.Command{
	Use:   "cp <source> <destination>",
	Short: "Copy files to/from Object Storage",
	Long: `Copy files between local filesystem and Object Storage.

Supported patterns:
  Local to Object Storage:   nhncloud object-storage cp ./file.txt obs://container/path/file.txt
  Object Storage to Local:   nhncloud object-storage cp obs://container/path/file.txt ./file.txt
  Between containers:        nhncloud object-storage cp obs://src/file.txt obs://dst/file.txt

Use obs:// prefix for Object Storage paths.`,
	Args: cobra.ExactArgs(2),
	Run:  runStorageCp,
}

var storageMvCmd = &cobra.Command{
	Use:   "mv <source> <destination>",
	Short: "Move files to/from Object Storage",
	Long: `Move files between local filesystem and Object Storage (copy + delete source).

Supported patterns:
  Local to Object Storage:   nhncloud object-storage mv ./file.txt obs://container/path/file.txt
  Object Storage to Local:   nhncloud object-storage mv obs://container/path/file.txt ./file.txt
  Between containers:        nhncloud object-storage mv obs://src/file.txt obs://dst/file.txt`,
	Args: cobra.ExactArgs(2),
	Run:  runStorageMv,
}

var storageSyncCmd = &cobra.Command{
	Use:   "sync <source> <destination>",
	Short: "Sync directories with Object Storage",
	Long: `Sync local directory with Object Storage container.

Supported patterns:
  Local to Object Storage:   nhncloud object-storage sync ./dir obs://container/prefix
  Object Storage to Local:   nhncloud object-storage sync obs://container/prefix ./dir`,
	Args: cobra.ExactArgs(2),
	Run:  runStorageSync,
}

func init() {
	objectStorageCmd.AddCommand(storageCpCmd)
	objectStorageCmd.AddCommand(storageMvCmd)
	objectStorageCmd.AddCommand(storageSyncCmd)

	storageCpCmd.Flags().Bool("recursive", false, "Copy directories recursively")
	storageCpCmd.Flags().BoolP("quiet", "q", false, "Suppress progress output")

	storageMvCmd.Flags().Bool("recursive", false, "Move directories recursively")
	storageMvCmd.Flags().BoolP("quiet", "q", false, "Suppress progress output")

	storageSyncCmd.Flags().Bool("delete", false, "Delete files in destination not in source")
	storageSyncCmd.Flags().BoolP("quiet", "q", false, "Suppress progress output")
	storageSyncCmd.Flags().String("exclude", "", "Exclude pattern (glob)")
}

func isObsPath(path string) bool {
	return strings.HasPrefix(path, obsPrefix)
}

func parseObsPath(path string) (container, objectPath string) {
	trimmed := strings.TrimPrefix(path, obsPrefix)
	parts := strings.SplitN(trimmed, "/", 2)
	container = parts[0]
	if len(parts) > 1 {
		objectPath = parts[1]
	}
	return
}

func runStorageCp(cmd *cobra.Command, args []string) {
	src, dst := args[0], args[1]
	recursive, _ := cmd.Flags().GetBool("recursive")
	quiet, _ := cmd.Flags().GetBool("quiet")

	srcIsObs := isObsPath(src)
	dstIsObs := isObsPath(dst)

	ctx := context.Background()

	switch {
	case !srcIsObs && dstIsObs:
		uploadFile(ctx, src, dst, recursive, quiet)
	case srcIsObs && !dstIsObs:
		downloadFile(ctx, src, dst, recursive, quiet)
	case srcIsObs && dstIsObs:
		copyBetweenObs(ctx, src, dst, recursive, quiet)
	default:
		exitWithError("At least one path must be an Object Storage path (obs://...)", nil)
	}
}

func runStorageMv(cmd *cobra.Command, args []string) {
	src, dst := args[0], args[1]
	recursive, _ := cmd.Flags().GetBool("recursive")
	quiet, _ := cmd.Flags().GetBool("quiet")

	srcIsObs := isObsPath(src)
	dstIsObs := isObsPath(dst)

	ctx := context.Background()

	switch {
	case !srcIsObs && dstIsObs:
		uploadFile(ctx, src, dst, recursive, quiet)
		if err := os.Remove(src); err != nil {
			exitWithError("Failed to remove source file after move", err)
		}
		if !quiet {
			fmt.Printf("Removed local file: %s\n", src)
		}
	case srcIsObs && !dstIsObs:
		downloadFile(ctx, src, dst, recursive, quiet)
		container, objPath := parseObsPath(src)
		client := getObjectStorageClient()
		if err := client.DeleteObject(ctx, container, objPath); err != nil {
			exitWithError("Failed to delete source object after move", err)
		}
		if !quiet {
			fmt.Printf("Deleted object: %s\n", src)
		}
	case srcIsObs && dstIsObs:
		copyBetweenObs(ctx, src, dst, recursive, quiet)
		container, objPath := parseObsPath(src)
		client := getObjectStorageClient()
		if err := client.DeleteObject(ctx, container, objPath); err != nil {
			exitWithError("Failed to delete source object after move", err)
		}
		if !quiet {
			fmt.Printf("Deleted object: %s\n", src)
		}
	default:
		exitWithError("At least one path must be an Object Storage path (obs://...)", nil)
	}
}

func uploadFile(ctx context.Context, localPath, obsPath string, recursive, quiet bool) {
	container, objPath := parseObsPath(obsPath)
	client := getObjectStorageClient()

	info, err := os.Stat(localPath)
	if err != nil {
		exitWithError("Failed to stat local file", err)
	}

	if info.IsDir() {
		if !recursive {
			exitWithError("Source is a directory, use --recursive to copy directories", nil)
		}
		uploadDirectory(ctx, client, localPath, container, objPath, quiet)
		return
	}

	file, err := os.Open(localPath)
	if err != nil {
		exitWithError("Failed to open local file", err)
	}
	defer file.Close()

	if objPath == "" || strings.HasSuffix(objPath, "/") {
		objPath = objPath + filepath.Base(localPath)
	}

	if !quiet {
		fmt.Printf("Uploading %s → obs://%s/%s ", localPath, container, objPath)
	}

	start := time.Now()
	_, err = client.PutObject(ctx, &object.PutObjectInput{
		Container:  container,
		ObjectName: objPath,
		Body:       file,
	})
	if err != nil {
		if !quiet {
			fmt.Println("FAILED")
		}
		exitWithError("Failed to upload file", err)
	}

	if !quiet {
		fmt.Printf("(%s, %s)\n", formatSize(info.Size()), time.Since(start).Round(time.Millisecond))
	}
}

func uploadDirectory(ctx context.Context, client *object.Client, localDir, container, prefix string, quiet bool) {
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
		exitWithError("Failed to walk directory", err)
	}

	if !quiet {
		fmt.Printf("Uploading %d files from %s to obs://%s/%s\n", len(files), localDir, container, prefix)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for _, filePath := range files {
		wg.Add(1)
		sem <- struct{}{}

		go func(fp string) {
			defer wg.Done()
			defer func() { <-sem }()

			relPath, _ := filepath.Rel(localDir, fp)
			objPath := prefix
			if objPath != "" && !strings.HasSuffix(objPath, "/") {
				objPath += "/"
			}
			objPath += relPath

			file, err := os.Open(fp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open %s: %v\n", fp, err)
				return
			}
			defer file.Close()

			_, err = client.PutObject(ctx, &object.PutObjectInput{
				Container:  container,
				ObjectName: objPath,
				Body:       file,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to upload %s: %v\n", fp, err)
				return
			}

			if !quiet {
				fmt.Printf("  Uploaded: %s\n", relPath)
			}
		}(filePath)
	}

	wg.Wait()
	if !quiet {
		fmt.Println("Upload complete.")
	}
}

func downloadFile(ctx context.Context, obsPath, localPath string, recursive, quiet bool) {
	container, objPath := parseObsPath(obsPath)
	client := getObjectStorageClient()

	if recursive {
		downloadDirectory(ctx, client, container, objPath, localPath, quiet)
		return
	}

	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	rawClient := createRawObjectClient(creds)

	if !quiet {
		fmt.Printf("Downloading obs://%s/%s → %s ", container, objPath, localPath)
	}

	start := time.Now()
	body, size, err := rawClient.GetObjectRaw(ctx, container, objPath)
	if err != nil {
		if !quiet {
			fmt.Println("FAILED")
		}
		exitWithError("Failed to download file", err)
	}
	defer body.Close()

	info, err := os.Stat(localPath)
	if err == nil && info.IsDir() {
		localPath = filepath.Join(localPath, filepath.Base(objPath))
	}

	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		exitWithError("Failed to create directory", err)
	}

	outFile, err := os.Create(localPath)
	if err != nil {
		exitWithError("Failed to create local file", err)
	}
	defer outFile.Close()

	written, err := io.Copy(outFile, body)
	if err != nil {
		exitWithError("Failed to write file", err)
	}

	if !quiet {
		fmt.Printf("(%s, %s)\n", formatSize(written), time.Since(start).Round(time.Millisecond))
	}
	_ = size
}

func downloadDirectory(ctx context.Context, client *object.Client, container, prefix, localDir string, quiet bool) {
	input := &object.ListObjectsInput{Prefix: prefix}
	result, err := client.ListObjects(ctx, container, input)
	if err != nil {
		exitWithError("Failed to list objects", err)
	}

	if !quiet {
		fmt.Printf("Downloading %d objects from obs://%s/%s to %s\n", len(result.Objects), container, prefix, localDir)
	}

	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	rawClient := createRawObjectClient(creds)

	for _, obj := range result.Objects {
		relPath := strings.TrimPrefix(obj.Name, prefix)
		relPath = strings.TrimPrefix(relPath, "/")
		localPath := filepath.Join(localDir, relPath)

		dir := filepath.Dir(localPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create directory %s: %v\n", dir, err)
			continue
		}

		body, _, err := rawClient.GetObjectRaw(ctx, container, obj.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to download %s: %v\n", obj.Name, err)
			continue
		}

		outFile, err := os.Create(localPath)
		if err != nil {
			body.Close()
			fmt.Fprintf(os.Stderr, "Failed to create %s: %v\n", localPath, err)
			continue
		}

		_, err = io.Copy(outFile, body)
		body.Close()
		outFile.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", localPath, err)
			continue
		}

		if !quiet {
			fmt.Printf("  Downloaded: %s\n", relPath)
		}
	}

	if !quiet {
		fmt.Println("Download complete.")
	}
}

func copyBetweenObs(ctx context.Context, srcPath, dstPath string, recursive, quiet bool) {
	srcContainer, srcObj := parseObsPath(srcPath)
	dstContainer, dstObj := parseObsPath(dstPath)
	client := getObjectStorageClient()

	if recursive {
		input := &object.ListObjectsInput{Prefix: srcObj}
		result, err := client.ListObjects(ctx, srcContainer, input)
		if err != nil {
			exitWithError("Failed to list source objects", err)
		}

		if !quiet {
			fmt.Printf("Copying %d objects from obs://%s/%s to obs://%s/%s\n",
				len(result.Objects), srcContainer, srcObj, dstContainer, dstObj)
		}

		for _, obj := range result.Objects {
			relPath := strings.TrimPrefix(obj.Name, srcObj)
			newObjPath := dstObj + relPath

			copyInput := &object.CopyObjectInput{
				SourceContainer:       srcContainer,
				SourceObjectName:      obj.Name,
				DestinationContainer:  dstContainer,
				DestinationObjectName: newObjPath,
			}

			if err := client.CopyObject(ctx, copyInput); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to copy %s: %v\n", obj.Name, err)
				continue
			}

			if !quiet {
				fmt.Printf("  Copied: %s\n", obj.Name)
			}
		}

		if !quiet {
			fmt.Println("Copy complete.")
		}
		return
	}

	if dstObj == "" || strings.HasSuffix(dstObj, "/") {
		dstObj = dstObj + filepath.Base(srcObj)
	}

	if !quiet {
		fmt.Printf("Copying obs://%s/%s → obs://%s/%s\n", srcContainer, srcObj, dstContainer, dstObj)
	}

	copyInput := &object.CopyObjectInput{
		SourceContainer:       srcContainer,
		SourceObjectName:      srcObj,
		DestinationContainer:  dstContainer,
		DestinationObjectName: dstObj,
	}

	if err := client.CopyObject(ctx, copyInput); err != nil {
		exitWithError("Failed to copy object", err)
	}

	if !quiet {
		fmt.Println("Copy complete.")
	}
}

func runStorageSync(cmd *cobra.Command, args []string) {
	src, dst := args[0], args[1]
	deleteExtra, _ := cmd.Flags().GetBool("delete")
	quiet, _ := cmd.Flags().GetBool("quiet")
	exclude, _ := cmd.Flags().GetString("exclude")

	srcIsObs := isObsPath(src)
	dstIsObs := isObsPath(dst)

	ctx := context.Background()

	switch {
	case !srcIsObs && dstIsObs:
		syncLocalToObs(ctx, src, dst, deleteExtra, quiet, exclude)
	case srcIsObs && !dstIsObs:
		syncObsToLocal(ctx, src, dst, deleteExtra, quiet, exclude)
	default:
		exitWithError("Sync requires one local path and one Object Storage path", nil)
	}
}

func syncLocalToObs(ctx context.Context, localDir, obsPath string, deleteExtra, quiet bool, exclude string) {
	container, prefix := parseObsPath(obsPath)
	client := getObjectStorageClient()

	localFiles := make(map[string]os.FileInfo)
	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if exclude != "" {
			if matched, _ := filepath.Match(exclude, filepath.Base(path)); matched {
				return nil
			}
		}
		relPath, _ := filepath.Rel(localDir, path)
		localFiles[relPath] = info
		return nil
	})
	if err != nil {
		exitWithError("Failed to walk local directory", err)
	}

	input := &object.ListObjectsInput{Prefix: prefix}
	result, err := client.ListObjects(ctx, container, input)
	if err != nil {
		exitWithError("Failed to list objects", err)
	}

	remoteFiles := make(map[string]object.Object)
	for _, obj := range result.Objects {
		relPath := strings.TrimPrefix(obj.Name, prefix)
		relPath = strings.TrimPrefix(relPath, "/")
		remoteFiles[relPath] = obj
	}

	var toUpload []string
	for relPath, info := range localFiles {
		remoteObj, exists := remoteFiles[relPath]
		if !exists || info.ModTime().After(parseOBSTime(remoteObj.LastModified)) {
			toUpload = append(toUpload, relPath)
		}
	}

	if !quiet {
		fmt.Printf("Syncing %s → obs://%s/%s\n", localDir, container, prefix)
		fmt.Printf("  Files to upload: %d\n", len(toUpload))
	}

	for _, relPath := range toUpload {
		localPath := filepath.Join(localDir, relPath)
		objPath := prefix
		if objPath != "" && !strings.HasSuffix(objPath, "/") {
			objPath += "/"
		}
		objPath += relPath

		file, err := os.Open(localPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open %s: %v\n", localPath, err)
			continue
		}

		_, err = client.PutObject(ctx, &object.PutObjectInput{
			Container:  container,
			ObjectName: objPath,
			Body:       file,
		})
		if err != nil {
			file.Close()
			fmt.Fprintf(os.Stderr, "Failed to upload %s: %v\n", relPath, err)
			continue
		}
		file.Close()

		if !quiet {
			fmt.Printf("  Uploaded: %s\n", relPath)
		}
	}

	if deleteExtra {
		var toDelete []string
		for relPath := range remoteFiles {
			if _, exists := localFiles[relPath]; !exists {
				toDelete = append(toDelete, relPath)
			}
		}

		if !quiet && len(toDelete) > 0 {
			fmt.Printf("  Files to delete: %d\n", len(toDelete))
		}

		for _, relPath := range toDelete {
			objPath := prefix
			if objPath != "" && !strings.HasSuffix(objPath, "/") {
				objPath += "/"
			}
			objPath += relPath

			if err := client.DeleteObject(ctx, container, objPath); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to delete %s: %v\n", relPath, err)
				continue
			}

			if !quiet {
				fmt.Printf("  Deleted: %s\n", relPath)
			}
		}
	}

	if !quiet {
		fmt.Println("Sync complete.")
	}
}

func syncObsToLocal(ctx context.Context, obsPath, localDir string, deleteExtra, quiet bool, exclude string) {
	container, prefix := parseObsPath(obsPath)
	client := getObjectStorageClient()

	input := &object.ListObjectsInput{Prefix: prefix}
	result, err := client.ListObjects(ctx, container, input)
	if err != nil {
		exitWithError("Failed to list objects", err)
	}

	remoteFiles := make(map[string]object.Object)
	for _, obj := range result.Objects {
		relPath := strings.TrimPrefix(obj.Name, prefix)
		relPath = strings.TrimPrefix(relPath, "/")
		if exclude != "" {
			if matched, _ := filepath.Match(exclude, filepath.Base(relPath)); matched {
				continue
			}
		}
		remoteFiles[relPath] = obj
	}

	localFiles := make(map[string]os.FileInfo)
	if _, err := os.Stat(localDir); err == nil {
		filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			relPath, _ := filepath.Rel(localDir, path)
			localFiles[relPath] = info
			return nil
		})
	}

	var toDownload []string
	for relPath, obj := range remoteFiles {
		localInfo, exists := localFiles[relPath]
		if !exists || parseOBSTime(obj.LastModified).After(localInfo.ModTime()) {
			toDownload = append(toDownload, relPath)
		}
	}

	if !quiet {
		fmt.Printf("Syncing obs://%s/%s → %s\n", container, prefix, localDir)
		fmt.Printf("  Files to download: %d\n", len(toDownload))
	}

	creds := credentials.NewStaticIdentity(getUsername(), getPassword(), getTenantID())
	rawClient := createRawObjectClient(creds)

	for _, relPath := range toDownload {
		obj := remoteFiles[relPath]
		localPath := filepath.Join(localDir, relPath)

		dir := filepath.Dir(localPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create directory %s: %v\n", dir, err)
			continue
		}

		body, _, err := rawClient.GetObjectRaw(ctx, container, obj.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to download %s: %v\n", relPath, err)
			continue
		}

		outFile, err := os.Create(localPath)
		if err != nil {
			body.Close()
			fmt.Fprintf(os.Stderr, "Failed to create %s: %v\n", localPath, err)
			continue
		}

		_, err = io.Copy(outFile, body)
		body.Close()
		outFile.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", localPath, err)
			continue
		}

		mtime := parseOBSTime(obj.LastModified)
		os.Chtimes(localPath, mtime, mtime)

		if !quiet {
			fmt.Printf("  Downloaded: %s\n", relPath)
		}
	}

	if deleteExtra {
		var toDelete []string
		for relPath := range localFiles {
			if _, exists := remoteFiles[relPath]; !exists {
				toDelete = append(toDelete, relPath)
			}
		}

		if !quiet && len(toDelete) > 0 {
			fmt.Printf("  Files to delete: %d\n", len(toDelete))
		}

		for _, relPath := range toDelete {
			localPath := filepath.Join(localDir, relPath)
			if err := os.Remove(localPath); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to delete %s: %v\n", localPath, err)
				continue
			}

			if !quiet {
				fmt.Printf("  Deleted: %s\n", relPath)
			}
		}
	}

	if !quiet {
		fmt.Println("Sync complete.")
	}
}

type rawObjectClient struct {
	region        string
	credentials   credentials.IdentityCredentials
	tokenEndpoint string
	token         string
	storageURL    string
}

func createRawObjectClient(creds credentials.IdentityCredentials) *rawObjectClient {
	return &rawObjectClient{
		region:      getRegion(),
		credentials: creds,
	}
}

func (c *rawObjectClient) authenticate(ctx context.Context) error {
	if c.token != "" {
		return nil
	}

	authURL := fmt.Sprintf("https://api-identity-infrastructure.%s.nhncloudservice.com/v2.0/tokens", c.region)

	authBody := fmt.Sprintf(`{
		"auth": {
			"tenantId": "%s",
			"passwordCredentials": {
				"username": "%s",
				"password": "%s"
			}
		}
	}`, c.credentials.GetTenantID(), c.credentials.GetUsername(), c.credentials.GetPassword())

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, strings.NewReader(authBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed: %s", string(body))
	}

	c.token = resp.Header.Get("X-Subject-Token")
	if c.token == "" {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "X-Auth-Token" {
				c.token = cookie.Value
				break
			}
		}
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	start := strings.Index(bodyStr, `"object-store"`)
	if start == -1 {
		return fmt.Errorf("object-store endpoint not found in catalog")
	}

	urlStart := strings.Index(bodyStr[start:], `"publicURL":"`) + start + 13
	urlEnd := strings.Index(bodyStr[urlStart:], `"`) + urlStart
	c.storageURL = bodyStr[urlStart:urlEnd]

	tokenStart := strings.Index(bodyStr, `"id":"`) + 6
	tokenEnd := strings.Index(bodyStr[tokenStart:], `"`) + tokenStart
	if c.token == "" {
		c.token = bodyStr[tokenStart:tokenEnd]
	}

	return nil
}

func (c *rawObjectClient) GetObjectRaw(ctx context.Context, container, objectName string) (io.ReadCloser, int64, error) {
	if err := c.authenticate(ctx); err != nil {
		return nil, 0, err
	}

	url := fmt.Sprintf("%s/%s/%s", c.storageURL, container, objectName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("X-Auth-Token", c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("failed to get object: status %d", resp.StatusCode)
	}

	return resp.Body, resp.ContentLength, nil
}
