// Package mailer cung cấp một service provider để gửi email
// với hỗ trợ template, xử lý hàng đợi và các tính năng mạnh mẽ khác.
//
// Package này cung cấp APIs để gửi email với nội dung văn bản thuần túy
// hoặc HTML, hỗ trợ render template, đính kèm file, nhúng hình ảnh và
// xử lý email qua hàng đợi. Nó được xây dựng dựa trên thư viện gopkg.in/gomail.v2
// và cung cấp một interface đơn giản, mạnh mẽ và linh hoạt.
//
// Ví dụ sử dụng cơ bản:
//
//	// Đăng ký service provider
//	app.Register(mailer.NewProvider(&mailer.Config{
//	    Host:       "smtp.example.com",
//	    Port:       587,
//	    Username:   "user@example.com",
//	    Password:   "password",
//	    Encryption: "tls",
//	    FromAddress: "noreply@example.com",
//	    FromName:   "Example App",
//	}))
//
//	// Sử dụng mailer để gửi email
//	m := app.Make("mailer").(mailer.Mailer)
//	message := m.NewMessage().
//	    To("recipient@example.com").
//	    Subject("Hello!").
//	    Text("This is a text message").
//	    HTML("<h1>Hello</h1><p>This is an HTML message</p>")
//
//	if err := m.Send(message); err != nil {
//	    log.Fatal(err)
//	}
//
// Gửi email với template:
//
//	data := map[string]interface{}{
//	    "Username": "John Doe",
//	    "ResetLink": "https://example.com/reset?token=abc123",
//	}
//
//	message := m.NewMessage().
//	    To("user@example.com").
//	    Subject("Password Reset").
//	    HTMLTemplate(`
//	        <h1>Password Reset</h1>
//	        <p>Hello {{.Username}},</p>
//	        <p>Please reset your password by clicking the link below:</p>
//	        <p><a href="{{.ResetLink}}">Reset Password</a></p>
//	    `).
//	    WithData(data)
//
//	if err := m.Send(message); err != nil {
//	    log.Fatal(err)
//	}
package mailer
