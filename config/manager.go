package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/go-fork/providers/config/formatter"
)

// Package config cung cấp khả năng quản lý cấu hình linh hoạt cho ứng dụng.
//
// Cho phép tải và truy xuất cấu hình từ nhiều nguồn khác nhau như file (YAML, JSON), biến môi trường, ...
// Các giá trị cấu hình được quản lý thông qua một Manager trung tâm, hỗ trợ truy vấn theo dot notation.
//
// Các thành phần chính:
//   - Manager: Interface quản lý cấu hình tổng quát
//   - DefaultManager: Implementation mặc định của Manager
//   - Các method hỗ trợ truy xuất, cập nhật, kiểm tra, nạp cấu hình
//
// Ví dụ sử dụng:
//   manager := config.NewManager()
//   manager.Load(formatter.NewYamlFormatter("config.yaml"))
//   appName := manager.GetString("app.name", "Default App")
//
// Các exception có thể phát sinh khi nạp file cấu hình hoặc khi key không tồn tại.
//
// Định nghĩa các interface và struct:

// Manager là interface chính để quản lý cấu hình ứng dụng.
//
// Cho phép truy xuất, cập nhật, kiểm tra, nạp cấu hình từ nhiều nguồn khác nhau.
// Các method:
//   - Get: Lấy giá trị cấu hình theo key, có thể truyền defaultValue nếu không tồn tại.
//   - Set: Đặt giá trị cấu hình cho key.
//   - Has: Kiểm tra key có tồn tại không.
//   - GetString, GetInt, GetBool: Lấy giá trị cấu hình với kiểu dữ liệu tương ứng.
//   - GetStringMap, GetStringSlice: Lấy giá trị dạng map hoặc slice.
//   - Unmarshal: Chuyển đổi giá trị cấu hình sang struct.
//   - Load: Nạp cấu hình từ một Formatter (nguồn cấu hình).
type Manager interface {
	// Get lấy một giá trị cấu hình theo khóa.
	// Params:
	//   - key: Tên khóa cấu hình (dot notation).
	//   - defaultValue: Giá trị mặc định nếu không tồn tại (tùy chọn).
	// Returns: Giá trị cấu hình hoặc defaultValue nếu không có.
	Get(key string, defaultValue ...interface{}) interface{}

	// Set đặt một giá trị cấu hình.
	// Params:
	//   - key: Tên khóa.
	//   - value: Giá trị cần đặt.
	Set(key string, value interface{})

	// Has kiểm tra xem một khóa có tồn tại không.
	// Params:
	//   - key: Tên khóa.
	// Returns: true nếu tồn tại, false nếu không.
	Has(key string) bool

	// GetString lấy một giá trị cấu hình dưới dạng string.
	// Params:
	//   - key: Tên khóa.
	//   - defaultValue: Giá trị mặc định nếu không tồn tại (tùy chọn).
	// Returns: Giá trị string hoặc defaultValue nếu không có.
	GetString(key string, defaultValue ...string) string

	// GetInt lấy một giá trị cấu hình dưới dạng int.
	// Params:
	//   - key: Tên khóa.
	//   - defaultValue: Giá trị mặc định nếu không tồn tại (tùy chọn).
	// Returns: Giá trị int hoặc defaultValue nếu không có.
	GetInt(key string, defaultValue ...int) int

	// GetBool lấy một giá trị cấu hình dưới dạng bool.
	// Params:
	//   - key: Tên khóa.
	//   - defaultValue: Giá trị mặc định nếu không tồn tại (tùy chọn).
	// Returns: Giá trị bool hoặc defaultValue nếu không có.
	GetBool(key string, defaultValue ...bool) bool

	// GetStringMap lấy một giá trị cấu hình dưới dạng map[string]interface{}.
	// Params:
	//   - key: Tên khóa.
	// Returns: map[string]interface{} hoặc map rỗng nếu không có.
	GetStringMap(key string) map[string]interface{}

	// GetStringSlice lấy một giá trị cấu hình dưới dạng []string.
	// Params:
	//   - key: Tên khóa.
	// Returns: []string hoặc slice rỗng nếu không có.
	GetStringSlice(key string) []string

	// Unmarshal chuyển đổi giá trị cấu hình sang struct.
	// Params:
	//   - key: Tên khóa.
	//   - out: Con trỏ struct đích.
	// Returns: error nếu không chuyển đổi được hoặc key không tồn tại.
	Unmarshal(key string, out interface{}) error

	// Load nạp cấu hình từ một nguồn Formatter.
	// Params:
	//   - formatter: Đối tượng Formatter (nguồn cấu hình).
	// Returns: error nếu không nạp được.
	Load(formatter formatter.Formatter) error
}

// DefaultManager là implementation mặc định của Manager.
//
// Lưu trữ các giá trị cấu hình, hỗ trợ thread-safe với sync.RWMutex.
// Có thể nạp nhiều nguồn Formatter, merge giá trị vào values.
type DefaultManager struct {
	formatters []formatter.Formatter  // Danh sách các nguồn cấu hình đã nạp
	values     map[string]interface{} // Map lưu trữ các giá trị cấu hình
	mu         sync.RWMutex           // Mutex để đảm bảo thread-safe
}

// NewManager tạo một DefaultManager mới.
// Returns: Đối tượng Manager.
func NewManager() Manager {
	return &DefaultManager{
		formatters: make([]formatter.Formatter, 0),
		values:     make(map[string]interface{}),
	}
}

// Get lấy một giá trị cấu hình theo key.
// Nếu không tồn tại trả về defaultValue nếu có, ngược lại trả về nil.
// Thread-safe.
func (m *DefaultManager) Get(key string, defaultValue ...interface{}) interface{} {
	if key == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.values[key]
	if ok {
		return val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

// Set đặt một giá trị cấu hình cho key.
// Thread-safe.
func (m *DefaultManager) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[key] = value
}

// Has kiểm tra key có tồn tại không.
// Returns: true nếu có, false nếu không.
// Thread-safe.
func (m *DefaultManager) Has(key string) bool {
	if key == "" {
		return false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.values[key]
	return ok
}

// GetString lấy giá trị cấu hình dạng string.
// Nếu không tồn tại trả về defaultValue nếu có, ngược lại trả về "".
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

// GetInt lấy giá trị cấu hình dạng int.
// Nếu không tồn tại trả về defaultValue nếu có, ngược lại trả về 0.
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

// GetBool lấy giá trị cấu hình dạng bool.
// Nếu không tồn tại trả về defaultValue nếu có, ngược lại trả về false.
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

// GetStringMap lấy giá trị cấu hình dạng map[string]interface{}.
// Nếu không tồn tại trả về map rỗng.
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

	// Thử chuyển đổi bằng JSON marshalling/unmarshalling
	jsonData, err := json.Marshal(val)
	if err == nil {
		var result map[string]interface{}
		if err := json.Unmarshal(jsonData, &result); err == nil {
			return result
		}
	}

	return make(map[string]interface{})
}

// GetStringSlice lấy giá trị cấu hình dạng []string.
// Nếu không tồn tại trả về slice rỗng.
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

// Unmarshal chuyển đổi giá trị cấu hình sang struct.
// Sử dụng encoding/json để chuyển đổi map sang struct.
// Returns: error nếu không chuyển đổi được hoặc key không tồn tại.
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
		// Nếu key rỗng, unmarshal toàn bộ cấu hình
		val = m.values
		exists = true
	} else {
		// Unmarshal một phần cấu hình dựa trên key
		val, exists = m.values[key]
	}

	if !exists {
		return fmt.Errorf("key '%s' not found in configuration", key)
	}

	// Sử dụng JSON marshalling để chuyển đổi
	jsonData, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := json.Unmarshal(jsonData, out); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return nil
}

// Load nạp cấu hình từ một Formatter.
// Thêm Formatter vào danh sách, merge các giá trị vào values.
// Returns: error nếu không nạp được.
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
