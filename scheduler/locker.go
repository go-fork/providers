package scheduler

import (
	"context"
	"errors"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
)

// RedisLockerOptions đã được di chuyển vào config.go

// redisLocker triển khai gocron.Locker interface sử dụng Redis làm backend.
type redisLocker struct {
	client  *redis.Client
	options RedisLockerOptionsTime
}

// redisLock triển khai gocron.Lock interface.
type redisLock struct {
	locker       *redisLocker
	key          string
	cancelRenew  context.CancelFunc
	renewContext context.Context
}

// NewRedisLocker tạo một Redis Locker mới để sử dụng với gocron.
// Nó có thể được chuyển vào phương thức WithDistributedLocker của scheduler.
//
// Example:
//
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: "localhost:6379",
//	})
//	locker, err := scheduler.NewRedisLocker(redisClient)
//	if err != nil {
//		log.Fatal(err)
//	}
//	sched.WithDistributedLocker(locker)
func NewRedisLocker(client *redis.Client, opts ...RedisLockerOptions) (gocron.Locker, error) {
	if client == nil {
		return nil, ErrRedisClientNil
	}

	// Sử dụng tùy chọn mặc định
	options := DefaultRedisLockerOptions()

	// Nếu có tùy chọn được cung cấp, sử dụng chúng
	if len(opts) > 0 {
		options = opts[0]

		// Validate các giá trị options
		if err := validateRedisLockerOptions(options); err != nil {
			return nil, err
		}
	}

	// Chuyển đổi options thành time.Duration
	timeOptions := options.ToTimeDuration()

	// Kiểm tra kết nối đến Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, ErrFailedToConnectToRedis
	}

	// Tạo locker
	locker := &redisLocker{
		client:  client,
		options: timeOptions,
	}

	return locker, nil
}

// Lock triển khai phương thức Lock của gocron.Locker interface.
func (r *redisLocker) Lock(ctx context.Context, key string) (gocron.Lock, error) {
	fullKey := r.options.KeyPrefix + key
	retries := 0

	for {
		// Cố gắng set key với expiration
		success, err := r.client.SetNX(ctx, fullKey, "locked", r.options.LockDuration).Result()

		// Nếu có lỗi không liên quan đến kết nối
		if err != nil && err != redis.ErrClosed && err != context.Canceled {
			return nil, err
		}

		// Nếu lock thành công
		if success {
			renewCtx, cancelFn := context.WithCancel(context.Background())
			lock := &redisLock{
				locker:       r,
				key:          key,
				renewContext: renewCtx,
				cancelRenew:  cancelFn,
			}

			// Bắt đầu quá trình tự động gia hạn khóa
			go lock.startRenewLoop()

			return lock, nil
		}

		// Nếu đã thử tối đa số lần
		if retries >= r.options.MaxRetries {
			return nil, ErrFailedToAcquireLock
		}

		// Chờ một khoảng thời gian trước khi thử lại
		select {
		case <-time.After(r.options.RetryDelay):
			retries++
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// startRenewLoop bắt đầu một goroutine để tự động gia hạn khóa trước khi hết hạn.
// Điều này ngăn khóa hết hạn trong khi job vẫn đang chạy.
func (r *redisLock) startRenewLoop() {
	renewInterval := r.locker.options.LockDuration / 3 * 2 // Gia hạn sau 2/3 thời gian hết hạn
	ticker := time.NewTicker(renewInterval)
	defer ticker.Stop()

	fullKey := r.locker.options.KeyPrefix + r.key

	for {
		select {
		case <-r.renewContext.Done():
			return
		case <-ticker.C:
			// Gia hạn khóa bằng cách đặt thời gian hết hạn mới
			// Sử dụng context với timeout để tránh block vô hạn
			ctx, cancel := context.WithTimeout(r.renewContext, 5*time.Second)
			err := r.locker.client.Expire(ctx, fullKey, r.locker.options.LockDuration).Err()
			cancel()
			if err != nil {
				// Log lỗi nếu cần thiết, nhưng không làm gián đoạn vòng lặp
				continue
			}
		}
	}
}

// Unlock triển khai phương thức Unlock của gocron.Lock interface.
func (r *redisLock) Unlock(ctx context.Context) error {
	// Dừng vòng lặp gia hạn trước
	if r.cancelRenew != nil {
		r.cancelRenew()
	}

	// Sau đó xóa khóa từ Redis
	fullKey := r.locker.options.KeyPrefix + r.key
	return r.locker.client.Del(ctx, fullKey).Err()
}

// validateRedisLockerOptions kiểm tra tính hợp lệ của các tùy chọn Redis Locker.
func validateRedisLockerOptions(options RedisLockerOptions) error {
	if options.LockDuration <= 0 {
		return ErrInvalidLockDuration
	}
	if options.MaxRetries < 0 {
		return ErrInvalidMaxRetries
	}
	if options.RetryDelay < 0 {
		return ErrInvalidRetryDelay
	}
	if options.KeyPrefix == "" {
		return ErrInvalidKeyPrefix
	}
	return nil
}

// Error constants
var (
	// ErrRedisClientNil được trả về khi Redis client nil.
	ErrRedisClientNil = errors.New("scheduler: redis client is nil")

	// ErrFailedToConnectToRedis được trả về khi không thể kết nối đến Redis.
	ErrFailedToConnectToRedis = errors.New("scheduler: failed to connect to redis")

	// ErrFailedToAcquireLock được trả về khi không thể lấy lock sau số lần thử tối đa.
	ErrFailedToAcquireLock = errors.New("scheduler: failed to acquire lock after maximum retries")

	// ErrInvalidLockDuration được trả về khi LockDuration không hợp lệ.
	ErrInvalidLockDuration = errors.New("scheduler: invalid lock duration")

	// ErrInvalidMaxRetries được trả về khi MaxRetries không hợp lệ.
	ErrInvalidMaxRetries = errors.New("scheduler: invalid max retries")

	// ErrInvalidRetryDelay được trả về khi RetryDelay không hợp lệ.
	ErrInvalidRetryDelay = errors.New("scheduler: invalid retry delay")

	// ErrInvalidKeyPrefix được trả về khi KeyPrefix không hợp lệ.
	ErrInvalidKeyPrefix = errors.New("scheduler: invalid key prefix")
)
