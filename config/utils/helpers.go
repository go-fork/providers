// Package utils cung cấp các tiện ích hỗ trợ cho package config.
//
// Package này chứa các hàm tiện ích để xử lý cấu trúc dữ liệu map phân cấp,
// hỗ trợ việc chuyển đổi giữa map phẳng (flat map) với key dạng dot notation
// và map lồng nhau (nested map) thể hiện cấu trúc phân cấp.
package utils

import "strings"

// flattenMapRecursive chuyển đổi map lồng nhau thành map phẳng với key dạng dot notation.
//
// Hàm này duyệt qua toàn bộ cấu trúc map lồng nhau và tạo ra một map phẳng
// trong đó mỗi key được ghép từ đường dẫn đến giá trị đó theo dot notation.
// Ví dụ, map lồng nhau:
//
//	{
//	  "database": {
//	    "host": "localhost",
//	    "port": 5432
//	  }
//	}
//
// sẽ được chuyển đổi thành map phẳng:
//
//	{
//	  "database.host": "localhost",
//	  "database.port": 5432
//	}
//
// Hàm này được sử dụng nội bộ và được gọi thông qua hàm exported FlattenMapRecursive.
//
// Params:
//   - result: map[string]interface{} - Map đích để lưu kết quả làm phẳng
//   - nested: map[string]interface{} - Map lồng nhau cần làm phẳng
//   - prefix: string - Tiền tố cho key kết quả (thường là rỗng ở lần gọi đầu tiên)
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

// unflattenDotMapRecursive chuyển đổi map phẳng với key dạng dot notation thành map lồng nhau.
//
// Hàm này thực hiện quá trình ngược lại với flattenMapRecursive, nhận một map phẳng
// với key dạng dot notation và tái cấu trúc nó thành một map lồng nhau. Các phần của key
// được phân tách bởi dấu chấm được sử dụng để xây dựng lại cấu trúc phân cấp.
//
// Ví dụ, map phẳng:
//
//	{
//	  "database.host": "localhost",
//	  "database.port": 5432,
//	  "database.credentials.username": "admin"
//	}
//
// sẽ được chuyển đổi thành map lồng nhau:
//
//	{
//	  "database": {
//	    "host": "localhost",
//	    "port": 5432,
//	    "credentials": {
//	      "username": "admin"
//	    }
//	  }
//	}
//
// Hàm này là private và được gọi thông qua hàm exported UnflattenDotMapRecursive.
//
// Params:
//   - flat: map[string]interface{} - Map phẳng với key dạng dot notation
//
// Returns:
//   - map[string]interface{}: Map lồng nhau được tái cấu trúc
func unflattenDotMapRecursive(flat map[string]interface{}) map[string]interface{} {
	nested := make(map[string]interface{})
	for k, v := range flat {
		parts := strings.SplitN(k, ".", 2)
		if len(parts) == 1 {
			nested[parts[0]] = v
		} else {
			sub, ok := nested[parts[0]].(map[string]interface{})
			if !ok {
				sub = make(map[string]interface{})
			}
			// merge submap recursively
			for kk, vv := range unflattenDotMapRecursive(map[string]interface{}{parts[1]: v}) {
				if exist, ok := sub[kk].(map[string]interface{}); ok {
					if vvmap, ok := vv.(map[string]interface{}); ok {
						// merge two maps
						for kkk, vvv := range vvmap {
							exist[kkk] = vvv
						}
						sub[kk] = exist
						continue
					}
				}
				sub[kk] = vv
			}
			nested[parts[0]] = sub
		}
	}
	return nested
}

// FlattenMapRecursive chuyển đổi map lồng nhau thành map phẳng với key dạng dot notation.
//
// Hàm này là wrapper công khai (exported) cho hàm flattenMapRecursive, cung cấp
// khả năng chuyển đổi từ cấu trúc dữ liệu phân cấp thành dạng phẳng với key được nối bởi dấu chấm.
//
// Quá trình làm phẳng sẽ xử lý cả map[string]interface{} và map[interface{}]interface{},
// tự động chuyển đổi key không phải string khi cần thiết.
//
// Params:
//   - result: map[string]interface{} - Map đích để lưu kết quả làm phẳng
//   - nested: map[string]interface{} - Map lồng nhau cần làm phẳng
//   - prefix: string - Tiền tố cho key kết quả (thường để trống khi gọi lần đầu)
//
// Examples:
//
//	nested := map[string]interface{}{
//	  "database": map[string]interface{}{
//	    "host": "localhost",
//	    "port": 5432,
//	  },
//	}
//	flat := make(map[string]interface{})
//	utils.FlattenMapRecursive(flat, nested, "")
//	// flat sẽ chứa: {"database.host": "localhost", "database.port": 5432}
func FlattenMapRecursive(result map[string]interface{}, nested map[string]interface{}, prefix string) {
	flattenMapRecursive(result, nested, prefix)
}

// UnflattenDotMapRecursive chuyển đổi map phẳng với key dạng dot notation thành map lồng nhau.
//
// Hàm này là wrapper công khai (exported) cho hàm unflattenDotMapRecursive, cung cấp
// khả năng chuyển đổi từ map phẳng với key dạng dot notation thành cấu trúc dữ liệu phân cấp.
//
// Hàm này thông minh trong việc xử lý các trường hợp xung đột và hợp nhất (merging) các key con
// trong cùng một nhánh phân cấp.
//
// Params:
//   - flat: map[string]interface{} - Map phẳng với key dạng dot notation
//
// Returns:
//   - map[string]interface{}: Map lồng nhau được tái cấu trúc
//
// Examples:
//
//	flat := map[string]interface{}{
//	  "app.name": "MyApp",
//	  "app.version": 1.0,
//	  "app.features.auth": true,
//	}
//	nested := utils.UnflattenDotMapRecursive(flat)
//	// nested sẽ có cấu trúc:
//	// {
//	//   "app": {
//	//     "name": "MyApp",
//	//     "version": 1.0,
//	//     "features": {
//	//       "auth": true
//	//     }
//	//   }
//	// }
func UnflattenDotMapRecursive(flat map[string]interface{}) map[string]interface{} {
	return unflattenDotMapRecursive(flat)
}
