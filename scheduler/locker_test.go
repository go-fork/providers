package scheduler

import (
	"context"
	"testing"

	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
)

func TestDefaultRedisLockerOptions(t *testing.T) {
	options := DefaultRedisLockerOptions()

	if options.KeyPrefix != "scheduler_lock:" {
		t.Errorf("Expected KeyPrefix 'scheduler_lock:', got '%s'", options.KeyPrefix)
	}

	if options.LockDuration != 30 {
		t.Errorf("Expected LockDuration 30, got %d", options.LockDuration)
	}

	if options.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries 3, got %d", options.MaxRetries)
	}

	if options.RetryDelay != 100 {
		t.Errorf("Expected RetryDelay 100, got %d", options.RetryDelay)
	}
}

func TestValidateRedisLockerOptions(t *testing.T) {
	tests := []struct {
		name    string
		options RedisLockerOptions
		wantErr bool
		errType error
	}{
		{
			name:    "valid options",
			options: DefaultRedisLockerOptions(),
			wantErr: false,
		},
		{
			name: "invalid lock duration - zero",
			options: RedisLockerOptions{
				KeyPrefix:    "test:",
				LockDuration: 0,
				MaxRetries:   3,
				RetryDelay:   100,
			},
			wantErr: true,
			errType: ErrInvalidLockDuration,
		},
		{
			name: "invalid lock duration - negative",
			options: RedisLockerOptions{
				KeyPrefix:    "test:",
				LockDuration: -1,
				MaxRetries:   3,
				RetryDelay:   100,
			},
			wantErr: true,
			errType: ErrInvalidLockDuration,
		},
		{
			name: "invalid max retries - negative",
			options: RedisLockerOptions{
				KeyPrefix:    "test:",
				LockDuration: 30,
				MaxRetries:   -1,
				RetryDelay:   100,
			},
			wantErr: true,
			errType: ErrInvalidMaxRetries,
		},
		{
			name: "invalid retry delay - negative",
			options: RedisLockerOptions{
				KeyPrefix:    "test:",
				LockDuration: 30,
				MaxRetries:   3,
				RetryDelay:   -1,
			},
			wantErr: true,
			errType: ErrInvalidRetryDelay,
		},
		{
			name: "invalid key prefix - empty",
			options: RedisLockerOptions{
				KeyPrefix:    "",
				LockDuration: 30,
				MaxRetries:   3,
				RetryDelay:   100,
			},
			wantErr: true,
			errType: ErrInvalidKeyPrefix,
		},
		{
			name: "max retries zero is valid",
			options: RedisLockerOptions{
				KeyPrefix:    "test:",
				LockDuration: 30,
				MaxRetries:   0,
				RetryDelay:   100,
			},
			wantErr: false,
		},
		{
			name: "retry delay zero is valid",
			options: RedisLockerOptions{
				KeyPrefix:    "test:",
				LockDuration: 30,
				MaxRetries:   3,
				RetryDelay:   0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRedisLockerOptions(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRedisLockerOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != tt.errType {
				t.Errorf("validateRedisLockerOptions() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestNewRedisLockerWithNilClient(t *testing.T) {
	_, err := NewRedisLocker(nil)
	if err != ErrRedisClientNil {
		t.Errorf("Expected ErrRedisClientNil, got %v", err)
	}
}

func TestNewRedisLockerWithInvalidOptions(t *testing.T) {
	// Skip test nếu không có Redis (vì chúng ta không thể tạo client thật)
	// Tạo mock client
	client := redis.NewClient(&redis.Options{
		Addr: "invalid:6379", // Invalid address to skip connection test in this case
	})

	invalidOptions := RedisLockerOptions{
		KeyPrefix:    "",
		LockDuration: -1,
		MaxRetries:   -1,
		RetryDelay:   -1,
	}

	_, err := NewRedisLocker(client, invalidOptions)
	if err == nil {
		t.Error("Expected error with invalid options")
	}
}

// Mock Redis Lock để test interface compliance
type mockRedisLock struct {
	unlocked bool
}

func (m *mockRedisLock) Unlock(ctx context.Context) error {
	m.unlocked = true
	return nil
}

// Mock Redis Locker để test interface compliance
type mockRedisLocker struct{}

func (m *mockRedisLocker) Lock(ctx context.Context, key string) (gocron.Lock, error) {
	return &mockRedisLock{}, nil
}

func TestRedisLockerInterface(t *testing.T) {
	var locker gocron.Locker = &mockRedisLocker{}

	ctx := context.Background()
	lock, err := locker.Lock(ctx, "test-key")
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}

	if lock == nil {
		t.Fatal("Lock is nil")
	}

	err = lock.Unlock(ctx)
	if err != nil {
		t.Fatalf("Failed to unlock: %v", err)
	}

	mockLock, ok := lock.(*mockRedisLock)
	if !ok {
		t.Fatal("Lock is not mockRedisLock")
	}

	if !mockLock.unlocked {
		t.Error("Lock was not unlocked")
	}
}
