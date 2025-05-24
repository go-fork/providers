package mailer

import (
	"encoding/json"
	"testing"
)

func TestMessage_MarshalJSON(t *testing.T) {
	msg := NewMessage().
		From("sender@example.com", "Sender Name").
		To("user1@example.com", "user2@example.com").
		CC("cc@example.com").
		BCC("bcc@example.com").
		ReplyTo("reply@example.com").
		Subject("Test Subject").
		Text("Plain text content").
		HTML("<p>HTML content</p>").
		WithData(map[string]interface{}{
			"Name": "John",
			"Age":  30,
		}).
		AttachBytes([]byte("file content"), "test.txt").
		EmbedBytes([]byte("image content"), "image.jpg", "img1")

	jsonData, err := msg.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() failed: %v", err)
	}

	if len(jsonData) == 0 {
		t.Fatal("MarshalJSON() returned empty data")
	}

	// Verify that the JSON can be parsed
	var messageJSON MessageJSON
	err = json.Unmarshal(jsonData, &messageJSON)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify basic fields
	if messageJSON.From != "sender@example.com" {
		t.Errorf("Expected From to be 'sender@example.com', got %s", messageJSON.From)
	}

	if messageJSON.FromName != "Sender Name" {
		t.Errorf("Expected FromName to be 'Sender Name', got %s", messageJSON.FromName)
	}

	if len(messageJSON.To) != 2 {
		t.Errorf("Expected 2 'To' recipients, got %d", len(messageJSON.To))
	}

	if messageJSON.To[0] != "user1@example.com" || messageJSON.To[1] != "user2@example.com" {
		t.Errorf("Expected specific 'To' recipients, got %v", messageJSON.To)
	}

	if len(messageJSON.Cc) != 1 || messageJSON.Cc[0] != "cc@example.com" {
		t.Errorf("Expected CC to contain 'cc@example.com', got %v", messageJSON.Cc)
	}

	if len(messageJSON.Bcc) != 1 || messageJSON.Bcc[0] != "bcc@example.com" {
		t.Errorf("Expected BCC to contain 'bcc@example.com', got %v", messageJSON.Bcc)
	}

	if len(messageJSON.ReplyTo) != 1 || messageJSON.ReplyTo[0] != "reply@example.com" {
		t.Errorf("Expected ReplyTo to contain 'reply@example.com', got %v", messageJSON.ReplyTo)
	}

	if messageJSON.Subject != "Test Subject" {
		t.Errorf("Expected Subject to be 'Test Subject', got %s", messageJSON.Subject)
	}

	if messageJSON.Text != "Plain text content" {
		t.Errorf("Expected Text to be 'Plain text content', got %s", messageJSON.Text)
	}

	if messageJSON.Html != "<p>HTML content</p>" {
		t.Errorf("Expected Html to be '<p>HTML content</p>', got %s", messageJSON.Html)
	}

	// Verify template data
	if len(messageJSON.TemplateData) != 2 {
		t.Errorf("Expected 2 template data entries, got %d", len(messageJSON.TemplateData))
	}

	if messageJSON.TemplateData["Name"] != "John" {
		t.Errorf("Expected TemplateData['Name'] to be 'John', got %v", messageJSON.TemplateData["Name"])
	}

	if messageJSON.TemplateData["Age"] != float64(30) { // JSON numbers become float64
		t.Errorf("Expected TemplateData['Age'] to be 30, got %v", messageJSON.TemplateData["Age"])
	}

	// Verify attachments
	if len(messageJSON.Attachments) != 1 {
		t.Errorf("Expected 1 attachment, got %d", len(messageJSON.Attachments))
	}

	attachment := messageJSON.Attachments[0]
	if attachment.Name != "test.txt" {
		t.Errorf("Expected attachment name to be 'test.txt', got %s", attachment.Name)
	}

	if string(attachment.Content) != "file content" {
		t.Errorf("Expected attachment content to be 'file content', got %s", string(attachment.Content))
	}

	if attachment.Path != "" {
		t.Errorf("Expected attachment path to be empty, got %s", attachment.Path)
	}

	// Verify embedded files
	if len(messageJSON.Embedded) != 1 {
		t.Errorf("Expected 1 embedded file, got %d", len(messageJSON.Embedded))
	}

	embedded := messageJSON.Embedded[0]
	if embedded.Name != "image.jpg" {
		t.Errorf("Expected embedded name to be 'image.jpg', got %s", embedded.Name)
	}

	if embedded.ContentID != "img1" {
		t.Errorf("Expected embedded ContentID to be 'img1', got %s", embedded.ContentID)
	}

	if string(embedded.Content) != "image content" {
		t.Errorf("Expected embedded content to be 'image content', got %s", string(embedded.Content))
	}
}

func TestMessage_MarshalJSON_EmptyMessage(t *testing.T) {
	msg := NewMessage()

	jsonData, err := msg.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() failed for empty message: %v", err)
	}

	if len(jsonData) == 0 {
		t.Fatal("MarshalJSON() returned empty data for empty message")
	}

	// Verify that the JSON can be parsed
	var messageJSON MessageJSON
	err = json.Unmarshal(jsonData, &messageJSON)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty message JSON: %v", err)
	}

	// Verify empty/nil slices
	if messageJSON.To == nil {
		t.Error("Expected To to be empty slice, not nil")
	}

	if messageJSON.Cc == nil {
		t.Error("Expected Cc to be empty slice, not nil")
	}

	if messageJSON.Bcc == nil {
		t.Error("Expected Bcc to be empty slice, not nil")
	}

	if messageJSON.ReplyTo == nil {
		t.Error("Expected ReplyTo to be empty slice, not nil")
	}

	if messageJSON.TemplateData == nil {
		t.Error("Expected TemplateData to be empty map, not nil")
	}
}

func TestMessage_UnmarshalJSON(t *testing.T) {
	// Create JSON data manually
	jsonData := `{
		"from": "sender@example.com",
		"from_name": "Sender Name",
		"to": ["user1@example.com", "user2@example.com"],
		"cc": ["cc@example.com"],
		"bcc": ["bcc@example.com"],
		"reply_to": ["reply@example.com"],
		"subject": "Test Subject",
		"text": "Plain text content",
		"html": "<p>HTML content</p>",
		"template_data": {
			"Name": "John",
			"Age": 30
		},
		"attachments": [{
			"name": "test.txt",
			"content": "ZmlsZSBjb250ZW50",
			"path": "",
			"content_id": ""
		}],
		"embedded": [{
			"name": "image.jpg",
			"content": "aW1hZ2UgY29udGVudA==",
			"path": "",
			"content_id": "img1"
		}]
	}`

	msg := NewMessage()
	err := msg.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("UnmarshalJSON() failed: %v", err)
	}

	// Verify basic fields
	if msg.from != "sender@example.com" {
		t.Errorf("Expected from to be 'sender@example.com', got %s", msg.from)
	}

	if msg.fromName != "Sender Name" {
		t.Errorf("Expected fromName to be 'Sender Name', got %s", msg.fromName)
	}

	if len(msg.to) != 2 {
		t.Errorf("Expected 2 'to' recipients, got %d", len(msg.to))
	}

	if msg.to[0] != "user1@example.com" || msg.to[1] != "user2@example.com" {
		t.Errorf("Expected specific 'to' recipients, got %v", msg.to)
	}

	if len(msg.cc) != 1 || msg.cc[0] != "cc@example.com" {
		t.Errorf("Expected cc to contain 'cc@example.com', got %v", msg.cc)
	}

	if len(msg.bcc) != 1 || msg.bcc[0] != "bcc@example.com" {
		t.Errorf("Expected bcc to contain 'bcc@example.com', got %v", msg.bcc)
	}

	if len(msg.replyTo) != 1 || msg.replyTo[0] != "reply@example.com" {
		t.Errorf("Expected replyTo to contain 'reply@example.com', got %v", msg.replyTo)
	}

	if msg.subject != "Test Subject" {
		t.Errorf("Expected subject to be 'Test Subject', got %s", msg.subject)
	}

	if msg.text != "Plain text content" {
		t.Errorf("Expected text to be 'Plain text content', got %s", msg.text)
	}

	if msg.html != "<p>HTML content</p>" {
		t.Errorf("Expected html to be '<p>HTML content</p>', got %s", msg.html)
	}

	// Verify template data
	if len(msg.templateData) != 2 {
		t.Errorf("Expected 2 template data entries, got %d", len(msg.templateData))
	}

	if msg.templateData["Name"] != "John" {
		t.Errorf("Expected templateData['Name'] to be 'John', got %v", msg.templateData["Name"])
	}

	if msg.templateData["Age"] != float64(30) { // JSON numbers become float64
		t.Errorf("Expected templateData['Age'] to be 30, got %v", msg.templateData["Age"])
	}

	// Verify attachments
	if len(msg.attachments) != 1 {
		t.Errorf("Expected 1 attachment, got %d", len(msg.attachments))
	}

	attachment := msg.attachments[0]
	if attachment.Name != "test.txt" {
		t.Errorf("Expected attachment name to be 'test.txt', got %s", attachment.Name)
	}

	if string(attachment.Content) != "file content" {
		t.Errorf("Expected attachment content to be 'file content', got %s", string(attachment.Content))
	}

	if attachment.Path != "" {
		t.Errorf("Expected attachment path to be empty, got %s", attachment.Path)
	}

	// Verify embedded files
	if len(msg.embedded) != 1 {
		t.Errorf("Expected 1 embedded file, got %d", len(msg.embedded))
	}

	embedded := msg.embedded[0]
	if embedded.Name != "image.jpg" {
		t.Errorf("Expected embedded name to be 'image.jpg', got %s", embedded.Name)
	}

	if embedded.ContentID != "img1" {
		t.Errorf("Expected embedded ContentID to be 'img1', got %s", embedded.ContentID)
	}

	if string(embedded.Content) != "image content" {
		t.Errorf("Expected embedded content to be 'image content', got %s", string(embedded.Content))
	}
}

func TestMessage_UnmarshalJSON_InvalidJSON(t *testing.T) {
	msg := NewMessage()

	// Test with invalid JSON
	err := msg.UnmarshalJSON([]byte("invalid json"))
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}
}

func TestMessage_UnmarshalJSON_EmptyJSON(t *testing.T) {
	msg := NewMessage()

	// Test with empty JSON object
	err := msg.UnmarshalJSON([]byte("{}"))
	if err != nil {
		t.Fatalf("UnmarshalJSON() failed for empty JSON: %v", err)
	}

	// Verify that fields are properly initialized (should be empty but not nil)
	if msg.to == nil {
		t.Error("Expected 'to' to be empty slice, not nil")
	}

	if msg.cc == nil {
		t.Error("Expected 'cc' to be empty slice, not nil")
	}

	if msg.bcc == nil {
		t.Error("Expected 'bcc' to be empty slice, not nil")
	}

	if msg.replyTo == nil {
		t.Error("Expected 'replyTo' to be empty slice, not nil")
	}

	if msg.templateData == nil {
		t.Error("Expected 'templateData' to be empty map, not nil")
	}
}

func TestMessage_JSONRoundTrip(t *testing.T) {
	// Create a complex message
	originalMsg := NewMessage().
		From("sender@example.com", "Sender Name").
		To("user1@example.com", "user2@example.com").
		CC("cc@example.com").
		BCC("bcc@example.com").
		ReplyTo("reply@example.com").
		Subject("Round Trip Test").
		Text("Plain text content").
		HTML("<p>HTML content</p>").
		WithData(map[string]interface{}{
			"Name":    "John",
			"Age":     30,
			"IsAdmin": true,
		}).
		AttachBytes([]byte("file content"), "test.txt").
		EmbedBytes([]byte("image content"), "image.jpg", "img1")

	// Marshal to JSON
	jsonData, err := originalMsg.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() failed: %v", err)
	}

	// Unmarshal from JSON
	restoredMsg := NewMessage()
	err = restoredMsg.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON() failed: %v", err)
	}

	// Compare basic fields
	if restoredMsg.from != originalMsg.from {
		t.Errorf("From mismatch: expected %s, got %s", originalMsg.from, restoredMsg.from)
	}

	if restoredMsg.fromName != originalMsg.fromName {
		t.Errorf("FromName mismatch: expected %s, got %s", originalMsg.fromName, restoredMsg.fromName)
	}

	if len(restoredMsg.to) != len(originalMsg.to) {
		t.Errorf("To length mismatch: expected %d, got %d", len(originalMsg.to), len(restoredMsg.to))
	}

	for i, addr := range originalMsg.to {
		if i >= len(restoredMsg.to) || restoredMsg.to[i] != addr {
			t.Errorf("To[%d] mismatch: expected %s, got %s", i, addr, restoredMsg.to[i])
		}
	}

	if restoredMsg.subject != originalMsg.subject {
		t.Errorf("Subject mismatch: expected %s, got %s", originalMsg.subject, restoredMsg.subject)
	}

	if restoredMsg.text != originalMsg.text {
		t.Errorf("Text mismatch: expected %s, got %s", originalMsg.text, restoredMsg.text)
	}

	if restoredMsg.html != originalMsg.html {
		t.Errorf("HTML mismatch: expected %s, got %s", originalMsg.html, restoredMsg.html)
	}

	// Compare template data
	if len(restoredMsg.templateData) != len(originalMsg.templateData) {
		t.Errorf("TemplateData length mismatch: expected %d, got %d", len(originalMsg.templateData), len(restoredMsg.templateData))
	}

	// Compare attachments
	if len(restoredMsg.attachments) != len(originalMsg.attachments) {
		t.Errorf("Attachments length mismatch: expected %d, got %d", len(originalMsg.attachments), len(restoredMsg.attachments))
	}

	if len(restoredMsg.attachments) > 0 && len(originalMsg.attachments) > 0 {
		if restoredMsg.attachments[0].Name != originalMsg.attachments[0].Name {
			t.Errorf("Attachment name mismatch: expected %s, got %s", originalMsg.attachments[0].Name, restoredMsg.attachments[0].Name)
		}

		if string(restoredMsg.attachments[0].Content) != string(originalMsg.attachments[0].Content) {
			t.Errorf("Attachment content mismatch: expected %s, got %s", string(originalMsg.attachments[0].Content), string(restoredMsg.attachments[0].Content))
		}
	}

	// Compare embedded files
	if len(restoredMsg.embedded) != len(originalMsg.embedded) {
		t.Errorf("Embedded length mismatch: expected %d, got %d", len(originalMsg.embedded), len(restoredMsg.embedded))
	}

	if len(restoredMsg.embedded) > 0 && len(originalMsg.embedded) > 0 {
		if restoredMsg.embedded[0].Name != originalMsg.embedded[0].Name {
			t.Errorf("Embedded name mismatch: expected %s, got %s", originalMsg.embedded[0].Name, restoredMsg.embedded[0].Name)
		}

		if restoredMsg.embedded[0].ContentID != originalMsg.embedded[0].ContentID {
			t.Errorf("Embedded ContentID mismatch: expected %s, got %s", originalMsg.embedded[0].ContentID, restoredMsg.embedded[0].ContentID)
		}

		if string(restoredMsg.embedded[0].Content) != string(originalMsg.embedded[0].Content) {
			t.Errorf("Embedded content mismatch: expected %s, got %s", string(originalMsg.embedded[0].Content), string(restoredMsg.embedded[0].Content))
		}
	}
}

func TestAttachmentJSON_Conversion(t *testing.T) {
	// Test Attachment to AttachmentJSON conversion
	attachment := Attachment{
		Path:      "/path/to/file.txt",
		Name:      "file.txt",
		Content:   []byte("file content"),
		ContentID: "file1",
	}

	attachmentJSON := AttachmentJSON(attachment)

	if attachmentJSON.Path != attachment.Path {
		t.Errorf("Path mismatch: expected %s, got %s", attachment.Path, attachmentJSON.Path)
	}

	if attachmentJSON.Name != attachment.Name {
		t.Errorf("Name mismatch: expected %s, got %s", attachment.Name, attachmentJSON.Name)
	}

	if string(attachmentJSON.Content) != string(attachment.Content) {
		t.Errorf("Content mismatch: expected %s, got %s", string(attachment.Content), string(attachmentJSON.Content))
	}

	if attachmentJSON.ContentID != attachment.ContentID {
		t.Errorf("ContentID mismatch: expected %s, got %s", attachment.ContentID, attachmentJSON.ContentID)
	}

	// Test AttachmentJSON to Attachment conversion
	convertedBack := Attachment(attachmentJSON)

	if convertedBack.Path != attachment.Path {
		t.Errorf("Converted Path mismatch: expected %s, got %s", attachment.Path, convertedBack.Path)
	}

	if convertedBack.Name != attachment.Name {
		t.Errorf("Converted Name mismatch: expected %s, got %s", attachment.Name, convertedBack.Name)
	}

	if string(convertedBack.Content) != string(attachment.Content) {
		t.Errorf("Converted Content mismatch: expected %s, got %s", string(attachment.Content), string(convertedBack.Content))
	}

	if convertedBack.ContentID != attachment.ContentID {
		t.Errorf("Converted ContentID mismatch: expected %s, got %s", attachment.ContentID, convertedBack.ContentID)
	}
}
