// Package formatter cung cấp các implementation cho việc nạp cấu hình từ các nguồn khác nhau.
package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-fork/providers/config/model"
)

// JSONFormatter đọc cấu hình từ file JSON.
type JSONFormatter struct {
	// path là đường dẫn tới file JSON
	path string
	// opts là các tùy chọn cho quá trình flatten
	opts FlattenOptions
}

// NewJSONFormatter tạo một JSONFormatter mới.
func NewJSONFormatter(path string) *JSONFormatter {
	return &JSONFormatter{
		path: path,
		opts: DefaultFlattenOptions(),
	}
}

// WithOptions cấu hình các tùy chọn cho JSONFormatter.
func (f *JSONFormatter) WithOptions(opts FlattenOptions) *JSONFormatter {
	f.opts = opts
	return f
}

// Load đọc nội dung từ file JSON.
func (f *JSONFormatter) Load() (interface{}, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return nil, fmt.Errorf("cannot read JSON file %s: %w", f.path, err)
	}
	return data, nil
}

// Parse chuyển đổi nội dung JSON thành cấu trúc dữ liệu Go.
func (f *JSONFormatter) Parse(data interface{}) (interface{}, error) {
	jsonData, ok := data.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected []byte but got %T", data)
	}

	var result interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, fmt.Errorf("cannot parse JSON: %w", err)
	}

	return result, nil
}

// Flatten chuyển đổi cấu trúc dữ liệu lồng nhau thành ConfigMap phẳng với dot notation.
func (f *JSONFormatter) Flatten(data interface{}, opts FlattenOptions) (model.ConfigMap, error) {
	result := model.ConfigMap{}
	if data == nil {
		return result, nil
	}

	err := f.flattenRecursive(data, "", result, opts)
	return result, err
}

// flattenRecursive đệ quy qua cấu trúc dữ liệu lồng nhau để tạo ConfigMap phẳng.
func (f *JSONFormatter) flattenRecursive(data interface{}, prefix string, result model.ConfigMap, opts FlattenOptions) error {
	if data == nil {
		if !opts.SkipNil {
			result[prefix] = model.ConfigValue{Value: nil, Type: model.TypeNil}
		}
		return nil
	}

	rv := reflect.ValueOf(data)

	switch rv.Kind() {
	case reflect.Map:
		// Lưu giá trị gốc nếu có prefix
		if prefix != "" {
			result[prefix] = model.ConfigValue{Value: data, Type: model.TypeMap}
		}

		for _, key := range rv.MapKeys() {
			strKey := fmt.Sprintf("%v", key.Interface())
			if !opts.CaseSensitive {
				strKey = strings.ToLower(strKey)
			}

			newPrefix := strKey
			if prefix != "" {
				newPrefix = prefix + opts.Separator + strKey
			}

			err := f.flattenRecursive(rv.MapIndex(key).Interface(), newPrefix, result, opts)
			if err != nil {
				return err
			}
		}

	case reflect.Slice, reflect.Array:
		// Lưu giá trị gốc
		if prefix != "" {
			result[prefix] = model.ConfigValue{Value: data, Type: model.TypeSlice}
		}

		length := rv.Len()
		for i := 0; i < length; i++ {
			newPrefix := prefix + opts.Separator + strconv.Itoa(i)
			if prefix == "" {
				newPrefix = strconv.Itoa(i)
			}

			err := f.flattenRecursive(rv.Index(i).Interface(), newPrefix, result, opts)
			if err != nil {
				return err
			}
		}

	case reflect.Bool:
		result[prefix] = model.ConfigValue{Value: rv.Bool(), Type: model.TypeBool}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result[prefix] = model.ConfigValue{Value: rv.Int(), Type: model.TypeInt}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result[prefix] = model.ConfigValue{Value: rv.Uint(), Type: model.TypeInt}

	case reflect.Float32, reflect.Float64:
		result[prefix] = model.ConfigValue{Value: rv.Float(), Type: model.TypeFloat}

	case reflect.String:
		result[prefix] = model.ConfigValue{Value: rv.String(), Type: model.TypeString}

	default:
		// Kiểu khác thì xử lý như string
		result[prefix] = model.ConfigValue{Value: fmt.Sprintf("%v", data), Type: model.TypeString}
	}

	return nil
}
