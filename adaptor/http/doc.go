/*
Package adapter cung cấp một adapter HTTP cho framework PAMM dựa trên thư viện net/http của Go.

Adapter này bao gồm:
  - Cấu hình server HTTP/HTTPS đầy đủ
  - Hỗ trợ middleware chain
  - Xử lý các route với phương thức HTTP cụ thể
  - Cấu hình TLS linh hoạt
  - Tùy chọn nén dữ liệu (compression)
  - Cấu hình giới hạn kích thước body

Ví dụ sử dụng cơ bản:

	// Tạo adapter với cấu hình mặc định
	httpAdapter := adapter.NewNetHTTPAdapter(nil)

	// Khởi tạo adapter với các middleware mặc định
	httpAdapter.Initialize()

	// Đăng ký route handler
	httpAdapter.HandleFunc("GET", "/hello", func(ctx httpCtx.Context) {
		ctx.Writer().Write([]byte("Hello, World!"))
	})

	// Khởi động server
	err := httpAdapter.Run()
	if err != nil {
		log.Fatal(err)
	}

Để cấu hình TLS:

	config := adapter.DefaultConfig()
	config.TLS.Enabled = true
	config.TLS.CertFile = "/path/to/cert.pem"
	config.TLS.KeyFile = "/path/to/key.pem"

	httpAdapter := adapter.NewNetHTTPAdapter(config)
	// Tiếp tục cấu hình và khởi động server...
*/
package adapter
