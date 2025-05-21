package formatter

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// YamlFormatter cung cấp Formatter cho phép nạp cấu hình từ file YAML.
//
// YamlFormatter hỗ trợ làm phẳng các key phân cấp thành dot notation.
type YamlFormatter struct {
	path string // Đường dẫn file YAML
}

// NewYamlFormatter khởi tạo một YamlFormatter mới.
// NewYamlFormatter nhận vào đường dẫn file YAML.
func NewYamlFormatter(path string) *YamlFormatter {
	return &YamlFormatter{
		path: path,
	}
}

// Load tải cấu hình từ file YAML.
// Load trả về map[string]interface{} chứa các giá trị cấu hình, hoặc error nếu có lỗi khi nạp file hoặc parse YAML.
func (p *YamlFormatter) Load() (map[string]interface{}, error) {
	// Kiểm tra file tồn tại
	if _, err := os.Stat(p.path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", p.path)
	}

	// Đọc file
	data, err := os.ReadFile(p.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Kiểm tra file rỗng
	if len(data) == 0 {
		return make(map[string]interface{}), nil
	}

	// Parse YAML
	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %v", err)
	}

	// Làm phẳng key phân cấp thành dot notation
	flattenMap := make(map[string]interface{})
	flattenMapRecursive(flattenMap, result, "")

	return flattenMap, nil
}

// Name trả về tên định danh của YamlFormatter.
// Name trả về string "yaml:<tên file>".
func (p *YamlFormatter) Name() string {
	return "yaml:" + filepath.Base(p.path)
}

// LoadFromDirectory nạp tất cả các file YAML trong thư mục.
// LoadFromDirectory nhận vào đường dẫn thư mục chứa file YAML, trả về map các giá trị cấu hình hoặc error nếu có lỗi khi đọc thư mục hoặc file.
func LoadFromDirectory(directory string) (map[string]interface{}, error) {
	if directory == "" {
		return nil, fmt.Errorf("empty directory path")
	}
	result := make(map[string]interface{})
	info, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("config directory not found: %s", directory)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", directory)
	}
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		filePath := filepath.Join(directory, file.Name())
		formatter := NewYamlFormatter(filePath)
		values, err := formatter.Load()
		if err != nil {
			return nil, fmt.Errorf("failed to load YAML file %s: %w", filePath, err)
		}
		for k, v := range values {
			result[k] = v
		}
	}
	return result, nil
}
