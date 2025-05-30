# Usage Guide - Config Package

## Cài đặt

```bash
go get go.fork.vn/providers/config@v0.1.0
```

## Import

```go
import "go.fork.vn/providers/config"
```

## Quick Start

### 1. Sử dụng cơ bản

```go
package main

import (
    "fmt"
    "log"
    
    "go.fork.vn/providers/config"
)

func main() {
    // Tạo config manager
    cfg := config.NewConfig()
    
    // Cấu hình file
    cfg.SetConfigName("app")
    cfg.AddConfigPath("./configs")
    cfg.SetConfigType("yaml")
    
    // Đặt giá trị mặc định
    cfg.SetDefault("server.port", 8080)
    cfg.SetDefault("app.debug", false)
    
    // Đọc cấu hình
    if err := cfg.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            fmt.Println("Không tìm thấy file config, sử dụng defaults")
        } else {
            log.Fatal("Lỗi đọc config:", err)
        }
    }
    
    // Sử dụng cấu hình
    if appName, exists := cfg.GetString("app.name"); exists {
        fmt.Printf("App: %s\n", appName)
    }
    
    if port, exists := cfg.GetInt("server.port"); exists {
        fmt.Printf("Port: %d\n", port)
    }
}
```

### 2. Với Dependency Injection

```go
package main

import (
    "fmt"
    "log"
    
    "go.fork.vn/di"
    "go.fork.vn/providers/config"
)

func main() {
    // Tạo DI container
    container := di.NewContainer()
    
    // Đăng ký config provider
    provider := &config.ServiceProvider{}
    if err := container.RegisterProvider(provider); err != nil {
        log.Fatal("Không thể đăng ký provider:", err)
    }
    
    // Resolve config manager
    cfg, err := di.Resolve[config.Manager](container)
    if err != nil {
        log.Fatal("Không thể resolve config:", err)
    }
    
    // Sử dụng như bình thường
    cfg.SetConfigName("app")
    cfg.AddConfigPath("./configs")
    cfg.ReadInConfig()
    
    if name, exists := cfg.GetString("app.name"); exists {
        fmt.Println("Application:", name)
    }
}
```

## Cấu hình File

### Cấu trúc thư mục

```
project/
├── configs/
│   ├── app.yaml
│   ├── app.production.yaml
│   └── app.development.yaml
├── main.go
└── go.mod
```

### File cấu hình mẫu (app.yaml)

```yaml
app:
  name: "My Application"
  version: "1.0.0"
  debug: false

server:
  host: "localhost"
  port: 8080
  timeout: 30s
  ssl:
    enabled: false
    cert_file: ""
    key_file: ""

database:
  host: "localhost"
  port: 5432
  name: "myapp"
  username: "postgres"
  password: "secret"
  max_connections: 10
  timeout: 5s
  replicas:
    - "replica1.example.com"
    - "replica2.example.com"

cache:
  enabled: true
  type: "redis"
  host: "localhost"
  port: 6379
  db: 0
  ttl: 3600

logging:
  level: "info"
  format: "json"
  output: "stdout"
  
features:
  user_registration: true
  email_verification: true
  social_login: false
```

### Multiple Environment Configs

```go
func setupConfig() config.Manager {
    cfg := config.NewConfig()
    
    // Determine environment
    env := os.Getenv("GO_ENV")
    if env == "" {
        env = "development"
    }
    
    // Base config first
    cfg.SetConfigName("app")
    cfg.AddConfigPath("./configs")
    
    // Environment-specific config
    if env != "development" {
        cfg.SetConfigName(fmt.Sprintf("app.%s", env))
        cfg.AddConfigPath(fmt.Sprintf("./configs/%s", env))
    }
    
    cfg.SetConfigType("yaml")
    
    // Set environment prefix
    cfg.SetEnvPrefix("MYAPP")
    cfg.AutomaticEnv()
    
    if err := cfg.ReadInConfig(); err != nil {
        log.Printf("Config error: %v", err)
    }
    
    return cfg
}
```

## Type-Safe Value Retrieval

### Basic Types

```go
// String values
if dbHost, exists := cfg.GetString("database.host"); exists {
    fmt.Printf("Database host: %s\n", dbHost)
} else {
    log.Fatal("database.host is required")
}

// Integer values
if port, exists := cfg.GetInt("server.port"); exists {
    server.Listen(port)
} else {
    server.Listen(8080) // default
}

// Boolean values
if debug, exists := cfg.GetBool("app.debug"); exists && debug {
    enableDebugMode()
}

// Float values
if ratio, exists := cfg.GetFloat("cache.hit_ratio"); exists {
    metrics.SetCacheRatio(ratio)
}

// Duration values
if timeout, exists := cfg.GetDuration("database.timeout"); exists {
    db.SetTimeout(timeout)
}

// Time values
if startTime, exists := cfg.GetTime("scheduler.start_time"); exists {
    scheduler.SetStartTime(startTime)
}
```

### Arrays and Slices

```go
// String slices
if hosts, exists := cfg.GetStringSlice("database.replicas"); exists {
    for _, host := range hosts {
        pool.AddReplica(host)
    }
}

// Integer slices
if ports, exists := cfg.GetIntSlice("server.ports"); exists {
    for _, port := range ports {
        listeners = append(listeners, createListener(port))
    }
}

// Generic slices
if items, exists := cfg.GetSlice("inventory.items"); exists {
    for _, item := range items {
        processItem(item)
    }
}
```

### Maps

```go
// Generic maps
if dbConfig, exists := cfg.GetMap("database"); exists {
    if host, ok := dbConfig["host"].(string); ok {
        fmt.Printf("DB Host: %s\n", host)
    }
}

// String maps
if settings, exists := cfg.GetStringMap("app.settings"); exists {
    for key, value := range settings {
        fmt.Printf("%s: %v\n", key, value)
    }
}

// String-to-String maps
if headers, exists := cfg.GetStringMapString("http.headers"); exists {
    for name, value := range headers {
        req.Header.Set(name, value)
    }
}

// String-to-StringSlice maps
if params, exists := cfg.GetStringMapStringSlice("routes.params"); exists {
    for route, paramList := range params {
        router.RegisterParams(route, paramList)
    }
}
```

## Struct Unmarshaling

### Định nghĩa Structs

```go
type AppConfig struct {
    Name    string `mapstructure:"name"`
    Version string `mapstructure:"version"`
    Debug   bool   `mapstructure:"debug"`
}

type ServerConfig struct {
    Host    string        `mapstructure:"host"`
    Port    int           `mapstructure:"port"`
    Timeout time.Duration `mapstructure:"timeout"`
    SSL     struct {
        Enabled  bool   `mapstructure:"enabled"`
        CertFile string `mapstructure:"cert_file"`
        KeyFile  string `mapstructure:"key_file"`
    } `mapstructure:"ssl"`
}

type DatabaseConfig struct {
    Host           string        `mapstructure:"host"`
    Port           int           `mapstructure:"port"`
    Name           string        `mapstructure:"name"`
    Username       string        `mapstructure:"username"`
    Password       string        `mapstructure:"password"`
    MaxConnections int           `mapstructure:"max_connections"`
    Timeout        time.Duration `mapstructure:"timeout"`
    Replicas       []string      `mapstructure:"replicas"`
}

type Config struct {
    App      AppConfig      `mapstructure:"app"`
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
}
```

### Unmarshal Operations

```go
// Unmarshal toàn bộ config
var fullConfig Config
if err := cfg.Unmarshal(&fullConfig); err != nil {
    log.Fatal("Không thể unmarshal config:", err)
}

fmt.Printf("App: %+v\n", fullConfig.App)
fmt.Printf("Server: %+v\n", fullConfig.Server)

// Unmarshal từng phần
var dbConfig DatabaseConfig
if err := cfg.UnmarshalKey("database", &dbConfig); err != nil {
    log.Fatal("Không thể unmarshal database config:", err)
}

// Sử dụng config
connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
    dbConfig.Host, dbConfig.Port, dbConfig.Name, 
    dbConfig.Username, dbConfig.Password)
```

## Environment Variables

### Automatic Environment Support

```go
cfg := config.NewConfig() // AutomaticEnv() đã được gọi sẵn

// Set prefix cho env variables
cfg.SetEnvPrefix("MYAPP")

// Giờ các env variables sau sẽ override config:
// MYAPP_DATABASE_HOST -> database.host
// MYAPP_SERVER_PORT -> server.port
// MYAPP_APP_DEBUG -> app.debug
```

### Explicit Binding

```go
// Bind specific keys
cfg.BindEnv("database.password", "DB_PASSWORD")
cfg.BindEnv("api.key", "API_SECRET_KEY")

// Bind với tên giống nhau
cfg.BindEnv("debug") // Tìm env var DEBUG
```

### Environment Variable Examples

```bash
# Set environment variables
export MYAPP_DATABASE_HOST=production-db.example.com
export MYAPP_SERVER_PORT=9000
export MYAPP_APP_DEBUG=true
export DB_PASSWORD=super-secret-password

# Run application
go run main.go
```

## Configuration Watching

### Basic Watching

```go
cfg := config.NewConfig()
cfg.SetConfigName("app")
cfg.AddConfigPath("./configs")
cfg.ReadInConfig()

// Setup watcher
cfg.OnConfigChange(func(e fsnotify.Event) {
    fmt.Printf("Config file changed: %s\n", e.Name)
    
    // Reload application components
    reloadDatabase(cfg)
    reloadCache(cfg)
    reloadLogger(cfg)
})

cfg.WatchConfig()

// Keep application running
select {}
```

### Advanced Watching với Graceful Reload

```go
type Application struct {
    cfg    config.Manager
    server *http.Server
    mu     sync.RWMutex
}

func (app *Application) setupConfigWatcher() {
    app.cfg.OnConfigChange(func(e fsnotify.Event) {
        log.Printf("Config changed: %s, reloading...", e.Name)
        
        app.mu.Lock()
        defer app.mu.Unlock()
        
        // Reload different components based on what changed
        if strings.Contains(e.Name, "app.yaml") {
            app.reloadServer()
            app.reloadDatabase()
        }
    })
    
    app.cfg.WatchConfig()
}

func (app *Application) reloadServer() {
    if port, exists := app.cfg.GetInt("server.port"); exists {
        if app.server.Addr != fmt.Sprintf(":%d", port) {
            // Port changed, restart server
            app.restartServer(port)
        }
    }
}
```

## Configuration Writing

### Basic Writing

```go
// Set values
cfg.Set("app.version", "2.0.0")
cfg.Set("features.new_feature", true)

// Write to current config file
if err := cfg.WriteConfig(); err != nil {
    log.Fatal("Cannot write config:", err)
}

// Write to specific file
if err := cfg.WriteConfigAs("config.new.yaml"); err != nil {
    log.Fatal("Cannot write config file:", err)
}

// Safe write (only if file doesn't exist)
if err := cfg.SafeWriteConfigAs("config.backup.yaml"); err != nil {
    log.Printf("Backup config exists or error: %v", err)
}
```

### Runtime Configuration Updates

```go
func updateFeatureFlag(cfg config.Manager, feature string, enabled bool) {
    key := fmt.Sprintf("features.%s", feature)
    
    if err := cfg.Set(key, enabled); err != nil {
        log.Printf("Cannot set %s: %v", key, err)
        return
    }
    
    // Persist change
    if err := cfg.WriteConfig(); err != nil {
        log.Printf("Cannot persist config: %v", err)
    }
    
    log.Printf("Feature %s set to %v", feature, enabled)
}
```

## Configuration Merging

### Merge from Reader

```go
// Additional config as JSON
additionalConfig := `{
  "features": {
    "beta_feature": true,
    "experimental": false
  },
  "cache": {
    "ttl": 7200
  }
}`

// Merge with current config
reader := strings.NewReader(additionalConfig)
if err := cfg.MergeConfig(reader); err != nil {
    log.Fatal("Cannot merge config:", err)
}
```

### Merge from File

```go
// Read base config
cfg.SetConfigFile("base.yaml")
cfg.ReadInConfig()

// Merge additional configs
cfg.SetConfigFile("features.yaml")
if err := cfg.MergeInConfig(); err != nil {
    log.Printf("Cannot merge features config: %v", err)
}

cfg.SetConfigFile("local.yaml")
if err := cfg.MergeInConfig(); err != nil {
    log.Printf("Cannot merge local config: %v", err)
}
```

## Validation Patterns

### Basic Validation

```go
func validateConfig(cfg config.Manager) error {
    required := []string{
        "database.host",
        "database.name",
        "database.username",
        "server.port",
    }
    
    for _, key := range required {
        if !cfg.Has(key) {
            return fmt.Errorf("required config key missing: %s", key)
        }
    }
    
    // Validate port range
    if port, exists := cfg.GetInt("server.port"); exists {
        if port < 1 || port > 65535 {
            return fmt.Errorf("invalid port: %d", port)
        }
    }
    
    return nil
}
```

### Advanced Validation

```go
type ConfigValidator struct {
    cfg config.Manager
}

func (v *ConfigValidator) ValidateDatabase() error {
    var dbConfig DatabaseConfig
    if err := v.cfg.UnmarshalKey("database", &dbConfig); err != nil {
        return fmt.Errorf("invalid database config: %w", err)
    }
    
    // Validate host
    if dbConfig.Host == "" {
        return errors.New("database host is required")
    }
    
    // Validate port
    if dbConfig.Port < 1 || dbConfig.Port > 65535 {
        return fmt.Errorf("invalid database port: %d", dbConfig.Port)
    }
    
    // Test connection
    if err := testDatabaseConnection(dbConfig); err != nil {
        return fmt.Errorf("cannot connect to database: %w", err)
    }
    
    return nil
}

func (v *ConfigValidator) ValidateAll() error {
    validators := []func() error{
        v.ValidateDatabase,
        v.ValidateServer,
        v.ValidateCache,
    }
    
    for _, validate := range validators {
        if err := validate(); err != nil {
            return err
        }
    }
    
    return nil
}
```

## Best Practices

### 1. Initialization Pattern

```go
func initConfig() config.Manager {
    cfg := config.NewConfig()
    
    // Setup file config
    cfg.SetConfigName("app")
    cfg.AddConfigPath("./configs")
    cfg.AddConfigPath("/etc/myapp/")
    cfg.AddConfigPath("$HOME/.myapp")
    cfg.SetConfigType("yaml")
    
    // Set reasonable defaults
    setDefaults(cfg)
    
    // Enable environment variables
    cfg.SetEnvPrefix("MYAPP")
    cfg.AutomaticEnv()
    
    // Read config file
    if err := cfg.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            log.Fatal("Config error:", err)
        }
        log.Println("No config file found, using defaults and env vars")
    }
    
    // Validate
    if err := validateConfig(cfg); err != nil {
        log.Fatal("Config validation failed:", err)
    }
    
    return cfg
}

func setDefaults(cfg config.Manager) {
    cfg.SetDefault("server.port", 8080)
    cfg.SetDefault("server.host", "localhost")
    cfg.SetDefault("app.debug", false)
    cfg.SetDefault("database.port", 5432)
    cfg.SetDefault("database.timeout", "30s")
    cfg.SetDefault("cache.enabled", true)
    cfg.SetDefault("cache.ttl", 3600)
}
```

### 2. Structured Configuration

```go
// Group related configs in structs
type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
}

type AppConfig struct {
    Database DatabaseConfig `mapstructure:"database"`
    Server   ServerConfig   `mapstructure:"server"`
    Cache    CacheConfig    `mapstructure:"cache"`
}

// Load once, use everywhere
func LoadConfig() *AppConfig {
    cfg := initConfig()
    
    var appConfig AppConfig
    if err := cfg.Unmarshal(&appConfig); err != nil {
        log.Fatal("Cannot unmarshal config:", err)
    }
    
    return &appConfig
}
```

### 3. Configuration Factory Pattern

```go
type ConfigManager struct {
    cfg config.Manager
}

func NewConfigManager() *ConfigManager {
    return &ConfigManager{
        cfg: initConfig(),
    }
}

func (cm *ConfigManager) GetDatabaseConfig() DatabaseConfig {
    var db DatabaseConfig
    cm.cfg.UnmarshalKey("database", &db)
    return db
}

func (cm *ConfigManager) GetServerConfig() ServerConfig {
    var server ServerConfig
    cm.cfg.UnmarshalKey("server", &server)
    return server
}

func (cm *ConfigManager) IsDebugMode() bool {
    debug, _ := cm.cfg.GetBool("app.debug")
    return debug
}
```

### 4. Error Handling

```go
func safeGetString(cfg config.Manager, key, defaultValue string) string {
    if value, exists := cfg.GetString(key); exists {
        return value
    }
    return defaultValue
}

func requireString(cfg config.Manager, key string) string {
    if value, exists := cfg.GetString(key); exists {
        return value
    }
    log.Fatalf("Required config key '%s' not found", key)
    return ""
}

func getIntWithValidation(cfg config.Manager, key string, min, max int) (int, error) {
    value, exists := cfg.GetInt(key)
    if !exists {
        return 0, fmt.Errorf("key '%s' not found", key)
    }
    
    if value < min || value > max {
        return 0, fmt.Errorf("key '%s' value %d out of range [%d, %d]", key, value, min, max)
    }
    
    return value, nil
}
```

## Common Issues & Solutions

### 1. File Not Found
```go
if err := cfg.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); ok {
        log.Println("Config file not found, using defaults")
        // Continue with defaults and env vars
    } else {
        log.Fatal("Error reading config:", err)
    }
}
```

### 2. Type Conversion Errors
```go
// Always check existence first
if port, exists := cfg.GetInt("server.port"); exists {
    if port <= 0 {
        log.Fatal("Invalid port value")
    }
    server.Listen(port)
} else {
    server.Listen(8080) // Use default
}
```

### 3. Environment Variable Override
```go
// Make sure AutomaticEnv is called
cfg.AutomaticEnv()
cfg.SetEnvPrefix("MYAPP")

// Env var names are case-sensitive and use underscores
// MYAPP_DATABASE_HOST maps to database.host
```

### 4. Watch Not Working
```go
// Must call ReadInConfig first
cfg.ReadInConfig()

// Then setup watching
cfg.OnConfigChange(func(e fsnotify.Event) {
    // Handle change
})
cfg.WatchConfig()
```
