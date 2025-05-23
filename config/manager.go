package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-fork/providers/config/formatter"
	"github.com/go-fork/providers/config/model"
	"github.com/go-fork/providers/config/utils"
)

// Định nghĩa các lỗi tiêu chuẩn
var (
	// ErrInvalidKey là lỗi khi key không hợp lệ (rỗng, ký tự đặc biệt)
	ErrInvalidKey = errors.New("invalid configuration key")

	// ErrKeyNotFound là lỗi khi key không tồn tại
	ErrKeyNotFound = errors.New("configuration key not found")

	// ErrTypeMismatch là lỗi khi kiểu dữ liệu không khớp
	ErrTypeMismatch = errors.New("configuration type mismatch")

	// ErrInvalidTarget là lỗi khi target không phải con trỏ hoặc struct
	ErrInvalidTarget = errors.New("invalid unmarshal target")

	// ErrParseFailed là lỗi khi parse dữ liệu thất bại
	ErrParseFailed = errors.New("configuration parse failed")
)

// Config là interface chính cho việc quản lý cấu hình.
type Config interface {
	// GetString trả về giá trị chuỗi cho key.
	GetString(key string) (string, bool)

	// GetInt trả về giá trị số nguyên cho key.
	GetInt(key string) (int, bool)

	// GetBool trả về giá trị boolean cho key.
	GetBool(key string) (bool, bool)

	// GetFloat trả về giá trị số thực cho key.
	GetFloat(key string) (float64, bool)

	// GetSlice trả về giá trị slice cho key.
	GetSlice(key string) ([]interface{}, bool)

	// GetMap trả về giá trị map cho key.
	GetMap(key string) (map[string]interface{}, bool)

	// Set cập nhật hoặc thêm một giá trị vào cấu hình.
	Set(key string, value interface{}) error

	// Has kiểm tra xem key có tồn tại hay không.
	Has(key string) bool

	// Unmarshal ánh xạ cấu hình vào một struct Go.
	Unmarshal(key string, target interface{}) error

	// Load tải cấu hình từ một formatter.
	Load(f formatter.Formatter) error
}

// DefaultConfig triển khai interface Config.
type DefaultConfig struct {
	formatters []formatter.Formatter    // Danh sách formatter
	configMap  model.ConfigMap          // Map phẳng
	mu         sync.RWMutex             // Thread-safe
	opts       formatter.FlattenOptions // Tùy chọn flatten
}

// NewConfig tạo một đối tượng Config mới.
func NewConfig() *DefaultConfig {
	return &DefaultConfig{
		formatters: []formatter.Formatter{},
		configMap:  make(model.ConfigMap),
		opts:       formatter.DefaultFlattenOptions(),
	}
}

// Load tải cấu hình từ một formatter, ghi đè các giá trị hiện có.
func (c *DefaultConfig) Load(f formatter.Formatter) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Load dữ liệu từ formatter
	data, err := f.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Parse dữ liệu
	parsedData, err := f.Parse(data)
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Flatten dữ liệu
	configMap, err := f.Flatten(parsedData, c.opts)
	if err != nil {
		return fmt.Errorf("failed to flatten configuration: %w", err)
	}

	// Merge với configMap hiện tại
	for key, value := range configMap {
		c.configMap[key] = value
	}

	return nil
}

// LoadAll tải cấu hình từ tất cả formatter đã đăng ký theo thứ tự ưu tiên.
func (c *DefaultConfig) LoadAll() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Xóa configMap hiện tại
	c.configMap = make(model.ConfigMap)

	// Duyệt qua các formatter theo thứ tự ưu tiên từ thấp đến cao
	for i := len(c.formatters) - 1; i >= 0; i-- {
		f := c.formatters[i]

		// Load dữ liệu
		data, err := f.Load()
		if err != nil {
			continue // Bỏ qua formatter nếu lỗi
		}

		// Parse dữ liệu
		parsedData, err := f.Parse(data)
		if err != nil {
			continue
		}

		// Flatten dữ liệu
		configMap, err := f.Flatten(parsedData, c.opts)
		if err != nil {
			continue
		}

		// Merge với configMap
		for key, value := range configMap {
			c.configMap[key] = value
		}
	}

	return nil
}

// GetString trả về giá trị chuỗi cho key.
func (c *DefaultConfig) GetString(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.configMap[key]; ok {
		switch val.Type {
		case model.TypeString:
			str, _ := val.Value.(string)
			return str, true
		case model.TypeInt, model.TypeFloat, model.TypeBool:
			str, _ := utils.ToString(val.Value)
			return str, true
		default:
			return "", false
		}
	}

	return "", false
}

// GetInt trả về giá trị số nguyên cho key.
func (c *DefaultConfig) GetInt(key string) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.configMap[key]; ok {
		switch val.Type {
		case model.TypeInt:
			if intVal, ok := val.Value.(int); ok {
				return intVal, true
			}
			if int64Val, ok := val.Value.(int64); ok {
				return int(int64Val), true
			}
		case model.TypeFloat:
			if floatVal, ok := val.Value.(float64); ok {
				return int(floatVal), true
			}
		case model.TypeString:
			if strVal, ok := val.Value.(string); ok {
				if intVal, err := utils.ToInt(strVal); err == nil {
					return intVal, true
				}
			}
		}
	}

	return 0, false
}

// GetBool trả về giá trị boolean cho key.
func (c *DefaultConfig) GetBool(key string) (bool, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.configMap[key]; ok {
		switch val.Type {
		case model.TypeBool:
			boolVal, _ := val.Value.(bool)
			return boolVal, true
		case model.TypeInt:
			if intVal, ok := val.Value.(int); ok {
				return intVal != 0, true
			}
			if int64Val, ok := val.Value.(int64); ok {
				return int64Val != 0, true
			}
		case model.TypeString:
			if strVal, ok := val.Value.(string); ok {
				if boolVal, err := utils.ToBool(strVal); err == nil {
					return boolVal, true
				}
			}
		}
	}

	return false, false
}

// GetFloat trả về giá trị số thực cho key.
func (c *DefaultConfig) GetFloat(key string) (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.configMap[key]; ok {
		switch val.Type {
		case model.TypeFloat:
			floatVal, _ := val.Value.(float64)
			return floatVal, true
		case model.TypeInt:
			if intVal, ok := val.Value.(int); ok {
				return float64(intVal), true
			}
			if int64Val, ok := val.Value.(int64); ok {
				return float64(int64Val), true
			}
		case model.TypeBool:
			if boolVal, ok := val.Value.(bool); ok {
				if boolVal {
					return 1.0, true
				}
				return 0.0, true
			}
		case model.TypeString:
			if strVal, ok := val.Value.(string); ok {
				if floatVal, err := utils.ToFloat(strVal); err == nil {
					return floatVal, true
				}
			}
		}
	}

	return 0, false
}

// GetSlice trả về giá trị slice cho key.
func (c *DefaultConfig) GetSlice(key string) ([]interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.configMap[key]; ok && val.Type == model.TypeSlice {
		if slice, ok := val.Value.([]interface{}); ok {
			return slice, true
		}
	}

	return nil, false
}

// GetMap trả về giá trị map cho key.
func (c *DefaultConfig) GetMap(key string) (map[string]interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.configMap[key]; ok && val.Type == model.TypeMap {
		if m, ok := val.Value.(map[string]interface{}); ok {
			return m, true
		}
	}

	return nil, false
}

// Set cập nhật hoặc thêm một giá trị vào cấu hình.
func (c *DefaultConfig) Set(key string, value interface{}) error {
	if key == "" && !c.opts.HandleEmptyKey {
		return fmt.Errorf("%w: key cannot be empty", ErrInvalidKey)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	valueType := model.TypeUnknown

	// Xác định kiểu dữ liệu của value
	switch value.(type) {
	case string:
		valueType = model.TypeString
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		valueType = model.TypeInt
	case float32, float64:
		valueType = model.TypeFloat
	case bool:
		valueType = model.TypeBool
	case []interface{}, []string, []int:
		valueType = model.TypeSlice
	case map[string]interface{}:
		valueType = model.TypeMap
	case nil:
		valueType = model.TypeNil
	default:
		// Kiểm tra kiểu phức tạp hơn
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			valueType = model.TypeSlice
		case reflect.Map:
			valueType = model.TypeMap
		default:
			// Chuyển đổi thành string nếu không xác định được kiểu
			value = fmt.Sprintf("%v", value)
			valueType = model.TypeString
		}
	}

	// Cập nhật ConfigMap
	c.configMap[key] = model.ConfigValue{
		Value: value,
		Type:  valueType,
	}

	return nil
}

// Has kiểm tra xem key có tồn tại hay không.
func (c *DefaultConfig) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.configMap[key]
	return ok
}

// Unmarshal ánh xạ cấu hình vào một struct Go.
func (c *DefaultConfig) Unmarshal(key string, target interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer: %w", ErrInvalidTarget)
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("target must point to a struct: %w", ErrInvalidTarget)
	}

	var data model.ConfigMap
	if key == "" {
		// Sử dụng toàn bộ ConfigMap
		data = c.configMap
	} else {
		// Lọc ConfigMap theo prefix
		data = make(model.ConfigMap)
		if val, exists := c.configMap[key]; exists {
			data[key] = val
		}

		// Thêm tất cả key con
		prefix := key + c.opts.Separator
		for k, v := range c.configMap {
			if strings.HasPrefix(k, prefix) {
				data[k] = v
			}
		}

		if len(data) == 0 {
			return fmt.Errorf("key %s not found: %w", key, ErrKeyNotFound)
		}
	}

	return c.unmarshalStruct(rv, data, key)
}

// toSnakeCase chuyển CamelCase/PascalCase về snake_case
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// unmarshalStruct ánh xạ cấu hình vào một struct.
func (c *DefaultConfig) unmarshalStruct(rv reflect.Value, data model.ConfigMap, prefix string) error {
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		fieldValue := rv.Field(i)

		// Bỏ qua các field không export
		if !fieldValue.CanSet() {
			continue
		}

		// Lấy tag "config" hoặc sử dụng tên field
		configTag := field.Tag.Get("config")
		if configTag == "-" {
			continue // Bỏ qua field nếu tag là "-"
		}

		fieldPath := configTag
		if fieldPath == "" {
			fieldPath = toSnakeCase(field.Name)
		}

		// Nếu có prefix, thêm vào path
		if prefix != "" && fieldPath != "" {
			fieldPath = prefix + c.opts.Separator + fieldPath
		}

		// --- Ưu tiên ánh xạ map/slice gốc nếu là struct/map/slice và không có tag config ---
		if configTag == "" && (fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Map || fieldValue.Kind() == reflect.Slice) {
			baseKey := toSnakeCase(field.Name)
			if prefix != "" {
				baseKey = prefix + c.opts.Separator + baseKey
			}
			if configVal, ok := data[baseKey]; ok && configVal.Type == model.TypeMap {
				if mapData, ok := configVal.Value.(map[string]interface{}); ok {
					// Nếu là struct, ánh xạ từng field con từ mapData
					if fieldValue.Kind() == reflect.Struct {
						subMap := make(model.ConfigMap)
						for k, v := range mapData {
							subMap[k] = model.ConfigValue{Value: v, Type: model.TypeUnknown}
						}
						if err := c.unmarshalStruct(fieldValue, subMap, ""); err != nil {
							return err
						}
						continue
					}
					// Nếu là slice, lấy trường con từ mapData
					if fieldValue.Kind() == reflect.Slice {
						sliceField := toSnakeCase(field.Name)
						if v, ok := mapData[sliceField]; ok {
							if sliceData, ok := v.([]interface{}); ok {
								if err := c.unmarshalSlice(fieldValue, sliceData, baseKey+c.opts.Separator+sliceField); err != nil {
									return err
								}
								continue
							}
						}
					}
					// Nếu là map, ánh xạ tiếp
					if fieldValue.Kind() == reflect.Map {
						if err := c.unmarshalMap(fieldValue, mapData, baseKey); err != nil {
							return err
						}
						continue
					}
				}
			}
		}

		// Xử lý struct lồng nhau (fallback dot notation)
		if fieldValue.Kind() == reflect.Struct {
			if err := c.unmarshalStruct(fieldValue, data, fieldPath); err != nil {
				return err
			}
			continue
		}

		// Xử lý slice
		if fieldValue.Kind() == reflect.Slice {
			if configVal, ok := data[fieldPath]; ok && configVal.Type == model.TypeSlice {
				if err := c.unmarshalSlice(fieldValue, configVal.Value, fieldPath); err != nil {
					return err
				}
			}
			continue
		}

		// Xử lý map
		if fieldValue.Kind() == reflect.Map {
			if configVal, ok := data[fieldPath]; ok && configVal.Type == model.TypeMap {
				if err := c.unmarshalMap(fieldValue, configVal.Value, fieldPath); err != nil {
					return err
				}
			}
			continue
		}

		// Xử lý kiểu dữ liệu cơ bản (string, int, bool, ...)
		if configVal, ok := data[fieldPath]; ok {
			if err := c.setField(fieldValue, configVal.Value, fieldPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// unmarshalSlice ánh xạ cấu hình vào một slice.
func (c *DefaultConfig) unmarshalSlice(field reflect.Value, value interface{}, path string) error {
	sliceValue, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("expected slice for %s but got %T", path, value)
	}

	sliceType := field.Type()
	// Không cần lưu elemType vì không sử dụng, chỉ trực tiếp sử dụng sliceType

	// Tạo slice mới với độ dài phù hợp
	newSlice := reflect.MakeSlice(sliceType, len(sliceValue), len(sliceValue))

	// Duyệt qua các phần tử của slice
	for i, val := range sliceValue {
		elemValue := newSlice.Index(i)

		// Set giá trị cho phần tử
		if err := c.setField(elemValue, val, path+c.opts.Separator+fmt.Sprint(i)); err != nil {
			return err
		}
	}

	// Set slice mới cho field
	field.Set(newSlice)
	return nil
}

// unmarshalMap ánh xạ cấu hình vào một map.
func (c *DefaultConfig) unmarshalMap(field reflect.Value, value interface{}, path string) error {
	mapValue, ok := value.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected map for %s but got %T", path, value)
	}

	mapType := field.Type()
	keyType := mapType.Key()
	elemType := mapType.Elem()

	// Kiểm tra kiểu key của map (chỉ hỗ trợ string)
	if keyType.Kind() != reflect.String {
		return fmt.Errorf("map key must be string for %s", path)
	}

	// Tạo map mới
	newMap := reflect.MakeMap(mapType)

	// Duyệt qua các phần tử của map
	for k, v := range mapValue {
		// Tạo key và value
		keyValue := reflect.ValueOf(k)
		elemValue := reflect.New(elemType).Elem()

		// Set giá trị cho phần tử
		if err := c.setField(elemValue, v, path+c.opts.Separator+k); err != nil {
			return err
		}

		// Thêm vào map
		newMap.SetMapIndex(keyValue, elemValue)
	}

	// Set map mới cho field
	field.Set(newMap)
	return nil
}

// setField set giá trị cho một field.
func (c *DefaultConfig) setField(field reflect.Value, value interface{}, path string) error {
	if !field.CanSet() {
		return nil
	}

	// Xử lý nil
	if value == nil {
		switch field.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
			field.Set(reflect.Zero(field.Type()))
			return nil
		default:
			return nil
		}
	}

	// Xử lý theo kiểu field
	switch field.Kind() {
	case reflect.String:
		str, err := utils.ToString(value)
		if err != nil {
			return fmt.Errorf("cannot convert %v to string for %s: %w", value, path, err)
		}
		field.SetString(str)

	case reflect.Bool:
		b, err := utils.ToBool(value)
		if err != nil {
			return fmt.Errorf("cannot convert %v to bool for %s: %w", value, path, err)
		}
		field.SetBool(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := utils.ToInt(value)
		if err != nil {
			return fmt.Errorf("cannot convert %v to int for %s: %w", value, path, err)
		}
		field.SetInt(int64(i))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := utils.ToInt(value)
		if err != nil || i < 0 {
			return fmt.Errorf("cannot convert %v to uint for %s: %w", value, path, err)
		}
		field.SetUint(uint64(i))

	case reflect.Float32, reflect.Float64:
		f, err := utils.ToFloat(value)
		if err != nil {
			return fmt.Errorf("cannot convert %v to float for %s: %w", value, path, err)
		}
		field.SetFloat(f)

	case reflect.Interface:
		field.Set(reflect.ValueOf(value))

	default:
		return fmt.Errorf("unsupported field type %s for %s", field.Type(), path)
	}

	return nil
}
