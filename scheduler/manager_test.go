package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/go-co-op/gocron"
)

func TestNewScheduler(t *testing.T) {
	scheduler := NewScheduler()
	if scheduler == nil {
		t.Fatal("NewScheduler() returned nil")
	}

	// Kiểm tra xem scheduler có implement Manager interface không
	var _ Manager = scheduler
}

func TestSchedulerFluentInterface(t *testing.T) {
	scheduler := NewScheduler()

	// Test fluent interface
	job, err := scheduler.Every(1).Second().Do(func() {
		// Empty job
	})

	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	if job == nil {
		t.Fatal("Job is nil")
	}
}

func TestSchedulerWithTag(t *testing.T) {
	scheduler := NewScheduler()

	// Tạo job với tag
	_, err := scheduler.Every(1).Second().Tag("test").Do(func() {
		// Empty job
	})

	if err != nil {
		t.Fatalf("Failed to create job with tag: %v", err)
	}

	// Tìm job theo tag
	jobs, err := scheduler.FindJobsByTag("test")
	if err != nil {
		t.Fatalf("Failed to find jobs by tag: %v", err)
	}

	if len(jobs) != 1 {
		t.Fatalf("Expected 1 job, got %d", len(jobs))
	}
}

func TestSchedulerCron(t *testing.T) {
	scheduler := NewScheduler()

	// Tạo job với cron expression
	_, err := scheduler.Cron("* * * * *").Do(func() {
		// Empty job
	})

	if err != nil {
		t.Fatalf("Failed to create cron job: %v", err)
	}
}

func TestSchedulerStartStop(t *testing.T) {
	scheduler := NewScheduler()

	// Kiểm tra ban đầu scheduler chưa chạy
	if scheduler.IsRunning() {
		t.Fatal("Scheduler should not be running initially")
	}

	// Start scheduler
	scheduler.StartAsync()

	// Chờ một chút để scheduler start
	time.Sleep(100 * time.Millisecond)

	// Kiểm tra scheduler đã chạy
	if !scheduler.IsRunning() {
		t.Fatal("Scheduler should be running after StartAsync()")
	}

	// Stop scheduler
	scheduler.Stop()

	// Chờ một chút để scheduler stop
	time.Sleep(100 * time.Millisecond)

	// Kiểm tra scheduler đã dừng
	if scheduler.IsRunning() {
		t.Fatal("Scheduler should not be running after Stop()")
	}
}

func TestSchedulerName(t *testing.T) {
	scheduler := NewScheduler()

	// Tạo job với tên
	job, err := scheduler.Every(1).Second().Name("test-job").Do(func() {
		// Empty job
	})

	if err != nil {
		t.Fatalf("Failed to create named job: %v", err)
	}

	if job == nil {
		t.Fatal("Job is nil")
	}
}

func TestSchedulerSingletonMode(t *testing.T) {
	scheduler := NewScheduler()

	// Tạo job với singleton mode
	_, err := scheduler.Every(1).Second().SingletonMode().Do(func() {
		// Empty job
	})

	if err != nil {
		t.Fatalf("Failed to create singleton job: %v", err)
	}
}

func TestSchedulerRemoveByTag(t *testing.T) {
	scheduler := NewScheduler()

	// Tạo job với tag
	_, err := scheduler.Every(1).Second().Tag("test", "remove").Do(func() {
		// Empty job
	})

	if err != nil {
		t.Fatalf("Failed to create job with tag: %v", err)
	}

	// Xóa job theo tag
	err = scheduler.RemoveByTag("test")
	if err != nil {
		t.Fatalf("Failed to remove job by tag: %v", err)
	}

	// Kiểm tra job đã bị xóa
	jobs, err := scheduler.FindJobsByTag("test")
	if err != nil && err.Error() != "gocron: no jobs found with given tag" {
		t.Fatalf("Failed to find jobs by tag: %v", err)
	}

	if len(jobs) != 0 {
		t.Fatalf("Expected 0 jobs after removal, got %d", len(jobs))
	}
}

func TestSchedulerClear(t *testing.T) {
	scheduler := NewScheduler()

	// Tạo một vài job
	_, err := scheduler.Every(1).Second().Do(func() {})
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	_, err = scheduler.Every(2).Second().Do(func() {})
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}

	// Clear tất cả jobs
	scheduler.Clear()

	// Kiểm tra không còn job nào
	gocronScheduler := scheduler.GetScheduler()
	if len(gocronScheduler.Jobs()) != 0 {
		t.Fatalf("Expected 0 jobs after Clear(), got %d", len(gocronScheduler.Jobs()))
	}
}

func TestSchedulerWithDistributedLocker(t *testing.T) {
	scheduler := NewScheduler()

	// Tạo mock locker
	mockLocker := &mockLocker{}

	// Test WithDistributedLocker trả về Manager
	result := scheduler.WithDistributedLocker(mockLocker)
	if result == nil {
		t.Fatal("WithDistributedLocker should return Manager")
	}

	// Kiểm tra result vẫn là cùng scheduler
	if result != scheduler {
		t.Fatal("WithDistributedLocker should return the same scheduler instance")
	}
}

func TestSchedulerRegisterEventListeners(t *testing.T) {
	scheduler := NewScheduler()

	// Tạo event listener function
	eventListener := func(j *gocron.Job) {
		// Mock event listener
	}

	// Test RegisterEventListeners không panic
	scheduler.RegisterEventListeners(eventListener)
}

// Mock implementations for testing
type mockLocker struct{}

func (m *mockLocker) Lock(ctx context.Context, key string) (gocron.Lock, error) {
	return &mockLock{}, nil
}

type mockLock struct{}

func (m *mockLock) Unlock(ctx context.Context) error {
	return nil
}
