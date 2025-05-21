package formatter

// Formatter định nghĩa interface cho các nguồn cấu hình (configuration source).
//
// Interface này cho phép trừu tượng hóa việc nạp dữ liệu cấu hình từ nhiều nguồn khác nhau (file, env, ...).
type Formatter interface {
	// Load tải dữ liệu cấu hình từ nguồn.
	// Load trả về map[string]interface{} chứa các giá trị cấu hình, hoặc error nếu có lỗi khi nạp.
	Load() (map[string]interface{}, error)

	// Name trả về tên định danh của Formatter.
	// Name trả về string là tên nguồn cấu hình (ví dụ: "env", "yaml:config.yaml").
	Name() string
}

// flattenMapRecursive làm phẳng map phân cấp thành dot notation.
// flattenMapRecursive nhận vào map kết quả, map phân cấp và prefix cho key.
func flattenMapRecursive(result map[string]interface{}, nested map[string]interface{}, prefix string) {
	if result == nil || nested == nil {
		return
	}

	for k, v := range nested {
		if k == "" {
			continue // Bỏ qua key rỗng
		}

		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch val := v.(type) {
		case map[string]interface{}:
			// Đệ quy với map con
			flattenMapRecursive(result, val, key)
		case map[interface{}]interface{}:
			// Chuyển đổi map[interface{}]interface{} thành map[string]interface{}
			stringMap := make(map[string]interface{}, len(val))
			for mk, mv := range val {
				if mkStr, ok := mk.(string); ok {
					stringMap[mkStr] = mv
				}
			}
			flattenMapRecursive(result, stringMap, key)
		default:
			result[key] = v
		}
	}
}

// Export flattenMapRecursive for test coverage
func FlattenMapRecursiveForTest(result map[string]interface{}, nested map[string]interface{}, prefix string) {
	flattenMapRecursive(result, nested, prefix)
}
