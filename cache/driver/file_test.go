package driver

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func init() {
	// Register complex types for gob encoding
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	gob.Register([]string{})
}

// getFilename là hàm tiện ích cho tests
// Áp dụng cùng logic của FileDriver.keyToFilename
func (d *FileDriver) getFilename(key string) string {
	filename, _ := d.keyToFilename(key)
	return filename
}

func createTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "file-driver-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

func TestFileDriverInitialization(t *testing.T) {
	t.Run("creates directory if it doesn't exist", func(t *testing.T) {
		// Arrange
		baseDir := createTempDir(t)
		cacheDir := filepath.Join(baseDir, "nonexistent-dir")

		// Verify the directory doesn't exist yet
		if _, err := os.Stat(cacheDir); !os.IsNotExist(err) {
			t.Fatalf("Expected directory not to exist before test")
		}

		// Act
		_, err := NewFileDriver(cacheDir)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify the directory was created
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			t.Errorf("Expected directory to be created, but it doesn't exist")
		}
	})

	t.Run("returns error with invalid path", func(t *testing.T) {
		// Arrange
		invalidPath := string([]byte{0}) // Invalid path character

		// Act
		_, err := NewFileDriver(invalidPath)

		// Assert
		if err == nil {
			t.Errorf("Expected error for invalid path, got nil")
		}
	})
}

func TestFileDriverGetSet(t *testing.T) {
	t.Run("set and get value", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()
		key := "test-key"
		value := "test-value"

		// Act - Set
		err = d.Set(ctx, key, value, 1*time.Hour)

		// Assert - Set
		if err != nil {
			t.Errorf("Expected no error for Set, got %v", err)
		}

		// Verify file was created
		filePath := d.getFilename(key)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected cache file to be created, but it doesn't exist")
		}

		// Act - Get
		result, found := d.Get(ctx, key)

		// Assert - Get
		if !found {
			t.Errorf("Expected to find key, but didn't")
		}
		if result != value {
			t.Errorf("Expected value to be %v, got %v", value, result)
		}
	})

	t.Run("returns not found when key doesn't exist", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()

		// Act
		_, found := d.Get(ctx, "nonexistent-key")

		// Assert
		if found {
			t.Errorf("Expected not to find key, but did")
		}
	})

	t.Run("returns not found when cache file is corrupted", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()
		key := "corrupted-key"

		// Create a corrupted cache file (invalid gob)
		filePath := d.getFilename(key)
		err = os.WriteFile(filePath, []byte("invalid-gob-data"), 0644)
		if err != nil {
			t.Fatalf("Failed to write corrupted file: %v", err)
		}

		// Act
		_, found := d.Get(ctx, key)

		// Assert
		if found {
			t.Errorf("Expected not to find corrupted key, but did")
		}
	})

	t.Run("complex value types", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()
		key := "complex-key"
		value := map[string]interface{}{
			"name":  "Test User",
			"age":   30,
			"roles": []string{"admin", "user"},
		}

		// Act - Set
		err = d.Set(ctx, key, value, 1*time.Hour)

		// Assert - Set
		if err != nil {
			t.Errorf("Expected no error for Set, got %v", err)
		}

		// Act - Get
		result, found := d.Get(ctx, key)

		// Assert - Get
		if !found {
			t.Errorf("Expected to find key, but didn't")
		}

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			t.Errorf("Expected result to be map[string]interface{}, got %T", result)
		} else {
			if resultMap["name"] != value["name"] {
				t.Errorf("Expected name to be %v, got %v", value["name"], resultMap["name"])
			}

			// Note: When encoding/decoding with gob, the type might be preserved
			age, ok := resultMap["age"]
			if !ok {
				t.Errorf("Expected age to be present in result map")
			} else if age != value["age"] {
				t.Errorf("Expected age to be %v, got %v", value["age"], age)
			}

			// Check roles array
			roles, ok := resultMap["roles"]
			if !ok {
				t.Errorf("Expected roles to be present in result map")
			} else {
				switch r := roles.(type) {
				case []string:
					// With gob encoding, the type might be preserved as []string
					rolesOriginal := value["roles"].([]string)
					if len(r) != len(rolesOriginal) {
						t.Errorf("Expected roles length to be %d, got %d", len(rolesOriginal), len(r))
					} else {
						for i, role := range rolesOriginal {
							if r[i] != role {
								t.Errorf("Expected role at index %d to be %v, got %v", i, role, r[i])
							}
						}
					}
				case []interface{}:
					// Alternatively, it might be []interface{}
					rolesOriginal := value["roles"].([]string)
					if len(r) != len(rolesOriginal) {
						t.Errorf("Expected roles length to be %d, got %d", len(rolesOriginal), len(r))
					} else {
						for i, role := range rolesOriginal {
							if r[i] != role {
								t.Errorf("Expected role at index %d to be %v, got %v", i, role, r[i])
							}
						}
					}
				default:
					t.Errorf("Expected roles to be []string or []interface{}, got %T", roles)
				}
			}
		}
	})
}

func TestFileDriverExpiration(t *testing.T) {
	t.Run("returns not found when item is expired", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()
		key := "expired-key"
		value := "expired-value"

		// Set with very short TTL
		err = d.Set(ctx, key, value, 50*time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Verify it exists initially
		if !d.Has(ctx, key) {
			t.Fatalf("Expected key to exist immediately after setting")
		}

		// Wait for expiration
		time.Sleep(100 * time.Millisecond)

		// Act & Assert
		if d.Has(ctx, key) {
			t.Errorf("Expected key to be expired, but it still exists")
		}
	})

	t.Run("never expires when TTL is negative", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()
		key := "forever-key"
		value := "forever-value"

		// Set with negative TTL for no expiration
		err = d.Set(ctx, key, value, -1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Act
		result, found := d.Get(ctx, key)

		// Assert
		if !found {
			t.Errorf("Expected to find key with negative TTL, but didn't")
		}
		if result != value {
			t.Errorf("Expected value to be %v, got %v", value, result)
		}
	})
}

func TestFileDriverHas(t *testing.T) {
	t.Run("returns true when key exists and not expired", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()
		key := "test-key"

		err = d.Set(ctx, key, "test-value", 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Act
		exists := d.Has(ctx, key)

		// Assert
		if !exists {
			t.Errorf("Expected key to exist, but it doesn't")
		}
	})

	t.Run("returns false when key doesn't exist", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()

		// Act
		exists := d.Has(ctx, "nonexistent-key")

		// Assert
		if exists {
			t.Errorf("Expected key not to exist, but it does")
		}
	})
}

func TestFileDriverDelete(t *testing.T) {
	t.Run("deletes existing key", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()
		key := "test-key"

		err = d.Set(ctx, key, "test-value", 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Verify file exists
		filePath := d.getFilename(key)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("Expected cache file to exist before delete")
		}

		// Act
		err = d.Delete(ctx, key)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify file was deleted
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Errorf("Expected cache file to be deleted, but it still exists")
		}
	})

	t.Run("doesn't error when key doesn't exist", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()

		// Act
		err = d.Delete(ctx, "nonexistent-key")

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestFileDriverFlush(t *testing.T) {
	t.Run("removes all cache files", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()

		// Set multiple keys
		keys := []string{"key1", "key2", "key3"}
		for _, key := range keys {
			err := d.Set(ctx, key, "value-"+key, 1*time.Hour)
			if err != nil {
				t.Fatalf("Failed to set value: %v", err)
			}
		}

		// Verify files exist
		for _, key := range keys {
			filePath := d.getFilename(key)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Fatalf("Expected cache file for %s to exist before flush", key)
			}
		}

		// Act
		err = d.Flush(ctx)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify directory still exists but is empty
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			t.Errorf("Expected cache directory to still exist after flush")
		}

		entries, err := os.ReadDir(cacheDir)
		if err != nil {
			t.Fatalf("Failed to read cache directory: %v", err)
		}

		if len(entries) != 0 {
			t.Errorf("Expected cache directory to be empty after flush, found %d entries", len(entries))
		}
	})

	t.Run("handles empty directory", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()

		// Act
		err = d.Flush(ctx)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("returns error if directory does not exist", func(t *testing.T) {
		d := &FileDriver{directory: "/path/does/not/exist"}
		err := d.Flush(context.Background())
		if err == nil {
			t.Errorf("Expected error when flushing non-existent directory, got nil")
		}
	})

	t.Run("returns error if cannot read directory", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d := &FileDriver{directory: cacheDir}
		// Đổi quyền thư mục để không đọc được
		os.Chmod(cacheDir, 0000)
		defer os.Chmod(cacheDir, 0755)
		err := d.Flush(context.Background())
		if err == nil {
			t.Errorf("Expected error when cannot read directory, got nil")
		}
	})

	t.Run("returns error if cannot remove file", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, _ := NewFileDriver(cacheDir)
		ctx := context.Background()
		key := "locked-file"
		d.Set(ctx, key, "value", 1*time.Hour)
		os.Chmod(cacheDir, 0000) // Không cho phép đọc/xóa file trong thư mục
		defer os.Chmod(cacheDir, 0755)
		err := d.Flush(ctx)
		if err == nil {
			t.Errorf("Expected error when cannot remove file in Flush, got nil")
		}

		// Tạo thư mục con không thể xóa file bên trong
		lockedDir := filepath.Join(cacheDir, "locked")
		os.Mkdir(lockedDir, 0700)
		lockedFile := filepath.Join(lockedDir, "locked-file")
		os.WriteFile(lockedFile, []byte("value"), 0644)
		os.Chmod(lockedDir, 0500) // Không cho phép ghi/xóa file bên trong
		defer os.Chmod(lockedDir, 0700)
		d.directory = cacheDir // Đảm bảo driver trỏ đúng thư mục
		err = d.Flush(ctx)
		if err == nil {
			t.Errorf("Expected error when cannot remove file in Flush, got nil")
		}
	})
}

func TestFileDriverMultipleOperations(t *testing.T) {
	t.Run("GetMultiple returns values and missing keys", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()

		// Set some keys
		d.Set(ctx, "key1", "value1", 1*time.Hour)
		d.Set(ctx, "key2", "value2", 1*time.Hour)
		// key3 will be missing

		// Act
		values, missing := d.GetMultiple(ctx, []string{"key1", "key2", "key3"})

		// Assert
		if len(values) != 2 {
			t.Errorf("Expected 2 values, got %d", len(values))
		}
		if values["key1"] != "value1" {
			t.Errorf("Expected key1 value to be 'value1', got %v", values["key1"])
		}
		if values["key2"] != "value2" {
			t.Errorf("Expected key2 value to be 'value2', got %v", values["key2"])
		}
		if len(missing) != 1 || missing[0] != "key3" {
			t.Errorf("Expected missing keys to be [key3], got %v", missing)
		}
	})

	t.Run("SetMultiple stores multiple values", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()
		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		// Act
		err = d.SetMultiple(ctx, values, 1*time.Hour)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify all files were created
		for key := range values {
			filePath := d.getFilename(key)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected cache file for %s to exist", key)
			}
		}

		// Verify values can be retrieved
		for key, expectedValue := range values {
			value, found := d.Get(ctx, key)
			if !found {
				t.Errorf("Expected to find key %s, but didn't", key)
			} else if value != expectedValue {
				t.Errorf("Expected value for key %s to be %v, got %v", key, expectedValue, value)
			}
		}
	})

	t.Run("returns error if any Set fails in SetMultiple", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}
		ctx := context.Background()
		// Tạo file không ghi được (ví dụ: file đã tồn tại và không có quyền ghi)
		badKey := string([]byte{0}) // Tên file không hợp lệ trên hầu hết hệ thống
		values := map[string]interface{}{
			"key1": "value1",
			badKey: "bad",
			"key2": "value2",
		}
		err = d.SetMultiple(ctx, values, 1*time.Hour)
		if err == nil {
			t.Errorf("Expected error when Set fails in SetMultiple, got nil")
		}
	})

	t.Run("DeleteMultiple removes multiple keys", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()

		// Set keys
		keysToSet := []string{"key1", "key2", "key3", "key4"}
		for _, key := range keysToSet {
			d.Set(ctx, key, "value-"+key, 1*time.Hour)
		}

		// Verify all keys exist
		for _, key := range keysToSet {
			if !d.Has(ctx, key) {
				t.Fatalf("Expected key %s to exist before delete", key)
			}
		}

		// Act
		keysToDelete := []string{"key1", "key3"}
		err = d.DeleteMultiple(ctx, keysToDelete)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify deleted keys are gone
		for _, key := range keysToDelete {
			if d.Has(ctx, key) {
				t.Errorf("Expected key %s to be deleted, but it still exists", key)
			}
		}

		// Verify other keys still exist
		if !d.Has(ctx, "key2") || !d.Has(ctx, "key4") {
			t.Errorf("Expected keys 'key2' and 'key4' to still exist")
		}
	})

	t.Run("returns error if any Delete fails in DeleteMultiple", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}
		ctx := context.Background()
		// Tạo key hợp lệ và key không hợp lệ
		d.Set(ctx, "ok", "value", 1*time.Hour)
		badKey := string([]byte{0}) // Tên file không hợp lệ
		err = d.DeleteMultiple(ctx, []string{"ok", badKey})
		if err == nil {
			t.Errorf("Expected error when Delete fails in DeleteMultiple, got nil")
		}
	})
}

func TestFileDriverRemember(t *testing.T) {
	t.Run("returns existing value when key exists", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()
		key := "remember-key"
		existingValue := "existing-value"

		d.Set(ctx, key, existingValue, 1*time.Hour)

		callbackCalled := false
		callback := func() (interface{}, error) {
			callbackCalled = true
			return "callback-value", nil
		}

		// Act
		value, err := d.Remember(ctx, key, 1*time.Hour, callback)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if value != existingValue {
			t.Errorf("Expected value to be %v, got %v", existingValue, value)
		}
		if callbackCalled {
			t.Errorf("Expected callback not to be called, but it was")
		}
	})

	t.Run("calls callback and stores when key doesn't exist", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()
		key := "missing-key"
		callbackValue := "callback-value"

		callbackCalled := false
		callback := func() (interface{}, error) {
			callbackCalled = true
			return callbackValue, nil
		}

		// Act
		value, err := d.Remember(ctx, key, 1*time.Hour, callback)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if value != callbackValue {
			t.Errorf("Expected value to be %v, got %v", callbackValue, value)
		}
		if !callbackCalled {
			t.Errorf("Expected callback to be called, but it wasn't")
		}

		// Verify value was stored
		storedValue, found := d.Get(ctx, key)
		if !found {
			t.Errorf("Expected key to be stored, but it wasn't")
		}
		if storedValue != callbackValue {
			t.Errorf("Expected stored value to be %v, got %v", callbackValue, storedValue)
		}
	})

	t.Run("returns error if callback fails", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}
		ctx := context.Background()
		key := "fail-callback"
		callback := func() (interface{}, error) {
			return nil, fmt.Errorf("callback error")
		}
		_, err = d.Remember(ctx, key, 1*time.Hour, callback)
		if err == nil {
			t.Errorf("Expected error from callback, got nil")
		}
	})

	t.Run("returns error if Set fails after callback", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}
		ctx := context.Background()
		badKey := string([]byte{0})
		callback := func() (interface{}, error) {
			return "value", nil
		}
		_, err = d.Remember(ctx, badKey, 1*time.Hour, callback)
		if err == nil {
			t.Errorf("Expected error from Set in Remember, got nil")
		}
	})
}

func TestFileDriverStats(t *testing.T) {
	t.Run("returns correct stats", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		ctx := context.Background()

		// Set some keys
		d.Set(ctx, "key1", "value1", 1*time.Hour)
		d.Set(ctx, "key2", "value2", 1*time.Hour)
		d.Set(ctx, "key3", "value3", 1*time.Hour)

		// Set one expired key (but file still exists)
		d.Set(ctx, "expired-key", "value", -1*time.Second)

		// Act
		stats := d.Stats(ctx)

		// Assert
		if stats["type"] != "file" {
			t.Errorf("Expected driver to be 'file', got %v", stats["type"])
		}

		// Should have 4 items (including expired one, since files still exist)
		if stats["count"] != 4 {
			t.Errorf("Expected 4 items, got %v", stats["count"])
		}

		if stats["path"] != cacheDir {
			t.Errorf("Expected path to be %v, got %v", cacheDir, stats["path"])
		}
	})
}

func TestFileDriverClose(t *testing.T) {
	t.Run("returns nil", func(t *testing.T) {
		// Arrange
		cacheDir := createTempDir(t)
		d, err := NewFileDriver(cacheDir)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}

		// Act
		err = d.Close()

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestFileDriverDeleteExpired(t *testing.T) {
	t.Run("removes expired files only", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, err := NewFileDriverWithOptions(cacheDir, 1*time.Second, 0)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}
		ctx := context.Background()

		// Set 1 expired (ttl=1ms), 1 valid, 1 no-expiry
		d.Set(ctx, "expired", "expired-value", 1*time.Millisecond)
		d.Set(ctx, "valid", "valid-value", 1*time.Hour)
		d.Set(ctx, "forever", "forever-value", -1)

		time.Sleep(10 * time.Millisecond) // Đảm bảo expired đã hết hạn
		d.deleteExpired()

		if _, found := d.Get(ctx, "expired"); found {
			t.Errorf("Expected 'expired' to be deleted, but it still exists")
		}
		if _, found := d.Get(ctx, "valid"); !found {
			t.Errorf("Expected 'valid' to exist, but it was deleted")
		}
		if _, found := d.Get(ctx, "forever"); !found {
			t.Errorf("Expected 'forever' to exist, but it was deleted")
		}
	})

	t.Run("handles corrupted file gracefully", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, err := NewFileDriverWithOptions(cacheDir, 1*time.Second, 0)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}
		// Tạo file rác không decode được
		file := filepath.Join(cacheDir, "corrupted")
		os.WriteFile(file, []byte("not-gob-data"), 0644)
		// Không panic, không xóa file hợp lệ
		d.Set(context.Background(), "valid", "ok", 1*time.Hour)
		d.deleteExpired()
		if _, found := d.Get(context.Background(), "valid"); !found {
			t.Errorf("Expected 'valid' to exist after deleteExpired with corrupted file")
		}
	})

	t.Run("deleteExpired handles error when opening file", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, err := NewFileDriverWithOptions(cacheDir, 1*time.Second, 0)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}
		// Tạo file không mở được (directory thay vì file)
		badFile := filepath.Join(cacheDir, "badfile")
		os.Mkdir(badFile, 0755)
		d.deleteExpired() // Không panic, không xóa file hợp lệ
	})

	t.Run("deleteExpired handles error when removing file", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, err := NewFileDriverWithOptions(cacheDir, 1*time.Millisecond, 0)
		if err != nil {
			t.Fatalf("Failed to create file driver: %v", err)
		}
		ctx := context.Background()
		d.Set(ctx, "expired", "expired-value", 1*time.Millisecond)
		time.Sleep(10 * time.Millisecond)
		filePath := d.getFilename("expired")
		os.Chmod(filePath, 0400) // Không xóa được
		defer os.Chmod(filePath, 0644)
		d.deleteExpired() // Không panic
	})
}

func TestFileDriver_ErrorPaths(t *testing.T) {
	t.Run("Set returns error if keyToFilename fails", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, _ := NewFileDriver(cacheDir)
		ctx := context.Background()
		badKey := string([]byte{0}) // invalid filename
		err := d.Set(ctx, badKey, "value", 1*time.Hour)
		if err == nil {
			t.Errorf("Expected error for invalid key in Set, got nil")
		}
	})

	t.Run("Get returns not found if keyToFilename fails", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, _ := NewFileDriver(cacheDir)
		ctx := context.Background()
		badKey := string([]byte{0})
		_, found := d.Get(ctx, badKey)
		if found {
			t.Errorf("Expected not found for invalid key in Get, got found")
		}
	})

	t.Run("Delete returns error if keyToFilename fails", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, _ := NewFileDriver(cacheDir)
		ctx := context.Background()
		badKey := string([]byte{0})
		err := d.Delete(ctx, badKey)
		if err == nil {
			t.Errorf("Expected error for invalid key in Delete, got nil")
		}
	})

	t.Run("Flush returns error if cannot remove file", func(t *testing.T) {
		cacheDir := createTempDir(t)
		d, _ := NewFileDriver(cacheDir)
		ctx := context.Background()
		key := "locked-file"
		d.Set(ctx, key, "value", 1*time.Hour)
		os.Chmod(cacheDir, 0000) // Không cho phép đọc/xóa file trong thư mục
		defer os.Chmod(cacheDir, 0755)
		err := d.Flush(ctx)
		if err == nil {
			t.Errorf("Expected error when cannot remove file in Flush, got nil")
		}

		// Tạo thư mục con không thể xóa file bên trong
		lockedDir := filepath.Join(cacheDir, "locked")
		os.Mkdir(lockedDir, 0700)
		lockedFile := filepath.Join(lockedDir, "locked-file")
		os.WriteFile(lockedFile, []byte("value"), 0644)
		os.Chmod(lockedDir, 0500) // Không cho phép ghi/xóa file bên trong
		defer os.Chmod(lockedDir, 0700)
		d.directory = cacheDir // Đảm bảo driver trỏ đúng thư mục
		err = d.Flush(ctx)
		if err == nil {
			t.Errorf("Expected error when cannot remove file in Flush, got nil")
		}
	})
}
