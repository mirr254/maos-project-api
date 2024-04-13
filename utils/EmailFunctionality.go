package utils

import (
	"net/smtp"
	"os"

	"github.com/sirupsen/logrus"
)

/*
SendEmail Sends an email to the user(toEmail)
   args: toEmail, subject, body
   returns: error
*/
func SendEmail(toEmail, subject string, body string) error {

	
	from := os.Getenv("FROM_EMAIL")
	pass := os.Getenv("EMAIL_PASSWORD")
	to := []string{toEmail}
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	logrus.Info("DETAILS: ", smtpHost, smtpPort, from, pass, toEmail, subject, body)

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
