// internals/services/EmailServiceInterface.go
package services

type EmailServiceInterface interface {
	SendTextEmail(to []string, subject, body string) error
	SendHTMLEmail(to []string, subject, body string) error
	SendEmailWithAttachment(to []string, subject, body, attachmentPath string) error
	SendEmailWithCCAndBCC(to []string, cc []string, bcc []string, subject, body string) error
}
