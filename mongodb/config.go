package mongodb

import (
	"time"
)

// Config holds MongoDB configuration.
type Config struct {
	// Connection URI for MongoDB
	// Format: mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
	URI string `yaml:"uri" mapstructure:"uri"`

	// Default database name
	Database string `yaml:"database" mapstructure:"database"`

	// Application name to identify the connection in MongoDB logs
	AppName string `yaml:"app_name" mapstructure:"app_name"`

	// Connection pool settings
	MaxPoolSize     uint64 `yaml:"max_pool_size" mapstructure:"max_pool_size"`           // Maximum number of connections in the connection pool
	MinPoolSize     uint64 `yaml:"min_pool_size" mapstructure:"min_pool_size"`           // Minimum number of connections in the connection pool
	MaxConnecting   uint64 `yaml:"max_connecting" mapstructure:"max_connecting"`         // Maximum number of connections being established concurrently
	MaxConnIdleTime uint64 `yaml:"max_conn_idle_time" mapstructure:"max_conn_idle_time"` // Maximum time (ms) a connection can remain idle

	// Timeout settings (all in milliseconds)
	ConnectTimeout         uint64 `yaml:"connect_timeout" mapstructure:"connect_timeout"`                   // Connection timeout
	ServerSelectionTimeout uint64 `yaml:"server_selection_timeout" mapstructure:"server_selection_timeout"` // Server selection timeout
	SocketTimeout          uint64 `yaml:"socket_timeout" mapstructure:"socket_timeout"`                     // Socket timeout (0 = no timeout)
	HeartbeatInterval      uint64 `yaml:"heartbeat_interval" mapstructure:"heartbeat_interval"`             // Heartbeat interval for monitoring server health
	LocalThreshold         uint64 `yaml:"local_threshold" mapstructure:"local_threshold"`                   // Local threshold for server selection
	Timeout                uint64 `yaml:"timeout" mapstructure:"timeout"`                                   // General operation timeout

	// TLS/SSL configuration
	TLS TLSConfig `yaml:"tls" mapstructure:"tls"`

	// Authentication configuration
	Auth AuthConfig `yaml:"auth" mapstructure:"auth"`

	// Read preference configuration
	// Modes: primary, primaryPreferred, secondary, secondaryPreferred, nearest
	ReadPreference ReadPreferenceConfig `yaml:"read_preference" mapstructure:"read_preference"`

	// Read concern configuration
	// Levels: local, available, majority, linearizable, snapshot
	ReadConcern ReadConcernConfig `yaml:"read_concern" mapstructure:"read_concern"`

	// Write concern configuration
	WriteConcern WriteConcernConfig `yaml:"write_concern" mapstructure:"write_concern"`

	// Retry configuration
	RetryWrites bool `yaml:"retry_writes" mapstructure:"retry_writes"` // Enable retryable writes
	RetryReads  bool `yaml:"retry_reads" mapstructure:"retry_reads"`   // Enable retryable reads

	// Compression configuration
	Compressors []string `yaml:"compressors" mapstructure:"compressors"` // Compression algorithms: snappy, zlib, zstd
	ZlibLevel   int      `yaml:"zlib_level" mapstructure:"zlib_level"`   // Compression level for zlib (1-9)
	ZstdLevel   int      `yaml:"zstd_level" mapstructure:"zstd_level"`   // Compression level for zstd (1-22)

	// Replica set configuration
	ReplicaSet string `yaml:"replica_set" mapstructure:"replica_set"` // Replica set name
	Direct     bool   `yaml:"direct" mapstructure:"direct"`           // Connect directly to a specific server

	// Load balancer configuration
	LoadBalanced bool `yaml:"load_balanced" mapstructure:"load_balanced"` // Enable load balanced mode

	// SRV configuration for DNS-based discovery
	SRV SRVConfig `yaml:"srv" mapstructure:"srv"`

	// Server API configuration
	ServerAPI ServerAPIConfig `yaml:"server_api" mapstructure:"server_api"`

	// Monitoring and logging
	ServerMonitoringMode     string `yaml:"server_monitoring_mode" mapstructure:"server_monitoring_mode"`           // Server monitoring mode: auto, stream, poll
	DisableOCSPEndpointCheck bool   `yaml:"disable_ocsp_endpoint_check" mapstructure:"disable_ocsp_endpoint_check"` // Disable OCSP endpoint check for TLS

	// BSON configuration
	BSON BSONConfig `yaml:"bson" mapstructure:"bson"`

	// Auto-encryption configuration (Enterprise/Atlas only)
	AutoEncryption AutoEncryptionConfig `yaml:"auto_encryption" mapstructure:"auto_encryption"`
}

// TLSConfig holds TLS/SSL configuration.
type TLSConfig struct {
	Enabled            bool   `yaml:"enabled" mapstructure:"enabled"`                           // Enable/disable TLS
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify" mapstructure:"insecure_skip_verify"` // Skip certificate verification (for testing only)
	CAFile             string `yaml:"ca_file" mapstructure:"ca_file"`                           // Path to CA certificate file
	CertFile           string `yaml:"cert_file" mapstructure:"cert_file"`                       // Path to client certificate file
	KeyFile            string `yaml:"key_file" mapstructure:"key_file"`                         // Path to client private key file
}

// AuthConfig holds authentication configuration.
type AuthConfig struct {
	Username                string            `yaml:"username" mapstructure:"username"`                                   // Username for authentication
	Password                string            `yaml:"password" mapstructure:"password"`                                   // Password for authentication
	AuthSource              string            `yaml:"auth_source" mapstructure:"auth_source"`                             // Authentication database
	AuthMechanism           string            `yaml:"auth_mechanism" mapstructure:"auth_mechanism"`                       // Authentication mechanism
	AuthMechanismProperties map[string]string `yaml:"auth_mechanism_properties" mapstructure:"auth_mechanism_properties"` // Additional authentication mechanism properties
}

// ReadPreferenceConfig holds read preference configuration.
type ReadPreferenceConfig struct {
	Mode         string              `yaml:"mode" mapstructure:"mode"`                   // Read preference mode
	TagSets      []map[string]string `yaml:"tag_sets" mapstructure:"tag_sets"`           // Tag sets for read preference
	MaxStaleness int                 `yaml:"max_staleness" mapstructure:"max_staleness"` // Maximum staleness in seconds
	HedgeEnabled bool                `yaml:"hedge_enabled" mapstructure:"hedge_enabled"` // Enable hedge reads for sharded clusters
}

// ReadConcernConfig holds read concern configuration.
type ReadConcernConfig struct {
	Level string `yaml:"level" mapstructure:"level"`
}

// WriteConcernConfig holds write concern configuration.
type WriteConcernConfig struct {
	W        interface{} `yaml:"w" mapstructure:"w"`                 // Can be int, string, or "majority"
	Journal  bool        `yaml:"journal" mapstructure:"journal"`     // Journal acknowledgment
	WTimeout int         `yaml:"w_timeout" mapstructure:"w_timeout"` // Write timeout in milliseconds
}

// SRVConfig holds SRV configuration for DNS-based discovery.
type SRVConfig struct {
	MaxHosts    int    `yaml:"max_hosts" mapstructure:"max_hosts"`       // Maximum number of hosts to connect to (0 = no limit)
	ServiceName string `yaml:"service_name" mapstructure:"service_name"` // SRV service name
}

// ServerAPIConfig holds Server API configuration.
type ServerAPIConfig struct {
	Version           string `yaml:"version" mapstructure:"version"`                       // API version ("1" is currently supported)
	Strict            bool   `yaml:"strict" mapstructure:"strict"`                         // Strict API version mode
	DeprecationErrors bool   `yaml:"deprecation_errors" mapstructure:"deprecation_errors"` // Return errors for deprecated features
}

// BSONConfig holds BSON configuration.
type BSONConfig struct {
	UseJSONStructTags     bool `yaml:"use_json_struct_tags" mapstructure:"use_json_struct_tags"`       // Use JSON struct tags for BSON marshaling
	ErrorOnInlineMap      bool `yaml:"error_on_inline_map" mapstructure:"error_on_inline_map"`         // Error on inline map fields
	AllowTruncatingFloats bool `yaml:"allow_truncating_floats" mapstructure:"allow_truncating_floats"` // Allow truncating floats when converting to integers
}

// AutoEncryptionConfig holds auto-encryption configuration (Enterprise/Atlas only).
type AutoEncryptionConfig struct {
	Enabled              bool                   `yaml:"enabled" mapstructure:"enabled"`                               // Enable auto-encryption
	KeyVaultNamespace    string                 `yaml:"key_vault_namespace" mapstructure:"key_vault_namespace"`       // Key vault namespace (database.collection)
	KMSProviders         map[string]interface{} `yaml:"kms_providers" mapstructure:"kms_providers"`                   // KMS providers configuration
	SchemaMap            map[string]interface{} `yaml:"schema_map" mapstructure:"schema_map"`                         // Schema map for automatic encryption
	BypassAutoEncryption bool                   `yaml:"bypass_auto_encryption" mapstructure:"bypass_auto_encryption"` // Bypass automatic encryption
	ExtraOptions         map[string]interface{} `yaml:"extra_options" mapstructure:"extra_options"`                   // Extra options for auto-encryption
}

// DefaultConfig returns default MongoDB configuration.
func DefaultConfig() *Config {
	return &Config{
		URI:                    "mongodb://localhost:27017",
		Database:               "myapp",
		AppName:                "golang-app",
		MaxPoolSize:            DefaultMaxPoolSize,
		MinPoolSize:            DefaultMinPoolSize,
		MaxConnecting:          DefaultMaxConnecting,
		MaxConnIdleTime:        DefaultConnIdleTime,
		ConnectTimeout:         DefaultConnectTimeout,
		ServerSelectionTimeout: DefaultServerTimeout,
		SocketTimeout:          0,
		HeartbeatInterval:      DefaultHeartbeat,
		LocalThreshold:         DefaultLocalThreshold,
		Timeout:                DefaultTimeout,
		TLS: TLSConfig{
			Enabled:            false,
			InsecureSkipVerify: false,
		},
		Auth: AuthConfig{
			AuthSource:              "admin",
			AuthMechanism:           "SCRAM-SHA-256",
			AuthMechanismProperties: make(map[string]string),
		},
		ReadPreference: ReadPreferenceConfig{
			Mode:         "primary",
			TagSets:      []map[string]string{},
			MaxStaleness: DefaultMaxStaleness,
			HedgeEnabled: false,
		},
		ReadConcern: ReadConcernConfig{
			Level: "majority",
		},
		WriteConcern: WriteConcernConfig{
			W:        "majority",
			Journal:  true,
			WTimeout: DefaultWTimeout,
		},
		RetryWrites:  true,
		RetryReads:   true,
		Compressors:  []string{},
		ZlibLevel:    DefaultZlibLevel,
		ZstdLevel:    DefaultZstdLevel,
		Direct:       false,
		LoadBalanced: false,
		SRV: SRVConfig{
			MaxHosts:    0,
			ServiceName: "mongodb",
		},
		ServerAPI: ServerAPIConfig{
			Version:           "1",
			Strict:            false,
			DeprecationErrors: false,
		},
		ServerMonitoringMode:     "auto",
		DisableOCSPEndpointCheck: false,
		BSON: BSONConfig{
			UseJSONStructTags:     false,
			ErrorOnInlineMap:      false,
			AllowTruncatingFloats: false,
		},
		AutoEncryption: AutoEncryptionConfig{
			Enabled:              false,
			KMSProviders:         make(map[string]interface{}),
			SchemaMap:            make(map[string]interface{}),
			BypassAutoEncryption: false,
			ExtraOptions:         make(map[string]interface{}),
		},
	}
}

// GetConnectTimeout returns connection timeout as time.Duration.
func (c *Config) GetConnectTimeout() time.Duration {
	return time.Duration(c.ConnectTimeout) * time.Millisecond
}

// GetServerSelectionTimeout returns server selection timeout as time.Duration.
func (c *Config) GetServerSelectionTimeout() time.Duration {
	return time.Duration(c.ServerSelectionTimeout) * time.Millisecond
}

// GetSocketTimeout returns socket timeout as time.Duration.
func (c *Config) GetSocketTimeout() time.Duration {
	if c.SocketTimeout == 0 {
		return 0 // No timeout
	}
	return time.Duration(c.SocketTimeout) * time.Millisecond
}

// GetHeartbeatInterval returns heartbeat interval as time.Duration.
func (c *Config) GetHeartbeatInterval() time.Duration {
	return time.Duration(c.HeartbeatInterval) * time.Millisecond
}

// GetLocalThreshold returns local threshold as time.Duration.
func (c *Config) GetLocalThreshold() time.Duration {
	return time.Duration(c.LocalThreshold) * time.Millisecond
}

// GetTimeout returns general timeout as time.Duration.
func (c *Config) GetTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Millisecond
}

// GetMaxConnIdleTime returns max connection idle time as time.Duration.
func (c *Config) GetMaxConnIdleTime() time.Duration {
	return time.Duration(c.MaxConnIdleTime) * time.Millisecond
}

// GetWTimeout returns write timeout as time.Duration.
func (c *Config) GetWTimeout() time.Duration {
	return time.Duration(c.WriteConcern.WTimeout) * time.Millisecond
}
