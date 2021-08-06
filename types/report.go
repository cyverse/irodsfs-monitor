package types

import "time"

// ReportInstance is a struct used to report an instance creation
type ReportInstance struct {
	Host                     string `json:"host"`
	Port                     int    `json:"port"`
	Zone                     string `json:"zone"`
	ClientUser               string `json:"client_user"`
	ProxyUser                string `json:"proxy_user"`
	AuthScheme               string `json:"auth_scheme"`
	ReadAheadMax             int    `json:"read_ahead_max"`
	OperationTimeout         string `json:"operation_timeout"`
	ConnectionIdleTimeout    string `json:"connection_idle_timeout"`
	ConnectionMax            int    `json:"connection_max"`
	MetadataCacheTimeout     string `json:"metadata_cache_timeout"`
	MetadataCacheCleanupTime string `json:"metadata_cache_cleanup_time"`
	FileBufferSizeMax        int64  `json:"file_buffer_size_max"`

	ClientHostName string `json:"client_hostname,omitempty"`
	InstanceID     string `json:"instance_id"`

	CreationTime time.Time `json:"creation_time"`
}

// FileBlock is an internal struct used in other structs
type FileBlock struct {
	Offset     int64     `json:"offset"`
	Length     int64     `json:"length"`
	AccessTime time.Time `json:"access_time"`
}

// ReportFileTransfer is a struct used to report file transfer information
type ReportFileTransfer struct {
	InstanceID string `json:"instance_id"`

	FilePath string `json:"file_path"`
	FileSize int64  `json:"file_size"`

	TransferOrder []FileBlock `json:"transfer_order"`
	TransferSize  int64       `json:"transfer_size"`

	FileOpenTime  time.Time `json:"file_open_time"`
	FileCloseTime time.Time `json:"file_close_time"`
}
