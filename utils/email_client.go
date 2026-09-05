package utils

import (
	"strconv"

	"gopkg.in/gomail.v2"
)

const EmailServiceNotEnabledErrorMessage = "email service is not enabled"

const (
	SampleEmailFailingRecipient = "failing-maintainer@example.com"
	SampleEmailSendFailedError  = "test email client failed to send email"
)

type EmailMessage struct {
	Subject string
	Body    string
}

var (
	SampleTestEmail = EmailMessage{
		Subject: "TEST admin email",
		Body:    "see new policy at https://quollix.org/legal/new-policy.md",
	}
	SampleMaintainersEmail = EmailMessage{
		Subject: "TEST maintainers email",
		Body:    "maintainers can review https://quollix.org/legal/new-policy.md",
	}
	SampleMaintainersFailureEmail = EmailMessage{
		Subject: "TEST maintainers partial failure",
		Body:    "maintainers can review https://quollix.org/legal/failing-policy.md",
	}
)

type EmailConfig struct {
	SMTPHost             string `validate:"host" json:"smtp_host"`
	SMTPPort             string `validate:"number" json:"smtp_port"`
	FromEmailAddress     string `validate:"email" json:"from_email_address"`
	EmailAccountUsername string `validate:"loose" json:"email_account_username"`
	EmailAccountPassword string `validate:"loose" json:"email_account_password"`
	IsEnabled            bool   `json:"is_enabled"`
}

var SampleEmailConfig = EmailConfig{ // #nosec G101 (CWE-798): Potential hardcoded credentials
	SMTPHost:             "smtp.test.example.com",
	SMTPPort:             "2525",
	FromEmailAddress:     "sender@test.example.com",
	EmailAccountUsername: "test-smtp-user",
	EmailAccountPassword: "test-smtp-password",
}

type EmailClient interface {
	SendEmail(emailConfig *EmailConfig, to, subject, body string) error
	CheckEmailServerConnection(emailConfig *EmailConfig) error
}

type EmailClientImpl struct{}

func (e *EmailClientImpl) SendEmail(emailConfig *EmailConfig, to, subject, body string) error {
	if !emailConfig.IsEnabled {
		return Logger.NewError(EmailServiceNotEnabledErrorMessage)
	}
	message := gomail.NewMessage()
	message.SetHeader("From", emailConfig.FromEmailAddress)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)
	port, err := strconv.Atoi(emailConfig.SMTPPort)
	if err != nil {
		return Logger.NewError(err.Error())
	}
	dialer := gomail.NewDialer(emailConfig.SMTPHost, port, emailConfig.EmailAccountUsername, emailConfig.EmailAccountPassword)
	err = dialer.DialAndSend(message)
	if err != nil {
		return Logger.NewError(err.Error())
	}
	return nil
}

func (e *EmailClientImpl) CheckEmailServerConnection(emailConfig *EmailConfig) error {
	port, err := strconv.Atoi(emailConfig.SMTPPort)
	if err != nil {
		return Logger.NewError(err.Error())
	}
	dialer := gomail.NewDialer(emailConfig.SMTPHost, port, emailConfig.EmailAccountUsername, emailConfig.EmailAccountPassword)
	sendCloser, err := dialer.Dial()
	if err != nil {
		return Logger.NewError(err.Error())
	}
	err = sendCloser.Close()
	if err != nil {
		return Logger.NewError(err.Error())
	}
	return err
}

type EmailClientMock struct{}

func (e EmailClientMock) SendEmail(emailConfig *EmailConfig, to, subject, body string) error {
	return nil
}

func (e EmailClientMock) CheckEmailServerConnection(emailConfig *EmailConfig) error {
	return nil
}
