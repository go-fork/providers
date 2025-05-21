package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// JsonFormatter cung cấp Formatter cho phép nạp cấu hình từ file JSON.
//
// JsonFormatter hỗ trợ làm phẳng các key phân cấp thành dot notation.
type JsonFormatter struct {
	path string // Đường dẫn file JSON
}

// NewJsonFormatter khởi tạo một JsonFormatter mới.
// NewJsonFormatter nhận vào đường dẫn file JSON.
func NewJsonFormatter(path string) *JsonFormatter {
	return &JsonFormatter{
		path: path,
	}
}

// Load tải cấu hình từ file JSON.
// Load trả về map[string]interface{} chứa các giá trị cấu hình, hoặc error nếu có lỗi khi nạp file hoặc parse JSON.
func (p *JsonFormatter) Load() (map[string]interface{}, error) {
	if p.path == "" {
		return nil, fmt.Errorf("empty file path")
	}

	// Đọc file
	data, err := os.ReadFile(p.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", p.path)
		}
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Kiểm tra file rỗng hoặc chỉ chứa {}
	if len(data) == 0 || (len(data) == 2 && string(data) == "{}") {
		return make(map[string]interface{}), nil
	}

	// Thử parse JSON thành map
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		// Nếu không thể parse thành map, thử parse thành các kiểu khác (array, primitive)
		var anyValue interface{}
		if errAny := json.Unmarshal(data, &anyValue); errAny == nil {
			return make(map[string]interface{}), nil
		}
		return nil, fmt.Errorf("failed to parse JSON config at %s: %v", p.path, err)
	}

	flattenMap := make(map[string]interface{})
	flattenMapRecursive(flattenMap, result, "")
	return flattenMap, nil
}

// Name trả về tên định danh của JsonFormatter.
// Name trả về string "json:<tên file>".
func (p *JsonFormatter) Name() string {
	return "json:" + filepath.Base(p.path)
}
