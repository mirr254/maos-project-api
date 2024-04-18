package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	models "maos-cloud-project-api/models"
	utils "maos-cloud-project-api/utils"

	gin "github.com/gin-gonic/gin"
)

var smtpHost = os.Getenv("SMTP_HOST")
var smtpPort = os.Getenv("SMTP_PORT")
var from = os.Getenv("FROM_EMAIL")
var pass = os.Getenv("EMAIL_PASSWORD")

type ErrorResponse struct {
    Error string `json:"error"`
}

func Signup(c *gin.Context) {

    var user models.User

	email_verification_token, err := utils.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "could not generate token"})
		return
	}
	
	db, err := models.InitDB()
	if err != nil {
		// Handle error
		logrus.Error(err)
		panic(err)
	}

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request"})
        return
    }

    // Check for empty email
    if user.Email == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "email must be provided"})
        return
    }

	if !govalidator.IsEmail(user.Email){
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid email address"})
		return
	}

    var existingUser models.User


    db.Where("email = ?", user.Email).First(&existingUser)

    if existingUser.ID != 0 {
        c.JSON(http.StatusConflict, ErrorResponse{Error: "user already exists"})
        return
    }

    hashedPassword, err := utils.GenerateHashPassword(user.Password)

    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "could not generate password hash"})
        return
    }

    user.Password = hashedPassword
	user.EmailVerificationToken = email_verification_token
	user.IsEmailVerified = false

    db.Create(&user)

	// Send email verification link
	subject := "Email Verification"
	body := "Click the link below to verify your email\n\n" + utils.CreateVerificationLink(email_verification_token)
	emailSendStatusChan := make(chan error)
	go func ()  {
		
		err := utils.SendEmail(smtpHost, smtpPort, from, pass,user.Email, subject, body)
		emailSendStatusChan <- err
	}()
	// TODO: Handle email sending error
	// err = <-emailSendStatusChan
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "could not send email"})
	// 	logrus.Error("Error sending email: ", err)
	// 	return
	// }
	// logrus.Info("Email sent")

    c.JSON(http.StatusCreated, gin.H{"success": "user created"})
	return 
}

func SendEmailVerification(emailSender utils.EmailSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		db, err := models.InitDB()

		if err != nil {
			// Handle error
			logrus.Error(err)
			panic(err)
		}

		//get current user
		currentUser := c.MustGet("user").(models.User)

		//get user email
		userEmail := currentUser.Email

		//check if user is already verified
		if currentUser.IsEmailVerified {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "email already verified"})
			return
		}

		// check if user is logged in
		cookie, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		_, err = utils.ParseToken(cookie)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		//generate token
		email_verification_token, err := utils.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "could not generate token"})
			return
		}

		//save token to user
		user.EmailVerificationToken = email_verification_token

		//update user
		db.Model(&currentUser).Update("email_verification_token", email_verification_token)
		//send email
		// Send email verification link
		subject := "Email Verification"
		body := "Click the link below to verify your email\n\n" + utils.CreateVerificationLink(email_verification_token)
		emailSendStatusChan := make(chan error)
		func () {
			err := emailSender.SendEmail(smtpHost, smtpPort, from, pass,userEmail, subject, body)
			emailSendStatusChan <- err
		}()
		
		c.JSON(http.StatusOK, gin.H{"success": "email verification link sent"})
	}
}

func VerifyEmail(c *gin.Context) {

	var user models.User
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid token"})
		return
	}

	db, err := models.InitDB()
	if err != nil {
		// Handle error
		logrus.Error(err)
		panic(err)
	}
	db.Where("email_verification_token = ?", token).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid token"})
		return
	}

	user.IsEmailVerified = true
	db.Save(&user)

	// Redirect to a success page
	c.JSON(http.StatusOK, gin.H{"success": "Email verified"})

	// http.Redirect(c.Writer, c.Request, "/api/v1/dashboard", http.StatusSeeOther)
}

func Login(c *gin.Context) {

	var user models.User
	db, err := models.InitDB()
	if err != nil {
		// Handle error
		panic(err)
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User

	db.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		Role: existingUser.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey, err := utils.GetSecretKey()
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}

	jwtKey := []byte(secretKey)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"success": "user logged in"})
}

func ResetPassword(c *gin.Context) {

	var user models.User

	db, err := models.InitDB()
	if err != nil {
		// Handle error
		panic(err)
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User

	db.Where("email = ?", user.Email).First(&existingUser)
	// models.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID == 0 {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)

	if errHash != nil {
		c.JSON(500, gin.H{"error": "could not generate password hash"})
		return
	}

	db.Model(&existingUser).Update("password", user.Password)

	c.JSON(http.StatusOK, gin.H{"success": "password updated"})
}

func Dashboard(c *gin.Context) {

	cookie, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, err := utils.ParseToken(cookie)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if claims.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "customer dashboard", "role": claims.Role})

}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"success": "user logged out"})
}
