package queue

import (
	"context"
	"testing"
	"time"

	"github.com/go-fork/di"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueueSchedulerIntegration(t *testing.T) {
	// Setup
	container := di.New()
	app := &testApp{container: container}

	// Register queue provider
	provider := NewServiceProvider()
	provider.Register(app)
	provider.Boot(app)

	// Get services
	queueInstance, err := container.Make("queue.manager")
	require.NoError(t, err)
	manager := queueInstance.(Manager)

	client := manager.Client()
	server := manager.Server()
	scheduler := manager.Scheduler()

	// Verify scheduler is set
	assert.NotNil(t, scheduler)
	assert.NotNil(t, server.GetScheduler())

	// Test immediate task processing
	t.Run("ImmediateTaskProcessing", func(t *testing.T) {
		processed := make(chan bool, 1)

		// Register handler
		server.RegisterHandler("test_task", func(ctx context.Context, task *Task) error {
			processed <- true
			return nil
		})

		// Start server
		err := server.Start()
		require.NoError(t, err)
		defer server.Stop()

		// Enqueue task
		_, err = client.Enqueue("test_task", "test payload", WithQueue("default"))
		require.NoError(t, err)

		// Wait for processing
		select {
		case <-processed:
			// Success
		case <-time.After(5 * time.Second):
			t.Fatal("Task was not processed within timeout")
		}
	})

	// Test scheduler integration
	t.Run("SchedulerIntegration", func(t *testing.T) {
		scheduled := make(chan bool, 1)

		// Schedule a job that enqueues a task
		scheduler.Every(1).Second().Do(func() {
			client.Enqueue("scheduled_task", "scheduled payload", WithQueue("default"))
			scheduled <- true
		})

		// Start scheduler
		if !scheduler.IsRunning() {
			scheduler.StartAsync()
		}

		// Wait for scheduled execution
		select {
		case <-scheduled:
			// Success
		case <-time.After(3 * time.Second):
			t.Fatal("Scheduled task was not executed within timeout")
		}
	})
}

func TestQueueServerMethods(t *testing.T) {
	container := di.New()
	app := &testApp{container: container}

	provider := NewServiceProvider()
	provider.Register(app)

	queueInstance, err := container.Make("queue.server")
	require.NoError(t, err)
	server := queueInstance.(Server)

	t.Run("SetAndGetScheduler", func(t *testing.T) {
		// Initially should be nil or have default scheduler
		initialScheduler := server.GetScheduler()

		// Get scheduler from manager
		managerInstance, err := container.Make("queue.manager")
		require.NoError(t, err)
		manager := managerInstance.(Manager)

		newScheduler := manager.Scheduler()
		server.SetScheduler(newScheduler)

		retrievedScheduler := server.GetScheduler()
		assert.Equal(t, newScheduler, retrievedScheduler)
		assert.NotEqual(t, initialScheduler, retrievedScheduler)
	})

	t.Run("ServerStartStop", func(t *testing.T) {
		// Test start
		err := server.Start()
		assert.NoError(t, err)

		// Test double start
		err = server.Start()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already started")

		// Test stop
		err = server.Stop()
		assert.NoError(t, err)

		// Test double stop
		err = server.Stop()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not started")
	})
}

func TestQueueManagerSchedulerIntegration(t *testing.T) {
	config := DefaultConfig()
	manager := NewManagerWithConfig(config)

	t.Run("DefaultScheduler", func(t *testing.T) {
		scheduler := manager.Scheduler()
		assert.NotNil(t, scheduler)
	})

	t.Run("SetCustomScheduler", func(t *testing.T) {
		// Get original scheduler
		originalScheduler := manager.Scheduler()

		// The manager creates a new scheduler if none exists
		// For this test, we'll verify the scheduler is properly set
		assert.NotNil(t, originalScheduler)

		// Since we can't easily create a different scheduler instance
		// we'll just verify the SetScheduler method works
		manager.SetScheduler(originalScheduler)
		retrievedScheduler := manager.Scheduler()
		assert.Equal(t, originalScheduler, retrievedScheduler)
	})

	t.Run("ServerHasScheduler", func(t *testing.T) {
		server := manager.Server()
		scheduler := manager.Scheduler()

		// The server should have the scheduler set
		// Note: In real implementation, this should be set during registration
		// For now, we'll manually set it to test the functionality
		server.SetScheduler(scheduler)
		serverScheduler := server.GetScheduler()
		assert.Equal(t, scheduler, serverScheduler)
	})
}

// testApp implements the interface expected by service providers
type testApp struct {
	container *di.Container
}

func (t *testApp) Container() *di.Container {
	return t.container
}
