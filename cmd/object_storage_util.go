package cmd

import (
	"fmt"
	"strings"
)

// OBSPath represents a parsed path which can be local or remote (OBS)
type OBSPath struct {
	Container string
	Object    string
	IsRemote  bool
	RawPath   string
}

// parseOBSPath parses a string into an OBSPath struct
// Remote paths must start with "obs://"
func parseOBSPath(path string) (*OBSPath, error) {
	if strings.HasPrefix(path, "obs://") {
		trimmed := strings.TrimPrefix(path, "obs://")
		// Split into Container and Object
		// usage: obs://container/object
		parts := strings.SplitN(trimmed, "/", 2)

		if len(parts) == 0 || parts[0] == "" {
			return nil, fmt.Errorf("invalid obs path: %s (missing container)", path)
		}

		container := parts[0]
		object := ""
		if len(parts) > 1 {
			object = parts[1]
		}

		return &OBSPath{
			Container: container,
			Object:    object,
			IsRemote:  true,
			RawPath:   path,
		}, nil
	}

	// Local path
	return &OBSPath{
		IsRemote: false,
		RawPath:  path,
	}, nil
}

// String returns the string representation
func (p *OBSPath) String() string {
	return p.RawPath
}
