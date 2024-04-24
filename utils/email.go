package utils

import (
	"net/smtp"
	"os"

	"fmt"
	"github.com/sirupsen/logrus"
)

type EmailSender interface {
	SendEmail( toEmail, subject, body string) error
}

type SMTPSender struct {}

func (s *SMTPSender) SendEmail( toEmail, subject, body string) error {

	return SendEmail( toEmail, subject, body)
}

/*
SendEmail Sends an email to the user(toEmail)
   args: toEmail, subject, body
   returns: error
*/
func SendEmail( toEmail, subject, body string) error {

	smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")
    from     := os.Getenv("FROM_EMAIL")
    pass     := os.Getenv("EMAIL_PASSWORD")

	to := []string{toEmail}
	message := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	var auth smtp.Auth
	if pass != "" {
		auth = smtp.PlainAuth("", from, pass, smtpHost)
	}

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		logrus.Error("Error sending email: ", err)
		return err
	}
	logrus.Info("EMAIL INFO: Email Sent")

	return nil

}

func getBaseUrl() string {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		//provide a default if not specified/localhost
		baseUrl = "http://localhost:8080"
	}
	return baseUrl
}
func CreateVerificationLink(route, email_token string) string {
	baseUrl := getBaseUrl()
	return fmt.Sprintf("%s/api/v1/%s?token=%s", baseUrl, route ,email_token)
}
