package block

type Volume struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Status           string             `json:"status"`
	Size             int                `json:"size"`
	VolumeType       string             `json:"volume_type"`
	Bootable         string             `json:"bootable"`
	Encrypted        bool               `json:"encrypted"`
	AvailabilityZone string             `json:"availability_zone"`
	SnapshotID       string             `json:"snapshot_id,omitempty"`
	SourceVolID      string             `json:"source_volid,omitempty"`
	Description      string             `json:"description,omitempty"`
	Attachments      []VolumeAttachment `json:"attachments,omitempty"`
	Metadata         map[string]string  `json:"metadata,omitempty"`
	CreatedAt        string             `json:"created_at"`
	UpdatedAt        string             `json:"updated_at,omitempty"`
}

type VolumeAttachment struct {
	ID         string `json:"id"`
	VolumeID   string `json:"volume_id"`
	ServerID   string `json:"server_id"`
	Device     string `json:"device"`
	AttachedAt string `json:"attached_at"`
}

type Snapshot struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Status      string            `json:"status"`
	Size        int               `json:"size"`
	VolumeID    string            `json:"volume_id"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at,omitempty"`
}

type VolumeType struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	ExtraSpecs  map[string]string `json:"extra_specs,omitempty"`
}

type ListVolumesOutput struct {
	Volumes []Volume `json:"volumes"`
}

type GetVolumeOutput struct {
	Volume Volume `json:"volume"`
}

type CreateVolumeInput struct {
	Name             string            `json:"name,omitempty"`
	Size             int               `json:"size"`
	VolumeType       string            `json:"volume_type,omitempty"`
	AvailabilityZone string            `json:"availability_zone,omitempty"`
	SnapshotID       string            `json:"snapshot_id,omitempty"`
	SourceVolID      string            `json:"source_volid,omitempty"`
	Description      string            `json:"description,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

type CreateVolumeOutput struct {
	Volume Volume `json:"volume"`
}

type UpdateVolumeInput struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type ExtendVolumeInput struct {
	NewSize int `json:"new_size"`
}

type AttachVolumeInput struct {
	ServerID string `json:"instance_uuid"`
	Device   string `json:"mountpoint,omitempty"`
}

type ListSnapshotsOutput struct {
	Snapshots []Snapshot `json:"snapshots"`
}

type GetSnapshotOutput struct {
	Snapshot Snapshot `json:"snapshot"`
}

type CreateSnapshotInput struct {
	Name        string `json:"name,omitempty"`
	VolumeID    string `json:"volume_id"`
	Description string `json:"description,omitempty"`
	Force       bool   `json:"force,omitempty"`
}

type CreateSnapshotOutput struct {
	Snapshot Snapshot `json:"snapshot"`
}

type ListVolumeTypesOutput struct {
	VolumeTypes []VolumeType `json:"volume_types"`
}
