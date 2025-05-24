// Package queue - Example integration with scheduler
package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-fork/di"
)

// ExampleQueueWithScheduler demonstrates how to use queue with scheduler integration
func ExampleQueueWithScheduler() {
	// Tạo DI container
	container := di.New()

	// Tạo mock app với container
	app := &exampleApp{container: container}

	// Đăng ký queue service provider (tự động đăng ký scheduler)
	queueProvider := NewServiceProvider()
	queueProvider.Register(app)
	queueProvider.Boot(app)

	// Lấy queue manager từ container
	queueInstance, err := container.Make("queue.manager")
	if err != nil {
		log.Fatal("Failed to get queue manager:", err)
	}
	manager := queueInstance.(Manager)

	// Lấy client và server
	client := manager.Client()
	server := manager.Server()

	// Đăng ký handlers cho server
	server.RegisterHandler("send_email", func(ctx context.Context, task *Task) error {
		log.Printf("Processing email task: %s", task.Payload)
		// Simulate email sending
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	server.RegisterHandler("process_data", func(ctx context.Context, task *Task) error {
		log.Printf("Processing data task: %s", task.Payload)
		// Simulate data processing
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	// Khởi động server
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
	defer server.Stop()

	// Demo 1: Enqueue immediate tasks
	log.Println("\n=== Demo 1: Immediate Tasks ===")
	tasks := []struct {
		name    string
		payload string
	}{
		{"send_email", "Welcome email for user 1"},
		{"send_email", "Welcome email for user 2"},
		{"process_data", "Process user registration data"},
	}

	for _, task := range tasks {
		if _, err := client.Enqueue(task.name, task.payload, WithQueue("default")); err != nil {
			log.Printf("Failed to enqueue task: %v", err)
		} else {
			log.Printf("Enqueued task: %s", task.name)
		}
	}

	// Demo 2: Schedule delayed tasks với scheduler
	log.Println("\n=== Demo 2: Scheduled Tasks ===")
	scheduler := manager.Scheduler()

	// Schedule một task chạy sau 5 giây
	scheduler.Every(5).Seconds().Do(func() {
		if _, err := client.Enqueue("send_email", "Scheduled email reminder", WithQueue("default")); err != nil {
			log.Printf("Failed to enqueue scheduled task: %v", err)
		} else {
			log.Printf("Enqueued scheduled task: send_email")
		}
	})

	// Schedule một task chạy mỗi 10 giây
	scheduler.Every(10).Seconds().Do(func() {
		if _, err := client.Enqueue("process_data", "Periodic data cleanup", WithQueue("default")); err != nil {
			log.Printf("Failed to enqueue periodic task: %v", err)
		} else {
			log.Printf("Enqueued periodic task: process_data")
		}
	})

	// Chờ để xem các task được xử lý
	log.Println("\n=== Processing Tasks (waiting 30 seconds) ===")
	time.Sleep(30 * time.Second)

	log.Println("\n=== Example completed ===")
}

// exampleApp implements the interface expected by service providers
type exampleApp struct {
	container *di.Container
}

func (m *exampleApp) Container() *di.Container {
	return m.container
}

// ExampleScheduledJobWithQueue demonstrates scheduling jobs that enqueue tasks
func ExampleScheduledJobWithQueue() {
	container := di.New()
	app := &exampleApp{container: container}

	// Register services
	queueProvider := NewServiceProvider()
	queueProvider.Register(app)
	queueProvider.Boot(app)

	// Get manager
	queueInstance, _ := container.Make("queue.manager")
	manager := queueInstance.(Manager)

	client := manager.Client()
	server := manager.Server()
	scheduler := manager.Scheduler()

	// Register job handlers
	server.RegisterHandler("daily_report", func(ctx context.Context, task *Task) error {
		log.Printf("Generating daily report: %s", task.Payload)
		return nil
	})

	server.RegisterHandler("weekly_backup", func(ctx context.Context, task *Task) error {
		log.Printf("Running weekly backup: %s", task.Payload)
		return nil
	})

	// Start server
	server.Start()
	defer server.Stop()

	// Schedule recurring jobs
	// Daily report at 9:00 AM
	scheduler.Every(1).Days().At("09:00").Do(func() {
		reportData := fmt.Sprintf("Report for %s", time.Now().Format("2006-01-02"))
		client.Enqueue("daily_report", reportData, WithQueue("reports"))
		log.Println("Scheduled daily report task")
	})

	// Weekly backup every Sunday at 2:00 AM
	scheduler.Every(1).Weeks().At("02:00").Do(func() {
		backupData := fmt.Sprintf("Backup for week %s", time.Now().Format("2006-W02"))
		client.Enqueue("weekly_backup", backupData, WithQueue("maintenance"))
		log.Println("Scheduled weekly backup task")
	})

	// For demo purposes, manually trigger some tasks
	log.Println("Manually triggering scheduled tasks for demo...")

	// Trigger daily report
	client.Enqueue("daily_report", "Manual daily report", WithQueue("reports"))

	// Trigger weekly backup
	client.Enqueue("weekly_backup", "Manual weekly backup", WithQueue("maintenance"))

	// Wait for processing
	time.Sleep(5 * time.Second)

	log.Println("Scheduled job example completed")
}
