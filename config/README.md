# Config

Một giải pháp quản lý cấu hình hoàn chỉnh và linh hoạt cho ứng dụng Go, dựa trên nền tảng của thư viện [Viper](https://github.com/spf13/viper).

## Giới thiệu

Package `config` cung cấp một wrapper tiện lợi cho thư viện Viper nổi tiếng, đồng thời mở rộng và chuẩn hóa API để dễ dàng tích hợp vào các ứng dụng thông qua Dependency Injection. Thư viện được thiết kế nhằm tối ưu quy trình quản lý cấu hình, đảm bảo tính nhất quán và linh hoạt cho các ứng dụng Go hiện đại.

Với thiết kế hướng interface rõ ràng và API sạch sẽ, package này cung cấp một lớp trừu tượng tiện lợi giúp làm việc với cấu hình trở nên đơn giản, đồng thời vẫn giữ được tất cả sức mạnh và linh hoạt của Viper.

Package này giải quyết các vấn đề phổ biến trong quản lý cấu hình:
- Đọc cấu hình từ nhiều nguồn (file, biến môi trường, flags)
- Tự động theo dõi và tải lại khi cấu hình thay đổi
- API an toàn về kiểu dữ liệu với cơ chế trả về giá trị kèm trạng thái
- Tích hợp liền mạch với hệ thống Dependency Injection

## Tính năng

| Phương thức | Mô tả | Tham số | Trả về/Lỗi |
|-------------|-------|---------|------------|
| `NewConfig()` | Tạo một đối tượng Manager mới sử dụng Viper làm backend | Không | `Manager`: Đối tượng Manager sẵn sàng sử dụng |
| `Get(key string)` | Trả về giá trị cho key theo kiểu gốc | `key`: Khóa cấu hình | `interface{}`: Giá trị gốc<br>`bool`: Tồn tại hay không |
| `GetString(key string)` | Trả về giá trị chuỗi cho key | `key`: Khóa cấu hình | `string`: Giá trị chuỗi<br>`bool`: Tồn tại hay không |
| `GetInt(key string)` | Trả về giá trị số nguyên cho key | `key`: Khóa cấu hình | `int`: Giá trị số nguyên<br>`bool`: Tồn tại hay không |
| `GetBool(key string)` | Trả về giá trị boolean cho key | `key`: Khóa cấu hình | `bool`: Giá trị boolean<br>`bool`: Tồn tại hay không |
| `GetFloat(key string)` | Trả về giá trị số thực cho key | `key`: Khóa cấu hình | `float64`: Giá trị số thực<br>`bool`: Tồn tại hay không |
| `GetDuration(key string)` | Trả về giá trị thời lượng cho key | `key`: Khóa cấu hình | `time.Duration`: Giá trị thời lượng<br>`bool`: Tồn tại hay không |
| `GetTime(key string)` | Trả về giá trị thời gian cho key | `key`: Khóa cấu hình | `time.Time`: Giá trị thời gian<br>`bool`: Tồn tại hay không |
| `GetSlice(key string)` | Trả về giá trị slice cho key | `key`: Khóa cấu hình | `[]interface{}`: Giá trị slice<br>`bool`: Tồn tại hay không |
| `GetStringSlice(key string)` | Trả về giá trị slice chuỗi cho key | `key`: Khóa cấu hình | `[]string`: Giá trị slice chuỗi<br>`bool`: Tồn tại hay không |
| `GetIntSlice(key string)` | Trả về giá trị slice số nguyên cho key | `key`: Khóa cấu hình | `[]int`: Giá trị slice số nguyên<br>`bool`: Tồn tại hay không |
| `GetMap(key string)` | Trả về giá trị map cho key | `key`: Khóa cấu hình | `map[string]interface{}`: Giá trị map<br>`bool`: Tồn tại hay không |
| `GetStringMap(key string)` | Trả về giá trị map của các interface cho key | `key`: Khóa cấu hình | `map[string]interface{}`: Giá trị map<br>`bool`: Tồn tại hay không |
| `GetStringMapString(key string)` | Trả về giá trị map của các chuỗi cho key | `key`: Khóa cấu hình | `map[string]string`: Giá trị map chuỗi<br>`bool`: Tồn tại hay không |
| `GetStringMapStringSlice(key string)` | Trả về giá trị map của các slice chuỗi cho key | `key`: Khóa cấu hình | `map[string][]string`: Giá trị map slice chuỗi<br>`bool`: Tồn tại hay không |
| `Set(key string, value interface{})` | Cập nhật hoặc thêm một giá trị vào cấu hình | `key`: Khóa cấu hình<br>`value`: Giá trị cần thiết lập | `error`: Lỗi nếu key rỗng |
| `SetDefault(key string, value interface{})` | Thiết lập giá trị mặc định cho key | `key`: Khóa cấu hình<br>`value`: Giá trị mặc định | Không |
| `Has(key string)` | Kiểm tra xem key có tồn tại hay không | `key`: Khóa cấu hình | `bool`: true nếu tồn tại |
| `AllSettings()` | Trả về tất cả cấu hình dưới dạng map | Không | `map[string]interface{}`: Tất cả cấu hình |
| `AllKeys()` | Trả về tất cả các khóa có giá trị | Không | `[]string`: Danh sách các khóa |
| `Unmarshal(key string, target interface{})` | Ánh xạ cấu hình vào một struct Go | `key`: Tiền tố khóa cấu hình<br>`target`: Con trỏ tới struct | `error`: Lỗi nếu có trong quá trình ánh xạ |
| `UnmarshalKey(key string, target interface{})` | Ánh xạ một khóa cụ thể vào struct | `key`: Khóa cấu hình<br>`target`: Con trỏ tới struct | `error`: Lỗi nếu có trong quá trình ánh xạ |
| `SetConfigFile(path string)` | Thiết lập đường dẫn tới file cấu hình | `path`: Đường dẫn đầy đủ tới file | Không |
| `SetConfigType(configType string)` | Thiết lập định dạng của file cấu hình | `configType`: Định dạng file | Không |
| `SetConfigName(name string)` | Thiết lập tên cho file cấu hình (không có phần mở rộng) | `name`: Tên file không có phần mở rộng | Không |
| `AddConfigPath(path string)` | Thêm đường dẫn để tìm kiếm file cấu hình | `path`: Đường dẫn | Không |
| `ReadInConfig()` | Tìm kiếm và đọc file cấu hình | Không | `error`: Lỗi nếu không thể đọc |
| `MergeInConfig()` | Gộp file cấu hình mới với cấu hình hiện tại | Không | `error`: Lỗi nếu không thể đọc |
| `WriteConfig()` | Ghi cấu hình hiện tại vào file | Không | `error`: Lỗi nếu không thể ghi |
| `SafeWriteConfig()` | Ghi cấu hình hiện tại vào file chỉ khi file không tồn tại | Không | `error`: Lỗi nếu file tồn tại |
| `WriteConfigAs(filename string)` | Ghi cấu hình vào file với tên được chỉ định | `filename`: Đường dẫn file | `error`: Lỗi nếu không thể ghi |
| `SafeWriteConfigAs(filename string)` | Ghi cấu hình vào file với tên được chỉ định chỉ khi file không tồn tại | `filename`: Đường dẫn file | `error`: Lỗi nếu file tồn tại |
| `WatchConfig()` | Theo dõi file cấu hình và tự động tải lại khi có thay đổi | Không | Không |
| `OnConfigChange(callback func(event fsnotify.Event))` | Thiết lập callback để chạy khi cấu hình thay đổi | `callback`: Hàm callback | Không |
| `SetEnvPrefix(prefix string)` | Thiết lập tiền tố cho biến môi trường | `prefix`: Tiền tố | Không |
| `AutomaticEnv()` | Kích hoạt tự động hỗ trợ biến môi trường | Không | Không |
| `BindEnv(input ...string)` | Ràng buộc một khóa Viper với biến môi trường | `input`: Khóa và tên biến môi trường | `error`: Lỗi nếu không thể ràng buộc |
| `MergeConfig(in io.Reader)` | Gộp cấu hình mới với cấu hình hiện tại | `in`: Reader chứa cấu hình | `error`: Lỗi nếu không thể gộp |

## Sử dụng mặc định

Sử dụng package `config` với Dependency Injection thông qua ServiceProvider:

```go
package main

import (
	"fmt"
	"log"

	"github.com/go-fork/di"
	"github.com/go-fork/providers/config"
)

func main() {
	// Khởi tạo ứng dụng với DI container
	app := di.New()
	
	// Đăng ký config service provider
	app.Register(config.NewServiceProvider())
	
	// Lấy config manager từ container
	container := app.Container()
	cfg := container.Get("config").(config.Manager)
	
	// Cấu hình và đọc file config
	cfg.SetConfigName("config")
	cfg.AddConfigPath(".")
	cfg.AddConfigPath("./configs")
	
	// Thiết lập một số giá trị mặc định
	cfg.SetDefault("app.name", "MyApp")
	cfg.SetDefault("app.port", 8080)
	cfg.SetDefault("database.timeout", "30s")
	
	// Đọc cấu hình từ file
	err := cfg.ReadInConfig()
	if err != nil {
		log.Printf("Cảnh báo: Không thể đọc file cấu hình: %v\n", err)
		log.Println("Tiếp tục với cấu hình mặc định và biến môi trường...")
	}
	
	// Kích hoạt hỗ trợ biến môi trường
	cfg.SetEnvPrefix("MYAPP")
	cfg.AutomaticEnv()
	
	// Sử dụng cấu hình
	if appName, ok := cfg.GetString("app.name"); ok {
		fmt.Printf("Tên ứng dụng: %s\n", appName)
	}
	
	if port, ok := cfg.GetInt("app.port"); ok {
		fmt.Printf("Cổng: %d\n", port)
	}
	
	// Khởi động các dịch vụ với cấu hình đã đọc...
}
```

## Sử dụng nâng cao (Advanced)

### Đọc cấu hình từ nhiều nguồn

```go
// Đọc từ file chính
cfg.SetConfigName("config")
cfg.AddConfigPath(".")
cfg.ReadInConfig()

// Gộp với file cấu hình môi trường
cfg.SetConfigName("config.development")
err := cfg.MergeInConfig()
if err != nil {
	log.Printf("Cảnh báo: Không tìm thấy file cấu hình môi trường: %v", err)
}

// Gộp với file cấu hình local (không theo dõi bởi git)
cfg.SetConfigName("config.local")
cfg.MergeInConfig() // Bỏ qua lỗi nếu không có
```

### Ánh xạ cấu hình vào struct

```go
// Định nghĩa struct cho cấu hình
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Server  struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		TLS  struct {
			Enabled  bool   `mapstructure:"enabled"`
			CertFile string `mapstructure:"cert_file"`
			KeyFile  string `mapstructure:"key_file"`
		} `mapstructure:"tls"`
	} `mapstructure:"server"`
	Database struct {
		Driver   string `mapstructure:"driver"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Options  struct {
			MaxOpenConns int           `mapstructure:"max_open_conns"`
			MaxIdleConns int           `mapstructure:"max_idle_conns"`
			MaxLifetime  time.Duration `mapstructure:"max_lifetime"`
		} `mapstructure:"options"`
	} `mapstructure:"database"`
}

// Ánh xạ toàn bộ cấu hình vào struct
var appConfig AppConfig
err := cfg.Unmarshal("", &appConfig)
if err != nil {
	log.Fatalf("Không thể ánh xạ cấu hình: %v", err)
}

fmt.Printf("Tên ứng dụng: %s\n", appConfig.Name)
fmt.Printf("Server sẽ chạy tại: %s:%d\n", appConfig.Server.Host, appConfig.Server.Port)

// Ánh xạ một phần cụ thể của cấu hình
var dbConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

err = cfg.Unmarshal("database", &dbConfig)
if err != nil {
	log.Fatalf("Không thể ánh xạ cấu hình database: %v", err)
}
```

### Theo dõi và tự động tải lại cấu hình

```go
// Thiết lập callback để xử lý khi cấu hình thay đổi
cfg.OnConfigChange(func(e fsnotify.Event) {
	fmt.Printf("Phát hiện thay đổi cấu hình: %s\n", e.Name)
	
	// Tải lại cấu hình vào struct
	var appConfig AppConfig
	err := cfg.Unmarshal("", &appConfig)
	if err != nil {
		log.Printf("Lỗi khi tải lại cấu hình: %v", err)
		return
	}
	
	// Tái khởi tạo các dịch vụ với cấu hình mới
	restartServices(appConfig)
	
	fmt.Println("Đã tải lại cấu hình thành công.")
})

// Bắt đầu theo dõi file cấu hình
cfg.WatchConfig()
```

### Xử lý cấu hình phân cấp

```go
// Ví dụ với cấu hình phân cấp trong YAML:
// 
// services:
//   auth:
//     enabled: true
//     providers:
//       - name: oauth
//         config:
//           client_id: '12345'
//           client_secret: 'secret'
//       - name: local
//         config:
//           password_policy:
//             min_length: 8
//             require_symbols: true

// Truy cập các giá trị cấu hình phân cấp
if authEnabled, ok := cfg.GetBool("services.auth.enabled"); ok && authEnabled {
	// Lấy danh sách các providers
	if providers, ok := cfg.GetSlice("services.auth.providers"); ok {
		for i, provider := range providers {
			if providerMap, ok := provider.(map[string]interface{}); ok {
				name := providerMap["name"].(string)
				fmt.Printf("Provider %d: %s\n", i+1, name)
				
				// Truy cập cấu hình cụ thể của provider
				configKey := fmt.Sprintf("services.auth.providers.%d.config", i)
				if providerConfig, ok := cfg.GetMap(configKey); ok {
					// Xử lý cấu hình provider
					fmt.Printf("  Cấu hình: %v\n", providerConfig)
				}
			}
		}
	}
}
```

## Lưu ý

1. **Kiểm tra tồn tại**: Luôn kiểm tra giá trị trả về thứ hai (boolean) khi đọc cấu hình để đảm bảo khóa tồn tại:
   ```go
   if value, ok := cfg.GetString("app.key"); ok {
       // Sử dụng value
   } else {
       // Xử lý trường hợp không có key
   }
   ```

2. **Độ ưu tiên**: Thứ tự ưu tiên của các nguồn cấu hình từ cao xuống thấp:
   - `Set()` gọi trực tiếp trong code
   - Biến môi trường
   - File cấu hình
   - Giá trị mặc định thông qua `SetDefault()`

3. **Kiểu dữ liệu an toàn**: Sử dụng phương thức phù hợp với kiểu dữ liệu cần lấy (`GetString`, `GetInt`, `GetBool`, v.v.) để tránh các lỗi ép kiểu.

4. **Ánh xạ struct**: Khi ánh xạ cấu hình vào struct:
   - Sử dụng tag `mapstructure` để chỉ định tên khóa
   - Struct phải được truyền vào dưới dạng con trỏ (`&myStruct`)
   - Cấu trúc phân cấp của struct phải tương ứng với cấu trúc cấu hình

5. **Biến môi trường**: Khi sử dụng biến môi trường:
   - Dấu chấm (`.`) trong key sẽ được chuyển thành dấu gạch dưới (`_`)
   - Chữ cái được chuyển thành chữ hoa
   - Tiền tố (nếu có) được thêm vào
   - Ví dụ: `database.host` với tiền tố `APP` trở thành `APP_DATABASE_HOST`

6. **Hot reload**: Khi sử dụng `WatchConfig()`, cần đảm bảo việc cập nhật cấu hình không gây ra race condition trong ứng dụng.

7. **Bảo mật**: Tránh lưu thông tin nhạy cảm (mật khẩu, khóa API) trong file cấu hình. Sử dụng biến môi trường hoặc các giải pháp lưu trữ bí mật (secret management) cho những thông tin này.
