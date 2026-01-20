package object

import (
	"io"
)

// Storage Classes
const (
	StorageClassStandard = "Standard"
	StorageClassEconomy  = "Economy"
)

// ==================== Account Types ====================
// AccountInfo represents storage account information
type AccountInfo struct {
	ContainerCount int64
	ObjectCount    int64
	BytesUsed      int64
}

// ==================== Container Types ====================
// Container represents a storage container
type Container struct {
	Name         string `json:"name"`
	Count        int64  `json:"count"`
	Bytes        int64  `json:"bytes"`
	LastModified string `json:"last_modified,omitempty"`
}

// ListContainersInput contains parameters for listing containers
type ListContainersInput struct {
	Format string // json or xml
	Marker string // pagination marker
	Limit  int    // max results
	Prefix string // filter by prefix
}

// ListContainersOutput contains the list of containers
type ListContainersOutput struct {
	Containers []Container
}

// ContainerInfo represents detailed container information
type ContainerInfo struct {
	Name                 string
	ObjectCount          int64
	BytesUsed            int64
	StoragePolicy        string
	ReadACL              string
	WriteACL             string
	ViewACL              string
	ObjectLifecycle      int
	ObjectTransferTo     string
	HistoryLocation      string
	VersionsRetention    int
	WormRetentionDay     int
	WebIndex             string
	WebError             string
	CORSAllowOrigin      string
	RfcCompliantEtags    bool
	IPACLAllowedList     string
	IPACLDeniedList      string
	DenyExtensionPolicy  string
	DenyKeywordPolicy    string
	AllowExtensionPolicy string
	AllowKeywordPolicy   string
	CustomMetadata       map[string]string
}

// CreateContainerInput contains parameters for creating a container
type CreateContainerInput struct {
	Name             string
	StoragePolicy    string // Standard or Economy
	WormRetentionDay int    // Object lock period in days
	ReadACL          string
	WriteACL         string
	Metadata         map[string]string
}

// UpdateContainerInput contains parameters for updating container settings
type UpdateContainerInput struct {
	Name                 string
	ReadACL              string
	WriteACL             string
	ViewACL              string
	IPACLAllowedList     string
	IPACLDeniedList      string
	ObjectLifecycle      *int   // days, nil to clear
	ObjectTransferTo     string // container for expired objects
	HistoryLocation      string // archive container for versioning
	VersionsRetention    *int   // days for version retention
	WebIndex             string // static website index document
	WebError             string // static website error suffix
	CORSAllowOrigin      string // CORS allowed origins
	RfcCompliantEtags    *bool  // RFC-compliant ETag format
	WormRetentionDay     *int   // Object lock period (can only extend)
	DenyExtensionPolicy  string // upload blacklist extensions
	DenyKeywordPolicy    string // upload blacklist keywords
	AllowExtensionPolicy string // upload whitelist extensions
	AllowKeywordPolicy   string // upload whitelist keywords
	Metadata             map[string]string
}

// ==================== Object Types ====================
// Object represents an object in a container
type Object struct {
	Name         string `json:"name"`
	Hash         string `json:"hash"`
	Bytes        int64  `json:"bytes"`
	ContentType  string `json:"content_type"`
	LastModified string `json:"last_modified"`
	Subdir       string `json:"subdir,omitempty"` // for delimiter queries
}

// ListObjectsInput contains parameters for listing objects
type ListObjectsInput struct {
	Prefix    string
	Delimiter string
	Marker    string
	Limit     int
	Format    string // json or xml
}

// ListObjectsOutput contains the list of objects
type ListObjectsOutput struct {
	Objects        []Object
	CommonPrefixes []string
}

// ObjectInfo represents detailed object information from HEAD request
type ObjectInfo struct {
	ContentType       string
	ContentLength     int64
	ETag              string
	LastModified      string
	Timestamp         int64
	DeleteAt          *int64
	WormRetainUntil   *int64
	ObjectManifest    string // DLO segment path
	StaticLargeObject bool   // SLO indicator
	ManifestETag      string // SLO manifest ETag
	CustomMetadata    map[string]string
}

// PutObjectInput contains parameters for uploading an object
type PutObjectInput struct {
	Container   string
	ObjectName  string
	Body        io.Reader
	ContentType string
	DeleteAt    *int64 // Unix timestamp for expiry
	DeleteAfter *int64 // TTL in seconds
	Metadata    map[string]string
}

// PutObjectOutput contains the result of object upload
type PutObjectOutput struct {
	ETag         string
	LastModified string
}

// GetObjectOutput contains the result of object download
type GetObjectOutput struct {
	Body          io.ReadCloser
	ContentType   string
	ContentLength int64
	ETag          string
	LastModified  string
	Metadata      map[string]string
}

// CopyObjectInput contains parameters for copying an object
type CopyObjectInput struct {
	SourceContainer       string
	SourceObjectName      string
	DestinationContainer  string
	DestinationObjectName string
}

// UpdateObjectMetadataInput contains parameters for updating object metadata
type UpdateObjectMetadataInput struct {
	Container       string
	ObjectName      string
	DeleteAt        *int64 // Unix timestamp for expiry
	DeleteAfter     *int64 // TTL in seconds
	WormRetainUntil *int64 // Lock expiry (can only extend)
	Metadata        map[string]string
}

// ==================== Multipart Upload Types ====================
// UploadSegmentInput contains parameters for uploading a segment
type UploadSegmentInput struct {
	Container    string
	ObjectName   string
	SegmentIndex int // sequence number (1, 2, 3, ...)
	Body         io.Reader
	ContentType  string
}

// CreateDLOManifestInput contains parameters for creating a DLO manifest
type CreateDLOManifestInput struct {
	Container        string
	ObjectName       string
	SegmentContainer string
	SegmentPrefix    string // segments path prefix
	ContentType      string
}

// SLOSegment represents a segment in an SLO manifest
type SLOSegment struct {
	Path      string `json:"path"`       // {container}/{object}
	ETag      string `json:"etag"`       // segment ETag
	SizeBytes int64  `json:"size_bytes"` // segment size
}

// CreateSLOManifestInput contains parameters for creating an SLO manifest
type CreateSLOManifestInput struct {
	Container   string
	ObjectName  string
	Segments    []SLOSegment
	ContentType string
}

// GetSLOManifestOutput contains the segments of an SLO
type GetSLOManifestOutput struct {
	Segments []SLOSegment
}

// ==================== Legacy Types (for backward compatibility) ====================
// ObjectMetadata for backward compatibility
type ObjectMetadata struct {
	ContentType        string
	ContentDisposition string
	ContentEncoding    string
	Metadata           map[string]string
}

// ContainerMetadata for backward compatibility
type ContainerMetadata struct {
	ReadACL          string
	WriteACL         string
	SyncTo           string
	SyncKey          string
	VersionsLocation string
	ObjectCount      int64
	BytesUsed        int64
	CustomMetadata   map[string]string
}
