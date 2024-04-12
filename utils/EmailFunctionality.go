package utils

import (
	"net/smtp"
	"os"
)

/*
SendEmail Sends an email to the user(toEmail)
   args: toEmail, subject, body
   returns: error
*/
func SendEmail(toEmail, subject string, body string) error {

	from := os.Getenv("EMAIL")
	pass := os.Getenv("EMAIL_PASSWORD")
	to := []string{toEmail}
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	message := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, pass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	return err

}
