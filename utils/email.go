package utils

import (
	"maos-cloud-project-api/config"
	"net/smtp"
	"os"

	"fmt"
	"github.com/sirupsen/logrus"
)

type EmailSender interface {
	SendEmail(cfg *config.Config, toEmail, subject, body string) error
}

type SMTPSender struct {}

func (s *SMTPSender) SendEmail(cfg *config.Config, toEmail, subject, body string) error {

	return SendEmail(cfg, toEmail, subject, body)
}

/*
SendEmail Sends an email to the user(toEmail)
   args: toEmail, subject, body
   returns: error
*/
func SendEmail(cfg *config.Config, toEmail, subject, body string) error {

	logrus.Info("SMTP CFG: ", cfg.SMTP_HOST)

	smtpHost := cfg.SMTP_HOST
    smtpPort := cfg.SMTP_PORT
    from     := cfg.FROM_EMAIL
    pass     := cfg.EMAIL_PASSWORD

	to := []string{toEmail}
	message := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	var auth smtp.Auth
	if pass != "" {
		auth = smtp.PlainAuth("", from, pass, smtpHost)
	}

	logrus.Info("PAAAASSSSS:", pass)

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
func CreateVerificationLink(route, email_token string) string {
	baseUrl := getBaseUrl()
	return fmt.Sprintf("%s/api/v1/%s?token=%s", baseUrl, route ,email_token)
}
