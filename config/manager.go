package config

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Manager là interface chính cho việc quản lý cấu hình, bao bọc và mở rộng chức năng
// của thư viện Viper. Interface này cung cấp các phương thức để truy xuất và quản lý
// cấu hình với API nhất quán và an toàn về kiểu dữ liệu.
type Manager interface {
	// GetString trả về giá trị chuỗi cho key.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất, hỗ trợ dot notation (vd: "app.name")
	//
	// Returns:
	//   - string: Giá trị chuỗi nếu key tồn tại
	//   - bool: true nếu key tồn tại, ngược lại là false
	//
	// Example:
	//
	//	if name, ok := cfg.GetString("app.name"); ok {
	//	    fmt.Printf("App name: %s\n", name)
	//	}
	GetString(key string) (string, bool)

	// GetInt trả về giá trị số nguyên cho key.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - int: Giá trị số nguyên nếu key tồn tại và giá trị hợp lệ
	//   - bool: true nếu key tồn tại và giá trị có thể chuyển đổi thành int, ngược lại là false
	//
	// Example:
	//
	//	if port, ok := cfg.GetInt("server.port"); ok {
	//	    server.Listen(port)
	//	}
	GetInt(key string) (int, bool)

	// GetBool trả về giá trị boolean cho key.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - bool: Giá trị boolean nếu key tồn tại và giá trị hợp lệ
	//   - bool: true nếu key tồn tại và giá trị có thể chuyển đổi thành boolean, ngược lại là false
	//
	// Lưu ý: Các giá trị "true", "1", "t", "yes", "y", "on" được coi là true,
	// các giá trị khác được coi là false
	//
	// Example:
	//
	//	if debug, ok := cfg.GetBool("app.debug"); ok && debug {
	//	    enableDebugging()
	//	}
	GetBool(key string) (bool, bool)

	// GetFloat trả về giá trị số thực cho key.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - float64: Giá trị số thực nếu key tồn tại và giá trị hợp lệ
	//   - bool: true nếu key tồn tại và giá trị có thể chuyển đổi thành float64, ngược lại là false
	//
	// Example:
	//
	//	if ratio, ok := cfg.GetFloat("scaling.ratio"); ok {
	//	    applyScaling(ratio)
	//	}
	GetFloat(key string) (float64, bool)

	// GetDuration trả về giá trị khoảng thời gian (Duration) cho key.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - time.Duration: Giá trị khoảng thời gian nếu key tồn tại và giá trị hợp lệ
	//   - bool: true nếu key tồn tại và giá trị có thể chuyển đổi thành Duration, ngược lại là false
	//
	// Lưu ý: Hỗ trợ chuỗi format như "5s", "2m", "1h30m" hoặc số nguyên (đơn vị là nanosecond)
	//
	// Example:
	//
	//	if timeout, ok := cfg.GetDuration("http.timeout"); ok {
	//	    client.Timeout = timeout
	//	}
	GetDuration(key string) (time.Duration, bool)

	// GetTime trả về giá trị thời gian (Time) cho key.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - time.Time: Giá trị thời gian nếu key tồn tại và giá trị hợp lệ
	//   - bool: true nếu key tồn tại và giá trị có thể chuyển đổi thành Time, ngược lại là false
	//
	// Example:
	//
	//	if startTime, ok := cfg.GetTime("schedule.start"); ok {
	//	    scheduler.SetStartTime(startTime)
	//	}
	GetTime(key string) (time.Time, bool)

	// GetSlice trả về giá trị slice (mảng các phần tử) cho key.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - []interface{}: Mảng các giá trị nếu key tồn tại và giá trị là slice
	//   - bool: true nếu key tồn tại và giá trị là slice, ngược lại là false
	//
	// Example:
	//
	//	if items, ok := cfg.GetSlice("inventory.items"); ok {
	//	    for _, item := range items {
	//	        processItem(item)
	//	    }
	//	}
	GetSlice(key string) ([]interface{}, bool)

	// GetStringSlice trả về giá trị slice chuỗi cho key.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - []string: Mảng các chuỗi nếu key tồn tại và giá trị có thể chuyển đổi thành slice chuỗi
	//   - bool: true nếu key tồn tại, ngược lại là false
	//
	// Lưu ý: Các phần tử không phải chuỗi sẽ được chuyển đổi thành chuỗi
	//
	// Example:
	//
	//	if hosts, ok := cfg.GetStringSlice("database.replicas"); ok {
	//	    for _, host := range hosts {
	//	        pool.AddReplica(host)
	//	    }
	//	}
	GetStringSlice(key string) ([]string, bool)

	// GetIntSlice trả về giá trị slice số nguyên cho key.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - []int: Mảng các số nguyên nếu key tồn tại và giá trị có thể chuyển đổi thành slice số nguyên
	//   - bool: true nếu key tồn tại, ngược lại là false
	//
	// Lưu ý: Các phần tử không phải số nguyên sẽ được chuyển đổi nếu có thể
	//
	// Example:
	//
	//	if ports, ok := cfg.GetIntSlice("server.ports"); ok {
	//	    for _, port := range ports {
	//	        listeners = append(listeners, listenOn(port))
	//	    }
	//	}
	GetIntSlice(key string) ([]int, bool)

	// GetMap trả về giá trị map (bảng băm) cho key.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - map[string]interface{}: Map các giá trị nếu key tồn tại và giá trị là map hoặc có thể xây dựng từ subkeys
	//   - bool: true nếu key tồn tại và giá trị là map, hoặc có subkeys, ngược lại là false
	//
	// Lưu ý: Cơ chế này có thể xây dựng map từ các subkey có cùng prefix
	//
	// Example:
	//
	//	if dbConfig, ok := cfg.GetMap("database"); ok {
	//	    host := dbConfig["host"].(string)
	//	    port := dbConfig["port"].(int)
	//	    connectToDatabase(host, port)
	//	}
	GetMap(key string) (map[string]interface{}, bool)

	// GetStringMap trả về giá trị liên kết với key dưới dạng map các interface{}.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - map[string]interface{}: Map các giá trị nếu key tồn tại
	//   - bool: true nếu key tồn tại, ngược lại là false
	//
	// Example:
	//
	//	if settings, ok := cfg.GetStringMap("app.settings"); ok {
	//	    for k, v := range settings {
	//	        fmt.Printf("%s: %v\n", k, v)
	//	    }
	//	}
	GetStringMap(key string) (map[string]interface{}, bool)

	// GetStringMapString trả về giá trị liên kết với key dưới dạng map các chuỗi.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - map[string]string: Map các chuỗi nếu key tồn tại
	//   - bool: true nếu key tồn tại, ngược lại là false
	//
	// Example:
	//
	//	if headers, ok := cfg.GetStringMapString("http.headers"); ok {
	//	    for name, value := range headers {
	//	        req.Header.Set(name, value)
	//	    }
	//	}
	GetStringMapString(key string) (map[string]string, bool)

	// GetStringMapStringSlice trả về giá trị liên kết với key dưới dạng map với giá trị là slice chuỗi.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - map[string][]string: Map của các slice chuỗi nếu key tồn tại
	//   - bool: true nếu key tồn tại, ngược lại là false
	//
	// Example:
	//
	//	if routeParams, ok := cfg.GetStringMapStringSlice("routes.parameters"); ok {
	//	    for route, params := range routeParams {
	//	        router.RegisterParams(route, params)
	//	    }
	//	}
	GetStringMapStringSlice(key string) (map[string][]string, bool)

	// Get trả về giá trị cho key theo kiểu gốc của nó.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần truy xuất
	//
	// Returns:
	//   - interface{}: Giá trị gốc nếu key tồn tại
	//   - bool: true nếu key tồn tại, ngược lại là false
	//
	// Example:
	//
	//	if value, ok := cfg.Get("complex.setting"); ok {
	//	    // Xử lý giá trị theo kiểu dữ liệu cụ thể
	//	    switch v := value.(type) {
	//	    case map[string]interface{}:
	//	        // Xử lý map
	//	    case []interface{}:
	//	        // Xử lý slice
	//	    }
	//	}
	Get(key string) (interface{}, bool)

	// Set cập nhật hoặc thêm một giá trị vào cấu hình.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần thiết lập, không được rỗng
	//   - value: interface{} - Giá trị cần thiết lập
	//
	// Returns:
	//   - error: Lỗi nếu key rỗng, nil nếu thành công
	//
	// Example:
	//
	//	err := cfg.Set("app.version", "1.0.0")
	//	if err != nil {
	//	    log.Fatalf("Không thể thiết lập cấu hình: %v", err)
	//	}
	Set(key string, value interface{}) error

	// SetDefault thiết lập giá trị mặc định cho key.
	//
	// Giá trị mặc định được sử dụng khi key không được thiết lập ở bất kỳ nguồn nào khác
	// (file, biến môi trường, Set trực tiếp).
	//
	// Params:
	//   - key: string - Khóa cấu hình cần thiết lập giá trị mặc định
	//   - value: interface{} - Giá trị mặc định
	//
	// Example:
	//
	//	cfg.SetDefault("server.port", 8080)
	//	cfg.SetDefault("log.level", "info")
	SetDefault(key string, value interface{})

	// Has kiểm tra xem key có tồn tại hay không.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần kiểm tra
	//
	// Returns:
	//   - bool: true nếu key tồn tại, ngược lại là false
	//
	// Example:
	//
	//	if cfg.Has("database.username") {
	//	    // Thiết lập kết nối database với thông tin xác thực
	//	} else {
	//	    // Sử dụng kết nối không xác thực hoặc mặc định
	//	}
	Has(key string) bool

	// AllSettings trả về tất cả cấu hình dưới dạng map.
	//
	// Returns:
	//   - map[string]interface{}: Tất cả cấu hình hiện tại
	//
	// Example:
	//
	//	settings := cfg.AllSettings()
	//	jsonConfig, _ := json.MarshalIndent(settings, "", "  ")
	//	fmt.Println("Cấu hình hiện tại:", string(jsonConfig))
	AllSettings() map[string]interface{}

	// AllKeys trả về tất cả các khóa có giá trị, bất kể chúng được thiết lập ở đâu.
	//
	// Returns:
	//   - []string: Danh sách tất cả các khóa đang có giá trị
	//
	// Example:
	//
	//	keys := cfg.AllKeys()
	//	fmt.Println("Tất cả các khóa cấu hình:")
	//	for _, key := range keys {
	//	    fmt.Printf("- %s\n", key)
	//	}
	AllKeys() []string

	// Unmarshal ánh xạ cấu hình vào một struct Go.
	//
	// Params:
	//   - key: string - Prefix của cấu hình cần ánh xạ, nếu rỗng thì ánh xạ toàn bộ cấu hình
	//   - target: interface{} - Con trỏ tới struct cần ánh xạ dữ liệu vào
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình ánh xạ, nil nếu thành công
	//
	// Example:
	//
	//	type DatabaseConfig struct {
	//	    Host     string `mapstructure:"host"`
	//	    Port     int    `mapstructure:"port"`
	//	    Username string `mapstructure:"username"`
	//	    Password string `mapstructure:"password"`
	//	}
	//
	//	var dbConfig DatabaseConfig
	//	err := cfg.UnmarshalKey("database", &dbConfig)
	//	if err != nil {
	//	    log.Fatalf("Không thể unmarshal cấu hình database: %v", err)
	//	}
	Unmarshal(target interface{}) error

	// UnmarshalKey ánh xạ một khóa cấu hình vào struct.
	//
	// Phương thức này tương tự Unmarshal nhưng chỉ hoạt động với một khóa duy nhất.
	//
	// Params:
	//   - key: string - Khóa cấu hình cần ánh xạ
	//   - target: interface{} - Con trỏ tới struct cần ánh xạ dữ liệu vào
	//
	// Returns:
	//   - error: Lỗi nếu có trong quá trình ánh xạ, nil nếu thành công
	//
	// Example:
	//
	//	type ServerConfig struct {
	//	    Port    int      `mapstructure:"port"`
	//	    Host    string   `mapstructure:"host"`
	//	    Domains []string `mapstructure:"domains"`
	//	}
	//
	//	var serverCfg ServerConfig
	//	err := cfg.UnmarshalKey("server", &serverCfg)
	UnmarshalKey(key string, target interface{}) error

	// SetConfigFile thiết lập đường dẫn tới file cấu hình.
	//
	// Params:
	//   - path: string - Đường dẫn đầy đủ tới file cấu hình
	//
	// Example:
	//
	//	cfg.SetConfigFile("/etc/myapp/config.yaml")
	//	err := cfg.ReadInConfig()
	SetConfigFile(path string)

	// SetConfigType thiết lập định dạng của file cấu hình.
	//
	// Định dạng này được sử dụng khi đọc cấu hình từ buffer hoặc
	// khi file cấu hình không có phần mở rộng để xác định định dạng.
	//
	// Params:
	//   - configType: string - Định dạng của file cấu hình (json, toml, yaml, yml, ini, hcl, env, props)
	//
	// Example:
	//
	//	cfg.SetConfigType("yaml")
	//	cfg.ReadConfig(bytes.NewBuffer(yamlConfig))
	SetConfigType(configType string)

	// SetConfigName thiết lập tên cho file cấu hình (không bao gồm phần mở rộng).
	//
	// Params:
	//   - name: string - Tên file cấu hình không có phần mở rộng
	//
	// Lưu ý: Phần mở rộng sẽ được tự động xác định dựa trên các định dạng hỗ trợ
	//
	// Example:
	//
	//	// Tìm file "config.yaml", "config.json", v.v. trong các đường dẫn đã thêm
	//	cfg.SetConfigName("config")
	SetConfigName(name string)

	// AddConfigPath thêm đường dẫn để tìm kiếm file cấu hình.
	//
	// Params:
	//   - path: string - Đường dẫn để Viper tìm kiếm file cấu hình
	//
	// Lưu ý: Có thể gọi nhiều lần để thêm nhiều đường dẫn
	//
	// Example:
	//
	//	cfg.SetConfigName("config")
	//	cfg.AddConfigPath(".")
	//	cfg.AddConfigPath("/etc/myapp/")
	//	cfg.AddConfigPath("$HOME/.myapp")
	//	err := cfg.ReadInConfig()
	AddConfigPath(path string)

	// ReadInConfig tìm kiếm và đọc file cấu hình từ đĩa hoặc key/value store,
	// từ một trong các đường dẫn đã định nghĩa.
	//
	// Returns:
	//   - error: Lỗi nếu không tìm thấy file cấu hình hoặc không thể đọc,
	//     nil nếu thành công
	//
	// Example:
	//
	//	cfg.SetConfigName("config")
	//	cfg.AddConfigPath(".")
	//	err := cfg.ReadInConfig()
	//	if err != nil {
	//	    if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	//	        // File không tồn tại, sử dụng giá trị mặc định
	//	    } else {
	//	        // Có lỗi khác khi đọc file
	//	        log.Fatalf("Không thể đọc cấu hình: %v", err)
	//	    }
	//	}
	ReadInConfig() error

	// MergeInConfig gộp file cấu hình mới với cấu hình hiện tại.
	//
	// Tìm kiếm và đọc file cấu hình từ đĩa hoặc key/value store,
	// sau đó gộp với cấu hình đã có.
	//
	// Returns:
	//   - error: Lỗi nếu không tìm thấy file cấu hình hoặc không thể đọc,
	//     nil nếu thành công
	//
	// Example:
	//
	//	// Đọc cấu hình chính
	//	cfg.SetConfigFile("config.yaml")
	//	cfg.ReadInConfig()
	//
	//	// Gộp với cấu hình bổ sung
	//	cfg.SetConfigFile("config.local.yaml")
	//	err := cfg.MergeInConfig()
	MergeInConfig() error

	// WriteConfig ghi cấu hình hiện tại vào file.
	//
	// Sẽ ghi đè lên file nếu nó đã tồn tại. Sẽ báo lỗi nếu không có file
	// cấu hình nào được thiết lập trước đó.
	//
	// Returns:
	//   - error: Lỗi trong quá trình ghi file, nil nếu thành công
	//
	// Example:
	//
	//	cfg.SetConfigFile("config.yaml")
	//	cfg.Set("database.host", "localhost")
	//	err := cfg.WriteConfig()
	WriteConfig() error

	// SafeWriteConfig ghi cấu hình hiện tại vào file chỉ khi file không tồn tại.
	//
	// Sẽ không ghi đè lên file nếu nó đã tồn tại. Sẽ báo lỗi nếu không có file
	// cấu hình nào được thiết lập trước đó.
	//
	// Returns:
	//   - error: Lỗi trong quá trình ghi file hoặc nếu file đã tồn tại, nil nếu thành công
	//
	// Example:
	//
	//	cfg.SetConfigFile("config.yaml")
	//	cfg.Set("app.name", "MyApp")
	//	err := cfg.SafeWriteConfig()
	SafeWriteConfig() error

	// WriteConfigAs ghi cấu hình hiện tại vào file với tên được chỉ định.
	//
	// Sẽ ghi đè lên file nếu nó đã tồn tại.
	//
	// Params:
	//   - filename: string - Đường dẫn tới file cấu hình để ghi
	//
	// Returns:
	//   - error: Lỗi trong quá trình ghi file, nil nếu thành công
	//
	// Example:
	//
	//	cfg.Set("version", "1.0.0")
	//	err := cfg.WriteConfigAs("/etc/myapp/config.yaml")
	WriteConfigAs(filename string) error

	// SafeWriteConfigAs ghi cấu hình hiện tại vào file với tên được chỉ định
	// chỉ khi file không tồn tại.
	//
	// Sẽ không ghi đè lên file nếu nó đã tồn tại.
	//
	// Params:
	//   - filename: string - Đường dẫn tới file cấu hình để ghi
	//
	// Returns:
	//   - error: Lỗi trong quá trình ghi file hoặc nếu file đã tồn tại, nil nếu thành công
	//
	// Example:
	//
	//	cfg.Set("version", "1.0.0")
	//	err := cfg.SafeWriteConfigAs("/etc/myapp/config.default.yaml")
	SafeWriteConfigAs(filename string) error

	// WatchConfig theo dõi file cấu hình và tự động tải lại khi có thay đổi.
	//
	// Phương thức này khởi tạo quá trình theo dõi file cấu hình trong nền.
	// Khi phát hiện thay đổi, cấu hình sẽ được tải lại tự động.
	//
	// Lưu ý: Hãy gọi OnConfigChange để đăng ký hàm callback xử lý sự kiện thay đổi
	//
	// Example:
	//
	//	cfg.OnConfigChange(func(e fsnotify.Event) {
	//	    fmt.Println("Cấu hình đã thay đổi:", e.Name)
	//	})
	//	cfg.WatchConfig()
	WatchConfig()

	// OnConfigChange thiết lập callback để chạy khi cấu hình thay đổi.
	//
	// Params:
	//   - callback: func(fsnotify.Event) - Hàm callback được gọi khi phát hiện thay đổi
	//
	// Lưu ý: Phải gọi WatchConfig() sau khi thiết lập callback này để kích hoạt theo dõi
	//
	// Example:
	//
	//	cfg.OnConfigChange(func(e fsnotify.Event) {
	//	    log.Printf("Phát hiện thay đổi cấu hình: %s", e.Name)
	//	    // Thực hiện các hành động cần thiết, ví dụ tải lại dịch vụ
	//	    if err := reloadServices(); err != nil {
	//	        log.Printf("Lỗi khi tải lại dịch vụ: %v", err)
	//	    }
	//	})
	//	cfg.WatchConfig()
	OnConfigChange(callback func(event fsnotify.Event))

	// SetEnvPrefix thiết lập tiền tố cho biến môi trường.
	//
	// Tất cả biến môi trường được sử dụng để ghi đè cấu hình sẽ được tìm kiếm
	// với tiền tố này, trừ khi được ràng buộc thông qua BindEnv.
	//
	// Params:
	//   - prefix: string - Tiền tố cho các biến môi trường
	//
	// Example:
	//
	//	// Thiết lập tiền tố "MYAPP_" cho biến môi trường
	//	cfg.SetEnvPrefix("MYAPP")
	//	cfg.AutomaticEnv()
	//	// Giờ biến môi trường MYAPP_DATABASE_HOST sẽ được ánh xạ tới database.host
	SetEnvPrefix(prefix string)

	// AutomaticEnv kích hoạt tự động hỗ trợ biến môi trường.
	//
	// Sau khi gọi phương thức này, mọi truy cập vào cấu hình sẽ kiểm tra xem
	// có biến môi trường phù hợp không và sử dụng giá trị từ đó nếu có.
	//
	// Tên biến môi trường được tạo thành từ:
	//   - Tiền tố (nếu có, từ SetEnvPrefix)
	//   - Tên khóa với dấu chấm (".") được thay thế bằng dấu gạch dưới ("_")
	//
	// Example:
	//
	//	cfg.SetEnvPrefix("MYAPP")
	//	cfg.AutomaticEnv()
	//	// Giờ đây, khi gọi cfg.GetString("database.host"), sẽ kiểm tra MYAPP_DATABASE_HOST
	AutomaticEnv()

	// BindEnv ràng buộc một khóa Viper với biến môi trường.
	//
	// Params:
	//   - input: ...string - Tham số biến đổi:
	//     - Nếu có 1 tham số: khóa Viper và tên biến môi trường được coi là giống nhau
	//     - Nếu có 2 tham số: tham số đầu là khóa Viper, tham số thứ hai là tên biến môi trường
	//
	// Returns:
	//   - error: Lỗi nếu không thể ràng buộc, nil nếu thành công
	//
	// Example:
	//
	//	// Ràng buộc "port" với biến môi trường "PORT"
	//	cfg.BindEnv("port")
	//
	//	// Ràng buộc "database.host" với biến môi trường "DB_HOST"
	//	cfg.BindEnv("database.host", "DB_HOST")
	BindEnv(input ...string) error

	// MergeConfig gộp cấu hình mới với cấu hình hiện tại.
	//
	// Params:
	//   - in: io.Reader - Reader chứa cấu hình cần gộp
	//
	// Returns:
	//   - error: Lỗi nếu không thể đọc hoặc gộp cấu hình, nil nếu thành công
	//
	// Example:
	//
	//	configData := []byte(`{"feature": {"enabled": true}}`)
	//	err := cfg.MergeConfig(bytes.NewBuffer(configData))
	//	if err != nil {
	//	    log.Fatalf("Không thể gộp cấu hình: %v", err)
	//	}
	MergeConfig(in io.Reader) error
}

// manager là struct triển khai interface Manager bằng cách nhúng Viper.
//
// manager sử dụng mô hình composition (nhúng struct) để kế thừa tất cả các phương thức
// của thư viện Viper, đồng thời mở rộng và chuẩn hóa API để phù hợp với interface Manager.
// Sử dụng mô hình này giúp dễ dàng tích hợp và mở rộng thư viện Viper mà không cần
// thay đổi mã nguồn gốc.
type manager struct {
	*viper.Viper // Nhúng Viper để thừa kế tính năng
}

// NewConfig tạo một đối tượng Manager mới sử dụng Viper làm backend.
//
// Hàm này khởi tạo một instance mới của Viper và cấu hình các thiết lập cơ bản,
// sau đó đóng gói trong struct manager để triển khai interface Manager.
//
// Returns:
//   - Manager: Đối tượng Manager cài đặt đầy đủ interface, sẵn sàng sử dụng
//
// Example:
//
//	cfg := config.NewConfig()
//	cfg.SetConfigFile("config.yaml")
//	err := cfg.ReadInConfig()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	appName, _ := cfg.GetString("app.name")
func NewConfig() Manager {
	v := viper.New()
	v.AutomaticEnv() // Tự động đọc từ biến môi trường

	return &manager{Viper: v}
}

// Get trả về giá trị cho key theo kiểu gốc của nó.
//
// Phương thức này trả về giá trị cấu hình theo kiểu dữ liệu gốc của nó
// (map, slice, string, number, boolean, v.v.) và một boolean cho biết key có tồn tại không.
//
// Params:
//   - key: string - Khóa cấu hình cần truy xuất
//
// Returns:
//   - interface{}: Giá trị gốc nếu key tồn tại, nil nếu không
//   - bool: true nếu key tồn tại, false nếu không
//
// Example:
//
//	if value, ok := cfg.Get("api.limits"); ok {
//	    // Kiểm tra kiểu dữ liệu của giá trị
//	    switch v := value.(type) {
//	    case map[string]interface{}:
//	        fmt.Printf("API limits có %d thiết lập\n", len(v))
//	        if maxRequests, exists := v["max_requests"]; exists {
//	            fmt.Printf("Max requests: %v\n", maxRequests)
//	        }
//	    case int:
//	        fmt.Printf("API limit: %d\n", v)
//	    default:
//	        fmt.Printf("API limits có kiểu không xác định: %T\n", v)
//	    }
//	} else {
//	    fmt.Println("Không tìm thấy cấu hình api.limits")
//	}
func (m *manager) Get(key string) (interface{}, bool) {
	if !m.Viper.IsSet(key) {
		return nil, false
	}
	return m.Viper.Get(key), true
}

// GetString trả về giá trị chuỗi cho key.
//
// Phương thức này trả về giá trị cấu hình dưới dạng chuỗi và một boolean
// cho biết key có tồn tại không. Viper sẽ tự động chuyển đổi các kiểu dữ liệu khác
// thành chuỗi nếu có thể.
//
// Params:
//   - key: string - Khóa cấu hình cần truy xuất
//
// Returns:
//   - string: Giá trị chuỗi nếu key tồn tại, chuỗi rỗng nếu không
//   - bool: true nếu key tồn tại, false nếu không
//
// Example:
//
//	if dbHost, ok := cfg.GetString("database.host"); ok {
//	    fmt.Printf("Kết nối tới database host: %s\n", dbHost)
//	    // Sử dụng dbHost để thiết lập kết nối
//	    conn := database.Connect(dbHost)
//	} else {
//	    fmt.Println("Không tìm thấy cấu hình database.host, sử dụng localhost")
//	    conn := database.Connect("localhost")
//	}
func (m *manager) GetString(key string) (string, bool) {
	if !m.Viper.IsSet(key) {
		return "", false
	}
	return m.Viper.GetString(key), true
}

// GetInt trả về giá trị số nguyên cho key.
func (m *manager) GetInt(key string) (int, bool) {
	if !m.Viper.IsSet(key) {
		return 0, false
	}
	return m.Viper.GetInt(key), true
}

// GetBool trả về giá trị boolean cho key.
func (m *manager) GetBool(key string) (bool, bool) {
	if !m.Viper.IsSet(key) {
		return false, false
	}
	return m.Viper.GetBool(key), true
}

// GetFloat trả về giá trị số thực cho key.
func (m *manager) GetFloat(key string) (float64, bool) {
	if !m.Viper.IsSet(key) {
		return 0, false
	}
	return m.Viper.GetFloat64(key), true
}

// GetDuration returns the duration value for a key.
func (m *manager) GetDuration(key string) (time.Duration, bool) {
	if !m.Viper.IsSet(key) {
		return 0, false
	}
	return m.Viper.GetDuration(key), true
}

// GetTime returns the time.Time value for a key.
func (m *manager) GetTime(key string) (time.Time, bool) {
	if !m.Viper.IsSet(key) {
		return time.Time{}, false
	}
	return m.Viper.GetTime(key), true
}

// GetSlice trả về giá trị slice cho key.
func (m *manager) GetSlice(key string) ([]interface{}, bool) {
	if !m.Viper.IsSet(key) {
		return nil, false
	}
	value := m.Viper.Get(key)
	if slice, ok := value.([]interface{}); ok {
		return slice, true
	}
	return nil, false
}

// GetStringSlice returns string slice value for a key.
func (m *manager) GetStringSlice(key string) ([]string, bool) {
	if !m.Viper.IsSet(key) {
		return nil, false
	}
	return m.Viper.GetStringSlice(key), true
}

// GetIntSlice returns int slice value for a key.
func (m *manager) GetIntSlice(key string) ([]int, bool) {
	if !m.Viper.IsSet(key) {
		return nil, false
	}
	return m.Viper.GetIntSlice(key), true
}

// GetMap trả về giá trị map cho key.
func (m *manager) GetMap(key string) (map[string]interface{}, bool) {
	// Kiểm tra nếu key không tồn tại
	if !m.Viper.IsSet(key) {
		return nil, false
	}

	// Kiểm tra nếu giá trị là map
	val := m.Viper.Get(key)
	if mapVal, ok := val.(map[string]interface{}); ok {
		return mapVal, true
	}

	// Sau đó kiểm tra các subkey
	return m.getMapFromSubKeys(key)
}

// getMapFromSubKeys tạo map từ các subkeys với tiền tố chung.
// Hàm này được tách ra để dễ dàng test và cải thiện độ bao phủ.
func (m *manager) getMapFromSubKeys(key string) (map[string]interface{}, bool) {
	result := make(map[string]interface{})
	prefix := key + "."
	subKeys := m.Viper.AllKeys()

	hasSubKey := false
	for _, subKey := range subKeys {
		if strings.HasPrefix(subKey, prefix) {
			hasSubKey = true
			shortKey := strings.TrimPrefix(subKey, prefix)
			result[shortKey] = m.Viper.Get(subKey)
		}
	}

	if hasSubKey {
		return result, true
	}
	return nil, false
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (m *manager) GetStringMap(key string) (map[string]interface{}, bool) {
	if !m.Viper.IsSet(key) {
		return nil, false
	}
	return m.Viper.GetStringMap(key), true
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (m *manager) GetStringMapString(key string) (map[string]string, bool) {
	if !m.Viper.IsSet(key) {
		return nil, false
	}
	return m.Viper.GetStringMapString(key), true
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (m *manager) GetStringMapStringSlice(key string) (map[string][]string, bool) {
	if !m.Viper.IsSet(key) {
		return nil, false
	}
	return m.Viper.GetStringMapStringSlice(key), true
}

// Set cập nhật hoặc thêm một giá trị vào cấu hình.
func (m *manager) Set(key string, value interface{}) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	m.Viper.Set(key, value)
	return nil
}

// SetDefault sets default value for a key.
func (m *manager) SetDefault(key string, value interface{}) {
	m.Viper.SetDefault(key, value)
}

// Has kiểm tra xem key có tồn tại hay không.
//
// Phương thức này xác định liệu một khóa cấu hình có được thiết lập hay không,
// bất kể nguồn cấu hình nào (file, biến môi trường, giá trị mặc định, v.v.).
//
// Params:
//   - key: string - Khóa cấu hình cần kiểm tra
//
// Returns:
//   - bool: true nếu key tồn tại, ngược lại là false
//
// Example:
//
//	if cfg.Has("database.credentials") {
//	    // Sử dụng thông tin xác thực từ cấu hình
//	    username, _ := cfg.GetString("database.credentials.username")
//	    password, _ := cfg.GetString("database.credentials.password")
//	    db.Connect(username, password)
//	} else {
//	    // Sử dụng xác thực mặc định
//	    db.ConnectWithDefaults()
//	}
func (m *manager) Has(key string) bool {
	return m.Viper.IsSet(key)
}

// AllSettings trả về tất cả cấu hình dưới dạng map.
//
// Phương thức này trả về tất cả các thiết lập cấu hình hiện tại dưới dạng map lồng nhau,
// phản ánh cấu trúc phân cấp của cấu hình.
//
// Returns:
//   - map[string]interface{}: Map lồng nhau chứa tất cả cấu hình hiện tại
//
// Example:
//
//	settings := cfg.AllSettings()
//
//	// In ra cấu hình dưới dạng JSON có định dạng
//	jsonData, _ := json.MarshalIndent(settings, "", "  ")
//	fmt.Println(string(jsonData))
//
//	// Truy cập trực tiếp vào cấu hình qua map
//	if app, ok := settings["app"].(map[string]interface{}); ok {
//	    fmt.Println("App name:", app["name"])
//	}
func (m *manager) AllSettings() map[string]interface{} {
	return m.Viper.AllSettings()
}

// AllKeys trả về tất cả các khóa đang giữ giá trị, bất kể chúng được thiết lập ở đâu.
//
// Phương thức này trả về danh sách tất cả các khóa cấu hình hiện có trong hệ thống,
// bao gồm các khóa từ file cấu hình, biến môi trường, giá trị mặc định, v.v.
// Các khóa được trả về dưới dạng đường dẫn phẳng (flat path) sử dụng dấu chấm làm dấu phân cách.
//
// Returns:
//   - []string: Slice chứa tất cả các khóa có giá trị
//
// Example:
//
//	keys := cfg.AllKeys()
//	fmt.Println("Tất cả các khóa cấu hình:")
//	for _, key := range keys {
//	    value, _ := cfg.Get(key)
//	    fmt.Printf("- %s = %v\n", key, value)
//	}
//
//	// Kiểm tra xem có khóa nào liên quan đến database không
//	dbKeys := []string{}
//	for _, key := range keys {
//	    if strings.HasPrefix(key, "database.") {
//	        dbKeys = append(dbKeys, key)
//	    }
//	}
//	fmt.Printf("Tìm thấy %d khóa database\n", len(dbKeys))
func (m *manager) AllKeys() []string {
	return m.Viper.AllKeys()
}

// Unmarshal ánh xạ toàn bộ cấu hình vào một struct Go.
//
// Phương thức này chuyển đổi toàn bộ cấu hình thành một đối tượng struct Go.
// Khác với UnmarshalKey, phương thức này luôn ánh xạ toàn bộ cấu hình vào struct.
//
// Params:
//   - target: interface{} - Con trỏ tới struct cần ánh xạ dữ liệu vào
//
// Returns:
//   - error: nil nếu ánh xạ thành công, lỗi nếu không thể ánh xạ giá trị vào struct
//
// Example:
//
//	// Ánh xạ toàn bộ cấu hình vào struct root
//	type AppConfig struct {
//	    AppName  string         `mapstructure:"app_name"`
//	    Version  string         `mapstructure:"version"`
//	    Database DatabaseConfig `mapstructure:"database"`
//	    Cache    struct {
//	        Enabled bool `mapstructure:"enabled"`
//	        TTL     int  `mapstructure:"ttl"`
//	    } `mapstructure:"cache"`
//	}
//
//	var appConfig AppConfig
//	err := cfg.Unmarshal(&appConfig)
//	if err != nil {
//	    log.Fatalf("Lỗi khi unmarshal toàn bộ cấu hình: %v", err)
//	}
//
//	// Để ánh xạ một phần cụ thể của cấu hình, hãy sử dụng UnmarshalKey thay thế
//	type DatabaseConfig struct {
//	    Host     string   `mapstructure:"host"`
//	    Port     int      `mapstructure:"port"`
//	    Username string   `mapstructure:"username"`
//	    Password string   `mapstructure:"password"`
//	    Replicas []string `mapstructure:"replicas"`
//	}
//
//	var dbConfig DatabaseConfig
//	err = cfg.UnmarshalKey("database", &dbConfig)
func (m *manager) Unmarshal(target interface{}) error {
	return m.Viper.Unmarshal(target)
}

// UnmarshalKey ánh xạ một khóa cấu hình duy nhất vào một struct.
//
// Phương thức này khác với Unmarshal ở chỗ nó luôn yêu cầu một khóa và
// không hỗ trợ ánh xạ toàn bộ cấu hình. UnmarshalKey thường được sử dụng khi
// bạn muốn ánh xạ một phần cụ thể của cấu hình.
//
// Params:
//   - key: string - Khóa cấu hình cần ánh xạ
//   - target: interface{} - Con trỏ tới struct cần ánh xạ dữ liệu vào
//
// Returns:
//   - error: Lỗi nếu có trong quá trình ánh xạ, nil nếu thành công
//
// Example:
//
//	// Ánh xạ cấu hình api.rate_limits vào struct RateLimits
//	type RateLimits struct {
//	    Enabled      bool `mapstructure:"enabled"`
//	    MaxRequests  int  `mapstructure:"max_requests"`
//	    PerTimeUnit  string `mapstructure:"per_time_unit"`
//	    BurstSize    int  `mapstructure:"burst_size"`
//	}
//
//	var limits RateLimits
//	err := cfg.UnmarshalKey("api.rate_limits", &limits)
//	if err != nil {
//	    log.Fatalf("Không thể ánh xạ cấu hình rate limits: %v", err)
//	}
//
//	if limits.Enabled {
//	    rateLimiter := middleware.NewRateLimiter(limits.MaxRequests, limits.PerTimeUnit, limits.BurstSize)
//	    app.Use(rateLimiter)
//	}
func (m *manager) UnmarshalKey(key string, target interface{}) error {
	return m.Viper.UnmarshalKey(key, target)
}

// SetConfigFile explicitly sets the path to a config file.
func (m *manager) SetConfigFile(path string) {
	m.Viper.SetConfigFile(path)
}

// SetConfigType sets the type of the configuration.
func (m *manager) SetConfigType(configType string) {
	m.Viper.SetConfigType(configType)
}

// SetConfigName sets name for the config file.
func (m *manager) SetConfigName(name string) {
	m.Viper.SetConfigName(name)
}

// AddConfigPath adds a path for Viper to search for the config file in.
func (m *manager) AddConfigPath(path string) {
	m.Viper.AddConfigPath(path)
}

// ReadInConfig will discover and load the configuration file from disk
// and key/value stores, searching in one of the defined paths.
func (m *manager) ReadInConfig() error {
	return m.Viper.ReadInConfig()
}

// MergeInConfig merges a new config file with the current configuration.
func (m *manager) MergeInConfig() error {
	return m.Viper.MergeInConfig()
}

// WriteConfig writes the current configuration to a file.
func (m *manager) WriteConfig() error {
	return m.Viper.WriteConfig()
}

// SafeWriteConfig writes the current configuration to file only if
// the file doesn't exist.
func (m *manager) SafeWriteConfig() error {
	return m.Viper.SafeWriteConfig()
}

// WriteConfigAs writes current configuration to a file with specified name.
func (m *manager) WriteConfigAs(filename string) error {
	return m.Viper.WriteConfigAs(filename)
}

// SafeWriteConfigAs writes the current configuration to file with specified
// name only if the file doesn't exist.
func (m *manager) SafeWriteConfigAs(filename string) error {
	return m.Viper.SafeWriteConfigAs(filename)
}

// WatchConfig watches the config file and reloads it when it changes.
func (m *manager) WatchConfig() {
	m.Viper.WatchConfig()
}

// OnConfigChange sets the callback to run when config changes.
func (m *manager) OnConfigChange(callback func(event fsnotify.Event)) {
	m.Viper.OnConfigChange(callback)
}

// SetEnvPrefix sets a prefix that ENV variables will use.
func (m *manager) SetEnvPrefix(prefix string) {
	m.Viper.SetEnvPrefix(prefix)
}

// AutomaticEnv enables automatic ENV variable support.
func (m *manager) AutomaticEnv() {
	m.Viper.AutomaticEnv()
}

// BindEnv binds a Viper key to a ENV variable.
func (m *manager) BindEnv(input ...string) error {
	return m.Viper.BindEnv(input...)
}

// MergeConfig merges a new config with the current config.
func (m *manager) MergeConfig(in io.Reader) error {
	return m.Viper.MergeConfig(in)
}
