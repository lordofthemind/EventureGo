package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
	"time"
)

type EmailService struct {
	smtpHost string
	smtpPort string
	from     string
	password string
}

// NewEmailService creates a new instance of EmailService with the given SMTP configurations.
func NewEmailService(smtpHost, smtpPort, from, password string) EmailServiceInterface {
	return &EmailService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		from:     from,
		password: password,
	}
}

func (es *EmailService) SendTextEmail(to []string, subject, body string) error {
	return es.sendMail(to, subject, body, "text/plain")
}

func (es *EmailService) SendHTMLEmail(to []string, subject, body string) error {
	return es.sendMail(to, subject, body, "text/html")
}

func (es *EmailService) SendEmailWithAttachment(to []string, subject, body, attachmentPath string) error {
	return es.sendEmailWithAttachmentHelper(to, subject, body, []string{attachmentPath})
}

func (es *EmailService) SendEmailWithMultipleAttachments(to []string, subject, body string, attachmentPaths []string) error {
	return es.sendEmailWithAttachmentHelper(to, subject, body, attachmentPaths)
}

func (es *EmailService) SendEmailWithInlineImages(to []string, subject, body string, imagePaths []string) error {
	// Similar to attachments, but inline images need to be handled differently with CID headers in the HTML content.
	return es.sendEmailWithAttachmentHelper(to, subject, body, imagePaths)
}

func (es *EmailService) SendEmailWithCCAndBCC(to []string, cc []string, bcc []string, subject, body string) error {
	recipients := append(to, append(cc, bcc...)...)
	return es.sendMail(recipients, subject, body, "text/plain")
}

func (es *EmailService) SendEmailWithHeaders(to []string, subject, body string, headers map[string]string) error {
	// Build a custom email message with the headers.
	msg := []byte(fmt.Sprintf("Subject: %s\n", subject))
	for key, value := range headers {
		msg = append(msg, []byte(fmt.Sprintf("%s: %s\n", key, value))...)
	}
	msg = append(msg, []byte(fmt.Sprintf("\n%s", body))...)
	return es.send(msg, to)
}

func (es *EmailService) SendPriorityEmail(to []string, subject, body string, priority string) error {
	headers := map[string]string{
		"X-Priority": priority,
	}
	return es.SendEmailWithHeaders(to, subject, body, headers)
}

func (es *EmailService) ScheduleEmail(to []string, subject, body string, sendAt time.Time) error {
	// Simple scheduler based on time difference
	time.Sleep(time.Until(sendAt))
	return es.SendTextEmail(to, subject, body)
}

func (es *EmailService) SendEmailWithReplyTo(to []string, subject, body, replyTo string) error {
	headers := map[string]string{
		"Reply-To": replyTo,
	}
	return es.SendEmailWithHeaders(to, subject, body, headers)
}

func (es *EmailService) SendBatchEmail(to []string, subject, body string) error {
	// You can enhance this to send emails in batches
	for _, recipient := range to {
		err := es.SendTextEmail([]string{recipient}, subject, body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (es *EmailService) SendEmailWithTracking(to []string, subject, body string, trackingID string) error {
	// Add tracking information (e.g., pixel or URL tracking)
	trackedBody := fmt.Sprintf("%s\n\nTracking ID: %s", body, trackingID)
	return es.SendTextEmail(to, subject, trackedBody)
}

func (es *EmailService) SendEmailWithAttachmentsAndInlineImages(to []string, subject, body string, attachmentPaths, imagePaths []string) error {
	// Add both attachments and inline images to the email
	return es.sendEmailWithAttachmentHelper(to, subject, body, append(attachmentPaths, imagePaths...))
}

// Helper function to send emails with attachments.
func (es *EmailService) sendEmailWithAttachmentHelper(to []string, subject, body string, attachmentPaths []string) error {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Email headers
	boundary := writer.Boundary()
	buffer.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	buffer.WriteString(fmt.Sprintf("MIME-version: 1.0;\nContent-Type: multipart/mixed; boundary=%s\n\n", boundary))

	// Plain text or HTML body part
	bodyPart, _ := writer.CreatePart(map[string][]string{"Content-Type": {"text/plain"}})
	bodyPart.Write([]byte(body))

	// Attach each file
	for _, attachmentPath := range attachmentPaths {
		fileContent, err := es.readFileAndEncodeBase64(attachmentPath)
		if err != nil {
			return err
		}
		attachmentPart, _ := writer.CreatePart(map[string][]string{
			"Content-Type":              {fmt.Sprintf("application/octet-stream; name=\"%s\"", filepath.Base(attachmentPath))},
			"Content-Disposition":       {fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(attachmentPath))},
			"Content-Transfer-Encoding": {"base64"},
		})
		attachmentPart.Write([]byte(fileContent))
	}

	writer.Close()

	// Send the email
	return es.send(buffer.Bytes(), to)
}

func (es *EmailService) sendMail(to []string, subject, body, contentType string) error {
	msg := []byte(fmt.Sprintf("Subject: %s\nMIME-version: 1.0;\nContent-Type: %s;\n\n%s", subject, contentType, body))
	return es.send(msg, to)
}

// Helper function to send the final email message.
func (es *EmailService) send(msg []byte, to []string) error {
	auth := smtp.PlainAuth("", es.from, es.password, es.smtpHost)
	return smtp.SendMail(es.smtpHost+":"+es.smtpPort, auth, es.from, to, msg)
}

// Helper function to read file and encode it in base64 format.
func (es *EmailService) readFileAndEncodeBase64(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileData), nil
}
