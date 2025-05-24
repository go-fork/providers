package queue

import (
	"context"
	"testing"
	"time"

	"github.com/go-fork/di"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderCleanupAndRetry(t *testing.T) {
	// Setup
	container := di.New()
	app := &testApp{container: container}

	// Register queue provider
	provider := NewServiceProvider()
	queueProvider := provider.(*ServiceProvider) // Cast to access internal methods
	queueProvider.Register(app)

	// Get services
	queueInstance, err := container.Make("queue.manager")
	require.NoError(t, err)
	manager := queueInstance.(Manager)

	adapter := manager.Adapter("")

	t.Run("CleanupOldDeadLetterTasks", func(t *testing.T) {
		ctx := context.Background()
		deadLetterQueueName := "default:dead"

		// Create old dead letter task (older than 7 days)
		oldTask := &DeadLetterTask{
			Task: Task{
				ID:        "old-task-1",
				Name:      "test_task",
				Queue:     "default",
				CreatedAt: time.Now().AddDate(0, 0, -8), // 8 days ago
			},
			Reason:   "test failure",
			FailedAt: time.Now().AddDate(0, 0, -8), // 8 days ago
		}

		// Create recent dead letter task (within 7 days)
		recentTask := &DeadLetterTask{
			Task: Task{
				ID:        "recent-task-1",
				Name:      "test_task",
				Queue:     "default",
				CreatedAt: time.Now().AddDate(0, 0, -1), // 1 day ago
			},
			Reason:   "test failure",
			FailedAt: time.Now().AddDate(0, 0, -1), // 1 day ago
		}

		// Enqueue both tasks
		err := adapter.Enqueue(ctx, deadLetterQueueName, oldTask)
		require.NoError(t, err)
		err = adapter.Enqueue(ctx, deadLetterQueueName, recentTask)
		require.NoError(t, err)

		// Verify both tasks are in the queue
		size, err := adapter.Size(ctx, deadLetterQueueName)
		require.NoError(t, err)
		assert.Equal(t, int64(2), size)

		// Run cleanup
		queueProvider.cleanupFailedJobs(manager)

		// Verify only recent task remains
		size, err = adapter.Size(ctx, deadLetterQueueName)
		require.NoError(t, err)
		assert.Equal(t, int64(1), size)

		// Verify the remaining task is the recent one
		var remainingTask DeadLetterTask
		err = adapter.Dequeue(ctx, deadLetterQueueName, &remainingTask)
		require.NoError(t, err)
		assert.Equal(t, "recent-task-1", remainingTask.Task.ID)

		// Clean up
		adapter.Clear(ctx, deadLetterQueueName)
	})

	t.Run("RetryReadyTasks", func(t *testing.T) {
		ctx := context.Background()
		retryQueueName := "default:retry"
		pendingQueueName := "default:pending"

		// Create task ready for retry (ProcessAt in the past)
		readyTask := &Task{
			ID:         "ready-task-1",
			Name:       "test_task",
			Queue:      "default",
			MaxRetry:   3,
			RetryCount: 1,
			CreatedAt:  time.Now().Add(-1 * time.Hour),
			ProcessAt:  time.Now().Add(-5 * time.Minute), // 5 minutes ago
		}

		// Create task not ready for retry (ProcessAt in the future)
		notReadyTask := &Task{
			ID:         "not-ready-task-1",
			Name:       "test_task",
			Queue:      "default",
			MaxRetry:   3,
			RetryCount: 1,
			CreatedAt:  time.Now().Add(-1 * time.Hour),
			ProcessAt:  time.Now().Add(5 * time.Minute), // 5 minutes from now
		}

		// Enqueue both tasks in retry queue
		err := adapter.Enqueue(ctx, retryQueueName, readyTask)
		require.NoError(t, err)
		err = adapter.Enqueue(ctx, retryQueueName, notReadyTask)
		require.NoError(t, err)

		// Verify both tasks are in retry queue
		retrySize, err := adapter.Size(ctx, retryQueueName)
		require.NoError(t, err)
		assert.Equal(t, int64(2), retrySize)

		// Verify pending queue is empty
		pendingSize, err := adapter.Size(ctx, pendingQueueName)
		require.NoError(t, err)
		assert.Equal(t, int64(0), pendingSize)

		// Run retry processing
		queueProvider.retryFailedJobs(manager)

		// Verify ready task moved to pending queue
		pendingSize, err = adapter.Size(ctx, pendingQueueName)
		require.NoError(t, err)
		assert.Equal(t, int64(1), pendingSize)

		// Verify not-ready task remains in retry queue
		retrySize, err = adapter.Size(ctx, retryQueueName)
		require.NoError(t, err)
		assert.Equal(t, int64(1), retrySize)

		// Verify the task in pending queue is the ready one
		var pendingTask Task
		err = adapter.Dequeue(ctx, pendingQueueName, &pendingTask)
		require.NoError(t, err)
		assert.Equal(t, "ready-task-1", pendingTask.ID)

		// Verify the task in retry queue is the not-ready one
		var retryTask Task
		err = adapter.Dequeue(ctx, retryQueueName, &retryTask)
		require.NoError(t, err)
		assert.Equal(t, "not-ready-task-1", retryTask.ID)

		// Clean up
		adapter.Clear(ctx, retryQueueName)
		adapter.Clear(ctx, pendingQueueName)
	})
}
