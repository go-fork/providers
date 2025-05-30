# Hướng dẫn sử dụng Mocks trong Cache Package

## Tổng quan

Package cache có hai loại mock objects để testing:

1. **MockManager** - Mock cho interface `Manager`
2. **MockDriver** - Mock cho interface `Driver`

Cả hai đều được tạo bởi [mockery](https://github.com/vektra/mockery) và sử dụng [testify/mock](https://github.com/stretchr/testify).

## MockManager

### Tạo Mock Instance

```go
import (
    "testing"
    "go.fork.vn/providers/cache/mocks"
)

func TestSomething(t *testing.T) {
    mockManager := mocks.NewMockManager(t)
    // ... setup expectations và test
}
```

### Ví dụ sử dụng MockManager

```go
func TestManagerOperations(t *testing.T) {
    t.Run("Test Get method", func(t *testing.T) {
        // Arrange
        mockManager := mocks.NewMockManager(t)
        mockManager.EXPECT().Get("test-key").Return("test-value", true)
        
        // Act
        value, found := mockManager.Get("test-key")
        
        // Assert
        assert.True(t, found)
        assert.Equal(t, "test-value", value)
    })
    
    t.Run("Test Set method", func(t *testing.T) {
        // Arrange
        mockManager := mocks.NewMockManager(t)
        mockManager.EXPECT().Set("key", "value", 5*time.Minute).Return(nil)
        
        // Act
        err := mockManager.Set("key", "value", 5*time.Minute)
        
        // Assert
        assert.NoError(t, err)
    })
}
```

## MockDriver

### Tạo Mock Instance

```go
func TestWithMockDriver(t *testing.T) {
    mockDriver := mocks.NewMockDriver(t)
    // ... setup expectations và test
}
```

### Ví dụ sử dụng MockDriver với Manager

```go
func TestManagerWithMockDriver(t *testing.T) {
    t.Run("Get returns value when default driver is set", func(t *testing.T) {
        // Arrange
        mockDriver := mocks.NewMockDriver(t)
        mockDriver.EXPECT().Get(context.Background(), "test-key").Return("test-value", true)

        manager := NewManager()
        manager.AddDriver("mock", mockDriver)
        manager.SetDefaultDriver("mock")

        // Act
        value, found := manager.Get("test-key")

        // Assert
        assert.True(t, found)
        assert.Equal(t, "test-value", value)
    })
}
```

## Các phương thức Mock có sẵn

### MockManager Methods
- `Get(key string) (interface{}, bool)`
- `Set(key string, value interface{}, ttl time.Duration) error`
- `Has(key string) bool`
- `Delete(key string) error`
- `DeleteMultiple(keys []string) error`
- `GetMultiple(keys []string) (map[string]interface{}, []string)`
- `SetMultiple(values map[string]interface{}, ttl time.Duration) error`
- `Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (interface{}, error)`
- `Flush() error`
- `Close() error`
- `Driver(name string) (driver.Driver, error)`
- `AddDriver(name string, driver driver.Driver)`
- `SetDefaultDriver(name string)`
- `Stats() map[string]map[string]interface{}`

### MockDriver Methods
- `Get(ctx context.Context, key string) (interface{}, bool)`
- `Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error`
- `Has(ctx context.Context, key string) bool`
- `Delete(ctx context.Context, key string) error`
- `DeleteMultiple(ctx context.Context, keys []string) error`
- `GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, []string)`
- `SetMultiple(ctx context.Context, values map[string]interface{}, ttl time.Duration) error`
- `Flush(ctx context.Context) error`
- `Close() error`
- `Stats() map[string]interface{}`

## Expecter Pattern

Mockery tạo ra các mock với expecter pattern, cho phép viết test dễ đọc hơn:

```go
// Setup expectation
mockManager.EXPECT().Get("key").Return("value", true)

// Setup multiple expectations
mockDriver.EXPECT().Set(mock.Anything, "key", "value", mock.Anything).Return(nil)
mockDriver.EXPECT().Get(mock.Anything, "key").Return("value", true)
```

## Chạy test

```bash
# Chạy tất cả tests
go test -v

# Chạy test cụ thể
go test -v -run TestManagerWithMockDriver
go test -v -run TestManagerWithMockery

# Chạy với coverage
go test -v -cover
```

## Regenerate Mocks

Để tạo lại mocks khi interface thay đổi:

```bash
mockery --config .mockery.yaml
```
