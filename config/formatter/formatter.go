// Package formatter cung cấp các interface và implementation để nạp cấu hình từ nhiều nguồn khác nhau.
//
// Package này định nghĩa interface Formatter để trừu tượng hóa việc nạp cấu hình từ các nguồn như
// file YAML, file JSON, biến môi trường, v.v. Các implementation của Formatter được cung cấp trong
// các file tương ứng (yaml.go, json.go, env.go).
package formatter

import (
	"github.com/go-fork/providers/config/model"
)

// FlattenOptions định nghĩa các tùy chọn cho quá trình flatten cấu trúc dữ liệu phân cấp.
type FlattenOptions struct {
	// Separator là ký tự phân tách được sử dụng trong dot notation (mặc định là ".")
	Separator string
	// SkipNil quyết định có bỏ qua các giá trị nil hay không
	SkipNil bool
	// HandleEmptyKey quyết định có xử lý key rỗng hay không
	HandleEmptyKey bool
	// CaseSensitive quyết định có phân biệt hoa thường trong key hay không
	CaseSensitive bool
}

// DefaultFlattenOptions trả về các tùy chọn mặc định cho quá trình flatten.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator:      ".",
		SkipNil:        true,
		HandleEmptyKey: false,
		CaseSensitive:  false,
	}
}

// Formatter định nghĩa interface cho các nguồn cấu hình (configuration source).
//
// Interface này cho phép trừu tượng hóa việc nạp dữ liệu cấu hình từ nhiều nguồn khác nhau
// như file YAML, file JSON, biến môi trường, v.v. Mỗi Formatter chịu trách nhiệm nạp dữ liệu
// từ một nguồn cụ thể và chuyển đổi nó thành một map[string]interface{} mà Manager có thể sử dụng.
//
// Các implementation phải đảm bảo trả về map với key dạng dot notation cho dữ liệu phân cấp.
type Formatter interface {
	Load() (interface{}, error)                                             // Tải dữ liệu
	Parse(data interface{}) (interface{}, error)                            // Phân tích cú pháp
	Flatten(data interface{}, opts FlattenOptions) (model.ConfigMap, error) // Flatten recursive
}
