// Package formatter cung cấp các định dạng (formatter) khác nhau để nạp cấu hình từ nhiều nguồn.
//
// Package này chứa các formatter cho các nguồn cấu hình phổ biến như YAML, JSON, và
// biến môi trường (environment variables). Các formatter thực thi interface Formatter,
// cho phép hệ thống quản lý cấu hình nạp dữ liệu từ nhiều nguồn khác nhau một cách thống nhất.
package formatter

import (
	"fmt"
	"os"

	"github.com/go-fork/providers/config/model"
	"gopkg.in/yaml.v3"
)

// YAMLFormatter đọc cấu hình từ file YAML.
type YAMLFormatter struct {
	// path là đường dẫn tới file YAML
	path string
	// opts là các tùy chọn cho quá trình flatten
	opts FlattenOptions
	// jsonFormatter là formatter sử dụng để flatten dữ liệu sau khi parse YAML
	jsonFormatter *JSONFormatter
}

// NewYAMLFormatter tạo một YAMLFormatter mới.
func NewYAMLFormatter(path string) *YAMLFormatter {
	return &YAMLFormatter{
		path:          path,
		opts:          DefaultFlattenOptions(),
		jsonFormatter: NewJSONFormatter(""), // Không sử dụng path trong JSONFormatter
	}
}

// WithOptions cấu hình các tùy chọn cho YAMLFormatter.
func (f *YAMLFormatter) WithOptions(opts FlattenOptions) *YAMLFormatter {
	f.opts = opts
	f.jsonFormatter.WithOptions(opts)
	return f
}

// Load đọc nội dung từ file YAML.
func (f *YAMLFormatter) Load() (interface{}, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return nil, fmt.Errorf("cannot read YAML file %s: %w", f.path, err)
	}
	return data, nil
}

// Parse chuyển đổi nội dung YAML thành cấu trúc dữ liệu Go.
func (f *YAMLFormatter) Parse(data interface{}) (interface{}, error) {
	yamlData, ok := data.([]byte)
	if !ok {
		return nil, fmt.Errorf("expected []byte but got %T", data)
	}

	var result interface{}
	if err := yaml.Unmarshal(yamlData, &result); err != nil {
		return nil, fmt.Errorf("cannot parse YAML: %w", err)
	}

	return result, nil
}

// Flatten chuyển đổi cấu trúc dữ liệu lồng nhau thành ConfigMap phẳng với dot notation.
// Sử dụng lại logic của JSONFormatter vì có cùng cấu trúc dữ liệu sau khi parse.
func (f *YAMLFormatter) Flatten(data interface{}, opts FlattenOptions) (model.ConfigMap, error) {
	return f.jsonFormatter.Flatten(data, opts)
}
