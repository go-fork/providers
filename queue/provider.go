package queue

import (
	"github.com/go-fork/di"
	"github.com/redis/go-redis/v9"
)

// RedisOptions chứa các tùy chọn cấu hình cho Redis.
type RedisOptions struct {
	// Addr là địa chỉ của Redis server.
	Addr string

	// Password là mật khẩu của Redis server.
	Password string

	// DB là số của Redis database.
	DB int

	// TLS xác định liệu có sử dụng TLS khi kết nối đến Redis server hay không.
	TLS bool

	// PoolSize là số lượng kết nối tối đa được duy trì trong pool.
	PoolSize int
}

// Provider là service provider cho queue.
type Provider struct {
	redisOptions RedisOptions
	serverOpts   ServerOptions
	client       Client
	server       Server
	redisClient  redis.UniversalClient
}

// NewServiceProvider tạo một provider mới với các tùy chọn cấu hình.
func NewServiceProvider(redisOpts RedisOptions, serverOpts ServerOptions) *Provider {
	return &Provider{
		redisOptions: redisOpts,
		serverOpts:   serverOpts,
	}
}

// Register đăng ký các dịch vụ vào container.
func (p *Provider) Register(c *di.Container) {
	c.Bind("queue.client", func(c *di.Container) interface{} {
		return p.GetClient()
	})

	c.Bind("queue.server", func(c *di.Container) interface{} {
		return p.GetServer()
	})

	c.Bind("queue.redis", func(c *di.Container) interface{} {
		return p.GetRedisClient()
	})
}

// Boot khởi động các dịch vụ.
func (p *Provider) Boot(c *di.Container) {
	// Không cần thực hiện gì đặc biệt trong boot
}

// GetRedisClient trả về Redis client.
func (p *Provider) GetRedisClient() redis.UniversalClient {
	if p.redisClient == nil {
		p.redisClient = redis.NewClient(&redis.Options{
			Addr:     p.redisOptions.Addr,
			Password: p.redisOptions.Password,
			DB:       p.redisOptions.DB,
			PoolSize: p.redisOptions.PoolSize,
		})
	}
	return p.redisClient
}

// GetClient trả về Client.
func (p *Provider) GetClient() Client {
	if p.client == nil {
		p.client = NewClientWithUniversalClient(p.GetRedisClient())
	}
	return p.client
}

// GetServer trả về Server.
func (p *Provider) GetServer() Server {
	if p.server == nil {
		p.server = NewServer(p.GetRedisClient(), p.serverOpts)
	}
	return p.server
}
