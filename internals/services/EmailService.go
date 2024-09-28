// internals/services/EmailService.go
package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/smtp"
	"path/filepath"
)

type EmailService struct {
	smtpHost string
	smtpPort string
	from     string
	password string
}

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
	// Create MIME message with attachment
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Email headers
	boundary := writer.Boundary()
	buffer.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	buffer.WriteString(fmt.Sprintf("MIME-version: 1.0;\nContent-Type: multipart/mixed; boundary=%s\n\n", boundary))

	// Plain text or HTML body part
	bodyPart, _ := writer.CreatePart(map[string][]string{"Content-Type": {"text/plain"}})
	bodyPart.Write([]byte(body))

	// File attachment part
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

	writer.Close()

	// Send email
	return es.send(buffer.Bytes(), to)
}

func (es *EmailService) SendEmailWithCCAndBCC(to []string, cc []string, bcc []string, subject, body string) error {
	recipients := append(to, append(cc, bcc...)...)
	return es.sendMail(recipients, subject, body, "text/plain")
}

func (es *EmailService) sendMail(to []string, subject, body, contentType string) error {
	msg := []byte(fmt.Sprintf("Subject: %s\nMIME-version: 1.0;\nContent-Type: %s;\n\n%s", subject, contentType, body))
	return es.send(msg, to)
}

func (es *EmailService) send(msg []byte, to []string) error {
	auth := smtp.PlainAuth("", es.from, es.password, es.smtpHost)
	return smtp.SendMail(es.smtpHost+":"+es.smtpPort, auth, es.from, to, msg)
}

func (es *EmailService) readFileAndEncodeBase64(filePath string) (string, error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileData), nil
}
