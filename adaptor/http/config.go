package adapter

import (
	"time"
)

var (
	// AdapterName là tên định danh của adapter HTTP
	AdapterName = "http"
	// ConfigPrefix là tiền tố được sử dụng trong cấu hình
	ConfigPrefix = "http.http"
)

type Config struct {
	// Addr optionally specifies the TCP address for the server to listen on,
	// in the form "host:port". If empty, ":http" (port 80) is used.
	// The service names are defined in RFC 6335 and assigned by IANA.
	// See net.Dial for details of the address format.
	Addr string `json:"addr" yaml:"addr"`

	// DisableGeneralOptionsHandler, if true, passes "OPTIONS *" requests to the Handler,
	// otherwise responds with 200 OK and Content-Length: 0.
	DisableGeneralOptionsHandler bool `json:"disable_general_options_handler" yaml:"disable_general_options_handler"`

	// TLSConfig optionally provides a TLS configuration for use
	// by ServeTLS and ListenAndServeTLS. Note that this value is
	// cloned by ServeTLS and ListenAndServeTLS, so it's not
	// possible to modify the configuration with methods like
	// tls.Config.SetSessionTicketKeys. To use
	// SetSessionTicketKeys, use Server.Serve with a TLS Listener
	// instead.
	TLSConfig *TLSConfig `json:"tls_config" yaml:"tls_config"`

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body. A zero or negative value means
	// there will be no timeout.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration `json:"read_timeout" yaml:"read_timeout"`

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body. If zero, the value of
	// ReadTimeout is used. If negative, or if zero and ReadTimeout
	// is zero or negative, there is no timeout.
	ReadHeaderTimeout time.Duration `json:"read_header_timeout" yaml:"read_header_timeout"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	// A zero or negative value means there will be no timeout.
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If zero, the value
	// of ReadTimeout is used. If negative, or if zero and ReadTimeout
	// is zero or negative, there is no timeout.
	IdleTimeout time.Duration `json:"idle_timeout" yaml:"idle_timeout"`

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes is used.
	MaxHeaderBytes int `json:"max_header_bytes" yaml:"max_header_bytes"`
}

type TLSConfig struct {
	// CertFile là đường dẫn đến file certificate
	// Bắt buộc nếu TLS được bật
	CertFile string `json:"cert_file" yaml:"cert_file"`

	// KeyFile là đường dẫn đến file private key
	// Bắt buộc nếu TLS được bật
	KeyFile string `json:"key_file" yaml:"key_file"`

	// MinVersion xác định phiên bản TLS tối thiểu được chấp nhận
	// Giá trị hợp lệ: "1.0", "1.1", "1.2", "1.3"
	// Mặc định: "1.2"
	MinVersion string `json:"min_version" yaml:"min_version"`

	// MaxVersion xác định phiên bản TLS tối đa được chấp nhận
	// Giá trị hợp lệ: "1.0", "1.1", "1.2", "1.3"
	// Mặc định: "1.3"
	MaxVersion string `json:"max_version" yaml:"max_version"`

	// PreferServerCipherSuites xác định việc ưu tiên sử dụng các cipher suites của server hay client
	// Mặc định: true
	PreferServerCipherSuites bool `json:"prefer_server_cipher_suites" yaml:"prefer_server_cipher_suites"`

	// CipherSuites là danh sách các cipher suites được chấp nhận
	// Nếu rỗng, sử dụng các cipher suites mặc định
	CipherSuites []string `json:"cipher_suites" yaml:"cipher_suites"`

	// CurvePreferences là danh sách các elliptic curves được ưu tiên sử dụng
	// Nếu rỗng, sử dụng các curve mặc định
	CurvePreferences []string `json:"curve_preferences" yaml:"curve_preferences"`
}

// DefaultConfig trả về cấu hình mặc định cho HTTP adapter
func DefaultConfig() *Config {
	return &Config{}
}
