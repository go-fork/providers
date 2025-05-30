# Config Package Documentation

## Tổng quan

Package `go.fork.vn/providers/config` cung cấp một wrapper tiện lợi và mạnh mẽ cho thư viện [Viper](https://github.com/spf13/viper), được thiết kế để tích hợp dễ dàng với hệ thống Dependency Injection `go.fork.vn/di`.

### Đặc điểm chính

- **Interface-based Design**: Sử dụng interface `Manager` để định nghĩa API rõ ràng
- **Type-safe Operations**: Tất cả getter methods trả về tuple (value, exists) để đảm bảo an toàn về kiểu
- **Dependency Injection Ready**: Tích hợp sẵn với `go.fork.vn/di` thông qua `ServiceProvider`
- **Viper Compatibility**: Tận dụng toàn bộ sức mạnh của thư viện Viper
- **Auto Environment Support**: Tự động hỗ trợ biến môi trường

## Kiến trúc

### Core Components

1. **Manager Interface**: Định nghĩa API chính cho việc quản lý cấu hình
2. **manager struct**: Implementation của Manager interface sử dụng Viper
3. **ServiceProvider**: Tích hợp với DI container

### Interface Manager

Interface `Manager` định nghĩa các nhóm phương thức sau:

#### Value Retrieval Methods
```go
GetString(key string) (string, bool)
GetInt(key string) (int, bool) 
GetBool(key string) (bool, bool)
GetFloat(key string) (float64, bool)
GetDuration(key string) (time.Duration, bool)
GetTime(key string) (time.Time, bool)
GetSlice(key string) ([]interface{}, bool)
GetStringSlice(key string) ([]string, bool)
GetIntSlice(key string) ([]int, bool)
Get(key string) (interface{}, bool)
```

#### Map Retrieval Methods
```go
GetMap(key string) (map[string]interface{}, bool)
GetStringMap(key string) (map[string]interface{}, bool)
GetStringMapString(key string) (map[string]string, bool)
GetStringMapStringSlice(key string) (map[string][]string, bool)
```

#### Configuration Management
```go
Set(key string, value interface{}) error
SetDefault(key string, value interface{})
Has(key string) bool
AllSettings() map[string]interface{}
AllKeys() []string
```

#### File Operations
```go
SetConfigFile(path string)
SetConfigType(configType string)
SetConfigName(name string)
AddConfigPath(path string)
ReadInConfig() error
MergeInConfig() error
WriteConfig() error
SafeWriteConfig() error
WriteConfigAs(filename string) error
SafeWriteConfigAs(filename string) error
```

#### Environment Variables
```go
SetEnvPrefix(prefix string)
AutomaticEnv()
BindEnv(input ...string) error
```

#### Watching & Merging
```go
WatchConfig()
OnConfigChange(callback func(event fsnotify.Event))
MergeConfig(in io.Reader) error
```

#### Unmarshaling
```go
Unmarshal(target interface{}) error
UnmarshalKey(key string, target interface{}) error
```

### Implementation Details

#### manager struct
```go
type manager struct {
    *viper.Viper // Embedding Viper instance
}
```

Struct `manager` sử dụng composition pattern để embed Viper instance, cho phép kế thừa tất cả functionality của Viper trong khi vẫn maintain interface riêng.

#### Constructor
```go
func NewConfig() Manager {
    v := viper.New()
    v.AutomaticEnv() // Tự động enable environment variables
    return &manager{Viper: v}
}
```

### Type Safety Features

Tất cả getter methods của interface Manager đều trả về tuple `(value, exists)` thay vì chỉ value như Viper gốc. Điều này đảm bảo:

1. **Explicit Existence Check**: Bạn luôn biết liệu key có tồn tại hay không
2. **Zero Value Disambiguation**: Phân biệt được zero value và không tồn tại
3. **Error Prevention**: Tránh sử dụng nhầm zero value khi key không tồn tại

#### Ví dụ Type Safety

```go
// Unsafe - Viper style
port := cfg.GetInt("server.port") // Trả về 0 nếu key không tồn tại

// Safe - Manager interface style
if port, exists := cfg.GetInt("server.port"); exists {
    server.Listen(port)
} else {
    log.Fatal("server.port configuration is required")
}
```

### ServiceProvider Integration

`ServiceProvider` implement interface `di.ServiceProvider` để tích hợp với DI container:

```go
type ServiceProvider struct{}

func (p *ServiceProvider) Register(container di.Container) error
func (p *ServiceProvider) Requires() []string
func (p *ServiceProvider) Providers() []string
func (p *ServiceProvider) Boot(app interface{})
```

#### Lifecycle

1. **Register Phase**: Tạo Manager instance và đăng ký vào container với key "config"
2. **Boot Phase**: No-op (không cần thực hiện gì)
3. **Runtime**: Manager có thể được resolve từ container

### Error Handling

#### Set Method
```go
func (m *manager) Set(key string, value interface{}) error {
    if key == "" {
        return fmt.Errorf("key cannot be empty")
    }
    m.Viper.Set(key, value)
    return nil
}
```

#### Existence Checking
Tất cả getter methods sử dụng `m.Viper.IsSet(key)` để kiểm tra sự tồn tại trước khi trả về value.

### Map Handling

Method `GetMap` có logic đặc biệt để xử lý hai trường hợp:

1. **Direct Map Value**: Key trực tiếp chứa map value
2. **Constructed from SubKeys**: Xây dựng map từ các subkey có cùng prefix

```go
func (m *manager) GetMap(key string) (map[string]interface{}, bool) {
    if !m.Viper.IsSet(key) {
        return nil, false
    }

    // Kiểm tra direct map value
    val := m.Viper.Get(key)
    if mapVal, ok := val.(map[string]interface{}); ok {
        return mapVal, true
    }

    // Xây dựng từ subkeys
    return m.getMapFromSubKeys(key)
}
```

### Environment Variable Support

Manager tự động enable environment variable support qua `AutomaticEnv()` trong constructor. Các method bổ sung:

- `SetEnvPrefix(prefix)`: Đặt prefix cho env vars
- `BindEnv(input...)`: Bind specific keys với env vars
- `AutomaticEnv()`: Enable automatic env var lookup

### Configuration Sources Priority

Viper sử dụng thứ tự ưu tiên sau (cao đến thấp):

1. Explicit Set calls
2. Command line flags  
3. Environment variables
4. Configuration files
5. Key/value stores
6. Default values

### Thread Safety

Viper (và do đó Manager) là thread-safe cho read operations nhưng **không thread-safe cho write operations**. Nếu cần concurrent writes, bạn cần implement locking riêng.

### Performance Considerations

- Getter methods có overhead nhỏ do existence checking
- Map construction từ subkeys có thể expensive với large configs
- File watching sử dụng fsnotify, có thể impact performance trên systems với nhiều files

### Best Practices

1. **Use Type-Safe Getters**: Luôn check exists boolean
2. **Set Defaults Early**: Sử dụng SetDefault trong initialization
3. **Organize Config Hierarchically**: Sử dụng dot notation cho nested config
4. **Handle File Not Found**: Check for ConfigFileNotFoundError specifically
5. **Validate Config After Load**: Implement validation logic sau ReadInConfig

### Common Patterns

#### Initialization Pattern
```go
cfg := config.NewConfig()
cfg.SetConfigName("app")
cfg.AddConfigPath("./configs")
cfg.SetConfigType("yaml")

// Set defaults
cfg.SetDefault("server.port", 8080)
cfg.SetDefault("app.debug", false)

// Read config
if err := cfg.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
        log.Fatal("Error reading config:", err)
    }
}
```

#### Struct Mapping Pattern
```go
type Config struct {
    App      AppConfig      `mapstructure:"app"`
    Database DatabaseConfig `mapstructure:"database"`
}

var config Config
if err := cfg.Unmarshal(&config); err != nil {
    log.Fatal("Cannot unmarshal config:", err)
}
```

#### Environment Override Pattern
```go
cfg.SetEnvPrefix("MYAPP")
cfg.AutomaticEnv()
// Env var MYAPP_DATABASE_HOST overrides database.host
```
