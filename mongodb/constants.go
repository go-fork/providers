package mongodb

import "time"

// Timeout constants for various operations
const (
	// DefaultPingTimeout is the default timeout for ping operations
	DefaultPingTimeout = 10 * time.Second

	// DefaultDisconnectTimeout is the default timeout for disconnect operations
	DefaultDisconnectTimeout = 10 * time.Second

	// DefaultHealthCheckTimeout is the default timeout for health check operations
	DefaultHealthCheckTimeout = 5 * time.Second

	// DefaultMonitoringTimeout is the default timeout for monitoring operations
	DefaultMonitoringTimeout = 10 * time.Second

	// DefaultStreamTimeout is the default timeout for stream operations
	DefaultStreamTimeout = 30 * time.Second

	// DefaultIndexNameConstant is the default constant for index field names
	DefaultIndexFieldName = "field_1"

	// Pool configuration constants
	DefaultMaxPoolSize    = 100
	DefaultMinPoolSize    = 5
	DefaultMaxConnecting  = 10
	DefaultConnIdleTime   = 600000
	DefaultConnectTimeout = 30000
	DefaultServerTimeout  = 30000
	DefaultHeartbeat      = 10000
	DefaultLocalThreshold = 15000
	DefaultTimeout        = 30000
	DefaultMaxStaleness   = 90
	DefaultWTimeout       = 30000
	DefaultZlibLevel      = 6
	DefaultZstdLevel      = 6
)
