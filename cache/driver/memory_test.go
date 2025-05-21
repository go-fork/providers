package driver

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMemoryDriverGet(t *testing.T) {
	t.Run("returns value when key exists and not expired", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"
		value := "test-value"

		err := d.Set(ctx, key, value, 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Act
		result, found := d.Get(ctx, key)

		// Assert
		if !found {
			t.Errorf("Expected to find key, but didn't")
		}
		if result != value {
			t.Errorf("Expected value to be %v, got %v", value, result)
		}
	})

	t.Run("returns not found when key doesn't exist", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Act
		_, found := d.Get(ctx, "nonexistent-key")

		// Assert
		if found {
			t.Errorf("Expected not to find key, but did")
		}
	})

	t.Run("returns not found when key is expired", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"
		value := "test-value"

		// Set with very short TTL to make it expire
		err := d.Set(ctx, key, value, 100*time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Wait for it to expire
		time.Sleep(200 * time.Millisecond)

		// Act
		_, found := d.Get(ctx, key)

		// Assert
		if found {
			t.Errorf("Expected not to find expired key, but did")
		}
	})
}

func TestMemoryDriverSet(t *testing.T) {
	t.Run("sets value successfully", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"
		value := "test-value"

		// Act
		err := d.Set(ctx, key, value, 1*time.Hour)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		result, found := d.Get(ctx, key)
		if !found {
			t.Errorf("Expected to find key, but didn't")
		}
		if result != value {
			t.Errorf("Expected value to be %v, got %v", value, result)
		}
	})

	t.Run("overwrites existing value", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"
		value1 := "test-value-1"
		value2 := "test-value-2"

		// Set initial value
		err := d.Set(ctx, key, value1, 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Act - overwrite
		err = d.Set(ctx, key, value2, 1*time.Hour)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		result, found := d.Get(ctx, key)
		if !found {
			t.Errorf("Expected to find key, but didn't")
		}
		if result != value2 {
			t.Errorf("Expected value to be %v, got %v", value2, result)
		}
	})

	t.Run("sets value with infinite TTL when negative", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"
		value := "test-value"

		// Act
		err := d.Set(ctx, key, value, -1*time.Hour)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		result, found := d.Get(ctx, key)
		if !found {
			t.Errorf("Expected to find key, but didn't")
		}
		if result != value {
			t.Errorf("Expected value to be %v, got %v", value, result)
		}
	})

	t.Run("sets complex value types", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"
		value := map[string]interface{}{
			"name":  "Test User",
			"age":   30,
			"roles": []string{"admin", "user"},
		}

		// Act
		err := d.Set(ctx, key, value, 1*time.Hour)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		result, found := d.Get(ctx, key)
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
			if resultMap["age"] != value["age"] {
				t.Errorf("Expected age to be %v, got %v", value["age"], resultMap["age"])
			}
		}
	})
}

func TestMemoryDriverHas(t *testing.T) {
	t.Run("returns true when key exists and not expired", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"

		err := d.Set(ctx, key, "test-value", 1*time.Hour)
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
		d := NewMemoryDriver()
		ctx := context.Background()

		// Act
		exists := d.Has(ctx, "nonexistent-key")

		// Assert
		if exists {
			t.Errorf("Expected key not to exist, but it does")
		}
	})

	t.Run("returns false when key is expired", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"

		// Set with short TTL to make it expire
		err := d.Set(ctx, key, "test-value", 100*time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Wait for it to expire
		time.Sleep(200 * time.Millisecond)

		// Act
		exists := d.Has(ctx, key)

		// Assert
		if exists {
			t.Errorf("Expected expired key not to exist, but it does")
		}
	})
}

func TestMemoryDriverDelete(t *testing.T) {
	t.Run("deletes existing key", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"

		err := d.Set(ctx, key, "test-value", 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Act
		err = d.Delete(ctx, key)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		exists := d.Has(ctx, key)
		if exists {
			t.Errorf("Expected key to be deleted, but it still exists")
		}
	})

	t.Run("doesn't error when key doesn't exist", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Act
		err := d.Delete(ctx, "nonexistent-key")

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestMemoryDriverFlush(t *testing.T) {
	t.Run("removes all keys", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Set multiple keys
		keys := []string{"key1", "key2", "key3"}
		for _, key := range keys {
			err := d.Set(ctx, key, "value-"+key, 1*time.Hour)
			if err != nil {
				t.Fatalf("Failed to set value: %v", err)
			}
		}

		// Verify keys exist
		for _, key := range keys {
			if !d.Has(ctx, key) {
				t.Fatalf("Expected key %s to exist before flush", key)
			}
		}

		// Act
		err := d.Flush(ctx)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify all keys are gone
		for _, key := range keys {
			if d.Has(ctx, key) {
				t.Errorf("Expected key %s to be deleted after flush", key)
			}
		}
	})

	t.Run("does nothing when cache is empty", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Act
		err := d.Flush(ctx)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestMemoryDriverGetMultiple(t *testing.T) {
	t.Run("gets multiple existing keys", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Set multiple keys
		keyValues := map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}

		for key, value := range keyValues {
			err := d.Set(ctx, key, value, 1*time.Hour)
			if err != nil {
				t.Fatalf("Failed to set value: %v", err)
			}
		}

		// Act
		values, missing := d.GetMultiple(ctx, []string{"key1", "key2", "key4"})

		// Assert
		if len(values) != 2 {
			t.Errorf("Expected 2 values, got %d", len(values))
		}

		if values["key1"] != "value1" {
			t.Errorf("Expected key1 value to be %v, got %v", "value1", values["key1"])
		}

		if values["key2"] != "value2" {
			t.Errorf("Expected key2 value to be %v, got %v", "value2", values["key2"])
		}

		if len(missing) != 1 || missing[0] != "key4" {
			t.Errorf("Expected missing keys to be [key4], got %v", missing)
		}
	})

	t.Run("returns empty map and all keys as missing when no keys exist", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		keys := []string{"key1", "key2", "key3"}

		// Act
		values, missing := d.GetMultiple(ctx, keys)

		// Assert
		if len(values) != 0 {
			t.Errorf("Expected 0 values, got %d", len(values))
		}

		if len(missing) != len(keys) {
			t.Errorf("Expected %d missing keys, got %d", len(keys), len(missing))
		}

		// Check all requested keys are in missing
		missingMap := make(map[string]bool)
		for _, key := range missing {
			missingMap[key] = true
		}

		for _, key := range keys {
			if !missingMap[key] {
				t.Errorf("Expected key %s to be in missing keys", key)
			}
		}
	})

	t.Run("skips expired keys", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Set one expired key
		err := d.Set(ctx, "expired", "value", 100*time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Wait for it to expire
		time.Sleep(200 * time.Millisecond)

		// Set one valid key
		err = d.Set(ctx, "valid", "value", 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Act
		values, missing := d.GetMultiple(ctx, []string{"expired", "valid"})

		// Assert
		if len(values) != 1 {
			t.Errorf("Expected 1 value, got %d", len(values))
		}

		if _, ok := values["valid"]; !ok {
			t.Errorf("Expected 'valid' key to be in values")
		}

		if len(missing) != 1 || missing[0] != "expired" {
			t.Errorf("Expected missing keys to be [expired], got %v", missing)
		}
	})
}

func TestMemoryDriverSetMultiple(t *testing.T) {
	t.Run("sets multiple values", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		values := map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": map[string]string{"nested": "value"},
		}

		// Act
		err := d.SetMultiple(ctx, values, 1*time.Hour)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Verify all keys were set
		for key, expectedValue := range values {
			value, found := d.Get(ctx, key)
			if !found {
				t.Errorf("Expected key %s to exist", key)
				continue
			}

			// For complex types like maps, we need to check more carefully
			switch key {
			case "key3":
				nestedMap, ok := value.(map[string]string)
				if !ok {
					t.Errorf("Expected key3 value to be map[string]string, got %T", value)
				} else if nestedMap["nested"] != "value" {
					t.Errorf("Expected nested value to be 'value', got %v", nestedMap["nested"])
				}
			default:
				if value != expectedValue {
					t.Errorf("Expected key %s value to be %v, got %v", key, expectedValue, value)
				}
			}
		}
	})

	t.Run("does nothing with empty map", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		values := map[string]interface{}{}

		// Act
		err := d.SetMultiple(ctx, values, 1*time.Hour)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("overwrites existing values", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Set initial values
		initial := map[string]interface{}{
			"key1": "old1",
			"key2": "old2",
		}

		err := d.SetMultiple(ctx, initial, 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to set initial values: %v", err)
		}

		// New values to overwrite
		updated := map[string]interface{}{
			"key1": "new1",
			"key3": "new3",
		}

		// Act
		err = d.SetMultiple(ctx, updated, 2*time.Hour)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// key1 should be updated
		value1, found := d.Get(ctx, "key1")
		if !found {
			t.Errorf("Expected key1 to exist")
		} else if value1 != "new1" {
			t.Errorf("Expected key1 value to be 'new1', got %v", value1)
		}

		// key2 should remain unchanged
		value2, found := d.Get(ctx, "key2")
		if !found {
			t.Errorf("Expected key2 to exist")
		} else if value2 != "old2" {
			t.Errorf("Expected key2 value to be 'old2', got %v", value2)
		}

		// key3 should be added
		value3, found := d.Get(ctx, "key3")
		if !found {
			t.Errorf("Expected key3 to exist")
		} else if value3 != "new3" {
			t.Errorf("Expected key3 value to be 'new3', got %v", value3)
		}
	})
}

func TestMemoryDriverDeleteMultiple(t *testing.T) {
	t.Run("deletes multiple existing keys", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Set multiple keys
		keys := []string{"key1", "key2", "key3"}
		for _, key := range keys {
			err := d.Set(ctx, key, "value-"+key, 1*time.Hour)
			if err != nil {
				t.Fatalf("Failed to set value: %v", err)
			}
		}

		// Verify keys exist
		for _, key := range keys {
			if !d.Has(ctx, key) {
				t.Fatalf("Expected key %s to exist before delete", key)
			}
		}

		// Act
		err := d.DeleteMultiple(ctx, []string{"key1", "key3"})

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// key1 and key3 should be deleted
		if d.Has(ctx, "key1") {
			t.Errorf("Expected key1 to be deleted")
		}
		if d.Has(ctx, "key3") {
			t.Errorf("Expected key3 to be deleted")
		}

		// key2 should still exist
		if !d.Has(ctx, "key2") {
			t.Errorf("Expected key2 to still exist")
		}
	})

	t.Run("doesn't error when keys don't exist", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Act
		err := d.DeleteMultiple(ctx, []string{"nonexistent1", "nonexistent2"})

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("doesn't error with empty key list", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Act
		err := d.DeleteMultiple(ctx, []string{})

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestMemoryDriverRemember(t *testing.T) {
	t.Run("returns existing value when key exists", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"
		existingValue := "existing-value"

		err := d.Set(ctx, key, existingValue, 1*time.Hour)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

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

	t.Run("calls callback and stores value when key doesn't exist", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"
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

		// Verify value was stored in cache
		cachedValue, found := d.Get(ctx, key)
		if !found {
			t.Errorf("Expected key to be stored in cache, but it wasn't")
		}
		if cachedValue != callbackValue {
			t.Errorf("Expected cached value to be %v, got %v", callbackValue, cachedValue)
		}
	})

	t.Run("returns error when callback returns error", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"
		expectedErr := errors.New("callback error")

		callback := func() (interface{}, error) {
			return nil, expectedErr
		}

		// Act
		_, err := d.Remember(ctx, key, 1*time.Hour, callback)

		// Assert
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}

		// Verify key was not stored
		_, found := d.Get(ctx, key)
		if found {
			t.Errorf("Expected key not to be stored when callback errors")
		}
	})

	t.Run("calls callback when key exists but is expired", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"
		existingValue := "existing-value"
		callbackValue := "callback-value"

		// Set with a short TTL to make it expire
		err := d.Set(ctx, key, existingValue, 100*time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Wait for it to expire
		time.Sleep(200 * time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

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
	})
}

func TestMemoryDriverStats(t *testing.T) {
	t.Run("returns correct stats", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Set some keys
		d.Set(ctx, "key1", "value1", 1*time.Hour)
		d.Set(ctx, "key2", "value2", 1*time.Hour)
		d.Set(ctx, "key3", "value3", -1*time.Second) // Expired immediately

		// Act
		stats := d.Stats(ctx)

		// Assert
		if stats["type"] != "memory" {
			t.Errorf("Expected driver to be 'memory', got %v", stats["type"])
		}

		// Should have 3 items (MemoryDriver doesn't automatically cleanup expired items until janitor runs)
		if stats["count"] != 3 {
			t.Errorf("Expected 3 items, got %v", stats["count"])
		}
	})

	t.Run("reflects changes in cache", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()

		// Initially empty
		initialStats := d.Stats(ctx)
		if initialStats["count"] != 0 {
			t.Errorf("Expected 0 items initially, got %v", initialStats["count"])
		}

		// Add some items
		d.Set(ctx, "key1", "value1", 1*time.Hour)
		d.Set(ctx, "key2", "value2", 1*time.Hour)

		// Check after adding
		afterAddStats := d.Stats(ctx)
		if afterAddStats["count"] != 2 {
			t.Errorf("Expected 2 items after adding, got %v", afterAddStats["count"])
		}

		// Delete an item
		d.Delete(ctx, "key1")

		// Check after deleting
		afterDeleteStats := d.Stats(ctx)
		if afterDeleteStats["count"] != 1 {
			t.Errorf("Expected 1 item after deleting, got %v", afterDeleteStats["count"])
		}

		// Flush
		d.Flush(ctx)

		// Check after flushing
		afterFlushStats := d.Stats(ctx)
		if afterFlushStats["count"] != 0 {
			t.Errorf("Expected 0 items after flushing, got %v", afterFlushStats["count"])
		}
	})
}

func TestMemoryDriverClose(t *testing.T) {
	t.Run("returns nil", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()

		// Act
		err := d.Close()

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestMemoryDriverExpiration(t *testing.T) {
	t.Run("auto-expires items", func(t *testing.T) {
		// Arrange
		d := NewMemoryDriver()
		ctx := context.Background()
		key := "test-key"
		value := "test-value"

		// Set with very short TTL
		err := d.Set(ctx, key, value, 50*time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to set value: %v", err)
		}

		// Verify key exists immediately
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
}
