# MongoDB Provider Documentation

## Tổng quan

MongoDB Provider là một package Go chuyên dụng cung cấp interface để tương tác với MongoDB database một cách type-safe và hiệu quả. Package này được thiết kế để tích hợp seamlessly với go.fork.vn dependency injection framework.

## Cài đặt

```bash
go get go.fork.vn/providers/mongodb@v0.1.0
```

## Import

```go
import "go.fork.vn/providers/mongodb"
```

## Kiến trúc

### Manager Interface

MongoDB Provider cung cấp `Manager` interface với các method chính:

```go
type Manager interface {
    GetClient(ctx context.Context) (*mongo.Client, error)
    GetDatabase(ctx context.Context, name string) (*mongo.Database, error)
    GetCollection(ctx context.Context, database, collection string) (*mongo.Collection, error)
    Ping(ctx context.Context) error
    Disconnect(ctx context.Context) error
    IsConnected(ctx context.Context) bool
    GetConnectionString() string
}
```

### ServiceProvider Implementation

`ServiceProvider` implement `di.ServiceProvider` interface để đăng ký MongoDB services:

```go
type ServiceProvider struct {
    configManager config.Manager
}
```

## Configuration

### Cấu trúc Config

```go
type Config struct {
    URI                 string        `mapstructure:"uri" yaml:"uri"`
    Database            string        `mapstructure:"database" yaml:"database"`
    MaxPoolSize         *uint64       `mapstructure:"max_pool_size" yaml:"max_pool_size"`
    MinPoolSize         *uint64       `mapstructure:"min_pool_size" yaml:"min_pool_size"`
    MaxConnIdleTime     time.Duration `mapstructure:"max_conn_idle_time" yaml:"max_conn_idle_time"`
    ServerSelectionTimeout time.Duration `mapstructure:"server_selection_timeout" yaml:"server_selection_timeout"`
    ConnectTimeout      time.Duration `mapstructure:"connect_timeout" yaml:"connect_timeout"`
    SocketTimeout       time.Duration `mapstructure:"socket_timeout" yaml:"socket_timeout"`
    SSL                 SSLConfig     `mapstructure:"ssl" yaml:"ssl"`
    Auth                AuthConfig    `mapstructure:"auth" yaml:"auth"`
}

type SSLConfig struct {
    Enabled            bool   `mapstructure:"enabled" yaml:"enabled"`
    CAFile             string `mapstructure:"ca_file" yaml:"ca_file"`
    CertificateFile    string `mapstructure:"certificate_file" yaml:"certificate_file"`
    PrivateKeyFile     string `mapstructure:"private_key_file" yaml:"private_key_file"`
    InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify" yaml:"insecure_skip_verify"`
}

type AuthConfig struct {
    Username string `mapstructure:"username" yaml:"username"`
    Password string `mapstructure:"password" yaml:"password"`
    AuthDB   string `mapstructure:"auth_db" yaml:"auth_db"`
}
```

### Environment Variables

MongoDB Provider hỗ trợ đọc configuration từ environment variables:

```bash
# Connection
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=myapp

# Pool settings
MONGODB_MAX_POOL_SIZE=100
MONGODB_MIN_POOL_SIZE=10
MONGODB_MAX_CONN_IDLE_TIME=30s

# Timeouts
MONGODB_SERVER_SELECTION_TIMEOUT=30s
MONGODB_CONNECT_TIMEOUT=10s
MONGODB_SOCKET_TIMEOUT=30s

# SSL Configuration
MONGODB_SSL_ENABLED=true
MONGODB_SSL_CA_FILE=/path/to/ca.pem
MONGODB_SSL_CERTIFICATE_FILE=/path/to/cert.pem
MONGODB_SSL_PRIVATE_KEY_FILE=/path/to/key.pem
MONGODB_SSL_INSECURE_SKIP_VERIFY=false

# Authentication
MONGODB_AUTH_USERNAME=admin
MONGODB_AUTH_PASSWORD=password123
MONGODB_AUTH_DB=admin
```

## Core Features

### 1. Connection Management

**Auto-reconnection và Connection Pooling:**
```go
manager := mongodb.NewManager(config)

// Get client với connection pooling
client, err := manager.GetClient(ctx)
if err != nil {
    return fmt.Errorf("failed to get client: %w", err)
}

// Check connection health
if err := manager.Ping(ctx); err != nil {
    return fmt.Errorf("mongodb ping failed: %w", err)
}
```

### 2. Database Operations

**Database và Collection Access:**
```go
// Get database instance
db, err := manager.GetDatabase(ctx, "myapp")
if err != nil {
    return fmt.Errorf("failed to get database: %w", err)
}

// Get collection instance
collection, err := manager.GetCollection(ctx, "myapp", "users")
if err != nil {
    return fmt.Errorf("failed to get collection: %w", err)
}
```

### 3. SSL/TLS Support

**Secure Connections:**
```go
config := &mongodb.Config{
    URI: "mongodb://localhost:27017",
    SSL: mongodb.SSLConfig{
        Enabled:            true,
        CAFile:             "/path/to/ca.pem",
        CertificateFile:    "/path/to/cert.pem",
        PrivateKeyFile:     "/path/to/key.pem",
        InsecureSkipVerify: false,
    },
}
```

### 4. Authentication

**Database Authentication:**
```go
config := &mongodb.Config{
    URI: "mongodb://localhost:27017",
    Auth: mongodb.AuthConfig{
        Username: "admin",
        Password: "password123",
        AuthDB:   "admin",
    },
}
```

## Dependency Injection Integration

### 1. Service Provider Registration

```go
package main

import (
    "go.fork.vn/di"
    "go.fork.vn/providers/mongodb"
    "go.fork.vn/providers/config"
)

func main() {
    container := di.NewContainer()
    
    // Register config provider
    configProvider := config.NewServiceProvider()
    container.RegisterServiceProvider(configProvider)
    
    // Register mongodb provider
    mongoProvider := mongodb.NewServiceProvider()
    container.RegisterServiceProvider(mongoProvider)
    
    // Build container
    if err := container.Build(); err != nil {
        panic(err)
    }
}
```

### 2. Service Resolution

```go
// Resolve MongoDB manager
var mongoManager mongodb.Manager
if err := container.Resolve(&mongoManager); err != nil {
    return fmt.Errorf("failed to resolve mongodb manager: %w", err)
}

// Use manager
client, err := mongoManager.GetClient(ctx)
if err != nil {
    return fmt.Errorf("failed to get mongodb client: %w", err)
}
```

## Error Handling

### Connection Errors

```go
manager := mongodb.NewManager(config)

// Check if connected
if !manager.IsConnected(ctx) {
    return errors.New("mongodb not connected")
}

// Ping database
if err := manager.Ping(ctx); err != nil {
    return fmt.Errorf("mongodb ping failed: %w", err)
}
```

### Graceful Shutdown

```go
// Disconnect when shutting down
defer func() {
    if err := manager.Disconnect(ctx); err != nil {
        log.Printf("Error disconnecting from MongoDB: %v", err)
    }
}()
```

## Testing

### 1. Mock Support

Package cung cấp mock implementations để testing:

```go
import "go.fork.vn/providers/mongodb/mocks"

func TestMyService(t *testing.T) {
    mockManager := &mocks.Manager{}
    
    // Setup expectations
    mockManager.On("GetClient", mock.Anything).Return(client, nil)
    mockManager.On("Ping", mock.Anything).Return(nil)
    
    // Test your service
    service := NewMyService(mockManager)
    err := service.DoSomething()
    
    assert.NoError(t, err)
    mockManager.AssertExpectations(t)
}
```

### 2. Integration Testing

```go
func TestMongoDBIntegration(t *testing.T) {
    config := &mongodb.Config{
        URI:      "mongodb://localhost:27017",
        Database: "test_db",
    }
    
    manager := mongodb.NewManager(config)
    defer manager.Disconnect(context.Background())
    
    // Test connection
    err := manager.Ping(context.Background())
    assert.NoError(t, err)
    
    // Test database operations
    db, err := manager.GetDatabase(context.Background(), "test_db")
    assert.NoError(t, err)
    assert.NotNil(t, db)
}
```

## Performance Tuning

### Connection Pool Settings

```go
config := &mongodb.Config{
    URI:                    "mongodb://localhost:27017",
    MaxPoolSize:            &[]uint64{100}[0],  // Max 100 connections
    MinPoolSize:            &[]uint64{10}[0],   // Min 10 connections
    MaxConnIdleTime:        30 * time.Second,   // Idle timeout
    ServerSelectionTimeout: 30 * time.Second,   // Selection timeout
    ConnectTimeout:         10 * time.Second,   // Connection timeout
    SocketTimeout:          30 * time.Second,   // Socket timeout
}
```

### Read Preferences

```go
// Configure read preferences for better performance
client, err := manager.GetClient(ctx)
if err != nil {
    return err
}

// Setup read preference
readPref := readpref.Secondary()
db := client.Database("myapp", options.Database().SetReadPreference(readPref))
```

## Migration Guide

### From v0.0.x to v0.1.0

**1. Update import statements:**
```go
// Old
import "github.com/go-fork/providers/mongodb"

// New  
import "go.fork.vn/providers/mongodb"
```

**2. Update go.mod:**
```bash
go get go.fork.vn/providers/mongodb@v0.1.0
```

**3. No API changes required** - tất cả existing code sẽ hoạt động bình thường.

## Best Practices

### 1. Connection Management
- Sử dụng single Manager instance trong toàn bộ application
- Implement proper graceful shutdown
- Monitor connection pool metrics

### 2. Error Handling
- Always check connection status trước khi operations
- Implement retry logic cho transient errors
- Log connection events để monitoring

### 3. Security
- Luôn sử dụng SSL/TLS trong production
- Store credentials trong environment variables
- Implement proper authentication

### 4. Performance
- Configure connection pool settings phù hợp với workload
- Sử dụng appropriate read preferences
- Monitor query performance

## Troubleshooting

### Common Issues

**1. Connection timeout:**
```bash
# Increase timeout settings
MONGODB_CONNECT_TIMEOUT=30s
MONGODB_SERVER_SELECTION_TIMEOUT=60s
```

**2. SSL certificate errors:**
```bash
# For development only
MONGODB_SSL_INSECURE_SKIP_VERIFY=true
```

**3. Authentication failures:**
```bash
# Check credentials and auth database
MONGODB_AUTH_USERNAME=correct_username
MONGODB_AUTH_PASSWORD=correct_password
MONGODB_AUTH_DB=admin
```

### Debug Mode

Enable debug logging để troubleshoot:

```go
// Enable MongoDB driver logging
mongodb.SetDebugMode(true)
```

## Contribution

Để contribute vào project:

1. Fork repository
2. Tạo feature branch
3. Implement changes với tests
4. Submit pull request

## License

MongoDB Provider được phân phối theo MIT License. Xem file LICENSE để biết thêm chi tiết.
