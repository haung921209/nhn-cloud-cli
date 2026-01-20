package image

import "time"

type Image struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	Visibility      string    `json:"visibility"`
	Protected       bool      `json:"protected"`
	Owner           string    `json:"owner"`
	Size            int64     `json:"size"`
	VirtualSize     *int64    `json:"virtual_size"`
	MinDisk         int       `json:"min_disk"`
	MinRAM          int       `json:"min_ram"`
	DiskFormat      string    `json:"disk_format"`
	ContainerFormat string    `json:"container_format"`
	OSType          string    `json:"os_type,omitempty"`
	OSDistro        string    `json:"os_distro,omitempty"`
	OSVersion       string    `json:"os_version,omitempty"`
	OSArchitecture  string    `json:"os_architecture,omitempty"`
	LoginUsername   string    `json:"login_username,omitempty"`
	HypervisorType  string    `json:"hypervisor_type,omitempty"`
	Description     string    `json:"description,omitempty"`
	Tags            []string  `json:"tags"`
	Checksum        string    `json:"checksum,omitempty"`
	Self            string    `json:"self"`
	File            string    `json:"file"`
	Schema          string    `json:"schema"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ListImagesOutput struct {
	Images []Image `json:"images"`
	Schema string  `json:"schema"`
	First  string  `json:"first,omitempty"`
	Next   string  `json:"next,omitempty"`
}

type ListImagesInput struct {
	Limit      int    `json:"limit,omitempty"`
	Name       string `json:"name,omitempty"`
	Owner      string `json:"owner,omitempty"`
	SizeMin    int64  `json:"size_min,omitempty"`
	SizeMax    int64  `json:"size_max,omitempty"`
	Status     string `json:"status,omitempty"`
	Visibility string `json:"visibility,omitempty"`
	OSType     string `json:"os_type,omitempty"`
	OSDistro   string `json:"os_distro,omitempty"`
	Marker     string `json:"marker,omitempty"`
	SortKey    string `json:"sort_key,omitempty"`
	SortDir    string `json:"sort_dir,omitempty"`
}

type CreateImageInput struct {
	Name            string   `json:"name"`
	ContainerFormat string   `json:"container_format,omitempty"`
	DiskFormat      string   `json:"disk_format,omitempty"`
	MinDisk         int      `json:"min_disk,omitempty"`
	MinRAM          int      `json:"min_ram,omitempty"`
	Protected       bool     `json:"protected,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	OSDistro        string   `json:"os_distro,omitempty"`
	OSVersion       string   `json:"os_version,omitempty"`
	OSType          string   `json:"os_type,omitempty"`
}

type UpdateImageOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type ImageMember struct {
	ImageID   string    `json:"image_id"`
	MemberID  string    `json:"member_id"`
	Status    string    `json:"status"`
	Schema    string    `json:"schema"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListImageMembersOutput struct {
	Members []ImageMember `json:"members"`
	Schema  string        `json:"schema"`
}

type CreateImageMemberInput struct {
	Member string `json:"member"`
}

type UpdateImageMemberInput struct {
	Status string `json:"status"`
}
