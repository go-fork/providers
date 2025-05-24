package mailer

import (
	"encoding/json"
)

// MessageJSON là cấu trúc dùng để serialize và deserialize Message
type MessageJSON struct {
	From         string                 `json:"from"`
	FromName     string                 `json:"from_name"`
	To           []string               `json:"to"`
	Cc           []string               `json:"cc"`
	Bcc          []string               `json:"bcc"`
	ReplyTo      []string               `json:"reply_to"`
	Subject      string                 `json:"subject"`
	Text         string                 `json:"text"`
	Html         string                 `json:"html"`
	Attachments  []AttachmentJSON       `json:"attachments"`
	Embedded     []AttachmentJSON       `json:"embedded"`
	TemplateData map[string]interface{} `json:"template_data"`
}

// AttachmentJSON là cấu trúc dùng để serialize và deserialize Attachment
type AttachmentJSON struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Content   []byte `json:"content"`
	ContentID string `json:"content_id"`
}

// MarshalJSON chuyển đổi Message thành JSON
func (m *Message) MarshalJSON() ([]byte, error) {
	mj := &MessageJSON{
		From:         m.from,
		FromName:     m.fromName,
		To:           m.to,
		Cc:           m.cc,
		Bcc:          m.bcc,
		ReplyTo:      m.replyTo,
		Subject:      m.subject,
		Text:         m.text,
		Html:         m.html,
		TemplateData: m.templateData,
	}

	// Chuyển đổi các attachment
	if len(m.attachments) > 0 {
		mj.Attachments = make([]AttachmentJSON, len(m.attachments))
		for i, att := range m.attachments {
			mj.Attachments[i] = AttachmentJSON(att)
		}
	}

	// Chuyển đổi các embedded
	if len(m.embedded) > 0 {
		mj.Embedded = make([]AttachmentJSON, len(m.embedded))
		for i, emb := range m.embedded {
			mj.Embedded[i] = AttachmentJSON(emb)
		}
	}

	return json.Marshal(mj)
}

// UnmarshalJSON chuyển đổi JSON thành Message
func (m *Message) UnmarshalJSON(data []byte) error {
	var mj MessageJSON
	err := json.Unmarshal(data, &mj)
	if err != nil {
		return err
	}

	// Gán các giá trị từ JSON vào Message
	m.from = mj.From
	m.fromName = mj.FromName

	// Initialize slices if nil
	if mj.To != nil {
		m.to = mj.To
	} else {
		m.to = []string{} // Initialize as empty slice instead of nil
	}

	if mj.Cc != nil {
		m.cc = mj.Cc
	} else {
		m.cc = []string{}
	}

	if mj.Bcc != nil {
		m.bcc = mj.Bcc
	} else {
		m.bcc = []string{}
	}

	if mj.ReplyTo != nil {
		m.replyTo = mj.ReplyTo
	} else {
		m.replyTo = []string{}
	}

	m.subject = mj.Subject
	m.text = mj.Text
	m.html = mj.Html

	if mj.TemplateData != nil {
		m.templateData = mj.TemplateData
	} else {
		m.templateData = make(map[string]interface{})
	}

	// Chuyển đổi các attachment
	if len(mj.Attachments) > 0 {
		m.attachments = make([]Attachment, len(mj.Attachments))
		for i, att := range mj.Attachments {
			m.attachments[i] = Attachment(att)
		}
	}

	// Chuyển đổi các embedded
	if len(mj.Embedded) > 0 {
		m.embedded = make([]Attachment, len(mj.Embedded))
		for i, emb := range mj.Embedded {
			m.embedded[i] = Attachment(emb)
		}
	}

	return nil
}
