package object

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/credentials"
	"github.com/haung921209/nhn-cloud-sdk-go/nhncloud/internal/client"
)

type Client struct {
	region        string
	credentials   credentials.IdentityCredentials
	httpClient    *http.Client
	tokenProvider *client.IdentityTokenProvider
	baseURL       string
	debug         bool
}

func NewClient(region string, creds credentials.IdentityCredentials, hc *http.Client, debug bool) *Client {
	c := &Client{
		region:      region,
		credentials: creds,
		debug:       debug,
	}

	if creds != nil {
		c.tokenProvider = client.NewIdentityTokenProvider(
			creds.GetTenantID(),
			creds.GetUsername(),
			creds.GetPassword(),
		)
	}

	return c
}

func (c *Client) ensureClient(ctx context.Context) error {
	if c.baseURL != "" {
		return nil
	}

	if c.tokenProvider == nil {
		return ErrNoCredentials
	}

	if _, err := c.tokenProvider.GetToken(ctx); err != nil {
		return fmt.Errorf("authenticate: %w", err)
	}

	baseURL, err := c.tokenProvider.GetServiceEndpoint("object-store", c.region)
	if err != nil || baseURL == "" {
		// Fallback to documented public endpoint pattern
		tenantID := c.credentials.GetTenantID()
		baseURL = fmt.Sprintf("https://%s-api-object-storage.nhncloudservice.com/v1/AUTH_%s", c.region, tenantID)
	}

	c.baseURL = baseURL
	return nil
}

func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	if err := c.ensureClient(ctx); err != nil {
		return nil, err
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	token, _ := c.tokenProvider.GetToken(ctx)
	req.Header.Set("X-Auth-Token", token)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Debug logging
	if c.debug {
		fmt.Printf("--- api request ---\n")
		fmt.Printf("%s %s\n", req.Method, req.URL.String())
		for k, v := range req.Header {
			fmt.Printf("%s: %s\n", k, strings.Join(v, ","))
		}
		fmt.Printf("-------------------\n")
	}

	httpClient := c.httpClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if c.debug {
		fmt.Printf("--- api response ---\n")
		fmt.Printf("Status: %s\n", resp.Status)
		for k, v := range resp.Header {
			fmt.Printf("%s: %s\n", k, strings.Join(v, ","))
		}
		// Note: Body reading would consume it, so we don't dump body here unless we buffer it.
		// For now header/status is enough to debug 403.
		fmt.Printf("--------------------\n")
	}

	return resp, nil
}

func (c *Client) GetAccountInfo(ctx context.Context) (*AccountInfo, error) {
	resp, err := c.doRequest(ctx, http.MethodHead, "", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get account info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("get account info: status %d", resp.StatusCode)
	}

	containerCount, _ := strconv.ParseInt(resp.Header.Get("X-Account-Container-Count"), 10, 64)
	objectCount, _ := strconv.ParseInt(resp.Header.Get("X-Account-Object-Count"), 10, 64)
	bytesUsed, _ := strconv.ParseInt(resp.Header.Get("X-Account-Bytes-Used"), 10, 64)

	return &AccountInfo{
		ContainerCount: containerCount,
		ObjectCount:    objectCount,
		BytesUsed:      bytesUsed,
	}, nil
}

func (c *Client) ListContainers(ctx context.Context, input *ListContainersInput) (*ListContainersOutput, error) {
	path := "?format=json"
	if input != nil {
		if input.Marker != "" {
			path += "&marker=" + input.Marker
		}
		if input.Prefix != "" {
			path += "&prefix=" + input.Prefix
		}
		if input.Limit > 0 {
			path += fmt.Sprintf("&limit=%d", input.Limit)
		}
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list containers: status %d", resp.StatusCode)
	}

	var containers []Container
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		return nil, fmt.Errorf("list containers: decode: %w", err)
	}

	return &ListContainersOutput{Containers: containers}, nil
}

func (c *Client) CreateContainer(ctx context.Context, input *CreateContainerInput) error {
	headers := make(map[string]string)

	if input.StoragePolicy != "" {
		headers["X-Storage-Policy"] = input.StoragePolicy
	}
	if input.WormRetentionDay > 0 {
		headers["X-Container-Worm-Retention-Day"] = strconv.Itoa(input.WormRetentionDay)
	}
	if input.ReadACL != "" {
		headers["X-Container-Read"] = input.ReadACL
	}
	if input.WriteACL != "" {
		headers["X-Container-Write"] = input.WriteACL
	}
	for k, v := range input.Metadata {
		headers["X-Container-Meta-"+k] = v
	}

	resp, err := c.doRequest(ctx, http.MethodPut, "/"+input.Name, nil, headers)
	if err != nil {
		return fmt.Errorf("create container %s: %w", input.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("create container %s: status %d", input.Name, resp.StatusCode)
	}

	return nil
}

func (c *Client) GetContainerInfo(ctx context.Context, containerName string) (*ContainerInfo, error) {
	resp, err := c.doRequest(ctx, http.MethodHead, "/"+containerName, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get container info %s: %w", containerName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("get container info %s: status %d", containerName, resp.StatusCode)
	}

	objectCount, _ := strconv.ParseInt(resp.Header.Get("X-Container-Object-Count"), 10, 64)
	bytesUsed, _ := strconv.ParseInt(resp.Header.Get("X-Container-Bytes-Used"), 10, 64)
	objectLifecycle, _ := strconv.Atoi(resp.Header.Get("X-Container-Object-Lifecycle"))
	versionsRetention, _ := strconv.Atoi(resp.Header.Get("X-Versions-Retention"))
	wormRetentionDay, _ := strconv.Atoi(resp.Header.Get("X-Container-Worm-Retention-Day"))

	info := &ContainerInfo{
		Name:                 containerName,
		ObjectCount:          objectCount,
		BytesUsed:            bytesUsed,
		StoragePolicy:        resp.Header.Get("X-Storage-Policy"),
		ReadACL:              resp.Header.Get("X-Container-Read"),
		WriteACL:             resp.Header.Get("X-Container-Write"),
		ViewACL:              resp.Header.Get("X-Container-View"),
		ObjectLifecycle:      objectLifecycle,
		ObjectTransferTo:     resp.Header.Get("X-Container-Object-Transfer-To"),
		HistoryLocation:      resp.Header.Get("X-History-Location"),
		VersionsRetention:    versionsRetention,
		WormRetentionDay:     wormRetentionDay,
		WebIndex:             resp.Header.Get("X-Container-Meta-Web-Index"),
		WebError:             resp.Header.Get("X-Container-Meta-Web-Error"),
		CORSAllowOrigin:      resp.Header.Get("X-Container-Meta-Access-Control-Allow-Origin"),
		RfcCompliantEtags:    resp.Header.Get("X-Container-Rfc-Compliant-Etags") == "true",
		IPACLAllowedList:     resp.Header.Get("X-Container-Ip-Acl-Allowed-List"),
		IPACLDeniedList:      resp.Header.Get("X-Container-Ip-Acl-Denied-List"),
		DenyExtensionPolicy:  resp.Header.Get("X-Container-Object-Deny-Extension-Policy"),
		DenyKeywordPolicy:    resp.Header.Get("X-Container-Object-Deny-Keyword-Policy"),
		AllowExtensionPolicy: resp.Header.Get("X-Container-Object-Allow-Extension-Policy"),
		AllowKeywordPolicy:   resp.Header.Get("X-Container-Object-Allow-Keyword-Policy"),
		CustomMetadata:       make(map[string]string),
	}

	for k, v := range resp.Header {
		if strings.HasPrefix(k, "X-Container-Meta-") && k != "X-Container-Meta-Web-Index" && k != "X-Container-Meta-Web-Error" && k != "X-Container-Meta-Access-Control-Allow-Origin" {
			info.CustomMetadata[strings.TrimPrefix(k, "X-Container-Meta-")] = v[0]
		}
	}

	return info, nil
}

func (c *Client) UpdateContainer(ctx context.Context, input *UpdateContainerInput) error {
	headers := make(map[string]string)

	if input.ReadACL != "" {
		headers["X-Container-Read"] = input.ReadACL
	}
	if input.WriteACL != "" {
		headers["X-Container-Write"] = input.WriteACL
	}
	if input.ViewACL != "" {
		headers["X-Container-View"] = input.ViewACL
	}
	if input.IPACLAllowedList != "" {
		headers["X-Container-Ip-Acl-Allowed-List"] = input.IPACLAllowedList
	}
	if input.IPACLDeniedList != "" {
		headers["X-Container-Ip-Acl-Denied-List"] = input.IPACLDeniedList
	}
	if input.ObjectLifecycle != nil {
		headers["X-Container-Object-Lifecycle"] = strconv.Itoa(*input.ObjectLifecycle)
	}
	if input.ObjectTransferTo != "" {
		headers["X-Container-Object-Transfer-To"] = input.ObjectTransferTo
	}
	if input.HistoryLocation != "" {
		headers["X-History-Location"] = input.HistoryLocation
	}
	if input.VersionsRetention != nil {
		headers["X-Versions-Retention"] = strconv.Itoa(*input.VersionsRetention)
	}
	if input.WebIndex != "" {
		headers["X-Container-Meta-Web-Index"] = input.WebIndex
	}
	if input.WebError != "" {
		headers["X-Container-Meta-Web-Error"] = input.WebError
	}
	if input.CORSAllowOrigin != "" {
		headers["X-Container-Meta-Access-Control-Allow-Origin"] = input.CORSAllowOrigin
	}
	if input.RfcCompliantEtags != nil {
		headers["X-Container-Rfc-Compliant-Etags"] = strconv.FormatBool(*input.RfcCompliantEtags)
	}
	if input.WormRetentionDay != nil {
		headers["X-Container-Worm-Retention-Day"] = strconv.Itoa(*input.WormRetentionDay)
	}
	if input.DenyExtensionPolicy != "" {
		headers["X-Container-Object-Deny-Extension-Policy"] = input.DenyExtensionPolicy
	}
	if input.DenyKeywordPolicy != "" {
		headers["X-Container-Object-Deny-Keyword-Policy"] = input.DenyKeywordPolicy
	}
	if input.AllowExtensionPolicy != "" {
		headers["X-Container-Object-Allow-Extension-Policy"] = input.AllowExtensionPolicy
	}
	if input.AllowKeywordPolicy != "" {
		headers["X-Container-Object-Allow-Keyword-Policy"] = input.AllowKeywordPolicy
	}
	for k, v := range input.Metadata {
		headers["X-Container-Meta-"+k] = v
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/"+input.Name, nil, headers)
	if err != nil {
		return fmt.Errorf("update container %s: %w", input.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("update container %s: status %d", input.Name, resp.StatusCode)
	}

	return nil
}

func (c *Client) DeleteContainer(ctx context.Context, containerName string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, "/"+containerName, nil, nil)
	if err != nil {
		return fmt.Errorf("delete container %s: %w", containerName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete container %s: status %d", containerName, resp.StatusCode)
	}

	return nil
}

func (c *Client) ListObjects(ctx context.Context, containerName string, input *ListObjectsInput) (*ListObjectsOutput, error) {
	path := "/" + containerName + "?format=xml"
	if input != nil {
		if input.Prefix != "" {
			path += "&prefix=" + input.Prefix
		}
		if input.Delimiter != "" {
			path += "&delimiter=" + input.Delimiter
		}
		if input.Marker != "" {
			path += "&marker=" + input.Marker
		}
		if input.Limit > 0 {
			path += fmt.Sprintf("&limit=%d", input.Limit)
		}
	}

	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("list objects in %s: %w", containerName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list objects in %s: status %d", containerName, resp.StatusCode)
	}

	var result struct {
		Container struct {
			Name string `xml:"name,attr"`
		} `xml:"container"`
		Objects []struct {
			Name         string `xml:"name"`
			Hash         string `xml:"hash"`
			Bytes        int64  `xml:"bytes"`
			ContentType  string `xml:"content_type"`
			LastModified string `xml:"last_modified"`
		} `xml:"object"`
		CommonPrefixes []struct {
			Prefix string `xml:"name,attr"`
		} `xml:"subdir"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("list objects %s: decode: %w", containerName, err)
	}

	objects := make([]Object, len(result.Objects))
	for i, o := range result.Objects {
		objects[i] = Object{
			Name:         o.Name,
			Hash:         o.Hash,
			Bytes:        o.Bytes,
			ContentType:  o.ContentType,
			LastModified: o.LastModified,
		}
	}

	prefixes := make([]string, len(result.CommonPrefixes))
	for i, p := range result.CommonPrefixes {
		prefixes[i] = p.Prefix
	}

	return &ListObjectsOutput{
		Objects:        objects,
		CommonPrefixes: prefixes,
	}, nil
}

func (c *Client) PutObject(ctx context.Context, input *PutObjectInput) (*PutObjectOutput, error) {
	headers := make(map[string]string)

	if input.ContentType != "" {
		headers["Content-Type"] = input.ContentType
	}
	if input.DeleteAt != nil {
		headers["X-Delete-At"] = strconv.FormatInt(*input.DeleteAt, 10)
	}
	if input.DeleteAfter != nil {
		headers["X-Delete-After"] = strconv.FormatInt(*input.DeleteAfter, 10)
	}
	for k, v := range input.Metadata {
		headers["X-Object-Meta-"+k] = v
	}

	path := "/" + input.Container + "/" + input.ObjectName
	resp, err := c.doRequest(ctx, http.MethodPut, path, input.Body, headers)
	if err != nil {
		return nil, fmt.Errorf("put object %s/%s: %w", input.Container, input.ObjectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("put object %s/%s: status %d", input.Container, input.ObjectName, resp.StatusCode)
	}

	return &PutObjectOutput{
		ETag: resp.Header.Get("Etag"),
	}, nil
}

func (c *Client) GetObjectInfo(ctx context.Context, containerName, objectName string) (*ObjectInfo, error) {
	path := "/" + containerName + "/" + objectName
	resp, err := c.doRequest(ctx, http.MethodHead, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get object info %s/%s: %w", containerName, objectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get object info %s/%s: status %d", containerName, objectName, resp.StatusCode)
	}

	contentLength, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	timestamp, _ := strconv.ParseInt(resp.Header.Get("X-Timestamp"), 10, 64)

	info := &ObjectInfo{
		ContentType:       resp.Header.Get("Content-Type"),
		ContentLength:     contentLength,
		ETag:              resp.Header.Get("Etag"),
		Timestamp:         timestamp,
		ObjectManifest:    resp.Header.Get("X-Object-Manifest"),
		StaticLargeObject: resp.Header.Get("X-Static-Large-Object") == "True",
		ManifestETag:      resp.Header.Get("X-Manifest-Etag"),
		CustomMetadata:    make(map[string]string),
	}

	if deleteAt := resp.Header.Get("X-Delete-At"); deleteAt != "" {
		v, _ := strconv.ParseInt(deleteAt, 10, 64)
		info.DeleteAt = &v
	}
	if wormRetain := resp.Header.Get("X-Object-Worm-Retain-Until"); wormRetain != "" {
		v, _ := strconv.ParseInt(wormRetain, 10, 64)
		info.WormRetainUntil = &v
	}

	for k, v := range resp.Header {
		if strings.HasPrefix(k, "X-Object-Meta-") {
			info.CustomMetadata[strings.TrimPrefix(k, "X-Object-Meta-")] = v[0]
		}
	}

	return info, nil
}

func (c *Client) GetObject(ctx context.Context, containerName, objectName string) (*GetObjectOutput, error) {
	path := "/" + containerName + "/" + objectName
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get object %s/%s: %w", containerName, objectName, err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("get object %s/%s: status %d", containerName, objectName, resp.StatusCode)
	}

	contentLength, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	metadata := make(map[string]string)
	for k, v := range resp.Header {
		if strings.HasPrefix(k, "X-Object-Meta-") {
			metadata[strings.TrimPrefix(k, "X-Object-Meta-")] = v[0]
		}
	}

	return &GetObjectOutput{
		Body:          resp.Body,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: contentLength,
		ETag:          resp.Header.Get("Etag"),
		Metadata:      metadata,
	}, nil
}

func (c *Client) CopyObject(ctx context.Context, input *CopyObjectInput) error {
	headers := map[string]string{
		"Destination": "/" + input.DestinationContainer + "/" + input.DestinationObjectName,
	}

	path := "/" + input.SourceContainer + "/" + input.SourceObjectName
	resp, err := c.doRequest(ctx, "COPY", path, nil, headers)
	if err != nil {
		return fmt.Errorf("copy object %s/%s: %w", input.SourceContainer, input.SourceObjectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("copy object %s/%s: status %d", input.SourceContainer, input.SourceObjectName, resp.StatusCode)
	}

	return nil
}

func (c *Client) UpdateObjectMetadata(ctx context.Context, input *UpdateObjectMetadataInput) error {
	headers := make(map[string]string)

	if input.DeleteAt != nil {
		headers["X-Delete-At"] = strconv.FormatInt(*input.DeleteAt, 10)
	}
	if input.DeleteAfter != nil {
		headers["X-Delete-After"] = strconv.FormatInt(*input.DeleteAfter, 10)
	}
	if input.WormRetainUntil != nil {
		headers["X-Object-Worm-Retain-Until"] = strconv.FormatInt(*input.WormRetainUntil, 10)
	}
	for k, v := range input.Metadata {
		headers["X-Object-Meta-"+k] = v
	}

	path := "/" + input.Container + "/" + input.ObjectName
	resp, err := c.doRequest(ctx, http.MethodPost, path, nil, headers)
	if err != nil {
		return fmt.Errorf("update object metadata %s/%s: %w", input.Container, input.ObjectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("update object metadata %s/%s: status %d", input.Container, input.ObjectName, resp.StatusCode)
	}

	return nil
}

func (c *Client) DeleteObject(ctx context.Context, containerName, objectName string) error {
	path := "/" + containerName + "/" + objectName
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return fmt.Errorf("delete object %s/%s: %w", containerName, objectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete object %s/%s: status %d", containerName, objectName, resp.StatusCode)
	}

	return nil
}

func (c *Client) UploadSegment(ctx context.Context, input *UploadSegmentInput) (*PutObjectOutput, error) {
	headers := make(map[string]string)
	if input.ContentType != "" {
		headers["Content-Type"] = input.ContentType
	}

	path := fmt.Sprintf("/%s/%s/%03d", input.Container, input.ObjectName, input.SegmentIndex)
	resp, err := c.doRequest(ctx, http.MethodPut, path, input.Body, headers)
	if err != nil {
		return nil, fmt.Errorf("upload segment %s/%s/%d: %w", input.Container, input.ObjectName, input.SegmentIndex, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("upload segment %s/%s/%d: status %d", input.Container, input.ObjectName, input.SegmentIndex, resp.StatusCode)
	}

	return &PutObjectOutput{
		ETag: resp.Header.Get("Etag"),
	}, nil
}

func (c *Client) CreateDLOManifest(ctx context.Context, input *CreateDLOManifestInput) error {
	headers := map[string]string{
		"X-Object-Manifest": input.SegmentContainer + "/" + input.SegmentPrefix + "/",
	}
	if input.ContentType != "" {
		headers["Content-Type"] = input.ContentType
	}

	path := "/" + input.Container + "/" + input.ObjectName
	resp, err := c.doRequest(ctx, http.MethodPut, path, strings.NewReader(""), headers)
	if err != nil {
		return fmt.Errorf("create DLO manifest %s/%s: %w", input.Container, input.ObjectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("create DLO manifest %s/%s: status %d", input.Container, input.ObjectName, resp.StatusCode)
	}

	return nil
}

func (c *Client) CreateSLOManifest(ctx context.Context, input *CreateSLOManifestInput) error {
	body, err := json.Marshal(input.Segments)
	if err != nil {
		return fmt.Errorf("create SLO manifest: marshal segments: %w", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	if input.ContentType != "" {
		headers["Content-Type"] = input.ContentType
	}

	path := "/" + input.Container + "/" + input.ObjectName + "?multipart-manifest=put"
	resp, err := c.doRequest(ctx, http.MethodPut, path, strings.NewReader(string(body)), headers)
	if err != nil {
		return fmt.Errorf("create SLO manifest %s/%s: %w", input.Container, input.ObjectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("create SLO manifest %s/%s: status %d", input.Container, input.ObjectName, resp.StatusCode)
	}

	return nil
}

func (c *Client) GetSLOManifest(ctx context.Context, containerName, objectName string) (*GetSLOManifestOutput, error) {
	path := "/" + containerName + "/" + objectName + "?multipart-manifest=get"
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get SLO manifest %s/%s: %w", containerName, objectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get SLO manifest %s/%s: status %d", containerName, objectName, resp.StatusCode)
	}

	var segments []SLOSegment
	if err := json.NewDecoder(resp.Body).Decode(&segments); err != nil {
		return nil, fmt.Errorf("get SLO manifest %s/%s: decode: %w", containerName, objectName, err)
	}

	return &GetSLOManifestOutput{Segments: segments}, nil
}
