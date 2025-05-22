// Package formatter cung cấp các interface và implementation để nạp cấu hình từ nhiều nguồn khác nhau.
//
// Package này định nghĩa interface Formatter để trừu tượng hóa việc nạp cấu hình từ các nguồn như
// file YAML, file JSON, biến môi trường, v.v. Các implementation của Formatter được cung cấp trong
// các file tương ứng (yaml.go, json.go, env.go).
package formatter

// Formatter định nghĩa interface cho các nguồn cấu hình (configuration source).
//
// Interface này cho phép trừu tượng hóa việc nạp dữ liệu cấu hình từ nhiều nguồn khác nhau
// như file YAML, file JSON, biến môi trường, v.v. Mỗi Formatter chịu trách nhiệm nạp dữ liệu
// từ một nguồn cụ thể và chuyển đổi nó thành một map[string]interface{} mà Manager có thể sử dụng.
//
// Các implementation phải đảm bảo trả về map với key dạng dot notation cho dữ liệu phân cấp.
type Formatter interface {
	// Load tải dữ liệu cấu hình từ nguồn.
	//
	// Phương thức này đọc và phân tích dữ liệu cấu hình từ nguồn cụ thể (file, env, ...)
	// và trả về kết quả dưới dạng map[string]interface{} với key dạng dot notation.
	//
	// Returns:
	//   - map[string]interface{}: Map chứa các cặp key-value cấu hình.
	//   - error: Lỗi nếu không thể nạp hoặc phân tích dữ liệu.
	//
	// Examples:
	//   values, err := formatter.Load()
	//   if err == nil {
	//     for k, v := range values {
	//       manager.Set(k, v)
	//     }
	//   }
	Load() (map[string]interface{}, error)

	// Name trả về tên định danh của Formatter.
	//
	// Phương thức này trả về một chuỗi mô tả nguồn cấu hình,
	// thường bao gồm loại formatter và nguồn dữ liệu cụ thể.
	//
	// Returns:
	//   - string: Tên định danh của Formatter (ví dụ: "yaml:config.yaml", "env:APP_").
	//
	// Examples:
	//   name := formatter.Name() // Ví dụ: "yaml:config.yaml"
	Name() string
}
