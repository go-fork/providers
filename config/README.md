# Go-Fork Config Provider

Gói `config` cung cấp giải pháp quản lý cấu hình (configuration) hiện đại, linh hoạt và mở rộng cho ứng dụng Go. Gói này được thiết kế để dễ dàng tích hợp với các ứng dụng Go hiện đại và hỗ trợ nhiều nguồn cấu hình khác nhau.

## Giới thiệu

Quản lý cấu hình là một phần quan trọng trong mọi ứng dụng. Gói `config` cung cấp một cách thống nhất để quản lý và truy xuất cấu hình từ nhiều nguồn khác nhau. Với khả năng truy xuất theo dot notation và các hàm helper an toàn về kiểu dữ liệu, gói này giúp đơn giản hóa việc làm việc với cấu hình trong ứng dụng Go.

## Tính năng nổi bật

- **Đa nguồn cấu hình**: Hỗ trợ nạp cấu hình từ file YAML, JSON, biến môi trường (env), hoặc custom provider.
- **Dot notation**: Truy xuất cấu hình phân cấp với cú pháp "a.b.c", tự động gom các key con.
- **Type-safe accessors**: Các phương thức truy xuất an toàn về kiểu dữ liệu (GetString, GetInt, GetBool, GetStringMap, GetStringSlice).
- **Giá trị mặc định**: Hỗ trợ giá trị mặc định khi key không tồn tại.
- **Unmarshal vào struct**: Dễ dàng chuyển đổi cấu hình thành các struct Go.
- **Thread-safe**: An toàn khi truy xuất và cập nhật đồng thời.
- **Tích hợp DI**: Dễ dàng tích hợp với Dependency Injection container.
- **Mở rộng**: Dễ dàng mở rộng với custom Formatter.

## Cấu trúc package

```
config/
  ├── doc.go                 # Tài liệu tổng quan về package
  ├── manager.go             # Định nghĩa interface Manager và DefaultManager
  ├── provider.go            # ServiceProvider tích hợp với DI container
  ├── formatter/
  │   ├── formatter.go       # Interface Formatter
  │   ├── yaml.go            # YAML formatter implementation
  │   ├── json.go            # JSON formatter implementation
  │   └── env.go             # Environment variables formatter
  └── utils/
      └── helpers.go         # Các hàm helper (flatten/unflatten maps)
```

## Cách hoạt động

### Đăng ký Service Provider

Service Provider cho phép tích hợp dễ dàng gói `config` vào ứng dụng sử dụng DI container:

```go
// Trong file bootstrap của ứng dụng
import "github.com/go-fork/providers/config"

func bootstrap(app interface{}) {
    // Đăng ký config provider
    configProvider := config.NewServiceProvider()
    configProvider.Register(app)
    
    // Boot các providers sau khi tất cả đã đăng ký
    configProvider.Boot(app)
}
```

ServiceProvider sẽ tự động:
1. Tạo một config manager mới
2. Nạp cấu hình từ biến môi trường (prefix APP_)
3. Nạp cấu hình từ thư mục `configs` (các file YAML)
4. Đăng ký manager vào container với key "config"

### Cách Load các formatter

Bạn có thể nạp cấu hình từ nhiều nguồn khác nhau thông qua các formatter:

```go
manager := config.NewManager()

// Nạp cấu hình từ file YAML
yamlFormatter := formatter.NewYamlFormatter("/path/to/config.yaml")
err := manager.Load(yamlFormatter)

// Nạp cấu hình từ file JSON
jsonFormatter := formatter.NewJsonFormatter("/path/to/config.json")
err := manager.Load(jsonFormatter)

// Nạp cấu hình từ biến môi trường với prefix 'APP_'
envFormatter := formatter.NewEnvFormatter("APP")
err := manager.Load(envFormatter)

// Nạp cấu hình từ thư mục chứa nhiều file YAML
values, err := formatter.LoadFromDirectory("/path/to/configs")
if err == nil {
    for k, v := range values {
        manager.Set(k, v)
    }
}
```

### Cách làm việc với các tính năng trong Manager

#### Truy xuất giá trị đơn giản

```go
// Lấy giá trị với kiểu tương ứng
appName := manager.GetString("app.name", "Default App") // Giá trị mặc định là "Default App"
port := manager.GetInt("app.port", 8080)               // Giá trị mặc định là 8080
debug := manager.GetBool("app.debug", false)           // Giá trị mặc định là false

// Kiểm tra key tồn tại
if manager.Has("app.secret") {
    // Xử lý khi key tồn tại
}

// Lấy giá trị không ép kiểu
value := manager.Get("app.value")
```

#### Làm việc với map và slice

```go
// Lấy cấu hình dạng map
dbConfig := manager.GetStringMap("database")
host := dbConfig["host"].(string)
port := dbConfig["port"].(float64) // JSON numbers được unmarshal thành float64

// Lấy cấu hình dạng slice
allowedHosts := manager.GetStringSlice("app.allowed_hosts")
for _, host := range allowedHosts {
    // Xử lý từng host
}
```

#### Unmarshal vào struct

```go
// Định nghĩa struct tương ứng với cấu hình
type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
    Options  map[string]interface{} `json:"options"`
}

// Unmarshal cấu hình vào struct
var dbConfig DatabaseConfig
err := manager.Unmarshal("database", &dbConfig)
if err != nil {
    // Xử lý lỗi
}

// Sử dụng struct đã unmarshal
fmt.Printf("Connecting to %s:%d\n", dbConfig.Host, dbConfig.Port)
```

#### Làm việc với dot notation

Dot notation cho phép làm việc với cấu hình phân cấp một cách dễ dàng:

```go
// Set giá trị với dot notation
manager.Set("database.connections.mysql.host", "localhost")
manager.Set("database.connections.mysql.port", 3306)

// Lấy tất cả cấu hình của một key cha
mysqlConfig := manager.GetStringMap("database.connections.mysql")

// Truy xuất trực tiếp key con
host := manager.GetString("database.connections.mysql.host")
```

Khi gọi `Get` hoặc `GetStringMap` với một key cha, manager sẽ tự động gom tất cả key con có cùng prefix thành một map lồng nhau.

---

Để biết thêm thông tin chi tiết và API reference, vui lòng xem tài liệu trong file `doc.go` hoặc chạy lệnh `go doc github.com/go-fork/providers/config`.
