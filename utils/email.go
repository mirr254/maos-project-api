package utils

import (
	"net/smtp"
	"os"

	"fmt"
	"github.com/sirupsen/logrus"
)

type EmailSender interface {
	SendEmail(smtpHost, smtpPort, from, pass, toEmail, subject, body string) error
}

type SMTPSender struct {}

func (s *SMTPSender) SendEmail(smtpHost, smtpPort, from, pass, toEmail, subject, body string) error {

	return SendEmail(smtpHost, smtpPort, from, pass, toEmail, subject, body)
}

/*
SendEmail Sends an email to the user(toEmail)
   args: toEmail, subject, body
   returns: error
*/
func SendEmail(smtpHost, smtpPort, from, pass, toEmail, subject, body string) error {

	logrus.Info("DETAILS: ", smtpHost, smtpPort, from, pass, toEmail, subject, body)
	
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
func CreateVerificationLink(email_token string) string {
	baseUrl := getBaseUrl()
	return fmt.Sprintf("%s/api/v1/verify-email?token=%s", baseUrl, email_token)
}
