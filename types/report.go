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
	BufferSizeMax            int64  `json:"buffer_size_max"`

	PoolAddress string `json:"pool_address"`

	ClientHostname string `json:"client_hostname,omitempty"`
	ClientHostIP   string `json:"client_host_ip,omitempty"` // filled by server
	InstanceID     string `json:"instance_id"`

	CreationTime     time.Time `json:"creation_time"`
	LastActivityTime time.Time `json:"last_activity_time,omitempty"`
	TerminationTime  time.Time `json:"termination_time,omitempty"` // may be empty
	Terminated       bool      `json:"terminated,omitempty"`
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

	FilePath     string `json:"file_path"`
	FileSize     int64  `json:"file_size"`
	FileOpenMode string `json:"file_open_mode"`

	TransferBlocks     []FileBlock `json:"transfer_blocks"`
	TransferSize       int64       `json:"transfer_size"`
	LargestBlockSize   int64       `json:"largest_block_size"`
	SmallestBlockSize  int64       `json:"smallest_block_size"`
	TransferBlockCount int64       `json:"transfer_block_count"`
	SequentialAccess   bool        `json:"sequential_access"`

	FileOpenTime  time.Time `json:"file_open_time"`
	FileCloseTime time.Time `json:"file_close_time"`
}
