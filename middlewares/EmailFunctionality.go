package middlewares

import (
	"net/http"
	"net/smtp"
)

/*
SendEmail Sends an email to the user(toEmail)
   args: toEmail, subject, body
   returns: error
*/
func SendEmail(toEmail, subject string, body string) error {
	from := "company-email@maos.com"
	pass := "password"
	to := []string{toEmail}
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	message := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, pass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	return err

}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {	
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return 
	}
	// Validate token and update email verification status in the database
    // This involves checking if the token exists, matches a user, and has not expired
    // For demonstration, let's assume a function `VerifyUserEmail(token string) error` does this
	err := VerifyUserEmail(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to a success page
	http.Redirect(w, r, "/api/v1/dashboard", http.StatusSeeOther)
}

func VerifyUserEmail(token string) error {
	// Check if token exists in the database
	// If it does, update the user's email verification status
	// If it doesn't, return an error

	

	return nil
}
