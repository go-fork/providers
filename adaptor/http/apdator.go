package adapter

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	httpAdapter "github.com/go-fork/pamm/pkg/infra/http/adapter"
	httpCtx "github.com/go-fork/pamm/pkg/infra/http/context"
)

// NetHTTPAdapter là adapter cho gói net/http của Go
type NetHTTPAdapter struct {
	*httpAdapter.BaseAdapter
	config      *Config
	mux         *http.ServeMux
	middlewares []func(httpCtx.Context)
	server      *http.Server
}

// NewNetHTTPAdapter tạo một adapter mới cho net/http
func NewNetHTTPAdapter(config *Config) *NetHTTPAdapter {
	if config == nil {
		config = DefaultConfig()
	}

	mux := http.NewServeMux()

	// Tạo địa chỉ từ cấu hình
	addr := fmt.Sprintf("%s:%d", config.Addr, config.Port)

	// Chỉ sử dụng các cấu hình được hỗ trợ bởi http.Server
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
		MaxHeaderBytes:    int(config.MaxHeaderBytes),
		ReadHeaderTimeout: config.ReadTimeout, // Sử dụng ReadTimeout khi không có ReadHeaderTimeout riêng
	}

	// Cấu hình TLS nếu được bật
	if config.TLS != nil && config.TLS.Enabled {
		server.TLSConfig = createTLSConfig(config.TLS)
	}

	adapter := &NetHTTPAdapter{
		BaseAdapter: httpAdapter.NewBaseAdapter(ConfigPrefix),
		config:      config,
		mux:         mux,
		middlewares: []func(httpCtx.Context){},
		server:      server,
	}

	adapter.SetHandler(mux)

	return adapter
}

// HandleFunc đăng ký một handler function với method và path
func (a *NetHTTPAdapter) HandleFunc(method, path string, handler func(ctx httpCtx.Context)) {
	a.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Tạo context cho request
		ctx := httpCtx.NewContext(w, r)

		// Thiết lập handlers chain bao gồm tất cả middlewares
		// và route handler cuối cùng
		handlers := make([]func(httpCtx.Context), 0, len(a.middlewares)+1)
		handlers = append(handlers, a.middlewares...)
		handlers = append(handlers, handler)

		ctx.SetHandlers(handlers)

		// Bắt đầu thực thi chuỗi middleware và handler
		ctx.Next()
	})
}

// Use thêm middleware vào adapter
func (a *NetHTTPAdapter) Use(middleware func(ctx httpCtx.Context)) {
	// Thêm middleware vào danh sách
	a.middlewares = append(a.middlewares, middleware)
}

// Run khởi động HTTP server với cấu hình từ config
func (a *NetHTTPAdapter) Run() error {
	// Áp dụng các middleware mặc định
	a.Initialize()

	// Kiểm tra nếu TLS được bật
	if a.config.TLS != nil && a.config.TLS.Enabled {
		return a.RunTLSWithConfig()
	}

	fmt.Printf("Server starting on http://%s:%d\n", a.config.Addr, a.config.Port)
	return a.server.ListenAndServe()
}

// RunTLSWithConfig khởi động HTTPS server với cấu hình TLS từ config
func (a *NetHTTPAdapter) RunTLSWithConfig() error {
	if a.config.TLS == nil || !a.config.TLS.Enabled {
		return fmt.Errorf("TLS not enabled in configuration")
	}

	fmt.Printf("Server starting on https://%s:%d\n", a.config.Addr, a.config.Port)
	return a.server.ListenAndServeTLS(a.config.TLS.CertFile, a.config.TLS.KeyFile)
}

// SetServer thiết lập http.Server cho adapter
func (a *NetHTTPAdapter) SetServer(server *http.Server) {
	a.server = server
	a.server.Handler = a.handler // Thiết lập handler từ BaseAdapter
}

// Initialize khởi tạo adapter với các middleware mặc định
func (a *NetHTTPAdapter) Initialize() {
	// Áp dụng các middleware mặc định theo cấu hình
	a.ApplyBodyLimit()
	a.ApplyCompression()

	// Thêm debug middleware nếu cần
	if a.config.Debug {
		a.Use(func(ctx httpCtx.Context) {
			start := time.Now()
			ctx.Next()
			elapsed := time.Since(start)

			req := ctx.Request()
			fmt.Printf("[DEBUG] %s %s - %v\n", req.Method, req.URL.Path, elapsed)
		})
	}
}

// ApplyCompression áp dụng compression middleware nếu được bật trong cấu hình
func (a *NetHTTPAdapter) ApplyCompression() {
	if a.config.Compression {
		a.Use(compressionMiddleware)
	}
}

// ApplyBodyLimit áp dụng body limit middleware theo cấu hình
func (a *NetHTTPAdapter) ApplyBodyLimit() {
	if a.config.BodyLimit > 0 {
		a.Use(func(ctx httpCtx.Context) {
			req := ctx.Request()
			if req.ContentLength > int64(a.config.BodyLimit) {
				ctx.Writer().WriteHeader(http.StatusRequestEntityTooLarge)
				ctx.Writer().Write([]byte("Request body too large"))
				return
			}
			ctx.Next()
		})
	}
}

// compressionMiddleware là middleware để nén response (gzip)
func compressionMiddleware(ctx httpCtx.Context) {
	// TODO: Implement compression middleware
	// Đây là placeholder, cần triển khai middleware nén thực tế
	ctx.Next()
}

// createTLSConfig tạo cấu hình TLS từ TLSConfig của adapter
func createTLSConfig(tlsConfig *TLSConfig) *tls.Config {
	config := &tls.Config{
		PreferServerCipherSuites: tlsConfig.PreferServerCipherSuites,
	}

	// Thiết lập phiên bản TLS
	minVersion, err := parseTLSVersion(tlsConfig.MinVersion)
	if err == nil {
		config.MinVersion = minVersion
	}

	maxVersion, err := parseTLSVersion(tlsConfig.MaxVersion)
	if err == nil {
		config.MaxVersion = maxVersion
	}

	// Thiết lập cipher suites nếu được cung cấp
	if len(tlsConfig.CipherSuites) > 0 {
		cipherSuites, err := parseCipherSuites(tlsConfig.CipherSuites)
		if err == nil {
			config.CipherSuites = cipherSuites
		}
	}

	// Thiết lập curve preferences nếu được cung cấp
	if len(tlsConfig.CurvePreferences) > 0 {
		curveIDs, err := parseCurvePreferences(tlsConfig.CurvePreferences)
		if err == nil {
			config.CurvePreferences = curveIDs
		}
	}

	// Đảm bảo NextProtos bao gồm HTTP/1.1
	config.NextProtos = []string{"http/1.1"}

	return config
}

// parseTLSVersion chuyển đổi chuỗi phiên bản TLS sang giá trị uint16
func parseTLSVersion(version string) (uint16, error) {
	switch version {
	case "1.0":
		return tls.VersionTLS10, nil
	case "1.1":
		return tls.VersionTLS11, nil
	case "1.2":
		return tls.VersionTLS12, nil
	case "1.3":
		return tls.VersionTLS13, nil
	default:
		return 0, fmt.Errorf("unsupported TLS version: %s", version)
	}
}

// parseCipherSuites chuyển đổi danh sách tên cipher suites thành giá trị uint16
func parseCipherSuites(names []string) ([]uint16, error) {
	// Ánh xạ từ tên đến giá trị cipher suite
	cipherMap := map[string]uint16{
		"TLS_RSA_WITH_RC4_128_SHA":                tls.TLS_RSA_WITH_RC4_128_SHA,
		"TLS_RSA_WITH_3DES_EDE_CBC_SHA":           tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		"TLS_RSA_WITH_AES_128_CBC_SHA":            tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		"TLS_RSA_WITH_AES_256_CBC_SHA":            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		"TLS_RSA_WITH_AES_128_CBC_SHA256":         tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		"TLS_RSA_WITH_AES_128_GCM_SHA256":         tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		"TLS_RSA_WITH_AES_256_GCM_SHA384":         tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_ECDSA_WITH_RC4_128_SHA":        tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA":    tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_RC4_128_SHA":          tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":     tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA":      tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256":   tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":   tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305":    tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305":  tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	}

	var result []uint16
	for _, name := range names {
		if val, ok := cipherMap[name]; ok {
			result = append(result, val)
		} else {
			// Thử chuyển đổi từ hexadecimal
			if strings.HasPrefix(name, "0x") {
				val, err := strconv.ParseUint(name[2:], 16, 16)
				if err != nil {
					return nil, fmt.Errorf("invalid cipher suite: %s", name)
				}
				result = append(result, uint16(val))
			} else {
				return nil, fmt.Errorf("unknown cipher suite: %s", name)
			}
		}
	}

	return result, nil
}

// parseCurvePreferences chuyển đổi danh sách tên curve preferences thành giá trị tls.CurveID
func parseCurvePreferences(names []string) ([]tls.CurveID, error) {
	// Ánh xạ từ tên đến giá trị curve
	curveMap := map[string]tls.CurveID{
		"P256":       tls.CurveP256,
		"P384":       tls.CurveP384,
		"P521":       tls.CurveP521,
		"X25519":     tls.X25519,
		"CURVE_P256": tls.CurveP256,
		"CURVE_P384": tls.CurveP384,
		"CURVE_P521": tls.CurveP521,
	}

	var result []tls.CurveID
	for _, name := range names {
		if val, ok := curveMap[name]; ok {
			result = append(result, val)
		} else {
			// Thử chuyển đổi từ số nguyên
			val, err := strconv.ParseUint(name, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("unknown curve name: %s", name)
			}
			result = append(result, tls.CurveID(val))
		}
	}

	return result, nil
}
