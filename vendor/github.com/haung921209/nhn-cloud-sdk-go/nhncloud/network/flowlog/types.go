package flowlog

import "time"

// Logger represents a flow log logger
type Logger struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description,omitempty"`
	ResourceType     string    `json:"resource_type"` // VPC, SUBNET, PORT
	ResourceID       string    `json:"resource_id"`
	FilterType       string    `json:"filter_type"`       // ALL, ACCEPT, DROP
	ConnectionAction string    `json:"connection_action"` // enable, disable
	StorageType      string    `json:"storage_type"`      // OBS
	StorageURL       string    `json:"storage_url"`
	LogFormat        string    `json:"log_format"`       // CSV, PARQUET
	CompressionType  string    `json:"compression_type"` // RAW, GZIP
	PartitionPeriod  string    `json:"partition_period"` // HOUR, DAY
	AdminStateUp     bool      `json:"admin_state_up"`
	State            string    `json:"state"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

// LoggingPort represents a logging port
type LoggingPort struct {
	ID        string    `json:"id"`
	LoggerID  string    `json:"logger_id"`
	PortID    string    `json:"port_id"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}

// ListLoggersOutput represents the response from listing loggers
type ListLoggersOutput struct {
	Loggers []Logger `json:"flowlog_loggers"`
}

// GetLoggerOutput represents the response containing a single logger
type GetLoggerOutput struct {
	Logger *Logger `json:"flowlog_logger"`
}

// CreateLoggerInput contains the logger creation data
type CreateLoggerInput struct {
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	ResourceType     string `json:"resource_type"`
	ResourceID       string `json:"resource_id"`
	FilterType       string `json:"filter_type"`
	ConnectionAction string `json:"connection_action,omitempty"`
	StorageType      string `json:"storage_type"`
	StorageURL       string `json:"storage_url"`
	LogFormat        string `json:"log_format"`
	CompressionType  string `json:"compression_type,omitempty"`
	PartitionPeriod  string `json:"partition_period,omitempty"`
	AdminStateUp     bool   `json:"admin_state_up,omitempty"`
}

// CreateLoggerRequest wraps the create input
type CreateLoggerRequest struct {
	Logger CreateLoggerInput `json:"flowlog_logger"`
}

// UpdateLoggerInput contains the logger update data
type UpdateLoggerInput struct {
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	ConnectionAction string `json:"connection_action,omitempty"`
	AdminStateUp     *bool  `json:"admin_state_up,omitempty"`
}

// UpdateLoggerRequest wraps the update input
type UpdateLoggerRequest struct {
	Logger UpdateLoggerInput `json:"flowlog_logger"`
}

// ListLoggingPortsOutput represents the response from listing logging ports
type ListLoggingPortsOutput struct {
	LoggingPorts []LoggingPort `json:"flowlog_logging_ports"`
}

// GetLoggingPortOutput represents the response containing a single logging port
type GetLoggingPortOutput struct {
	LoggingPort *LoggingPort `json:"flowlog_logging_port"`
}
