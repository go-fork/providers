package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	texttemplate "text/template"

	"gopkg.in/gomail.v2"
)

// Message đại diện cho một email với đầy đủ thông tin
type Message struct {
	// Các thuộc tính cơ bản của email
	from        string
	fromName    string
	to          []string
	cc          []string
	bcc         []string
	replyTo     []string
	subject     string
	text        string
	html        string
	attachments []Attachment
	embedded    []Attachment

	// Thuộc tính cho template rendering
	textTemplate      string
	htmlTemplate      string
	templateData      map[string]interface{}
	templateFuncs     template.FuncMap
	textTemplateFuncs texttemplate.FuncMap
}

// Attachment chứa thông tin về tệp đính kèm
type Attachment struct {
	Path      string
	Name      string
	Content   []byte
	ContentID string // Dùng cho embedded images
}

// NewMessage tạo một message mới
func NewMessage() *Message {
	return &Message{
		to:                make([]string, 0),
		cc:                make([]string, 0),
		bcc:               make([]string, 0),
		replyTo:           make([]string, 0),
		attachments:       make([]Attachment, 0),
		embedded:          make([]Attachment, 0),
		templateData:      make(map[string]interface{}),
		templateFuncs:     template.FuncMap{},
		textTemplateFuncs: texttemplate.FuncMap{},
	}
}

// From đặt địa chỉ người gửi
func (m *Message) From(address string, name ...string) *Message {
	m.from = address
	if len(name) > 0 {
		m.fromName = name[0]
	}
	return m
}

// To thêm một hoặc nhiều người nhận
func (m *Message) To(addresses ...string) *Message {
	m.to = append(m.to, addresses...)
	return m
}

// CC thêm một hoặc nhiều người nhận carbon copy
func (m *Message) CC(addresses ...string) *Message {
	m.cc = append(m.cc, addresses...)
	return m
}

// BCC thêm một hoặc nhiều người nhận blind carbon copy
func (m *Message) BCC(addresses ...string) *Message {
	m.bcc = append(m.bcc, addresses...)
	return m
}

// ReplyTo đặt địa chỉ reply-to
func (m *Message) ReplyTo(addresses ...string) *Message {
	m.replyTo = append(m.replyTo, addresses...)
	return m
}

// Subject đặt tiêu đề email
func (m *Message) Subject(subject string) *Message {
	m.subject = subject
	return m
}

// Text đặt nội dung văn bản thuần túy
func (m *Message) Text(content string) *Message {
	m.text = content
	return m
}

// HTML đặt nội dung HTML
func (m *Message) HTML(content string) *Message {
	m.html = content
	return m
}

// TextTemplate đặt template cho nội dung văn bản thuần túy
func (m *Message) TextTemplate(content string) *Message {
	m.textTemplate = content
	return m
}

// HTMLTemplate đặt template cho nội dung HTML
func (m *Message) HTMLTemplate(content string) *Message {
	m.htmlTemplate = content
	return m
}

// WithTextTemplateFuncs đặt các template functions cho text template
func (m *Message) WithTextTemplateFuncs(funcs texttemplate.FuncMap) *Message {
	for name, fn := range funcs {
		m.textTemplateFuncs[name] = fn
	}
	return m
}

// WithHTMLTemplateFuncs đặt các template functions cho HTML template
func (m *Message) WithHTMLTemplateFuncs(funcs template.FuncMap) *Message {
	for name, fn := range funcs {
		m.templateFuncs[name] = fn
	}
	return m
}

// WithData đặt dữ liệu cho template rendering
func (m *Message) WithData(data map[string]interface{}) *Message {
	for k, v := range data {
		m.templateData[k] = v
	}
	return m
}

// Attach đính kèm tệp từ đường dẫn
func (m *Message) Attach(path string, name ...string) *Message {
	filename := filepath.Base(path)
	if len(name) > 0 && name[0] != "" {
		filename = name[0]
	}

	m.attachments = append(m.attachments, Attachment{
		Path: path,
		Name: filename,
	})
	return m
}

// AttachBytes đính kèm dữ liệu binary
func (m *Message) AttachBytes(content []byte, name string) *Message {
	m.attachments = append(m.attachments, Attachment{
		Content: content,
		Name:    name,
	})
	return m
}

// Embed nhúng một hình ảnh vào nội dung HTML
func (m *Message) Embed(path string, cid string) *Message {
	m.embedded = append(m.embedded, Attachment{
		Path:      path,
		ContentID: cid,
	})
	return m
}

// EmbedBytes nhúng dữ liệu binary vào nội dung HTML
func (m *Message) EmbedBytes(content []byte, name string, cid string) *Message {
	m.embedded = append(m.embedded, Attachment{
		Content:   content,
		Name:      name,
		ContentID: cid,
	})
	return m
}

// renderTextTemplate renders the text template with provided data
func (m *Message) renderTextTemplate() (string, error) {
	if m.textTemplate == "" {
		return m.text, nil
	}

	// Set strict option to catch nested field errors
	tmpl, err := texttemplate.New("text").Funcs(m.textTemplateFuncs).Option("missingkey=error").Parse(m.textTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse text template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, m.templateData); err != nil {
		return "", fmt.Errorf("failed to execute text template: %w", err)
	}

	return buf.String(), nil
}

// renderHTMLTemplate renders the HTML template with provided data
func (m *Message) renderHTMLTemplate() (string, error) {
	if m.htmlTemplate == "" {
		return m.html, nil
	}

	// Use template.HTML to prevent auto-escaping of HTML content and set strict option
	tmpl, err := template.New("html").Funcs(m.templateFuncs).Option("missingkey=error").Parse(m.htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, m.templateData); err != nil {
		return "", fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return buf.String(), nil
}

// Validate kiểm tra message có hợp lệ không
func (m *Message) Validate() error {
	if len(m.to) == 0 && len(m.cc) == 0 && len(m.bcc) == 0 {
		return errors.New("email must have at least one recipient (to, cc, or bcc)")
	}

	if m.text == "" && m.html == "" && m.textTemplate == "" && m.htmlTemplate == "" {
		return errors.New("email must have either text or HTML content")
	}

	return nil
}

// BuildGoMailMessage chuyển đổi Message thành gomail.Message
func (m *Message) BuildGoMailMessage(defaultFrom, defaultFromName string) (*gomail.Message, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}

	// Render templates nếu có
	renderedText, err := m.renderTextTemplate()
	if err != nil {
		return nil, err
	}

	renderedHTML, err := m.renderHTMLTemplate()
	if err != nil {
		return nil, err
	}

	msg := gomail.NewMessage()

	// Đặt người gửi
	from := m.from
	if from == "" {
		from = defaultFrom
	}

	fromName := m.fromName
	if fromName == "" {
		fromName = defaultFromName
	}

	if fromName != "" {
		msg.SetAddressHeader("From", from, fromName)
	} else {
		msg.SetHeader("From", from)
	}

	// Đặt người nhận
	if len(m.to) > 0 {
		msg.SetHeader("To", m.to...)
	}

	if len(m.cc) > 0 {
		msg.SetHeader("Cc", m.cc...)
	}

	if len(m.bcc) > 0 {
		msg.SetHeader("Bcc", m.bcc...)
	}

	if len(m.replyTo) > 0 {
		msg.SetHeader("Reply-To", m.replyTo...)
	}

	// Đặt tiêu đề
	msg.SetHeader("Subject", m.subject)

	// Đặt nội dung
	if renderedHTML != "" && renderedText != "" {
		msg.SetBody("text/plain", renderedText)
		msg.AddAlternative("text/html", renderedHTML)
	} else if renderedHTML != "" {
		msg.SetBody("text/html", renderedHTML)
	} else if renderedText != "" {
		msg.SetBody("text/plain", renderedText)
	}

	// Thêm tệp đính kèm
	for _, attachment := range m.attachments {
		if attachment.Path != "" {
			msg.Attach(attachment.Path, gomail.Rename(attachment.Name))
		} else if len(attachment.Content) > 0 {
			// Sử dụng cách tiếp cận khác vì gomail không có phương thức AttachReader
			// Chúng ta cần viết dữ liệu tạm ra file và đính kèm
			tmpFile, err := os.CreateTemp("", "mail-attachment-*")
			if err != nil {
				return nil, fmt.Errorf("failed to create temp file for attachment: %w", err)
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			if _, err := tmpFile.Write(attachment.Content); err != nil {
				return nil, fmt.Errorf("failed to write attachment content: %w", err)
			}
			tmpFile.Close() // Đóng file trước khi đính kèm

			msg.Attach(tmpFile.Name(), gomail.Rename(attachment.Name))
		}
	}

	// Nhúng tệp (embedded files)
	for _, embedded := range m.embedded {
		headers := map[string][]string{"Content-ID": {fmt.Sprintf("<%s>", embedded.ContentID)}}
		if embedded.Path != "" {
			msg.Embed(embedded.Path, gomail.SetHeader(headers))
		} else if len(embedded.Content) > 0 {
			// Sử dụng cách tiếp cận tương tự như với attachments
			tmpFile, err := os.CreateTemp("", "mail-embedded-*")
			if err != nil {
				return nil, fmt.Errorf("failed to create temp file for embedded content: %w", err)
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			if _, err := tmpFile.Write(embedded.Content); err != nil {
				return nil, fmt.Errorf("failed to write embedded content: %w", err)
			}
			tmpFile.Close() // Đóng file trước khi nhúng

			msg.Embed(tmpFile.Name(), gomail.SetHeader(headers))
		}
	}

	return msg, nil
}
