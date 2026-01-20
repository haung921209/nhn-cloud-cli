package compute

type Server struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Status           string               `json:"status"`
	TenantID         string               `json:"tenant_id"`
	UserID           string               `json:"user_id"`
	KeyName          string               `json:"key_name,omitempty"`
	ImageID          string               `json:"imageId,omitempty"`
	FlavorID         string               `json:"flavor>id,omitempty"`
	AvailabilityZone string               `json:"OS-EXT-AZ:availability_zone,omitempty"`
	Created          string               `json:"created"`
	Updated          string               `json:"updated,omitempty"`
	Addresses        map[string][]Address `json:"addresses,omitempty"`
	Metadata         map[string]string    `json:"metadata,omitempty"`
	SecurityGroups   []SecurityGroup      `json:"security_groups,omitempty"`
}

type Address struct {
	Addr    string `json:"addr"`
	Version int    `json:"version"`
	Type    string `json:"OS-EXT-IPS:type,omitempty"`
}

type SecurityGroup struct {
	Name string `json:"name"`
}

type Flavor struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	RAM   int    `json:"ram"`
	VCPUs int    `json:"vcpus"`
	Disk  int    `json:"disk"`
}

type Image struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	MinDisk int    `json:"minDisk"`
	MinRAM  int    `json:"minRam"`
	Created string `json:"created"`
	Updated string `json:"updated,omitempty"`
}

type KeyPair struct {
	Name        string `json:"name"`
	PublicKey   string `json:"public_key"`
	Fingerprint string `json:"fingerprint"`
}

type ListServersOutput struct {
	Servers []Server `json:"servers"`
}

type GetServerOutput struct {
	Server Server `json:"server"`
}

type CreateServerInput struct {
	Name               string               `json:"name"`
	ImageRef           string               `json:"imageRef"`
	FlavorRef          string               `json:"flavorRef"`
	KeyName            string               `json:"key_name,omitempty"`
	AvailabilityZone   string               `json:"availability_zone,omitempty"`
	Networks           []ServerNetwork      `json:"networks,omitempty"`
	SecurityGroups     []SecurityGroup      `json:"security_groups,omitempty"`
	Metadata           map[string]string    `json:"metadata,omitempty"`
	UserData           string               `json:"user_data,omitempty"`
	BlockDeviceMapping []BlockDeviceMapping `json:"block_device_mapping_v2,omitempty"`
}

type ServerNetwork struct {
	UUID    string `json:"uuid,omitempty"`
	Port    string `json:"port,omitempty"`
	FixedIP string `json:"fixed_ip,omitempty"`
}

type BlockDeviceMapping struct {
	BootIndex           int    `json:"boot_index"`
	UUID                string `json:"uuid,omitempty"`
	SourceType          string `json:"source_type"`
	DestinationType     string `json:"destination_type"`
	VolumeSize          int    `json:"volume_size,omitempty"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
}

type CreateServerOutput struct {
	Server Server `json:"server"`
}

type ListFlavorsOutput struct {
	Flavors []Flavor `json:"flavors"`
}

type ListImagesOutput struct {
	Images []Image `json:"images"`
}

type ListKeyPairsOutput struct {
	KeyPairs []struct {
		KeyPair KeyPair `json:"keypair"`
	} `json:"keypairs"`
}

type CreateKeyPairInput struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key,omitempty"`
}

type CreateKeyPairOutput struct {
	KeyPair struct {
		Name        string `json:"name"`
		PublicKey   string `json:"public_key"`
		PrivateKey  string `json:"private_key,omitempty"`
		Fingerprint string `json:"fingerprint"`
	} `json:"keypair"`
}

type ActionInput struct {
	Action string `json:"action,omitempty"`
}

type RebootInput struct {
	Reboot struct {
		Type string `json:"type"`
	} `json:"reboot"`
}

type ResizeInput struct {
	Resize struct {
		FlavorRef string `json:"flavorRef"`
	} `json:"resize"`
}

type AvailabilityZone struct {
	ZoneName  string `json:"zoneName"`
	ZoneState struct {
		Available bool `json:"available"`
	} `json:"zoneState"`
}

type ListAvailabilityZonesOutput struct {
	AvailabilityZoneInfo []AvailabilityZone `json:"availabilityZoneInfo"`
}
