// File manager.go định nghĩa các giao diện và cài đặt cốt lõi cho hệ thống quản lý cấu hình.
//
// File này bao gồm:
//   1. Giao diện Manager: Định nghĩa các phương thức API cho việc quản lý cấu hình
//   2. Cấu trúc DefaultManager: Cài đặt an toàn đa luồng của giao diện Manager
//   3. Các phương thức truy xuất và cập nhật cho việc CRUD và biến đổi giá trị cấu hình
//
// Kiến trúc tổng quan:
//   - Kho lưu trữ cấu hình tập trung với cấu trúc dữ liệu khóa-giá trị phẳng
//   - Mẫu truy cập dạng ký hiệu chấm (dot notation) cho dữ liệu phân cấp
//   - Các thao tác an toàn đa luồng với cơ chế khóa chi tiết
//   - Truy xuất an toàn về kiểu dữ liệu với khả năng chuyển đổi kiểu mạnh mẽ
//   - Kiến trúc định dạng có thể mở rộng cho nhiều nguồn cấu hình khác nhau

// Gói config cung cấp hệ thống quản lý cấu hình linh hoạt và dễ mở rộng.
//
// Hệ thống cho phép tổng hợp và chuẩn hóa dữ liệu cấu hình từ nhiều nguồn khác nhau
// (YAML, JSON, biến môi trường, v.v.) vào một kho lưu trữ tập trung, an toàn đa luồng.
// Các khóa cấu hình được quản lý dưới dạng ký hiệu chấm, hỗ trợ tra cứu phân cấp và
// tự động chuyển đổi cấu trúc lồng nhau.
//
// Các thành phần chính:
//   - Manager: Giao diện chính định nghĩa API truy cập cấu hình
//   - DefaultManager: Cài đặt an toàn đa luồng với khả năng tra cứu tối ưu
//   - Các phương thức truy xuất theo kiểu dữ liệu cụ thể (GetString, GetInt, GetBool, v.v.)
//   - Hỗ trợ chuyển đổi cấu hình thành struct với liên kết JSON
//
// Ví dụ sử dụng:
//   manager := config.NewManager()
//   manager.Load(formatter.NewYamlFormatter("config.yaml"))
//   appName := manager.GetString("app.name", "Default App")
//
// Xử lý ngoại lệ:
//   - Lỗi tải cấu hình: Lỗi đọc/ghi file, vấn đề xác thực định dạng
//   - Tình huống không tìm thấy khóa: Thay thế bằng giá trị mặc định
//   - Lỗi chuyển đổi kiểu dữ liệu: Sử dụng giá trị dự phòng an toàn
//
// Các khái niệm cốt lõi:

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/go-fork/providers/config/formatter"
	"github.com/go-fork/providers/config/utils"
)

// Manager là interface chính để quản lý cấu hình ứng dụng.
//
// Interface này định nghĩa các phương thức cần thiết để truy xuất, cập nhật, kiểm tra
// và nạp cấu hình từ nhiều nguồn khác nhau. Manager cung cấp các phương thức type-safe
// để lấy giá trị cấu hình với kiểu dữ liệu mong muốn và hỗ trợ giá trị mặc định.
//
// Các phương thức:
//   - Get: Lấy giá trị cấu hình theo key, có thể truyền defaultValue nếu không tồn tại.
//   - Set: Đặt giá trị cấu hình cho key.
//   - Has: Kiểm tra key có tồn tại không.
//   - GetString, GetInt, GetBool: Lấy giá trị cấu hình với kiểu dữ liệu tương ứng.
//   - GetStringMap, GetStringSlice: Lấy giá trị dạng map hoặc slice.
//   - Unmarshal: Chuyển đổi giá trị cấu hình sang struct.
//   - Load: Nạp cấu hình từ một Formatter (nguồn cấu hình).
type Manager interface {
	// Get lấy một giá trị cấu hình theo khóa.
	//
	// Phương thức này tìm kiếm giá trị cấu hình theo key được cung cấp.
	// Nếu key không tồn tại và defaultValue được cung cấp, trả về defaultValue đầu tiên.
	// Nếu key là key cha (ví dụ: "database"), phương thức sẽ gom các key con thành map lồng nhau.
	//
	// Params:
	//   - key: string - Tên khóa cấu hình (hỗ trợ dot notation, ví dụ: "database.host").
	//   - defaultValue: ...interface{} - Giá trị mặc định nếu key không tồn tại (tùy chọn).
	//
	// Returns:
	//   - interface{}: Giá trị cấu hình hoặc defaultValue nếu không tồn tại.
	//
	// Examples:
	//   manager.Get("app.name") // Trả về giá trị của key "app.name" hoặc nil nếu không tồn tại
	//   manager.Get("app.name", "Default App") // Trả về giá trị hoặc "Default App" nếu không tồn tại
	//   manager.Get("database") // Trả về map lồng nhau gồm tất cả key con của "database"
	Get(key string, defaultValue ...interface{}) interface{}

	// Set đặt một giá trị cấu hình cho key.
	//
	// Phương thức này lưu giá trị vào bộ nhớ với key tương ứng. Key hỗ trợ
	// dot notation để tổ chức cấu hình phân cấp.
	//
	// Params:
	//   - key: string - Tên khóa cấu hình (hỗ trợ dot notation).
	//   - value: interface{} - Giá trị cần đặt cho key.
	//
	// Examples:
	//   manager.Set("app.name", "My Application")
	//   manager.Set("database.host", "localhost")
	//   manager.Set("features.enabled", true)
	Set(key string, value interface{})

	// Has kiểm tra một khóa có tồn tại trong cấu hình hay không.
	//
	// Phương thức này kiểm tra key được cung cấp có tồn tại không. Đối với key cha,
	// nó cũng kiểm tra sự tồn tại của bất kỳ key con nào.
	//
	// Params:
	//   - key: string - Tên khóa cấu hình cần kiểm tra.
	//
	// Returns:
	//   - bool: true nếu key tồn tại (hoặc có key con), ngược lại false.
	//
	// Examples:
	//   manager.Has("app.name") // true nếu key "app.name" tồn tại
	//   manager.Has("database") // true nếu có bất kỳ key nào bắt đầu bằng "database."
	Has(key string) bool

	// GetString lấy một giá trị cấu hình dưới dạng string.
	//
	// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành string.
	// Hỗ trợ giá trị mặc định nếu key không tồn tại hoặc không thể chuyển đổi thành string.
	//
	// Params:
	//   - key: string - Tên khóa cấu hình.
	//   - defaultValue: ...string - Giá trị mặc định nếu key không tồn tại (tùy chọn).
	//
	// Returns:
	//   - string: Giá trị string hoặc defaultValue nếu không thể lấy giá trị.
	//
	// Examples:
	//   manager.GetString("app.name") // Trả về string hoặc "" nếu không tồn tại
	//   manager.GetString("app.name", "Default") // Trả về string hoặc "Default" nếu không tồn tại
	GetString(key string, defaultValue ...string) string

	// GetInt lấy một giá trị cấu hình dưới dạng int.
	//
	// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành int.
	// Hỗ trợ giá trị mặc định nếu key không tồn tại hoặc không thể chuyển đổi thành int.
	//
	// Params:
	//   - key: string - Tên khóa cấu hình.
	//   - defaultValue: ...int - Giá trị mặc định nếu key không tồn tại (tùy chọn).
	//
	// Returns:
	//   - int: Giá trị int hoặc defaultValue nếu không thể lấy giá trị.
	//
	// Examples:
	//   manager.GetInt("app.port") // Trả về int hoặc 0 nếu không tồn tại
	//   manager.GetInt("app.port", 8080) // Trả về int hoặc 8080 nếu không tồn tại
	GetInt(key string, defaultValue ...int) int

	// GetBool lấy một giá trị cấu hình dưới dạng bool.
	//
	// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành bool.
	// Hỗ trợ giá trị mặc định nếu key không tồn tại hoặc không thể chuyển đổi thành bool.
	// Đối với string, các giá trị "true", "yes", "1" được coi là true.
	//
	// Params:
	//   - key: string - Tên khóa cấu hình.
	//   - defaultValue: ...bool - Giá trị mặc định nếu key không tồn tại (tùy chọn).
	//
	// Returns:
	//   - bool: Giá trị bool hoặc defaultValue nếu không thể lấy giá trị.
	//
	// Examples:
	//   manager.GetBool("app.debug") // Trả về bool hoặc false nếu không tồn tại
	//   manager.GetBool("app.debug", true) // Trả về bool hoặc true nếu không tồn tại
	GetBool(key string, defaultValue ...bool) bool

	// GetStringMap lấy một giá trị cấu hình dưới dạng map[string]interface{}.
	//
	// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành map[string]interface{}.
	// Nếu key là key cha, nó sẽ gom tất cả key con thành một map lồng nhau.
	//
	// Params:
	//   - key: string - Tên khóa cấu hình.
	//
	// Returns:
	//   - map[string]interface{}: Map chứa các giá trị hoặc map rỗng nếu không tồn tại.
	//
	// Examples:
	//   manager.GetStringMap("database") // Trả về map chứa tất cả cấu hình database
	GetStringMap(key string) map[string]interface{}

	// GetStringSlice lấy một giá trị cấu hình dưới dạng []string.
	//
	// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành []string.
	// Hỗ trợ nhiều định dạng đầu vào khác nhau: []string, []interface{}, chuỗi JSON.
	//
	// Params:
	//   - key: string - Tên khóa cấu hình.
	//
	// Returns:
	//   - []string: Slice chuỗi hoặc slice rỗng nếu không tồn tại.
	//
	// Examples:
	//   manager.GetStringSlice("app.allowed_hosts") // Trả về []string từ cấu hình
	GetStringSlice(key string) []string

	// Unmarshal chuyển đổi giá trị cấu hình thành struct.
	//
	// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành struct
	// thông qua JSON marshaling/unmarshaling. Hỗ trợ đầy đủ các tag `json` trong struct.
	// Nếu key là key cha, nó sẽ gom tất cả key con thành một struct có cấu trúc phù hợp.
	//
	// Params:
	//   - key: string - Tên khóa cấu hình.
	//   - out: interface{} - Con trỏ đến struct đích.
	//
	// Returns:
	//   - error: Lỗi nếu key không tồn tại, không thể marshal/unmarshal, hoặc out không phải con trỏ.
	//
	// Examples:
	//   var dbConfig DatabaseConfig
	//   err := manager.Unmarshal("database", &dbConfig)
	Unmarshal(key string, out interface{}) error

	// Load nạp cấu hình từ một nguồn Formatter.
	//
	// Phương thức này sử dụng một đối tượng Formatter để nạp cấu hình từ nguồn bên ngoài
	// (như file YAML, JSON, biến môi trường, ...). Cấu hình mới được merge vào cấu hình hiện tại.
	//
	// Params:
	//   - formatter: formatter.Formatter - Đối tượng Formatter để nạp cấu hình.
	//
	// Returns:
	//   - error: Lỗi nếu không thể nạp cấu hình hoặc formatter là nil.
	//
	// Examples:
	//   err := manager.Load(formatter.NewYamlFormatter("config.yaml"))
	//   err := manager.Load(formatter.NewEnvFormatter("APP_"))
	Load(formatter formatter.Formatter) error
}

// DefaultManager là implementation mặc định của Manager.
//
// DefaultManager triển khai đầy đủ các phương thức của Manager một cách thread-safe
// sử dụng sync.RWMutex. Struct này lưu trữ cấu hình trong một map phẳng với key
// dạng dot notation và hỗ trợ truy vấn các key cha bằng cách gom các key con.
//
// Thuộc tính:
//   - formatters: Danh sách các formatter đã được nạp vào manager.
//   - values: Map lưu trữ các cặp key-value cấu hình.
//   - mu: Mutex để đảm bảo thread-safe khi truy xuất và cập nhật cấu hình.
type DefaultManager struct {
	formatters []formatter.Formatter  // Danh sách các nguồn cấu hình đã nạp
	values     map[string]interface{} // Map lưu trữ các giá trị cấu hình
	mu         sync.RWMutex           // Mutex để đảm bảo thread-safe
}

// NewManager tạo và trả về một đối tượng DefaultManager mới.
//
// Phương thức này khởi tạo một DefaultManager mới với slice formatters
// và map values rỗng, sẵn sàng để sử dụng.
//
// Returns:
//   - Manager: Một đối tượng Manager mới (thực tế là *DefaultManager).
//
// Examples:
//
//	manager := config.NewManager()
//	manager.Set("app.name", "My App")
func NewManager() Manager {
	return &DefaultManager{
		formatters: make([]formatter.Formatter, 0),
		values:     make(map[string]interface{}),
	}
}

// Get lấy một giá trị cấu hình theo key.
//
// Phương thức này tìm kiếm giá trị trong map cấu hình theo key được cung cấp.
// Nếu key tồn tại chính xác, trả về giá trị tương ứng.
// Nếu key không tồn tại:
//   - Nếu key là key cha (ví dụ "database"), phương thức sẽ tự động gom các key con
//     (như "database.host", "database.port") thành một map lồng nhau.
//   - Nếu không tìm thấy key và key con, trả về defaultValue nếu có, ngược lại trả về nil.
//
// Thread-safe nhờ sử dụng RWMutex để đảm bảo an toàn khi đọc cấu hình.
//
// Params:
//   - key: string - Tên khóa cấu hình cần truy vấn.
//   - defaultValue: ...interface{} - Giá trị mặc định nếu key không tồn tại.
//
// Returns:
//   - interface{}: Giá trị cấu hình, map lồng nhau (nếu là key cha),
//     defaultValue (nếu key không tồn tại), hoặc nil.
func (m *DefaultManager) Get(key string, defaultValue ...interface{}) interface{} {
	if key == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Tìm key chính xác
	if val, ok := m.values[key]; ok {
		return val
	}

	// Gom các key con nếu có
	prefix := key + "."
	result := make(map[string]interface{})
	for k, v := range m.values {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			subKey := k[len(prefix):]
			result[subKey] = v
		}
	}
	if len(result) > 0 {
		return utils.UnflattenDotMapRecursive(result)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

// Set đặt một giá trị cấu hình cho key.
//
// Phương thức này lưu trữ một giá trị trong map cấu hình với key tương ứng.
// Nếu key đã tồn tại, giá trị sẽ được ghi đè. Key hỗ trợ dot notation để
// tổ chức cấu hình phân cấp (ví dụ: "database.host", "database.port").
//
// Thread-safe nhờ sử dụng Mutex để đảm bảo an toàn khi cập nhật cấu hình.
//
// Params:
//   - key: string - Tên khóa cấu hình để lưu giá trị.
//   - value: interface{} - Giá trị cần lưu trữ.
func (m *DefaultManager) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[key] = value
}

// Has kiểm tra một key có tồn tại trong cấu hình hay không.
//
// Phương thức này kiểm tra xem key được cung cấp có tồn tại trong map cấu hình không.
// Đối với key cha (ví dụ: "database"), phương thức sẽ kiểm tra sự tồn tại của bất kỳ
// key con nào có tiền tố là key đó (như "database.host", "database.port").
//
// Thread-safe nhờ sử dụng RWMutex để đảm bảo an toàn khi đọc cấu hình.
//
// Params:
//   - key: string - Tên khóa cấu hình cần kiểm tra.
//
// Returns:
//   - bool: true nếu key tồn tại hoặc có key con, false nếu không.
func (m *DefaultManager) Has(key string) bool {
	if key == "" {
		return false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if _, ok := m.values[key]; ok {
		return true
	}
	prefix := key + "."
	for k := range m.values {
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// GetString lấy giá trị cấu hình dưới dạng string.
//
// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành string.
// Hỗ trợ các kiểu dữ liệu:
//   - string: Trả về trực tiếp
//   - fmt.Stringer: Gọi phương thức String()
//
// Nếu key không tồn tại hoặc không thể chuyển đổi thành string, trả về
// defaultValue nếu có, ngược lại trả về chuỗi rỗng "".
//
// Params:
//   - key: string - Tên khóa cấu hình.
//   - defaultValue: ...string - Giá trị mặc định nếu key không tồn tại (tùy chọn).
//
// Returns:
//   - string: Giá trị string, defaultValue, hoặc chuỗi rỗng.
func (m *DefaultManager) GetString(key string, defaultValue ...string) string {
	val := m.Get(key)
	switch v := val.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// GetInt lấy giá trị cấu hình dưới dạng int.
//
// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành int.
// Hỗ trợ nhiều kiểu dữ liệu:
//   - int, int32, int64: Chuyển đổi sang int
//   - float32, float64: Chuyển đổi sang int (cắt phần thập phân)
//   - string: Parse thành int nếu có thể
//
// Nếu key không tồn tại hoặc không thể chuyển đổi thành int, trả về
// defaultValue nếu có, ngược lại trả về 0.
//
// Params:
//   - key: string - Tên khóa cấu hình.
//   - defaultValue: ...int - Giá trị mặc định nếu key không tồn tại (tùy chọn).
//
// Returns:
//   - int: Giá trị int, defaultValue, hoặc 0.
func (m *DefaultManager) GetInt(key string, defaultValue ...int) int {
	val := m.Get(key)
	if val == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		var intVal int
		if _, err := fmt.Sscanf(v, "%d", &intVal); err == nil {
			return intVal
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

// GetBool lấy giá trị cấu hình dưới dạng bool.
//
// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành bool.
// Hỗ trợ các kiểu dữ liệu:
//   - bool: Trả về trực tiếp
//   - string: Chuyển đổi dựa trên các quy tắc sau:
//   - "true", "1", "yes", "on" -> true
//   - "false", "0", "no", "off" -> false
//
// Nếu key không tồn tại hoặc không thể chuyển đổi thành bool, trả về
// defaultValue nếu có, ngược lại trả về false.
//
// Params:
//   - key: string - Tên khóa cấu hình.
//   - defaultValue: ...bool - Giá trị mặc định nếu key không tồn tại (tùy chọn).
//
// Returns:
//   - bool: Giá trị bool, defaultValue, hoặc false.
func (m *DefaultManager) GetBool(key string, defaultValue ...bool) bool {
	val := m.Get(key)
	switch v := val.(type) {
	case bool:
		return v
	case string:
		if v == "true" || v == "1" || v == "yes" || v == "on" {
			return true
		}
		if v == "false" || v == "0" || v == "no" || v == "off" {
			return false
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

// GetStringMap lấy giá trị cấu hình dưới dạng map[string]interface{}.
//
// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành map[string]interface{}.
// Hỗ trợ nhiều kiểu dữ liệu:
//   - map[string]interface{}: Trả về trực tiếp
//   - map[interface{}]interface{}: Chuyển đổi key thành string khi có thể
//   - Các kiểu khác: Thử chuyển đổi bằng JSON marshaling/unmarshaling
//
// Nếu key là key cha, nó sẽ gom tất cả key con thành một map lồng nhau.
// Nếu key không tồn tại hoặc không thể chuyển đổi thành map, trả về map rỗng.
//
// Params:
//   - key: string - Tên khóa cấu hình.
//
// Returns:
//   - map[string]interface{}: Map chứa các giá trị hoặc map rỗng.
func (m *DefaultManager) GetStringMap(key string) map[string]interface{} {
	if key == "" {
		return make(map[string]interface{})
	}

	val := m.Get(key)
	if val == nil {
		return make(map[string]interface{})
	}

	switch v := val.(type) {
	case map[string]interface{}:
		return v
	case map[interface{}]interface{}:
		// Chuyển đổi map[interface{}]interface{} thành map[string]interface{}
		result := make(map[string]interface{}, len(v))
		for mk, mv := range v {
			if strKey, ok := mk.(string); ok {
				result[strKey] = mv
			}
		}
		return result
	}

	// Thử chuyển đổi bằng JSON marshalling/unmarshaling
	jsonData, err := json.Marshal(val)
	if err == nil {
		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err == nil {
			return result
		}
	}

	return make(map[string]interface{})
}

// GetStringSlice lấy giá trị cấu hình dưới dạng []string.
//
// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành []string.
// Hỗ trợ nhiều định dạng đầu vào khác nhau:
//   - []string: Trả về trực tiếp
//   - []interface{}: Chuyển đổi mỗi phần tử thành string
//   - string: Thử parse như JSON array, nếu không được thì trả về slice với 1 phần tử
//   - []byte: Thử parse như JSON array hoặc []interface{}
//
// Phương thức cố gắng chuyển đổi các kiểu dữ liệu khác thành string slice
// thông qua việc sử dụng JSON marshaling/unmarshaling khi cần thiết.
//
// Nếu key không tồn tại hoặc không thể chuyển đổi thành []string, trả về slice rỗng.
//
// Params:
//   - key: string - Tên khóa cấu hình.
//
// Returns:
//   - []string: Slice chuỗi hoặc slice rỗng.
func (m *DefaultManager) GetStringSlice(key string) []string {
	if key == "" {
		return []string{}
	}

	val := m.Get(key)
	if val == nil {
		return []string{}
	}

	// Kiểm tra các trường hợp chính xác
	switch v := val.(type) {
	case []string:
		return v
	case []interface{}:
		// Chuyển đổi []interface{} thành []string
		result := make([]string, 0, len(v))
		for _, item := range v {
			switch sv := item.(type) {
			case string:
				result = append(result, sv)
			case fmt.Stringer:
				result = append(result, sv.String())
			default:
				// Chuyển đổi các kiểu khác thành string
				result = append(result, fmt.Sprintf("%v", item))
			}
		}
		return result
	case string:
		// Thử parse string như JSON array
		var result []string
		if err := json.Unmarshal([]byte(v), &result); err == nil {
			return result
		}
		// Nếu không phải JSON array, trả về slice với 1 phần tử
		return []string{v}
	}

	// Thử chuyển đổi bằng JSON
	jsonData, err := json.Marshal(val)
	if err == nil {
		var result []string
		if err := json.Unmarshal(jsonData, &result); err == nil {
			return result
		}
		// Nếu không thể unmarshal thành []string, thử unmarshal thành []interface{}
		var resultInterface []interface{}
		if err := json.Unmarshal(jsonData, &resultInterface); err == nil {
			result := make([]string, 0, len(resultInterface))
			for _, item := range resultInterface {
				if str, ok := item.(string); ok {
					result = append(result, str)
				} else {
					result = append(result, fmt.Sprintf("%v", item))
				}
			}
			return result
		}
		// Nếu là struct, trả về []string với 1 phần tử là fmt.Sprintf("%v", val)
		if reflect.TypeOf(val).Kind() == reflect.Struct {
			return []string{fmt.Sprintf("%v", val)}
		}
	}

	// Kiểm tra nếu là []byte thì thử parse như JSON array hoặc interface
	if b, ok := val.([]byte); ok {
		var result []string
		if err := json.Unmarshal(b, &result); err == nil {
			return result
		}
		var resultInterface []interface{}
		if err := json.Unmarshal(b, &resultInterface); err == nil {
			result := make([]string, 0, len(resultInterface))
			for _, item := range resultInterface {
				if str, ok := item.(string); ok {
					result = append(result, str)
				} else {
					result = append(result, fmt.Sprintf("%v", item))
				}
			}
			return result
		}
		// Nếu không phải JSON, trả về slice với 1 phần tử là fmt.Sprintf("%v", val)
		return []string{fmt.Sprintf("%v", val)}
	}
	return []string{}
}

// Unmarshal chuyển đổi giá trị cấu hình thành struct.
//
// Phương thức này tìm kiếm giá trị theo key và chuyển đổi nó thành struct
// thông qua JSON marshaling/unmarshaling. Hỗ trợ đầy đủ các tag `json` trong struct.
//
// Quy trình xử lý:
//  1. Kiểm tra tính hợp lệ của tham số out (phải là non-nil pointer)
//  2. Tìm kiếm giá trị theo key (chính xác hoặc gom các key con)
//  3. Chuyển đổi giá trị thành dạng JSON
//  4. Giải mã JSON vào struct đích
//
// Nếu key là key cha, nó sẽ gom tất cả key con thành một struct có cấu trúc phù hợp.
//
// Thread-safe nhờ sử dụng RWMutex để đảm bảo an toàn khi đọc cấu hình.
//
// Params:
//   - key: string - Tên khóa cấu hình.
//   - out: interface{} - Con trỏ đến struct đích.
//
// Returns:
//   - error: Lỗi nếu key không tồn tại, không thể marshal/unmarshal, hoặc out không phải con trỏ.
//
// Exceptions:
//   - "output pointer cannot be nil": out là nil
//   - "output must be a non-nil pointer": out không phải con trỏ hoặc là con trỏ nil
//   - "key 'x' not found in configuration": key không tồn tại trong cấu hình
//   - "failed to marshal configuration": không thể chuyển đổi giá trị thành JSON
//   - "failed to unmarshal configuration": không thể giải mã JSON vào struct
func (m *DefaultManager) Unmarshal(key string, out interface{}) error {
	if out == nil {
		return fmt.Errorf("output pointer cannot be nil")
	}

	// Kiểm tra xem out có phải là con trỏ không
	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr || outValue.IsNil() {
		return fmt.Errorf("output must be a non-nil pointer")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var val interface{}
	var exists bool

	if key == "" {
		val = m.values
		exists = true
	} else if v, ok := m.values[key]; ok {
		val = v
		exists = true
	} else {
		// Gom các key con
		prefix := key + "."
		result := make(map[string]interface{})
		for k, v := range m.values {
			if len(k) > len(prefix) && k[:len(prefix)] == prefix {
				subKey := k[len(prefix):]
				result[subKey] = v
			}
		}
		if len(result) > 0 {
			val = utils.UnflattenDotMapRecursive(result)
			exists = true
		}
	}

	if !exists {
		return fmt.Errorf("key '%s' not found in configuration", key)
	}

	jsonData, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}
	if err := json.Unmarshal(jsonData, out); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}
	return nil
}

// Load nạp cấu hình từ một nguồn thông qua Formatter.
//
// Phương thức này thực hiện các bước sau:
//  1. Kiểm tra tính hợp lệ của formatter (không được nil)
//  2. Gọi phương thức Load của formatter để lấy các giá trị cấu hình
//  3. Thêm formatter vào danh sách formatters đã nạp
//  4. Merge các cặp key-value nhận được vào map values
//
// Đây là phương thức chính để nạp cấu hình từ các nguồn như file YAML, JSON,
// biến môi trường, hoặc bất kỳ nguồn dữ liệu nào khác thông qua các formatter.
//
// Thread-safe nhờ sử dụng Mutex để đảm bảo an toàn khi cập nhật cấu hình.
//
// Params:
//   - formatter: formatter.Formatter - Đối tượng Formatter để nạp cấu hình.
//
// Returns:
//   - error: Lỗi nếu formatter là nil, hoặc formatter.Load() thất bại.
//
// Exceptions:
//   - "formatter cannot be nil": formatter là nil
//   - "failed to load config from X": Lỗi khi gọi formatter.Load()
//
// Examples:
//
//	err := manager.Load(formatter.NewYamlFormatter("config/app.yaml"))
//	err := manager.Load(formatter.NewEnvFormatter("APP_"))
func (m *DefaultManager) Load(formatter formatter.Formatter) error {
	if formatter == nil {
		return errors.New("formatter cannot be nil")
	}

	values, err := formatter.Load()
	if err != nil {
		return fmt.Errorf("failed to load config from %s: %w", formatter.Name(), err)
	}

	if values == nil {
		return nil // Không có gì để nạp
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Thêm formatter vào danh sách đã nạp
	m.formatters = append(m.formatters, formatter)

	// Merge giá trị vào map values
	for k, v := range values {
		if k != "" { // Bỏ qua key rỗng
			m.values[k] = v
		}
	}

	return nil
}
