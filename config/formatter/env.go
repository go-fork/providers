// Package formatter cung cấp các implementation cho việc nạp cấu hình từ các nguồn khác nhau.
package formatter

import (
	"os"
	"reflect"
	"strings"

	"github.com/go-fork/providers/config/model"
)

// EnvFormatter đọc cấu hình từ biến môi trường.
// Formatter này chuyển đổi tên biến từ dạng UPPER_CASE sang lower.case (dot notation).
type EnvFormatter struct {
	// prefix là tiền tố của biến môi trường, ví dụ "APP_" (tùy chọn)
	prefix string
	// opts là các tùy chọn cho quá trình flatten
	opts FlattenOptions
}

// NewEnvFormatter tạo một EnvFormatter mới với prefix tùy chọn.
func NewEnvFormatter(prefix string) *EnvFormatter {
	return &EnvFormatter{
		prefix: prefix,
		opts:   DefaultFlattenOptions(),
	}
}

// WithOptions cấu hình các tùy chọn cho EnvFormatter.
func (f *EnvFormatter) WithOptions(opts FlattenOptions) *EnvFormatter {
	f.opts = opts
	return f
}

// Load đọc tất cả biến môi trường từ hệ thống.
func (f *EnvFormatter) Load() (interface{}, error) {
	return os.Environ(), nil
}

// Parse chuyển đổi các biến môi trường thành một map.
// Chuyển đổi tên biến từ dạng UPPER_CASE sang lower.case (dot notation).
func (f *EnvFormatter) Parse(data interface{}) (interface{}, error) {
	envVars, ok := data.([]string)
	if !ok {
		return nil, nil
	}

	result := make(map[string]interface{})

	for _, env := range envVars {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]

		// Kiểm tra prefix nếu có
		if f.prefix != "" && !strings.HasPrefix(key, f.prefix) {
			continue
		}

		// Loại bỏ prefix
		if f.prefix != "" {
			key = key[len(f.prefix):]
		}

		// Chuyển đổi key từ UPPER_CASE sang lower.case
		if !f.opts.CaseSensitive {
			key = strings.ToLower(key)
		}

		// Chuyển đổi underscore (_) thành separator (thường là dấu chấm)
		key = strings.ReplaceAll(key, "_", f.opts.Separator)

		result[key] = value
	}

	return result, nil
}

// Flatten chuyển đổi map từ Parse thành ConfigMap phẳng.
func (f *EnvFormatter) Flatten(data interface{}, opts FlattenOptions) (model.ConfigMap, error) {
	if data == nil {
		return model.ConfigMap{}, nil
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return model.ConfigMap{}, nil
	}

	result := model.ConfigMap{}

	for key, value := range dataMap {
		if f.opts.SkipNil && value == nil {
			continue
		}

		valueType := determineValueType(value)
		result[key] = model.ConfigValue{
			Value: value,
			Type:  valueType,
		}
	}

	return result, nil
}

// determineValueType xác định kiểu của một giá trị.
func determineValueType(value interface{}) model.ValueType {
	if value == nil {
		return model.TypeNil
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return model.TypeString
	case reflect.Bool:
		return model.TypeBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return model.TypeInt // Biến môi trường thường là string, có thể chuyển đổi sau
	case reflect.Slice, reflect.Array:
		return model.TypeSlice
	case reflect.Map:
		return model.TypeMap
	default:
		return model.TypeString
	}
}
