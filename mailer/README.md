# Mailer Provider

Mailer Provider là giải pháp gửi email đơn giản và mạnh mẽ cho ứng dụng Go, được xây dựng dựa trên thư viện [gopkg.in/gomail.v2](https://gopkg.in/gomail.v2).

## Tính năng nổi bật

- Tích hợp đầy đủ với DI container của ứng dụng, tương thích với di v0.0.5
- API fluent và dễ sử dụng
- Hỗ trợ gửi cả email văn bản thuần túy và HTML
- Hỗ trợ render template từ `text/template` và `html/template`
- Hỗ trợ file đính kèm và nhúng hình ảnh
- Hỗ trợ xử lý hàng đợi (queue) cho email
- Dễ dàng test với MockMailer
- Triển khai đầy đủ các methods Requires và Providers cho di v0.0.5

## Cài đặt

Để cài đặt Mailer Provider, bạn có thể sử dụng lệnh go get:

```bash
go get github.com/go-fork/providers/mailer
```

## Cách sử dụng

### 1. Đăng ký Service Provider

```go
package main

import (
    "github.com/go-fork/di"
    "github.com/go-fork/providers/mailer"
)

func main() {
    app := di.New()
    
    // Đăng ký service provider với cấu hình SMTP cơ bản
    app.Register(mailer.NewProvider(&mailer.Config{
        Host:       "smtp.example.com",
        Port:       587,
        Username:   "user@example.com",
        Password:   "password",
        Encryption: "tls",
        FromAddress: "noreply@example.com",
        FromName:   "Example App",
    }))
    
    // Hoặc với hỗ trợ queue
    queueSettings := mailer.DefaultQueueSettings()
    queueSettings.QueueEnabled = true
    queueSettings.Workers = 2

    app.Register(mailer.NewProvider(&mailer.Config{
        Host:       "smtp.example.com",
        Port:       587,
        Username:   "user@example.com",
        Password:   "password",
        Encryption: "tls",
        FromAddress: "noreply@example.com",
        FromName:   "Example App",
    }).WithQueue(queueSettings))
    
    // Khởi động ứng dụng
    app.Boot()
}
```

### 2. Gửi Email Đơn Giản

```go
// Lấy mailer từ container
container := app.Container()
m := container.MustMake("mailer").(mailer.Mailer)

// Tạo và gửi email
message := m.NewMessage().
    To("recipient@example.com").
    Subject("Hello from Go-Fork!").
    Text("This is a plain text message").
    HTML("<h1>Hello</h1><p>This is an HTML message</p>")

if err := m.Send(message); err != nil {
    log.Fatal(err)
}
```

### 3. Gửi Email với Template

```go
// Tạo và gửi email với template
data := map[string]interface{}{
    "Username": "John Doe",
    "ResetLink": "https://example.com/reset?token=abc123",
}

message := m.NewMessage().
    To("user@example.com").
    Subject("Password Reset").
    TextTemplate("Hello {{.Username}}, please reset your password using this link: {{.ResetLink}}").
    HTMLTemplate(`
        <h1>Password Reset</h1>
        <p>Hello {{.Username}},</p>
        <p>Please reset your password by clicking the link below:</p>
        <p><a href="{{.ResetLink}}">Reset Password</a></p>
    `).
    WithData(data)

if err := m.Send(message); err != nil {
    log.Fatal(err)
}
```

### 4. Gửi Email với File Đính Kèm

```go
message := m.NewMessage().
    To("user@example.com").
    Subject("Your Invoice").
    Text("Please find your invoice attached").
    Attach("/path/to/invoice.pdf", "Invoice-May-2025.pdf")

if err := m.Send(message); err != nil {
    log.Fatal(err)
}
```

### 5. Nhúng Hình Ảnh vào Email

```go
message := m.NewMessage().
    To("user@example.com").
    Subject("Check out our logo").
    HTML(`<h1>Our Logo</h1><p><img src="cid:logo.png" /></p>`).
    Embed("/path/to/logo.png", "logo.png")

if err := m.Send(message); err != nil {
    log.Fatal(err)
}
```

### 6. Sử Dụng Hàng Đợi Email

```go
// Lấy mailer từ container
container := app.Container()
m := container.Get("mailer").(mailer.Mailer)

// Gửi email qua hàng đợi (tự động xử lý nếu QueueEnabled = true)
message := m.NewMessage().
    To("recipient@example.com").
    Subject("Queued Email").
    Text("This email is being sent through the queue")

if err := m.Send(message); err != nil {
    log.Fatal(err)
}

// Nếu muốn truy cập trực tiếp vào hàng đợi
if qm, ok := m.(*mailer.QueuedMailer); ok {
    queue := qm.GetQueue()
    
    // Liệt kê các email trong hàng đợi
    emails, _ := queue.List()
    fmt.Printf("Có %d email đang trong hàng đợi\n", len(emails))
    
    // Xử lý hàng đợi thủ công nếu cần
    queue.Process(qm.GetMailer())
}
```

### 7. Sử Dụng Template Functions

```go
// Tạo template functions
htmlFuncs := template.FuncMap{
    "formatDate": func(t time.Time) string {
        return t.Format("02/01/2006")
    },
}

textFuncs := texttemplate.FuncMap{
    "formatDate": func(t time.Time) string {
        return t.Format("02/01/2006")
    },
}

// Sử dụng trong template
data := map[string]interface{}{
    "Username": "John",
    "Date": time.Now(),
}

message := m.NewMessage().
    To("user@example.com").
    Subject("Welcome").
    HTMLTemplate(`<p>Welcome {{.Username}}! Today is {{formatDate .Date}}</p>`).
    WithHTMLTemplateFuncs(htmlFuncs).
    TextTemplate(`Welcome {{.Username}}! Today is {{formatDate .Date}}`).
    WithTextTemplateFuncs(textFuncs).
    WithData(data)

if err := m.Send(message); err != nil {
    log.Fatal(err)
}
```

## Testing

Package này cung cấp MockMailer để dễ dàng test ứng dụng:

```go
// Tạo mock mailer
mockMailer := mailer.NewMockMailer()

// Đưa vào provider
provider := mailer.NewProvider(config)
provider.SetMailer(mockMailer)

// Đăng ký provider
app.Register(provider)

// Sau khi gửi email, kiểm tra
if len(mockMailer.SentMessages) != 1 {
    t.Fatalf("Expected 1 message, got %d", len(mockMailer.SentMessages))
}
sentMessage := mockMailer.SentMessages[0]
```

## Cấu trúc

- **Config**: Đối tượng cấu hình cho SMTP server
- **Message**: Đối tượng đại diện cho một email
- **Mailer**: Interface với các phương thức để gửi email
- **SMTPMailer**: Triển khai Mailer interface với SMTP
- **MockMailer**: Triển khai Mailer interface cho testing
- **Queue**: Interface cho hàng đợi email
- **MemoryQueue**: Triển khai Queue interface với lưu trữ trong bộ nhớ
- **QueuedMailer**: Wrapper cho Mailer với khả năng queue
- **Provider**: Service provider để tích hợp với DI container

## Yêu cầu hệ thống

- Go 1.18 trở lên

## Giấy phép

Mã nguồn này được phân phối dưới giấy phép MIT.
