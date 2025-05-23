// Package utils cung cấp các hàm tiện ích cho hệ thống quản lý cấu hình.
package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// IsNil kiểm tra xem một giá trị có phải là nil không, xử lý đúng các kiểu dữ liệu interface, map, slice, ...
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

// ToString chuyển đổi một giá trị bất kỳ thành chuỗi.
func ToString(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}

	switch val := v.(type) {
	case string:
		return val, nil
	case bool:
		return strconv.FormatBool(val), nil
	case int:
		return strconv.Itoa(val), nil
	case int64:
		return strconv.FormatInt(val, 10), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case []byte:
		return string(val), nil
	default:
		return fmt.Sprintf("%v", val), nil
	}
}

// ToBool chuyển đổi một giá trị thành boolean an toàn.
func ToBool(v interface{}) (bool, error) {
	if v == nil {
		return false, fmt.Errorf("cannot convert nil to bool")
	}

	switch val := v.(type) {
	case bool:
		return val, nil
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(val).Int() != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(val).Uint() != 0, nil
	case float32, float64:
		return reflect.ValueOf(val).Float() != 0, nil
	case string:
		val = strings.ToLower(strings.TrimSpace(val))
		switch val {
		case "true", "yes", "1", "on", "t", "y":
			return true, nil
		case "false", "no", "0", "off", "f", "n", "":
			return false, nil
		default:
			return false, fmt.Errorf("cannot convert string %q to bool", val)
		}
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

// ToInt chuyển đổi một giá trị thành int an toàn.
func ToInt(v interface{}) (int, error) {
	if v == nil {
		return 0, fmt.Errorf("cannot convert nil to int")
	}

	switch val := v.(type) {
	case int:
		return val, nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case int64:
		return int(val), nil
	case uint:
		return int(val), nil
	case uint8:
		return int(val), nil
	case uint16:
		return int(val), nil
	case uint32:
		return int(val), nil
	case uint64:
		return int(val), nil
	case float32:
		return int(val), nil
	case float64:
		return int(val), nil
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case string:
		return strconv.Atoi(strings.TrimSpace(val))
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

// ToFloat chuyển đổi một giá trị thành float64 an toàn.
func ToFloat(v interface{}) (float64, error) {
	if v == nil {
		return 0, fmt.Errorf("cannot convert nil to float64")
	}

	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(strings.TrimSpace(val), 64)
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

// MergeConfigMaps hợp nhất nhiều ConfigMap với thứ tự ưu tiên.
// ConfigMap đầu tiên có độ ưu tiên cao nhất.
func MergeConfigMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Duyệt qua các map theo thứ tự ưu tiên từ thấp đến cao
	for i := len(maps) - 1; i >= 0; i-- {
		for key, val := range maps[i] {
			result[key] = val
		}
	}

	return result
}

// GetStructTag lấy giá trị tag "config" từ một struct field.
func GetStructTag(field reflect.StructField) string {
	return field.Tag.Get("config")
}
