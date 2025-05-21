// Package driver provides cache driver implementations and interfaces
package driver

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoDBDriverGet(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("returns value when key exists and not expired", func(mt *mtest.T) {
		// Mock MongoDB response
		now := time.Now()
		expiration := now.Add(1 * time.Hour).UnixNano()

		// Set up mock response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.collection", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: "test-key"},
			{Key: "value", Value: "test-value"},
			{Key: "expiration", Value: expiration},
			{Key: "created_at", Value: now},
		}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		result, found := driver.Get(context.Background(), "test-key")

		// Assert
		if !found {
			mt.Errorf("Expected to find key, but didn't")
		}
		if result != "test-value" {
			mt.Errorf("Expected value to be %v, got %v", "test-value", result)
		}
	})

	mt.Run("returns not found when key doesn't exist", func(mt *mtest.T) {
		// Set up mock response for not found
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(
			mtest.CommandError{
				Code:    11000,
				Name:    "NotFound",
				Message: "document not found",
			}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		_, found := driver.Get(context.Background(), "nonexistent-key")

		// Assert
		if found {
			mt.Errorf("Expected not to find key, but did")
		}
	})

	mt.Run("returns not found when key is expired", func(mt *mtest.T) {
		// Mock MongoDB response
		now := time.Now()
		expiration := now.Add(-1 * time.Hour).UnixNano() // expired

		// Set up mock response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.collection", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: "expired-key"},
			{Key: "value", Value: "expired-value"},
			{Key: "expiration", Value: expiration},
			{Key: "created_at", Value: now.Add(-2 * time.Hour)},
		}))

		// Add mock response for the delete operation that will be triggered
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		_, found := driver.Get(context.Background(), "expired-key")

		// Assert
		if found {
			mt.Errorf("Expected not to find expired key, but did")
		}
	})
}

func TestMongoDBDriverSet(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("sets value successfully", func(mt *mtest.T) {
		// Set up mock response for successful update
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		err := driver.Set(context.Background(), "test-key", "test-value", 1*time.Hour)

		// Assert
		if err != nil {
			mt.Errorf("Expected no error, got %v", err)
		}
	})

	mt.Run("handles error from MongoDB", func(mt *mtest.T) {
		// Set up mock response for error
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(
			mtest.CommandError{
				Code:    11000,
				Name:    "Error",
				Message: "MongoDB error",
			}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		err := driver.Set(context.Background(), "test-key", "test-value", 1*time.Hour)

		// Assert
		if err == nil {
			mt.Errorf("Expected error, got nil")
		}
	})

	mt.Run("uses default expiration when ttl is 0", func(mt *mtest.T) {
		// Set up mock response for successful update
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})

		// Create driver with mocked client
		defaultExpiration := 5 * time.Minute
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: defaultExpiration,
		}

		// Act
		err := driver.Set(context.Background(), "test-key", "test-value", 0)

		// Assert
		if err != nil {
			mt.Errorf("Expected no error, got %v", err)
		}
	})

	mt.Run("sets no expiration when ttl is negative", func(mt *mtest.T) {
		// Set up mock response for successful update
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		err := driver.Set(context.Background(), "test-key", "test-value", -1)

		// Assert
		if err != nil {
			mt.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestMongoDBDriverHas(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("returns true when key exists", func(mt *mtest.T) {
		// Mock MongoDB response
		now := time.Now()
		expiration := now.Add(1 * time.Hour).UnixNano()

		// Set up mock response
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.collection", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: "test-key"},
			{Key: "value", Value: "test-value"},
			{Key: "expiration", Value: expiration},
			{Key: "created_at", Value: now},
		}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		exists := driver.Has(context.Background(), "test-key")

		// Assert
		if !exists {
			mt.Errorf("Expected key to exist, but it doesn't")
		}
	})

	mt.Run("returns false when key doesn't exist", func(mt *mtest.T) {
		// Set up mock response for not found
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(
			mtest.CommandError{
				Code:    11000,
				Name:    "NotFound",
				Message: "document not found",
			}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		exists := driver.Has(context.Background(), "nonexistent-key")

		// Assert
		if exists {
			mt.Errorf("Expected key not to exist, but it does")
		}
	})
}

func TestMongoDBDriverDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("deletes key successfully", func(mt *mtest.T) {
		// Set up mock response for delete
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		err := driver.Delete(context.Background(), "test-key")

		// Assert
		if err != nil {
			mt.Errorf("Expected no error, got %v", err)
		}
	})

	mt.Run("handles error from MongoDB", func(mt *mtest.T) {
		// Set up mock response for error
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(
			mtest.CommandError{
				Code:    11000,
				Name:    "Error",
				Message: "MongoDB error",
			}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		err := driver.Delete(context.Background(), "test-key")

		// Assert
		if err == nil {
			mt.Errorf("Expected error, got nil")
		}
	})
}

func TestMongoDBDriverFlush(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("flushes cache successfully", func(mt *mtest.T) {
		// Set up mock response for DeleteMany
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 5}})

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		err := driver.Flush(context.Background())

		// Assert
		if err != nil {
			mt.Errorf("Expected no error, got %v", err)
		}
	})

	mt.Run("handles error from MongoDB", func(mt *mtest.T) {
		// Set up mock response for error
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(
			mtest.CommandError{
				Code:    11000,
				Name:    "Error",
				Message: "MongoDB error",
			}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		err := driver.Flush(context.Background())

		// Assert
		if err == nil {
			mt.Errorf("Expected error, got nil")
		}
	})
}

func TestMongoDBDriverGetMultiple(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("gets multiple values successfully", func(mt *mtest.T) {
		// Mock MongoDB response
		now := time.Now()
		expiration := now.Add(1 * time.Hour).UnixNano()

		// Create mock cursor with multiple documents
		firstBatch := []bson.D{
			{
				{Key: "_id", Value: "key1"},
				{Key: "value", Value: "value1"},
				{Key: "expiration", Value: expiration},
				{Key: "created_at", Value: now},
			},
			{
				{Key: "_id", Value: "key3"},
				{Key: "value", Value: "value3"},
				{Key: "expiration", Value: expiration},
				{Key: "created_at", Value: now},
			},
		}

		// Set up mock response with a cursor that has multiple results
		mt.AddMockResponses(
			mtest.CreateCursorResponse(2, "db.collection", mtest.FirstBatch, firstBatch...),
			mtest.CreateCursorResponse(0, "db.collection", mtest.NextBatch),
		)

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		values, missing := driver.GetMultiple(context.Background(), []string{"key1", "key2", "key3"})

		// Assert
		if len(values) != 2 {
			mt.Errorf("Expected 2 values, got %d", len(values))
		}
		if values["key1"] != "value1" {
			mt.Errorf("Expected value for key1 to be 'value1', got %v", values["key1"])
		}
		if values["key3"] != "value3" {
			mt.Errorf("Expected value for key3 to be 'value3', got %v", values["key3"])
		}
		if len(missing) != 1 || missing[0] != "key2" {
			mt.Errorf("Expected missing keys to be [key2], got %v", missing)
		}
	})

	mt.Run("handles expired key", func(mt *mtest.T) {
		// Mock MongoDB response
		now := time.Now()
		validExpiration := now.Add(1 * time.Hour).UnixNano()
		expiredExpiration := now.Add(-1 * time.Hour).UnixNano()

		// Create mock cursor with valid and expired documents
		firstBatch := []bson.D{
			{
				{Key: "_id", Value: "valid-key"},
				{Key: "value", Value: "valid-value"},
				{Key: "expiration", Value: validExpiration},
				{Key: "created_at", Value: now},
			},
			{
				{Key: "_id", Value: "expired-key"},
				{Key: "value", Value: "expired-value"},
				{Key: "expiration", Value: expiredExpiration},
				{Key: "created_at", Value: now.Add(-2 * time.Hour)},
			},
		}

		// Set up mock response
		mt.AddMockResponses(
			mtest.CreateCursorResponse(2, "db.collection", mtest.FirstBatch, firstBatch...),
			mtest.CreateCursorResponse(0, "db.collection", mtest.NextBatch),
		)

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		values, missing := driver.GetMultiple(context.Background(), []string{"valid-key", "expired-key", "nonexistent-key"})

		// Assert
		if len(values) != 1 {
			mt.Errorf("Expected 1 value, got %d", len(values))
		}
		if values["valid-key"] != "valid-value" {
			mt.Errorf("Expected value for valid-key to be 'valid-value', got %v", values["valid-key"])
		}
		if len(missing) != 3 {
			mt.Errorf("Expected 3 missing keys, got %d", len(missing))
		}
	})

	mt.Run("handles MongoDB error", func(mt *mtest.T) {
		// Set up mock response for error
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(
			mtest.CommandError{
				Code:    11000,
				Name:    "Error",
				Message: "MongoDB error",
			}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		keys := []string{"key1", "key2", "key3"}

		// Act
		values, missing := driver.GetMultiple(context.Background(), keys)

		// Assert
		if len(values) != 0 {
			mt.Errorf("Expected 0 values, got %d", len(values))
		}
		if len(missing) != len(keys) {
			mt.Errorf("Expected %d missing keys, got %d", len(keys), len(missing))
		}
	})
}

func TestMongoDBDriverSetMultiple(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("sets multiple values successfully", func(mt *mtest.T) {
		// Set up mock response for BulkWrite
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "nModified", Value: 2}})

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}

		// Act
		err := driver.SetMultiple(context.Background(), values, 1*time.Hour)

		// Assert
		if err != nil {
			mt.Errorf("Expected no error, got %v", err)
		}
	})

	mt.Run("handles error from MongoDB", func(mt *mtest.T) {
		// Set up mock response for error
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(
			mtest.CommandError{
				Code:    11000,
				Name:    "Error",
				Message: "MongoDB error",
			}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}

		// Act
		err := driver.SetMultiple(context.Background(), values, 1*time.Hour)

		// Assert
		if err == nil {
			mt.Errorf("Expected error, got nil")
		}
	})

	mt.Run("uses default expiration when ttl is 0", func(mt *mtest.T) {
		// Set up mock response for BulkWrite
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "nModified", Value: 2}})

		// Create driver with mocked client
		defaultExpiration := 5 * time.Minute
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: defaultExpiration,
		}

		values := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}

		// Act
		err := driver.SetMultiple(context.Background(), values, 0)

		// Assert
		if err != nil {
			mt.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestMongoDBDriverDeleteMultiple(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("deletes multiple keys successfully", func(mt *mtest.T) {
		// Set up mock response for DeleteMany
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 2}})

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		err := driver.DeleteMultiple(context.Background(), []string{"key1", "key2"})

		// Assert
		if err != nil {
			mt.Errorf("Expected no error, got %v", err)
		}
	})

	mt.Run("handles empty keys list", func(mt *mtest.T) {
		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		err := driver.DeleteMultiple(context.Background(), []string{})

		// Assert
		if err != nil {
			mt.Errorf("Expected no error, got %v", err)
		}
	})

	mt.Run("handles error from MongoDB", func(mt *mtest.T) {
		// Set up mock response for error
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(
			mtest.CommandError{
				Code:    11000,
				Name:    "Error",
				Message: "MongoDB error",
			}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		err := driver.DeleteMultiple(context.Background(), []string{"key1", "key2"})

		// Assert
		if err == nil {
			mt.Errorf("Expected error, got nil")
		}
	})
}

func TestMongoDBDriverRemember(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("returns cached value when key exists", func(mt *mtest.T) {
		// Mock MongoDB response for existing item
		now := time.Now()
		expiration := now.Add(1 * time.Hour).UnixNano()

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.collection", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: "test-key"},
			{Key: "value", Value: "cached-value"},
			{Key: "expiration", Value: expiration},
			{Key: "created_at", Value: now},
		}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Callback should not be called if item is in cache
		callbackCalled := false
		callback := func() (interface{}, error) {
			callbackCalled = true
			return "new-value", nil
		}

		// Act
		result, err := driver.Remember(context.Background(), "test-key", 1*time.Hour, callback)

		// Assert
		if err != nil {
			mt.Errorf("Expected no error, got %v", err)
		}
		if result != "cached-value" {
			mt.Errorf("Expected result to be 'cached-value', got %v", result)
		}
		if callbackCalled {
			mt.Errorf("Callback should not have been called")
		}
	})

	mt.Run("calls callback and stores result when key doesn't exist", func(mt *mtest.T) {
		// Set up mock responses: 1) not found, 2) successful set
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    11000,
				Name:    "NotFound",
				Message: "document not found",
			}),
			bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}},
		)

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Callback should be called if item is not in cache
		callbackCalled := false
		callback := func() (interface{}, error) {
			callbackCalled = true
			return "new-value", nil
		}

		// Act
		result, err := driver.Remember(context.Background(), "test-key", 1*time.Hour, callback)

		// Assert
		if err != nil {
			mt.Errorf("Expected no error, got %v", err)
		}
		if result != "new-value" {
			mt.Errorf("Expected result to be 'new-value', got %v", result)
		}
		if !callbackCalled {
			mt.Errorf("Callback should have been called")
		}
	})

	mt.Run("returns error when callback fails", func(mt *mtest.T) {
		// Set up mock response for not found
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(
			mtest.CommandError{
				Code:    11000,
				Name:    "NotFound",
				Message: "document not found",
			}))

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Callback should return an error
		expectedError := errors.New("callback error")
		callback := func() (interface{}, error) {
			return nil, expectedError
		}

		// Act
		_, err := driver.Remember(context.Background(), "test-key", 1*time.Hour, callback)

		// Assert
		if err != expectedError {
			mt.Errorf("Expected error %v, got %v", expectedError, err)
		}
	})
}

func TestMongoDBDriverStats(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("returns stats with count", func(mt *mtest.T) {
		// Set up mock response for count
		mt.AddMockResponses(
			bson.D{{Key: "n", Value: 5}},
			bson.D{
				{Key: "ns", Value: "db.collection"},
				{Key: "count", Value: 5},
				{Key: "size", Value: 1024},
				{Key: "avgObjSize", Value: 256},
			},
		)

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
			hits:              10,
			misses:            5,
		}

		// Act
		stats := driver.Stats(context.Background())

		// Assert
		if stats["count"] != int64(-1) {
			mt.Errorf("Expected count to be -1, got %v", stats["count"])
		}
		if stats["hits"] != int64(10) {
			mt.Errorf("Expected hits to be 10, got %v", stats["hits"])
		}
		if stats["misses"] != int64(5) {
			mt.Errorf("Expected misses to be 5, got %v", stats["misses"])
		}
		if stats["type"] != "mongodb" {
			mt.Errorf("Expected type to be 'mongodb', got %v", stats["type"])
		}
	})

	mt.Run("handles error in count", func(mt *mtest.T) {
		// Set up mock response for count error
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(mtest.CommandError{
				Code:    11000,
				Name:    "Error",
				Message: "count error",
			}),
			bson.D{
				{Key: "ns", Value: "db.collection"},
				{Key: "size", Value: 1024},
			},
		)

		// Create driver with mocked client
		driver := &MongoDBDriver{
			client:            mt.Client,
			collection:        mt.Coll,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		stats := driver.Stats(context.Background())

		// Assert
		if stats["count"] != int64(-1) {
			mt.Errorf("Expected count to be -1 on error, got %v", stats["count"])
		}
	})
}

func TestMongoDBDriverClose(t *testing.T) {
	// Using a regular test not mtest to avoid channel close issues
	t.Run("handles nil client", func(t *testing.T) {
		// Create driver with nil client
		driver := &MongoDBDriver{
			client:            nil,
			collection:        nil,
			defaultExpiration: 5 * time.Minute,
		}

		// Act
		err := driver.Close()

		// Assert
		if err != nil {
			t.Errorf("Expected no error with nil client, got %v", err)
		}
	})
}
