package router

import (
	"maos-cloud-project-api/controllers"
	"maos-cloud-project-api/middlewares"
	"maos-cloud-project-api/mocks"
	"maos-cloud-project-api/utils"
	"maos-cloud-project-api/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AuthRoutes(r *gin.Engine, cfg *config.Config) {
    v1 := r.Group("/api/v1")

    logrus.Info("MOCK_TESTS: ", cfg.MOCK_TESTS)

    var emailSender utils.EmailSender
    if cfg.MOCK_TESTS == "true" {
        logrus.Info("MOOOOOOCK: Using Mock Email Sender")
        emailSender = &mocks.MockEmailSender{}
    } else {
        logrus.Info("NOOOOMMOOCK: Using SMTP Email Sender")
        emailSender = &utils.SMTPSender{}
    }

    v1.POST("/login", controllers.Login)
    v1.POST("/signup", controllers.Signup(emailSender))
    v1.GET("/dashboard",middlewares.IsAuthorized(), controllers.Dashboard)
    v1.GET("/logout", controllers.Logout)
    v1.POST("/resetpassword", controllers.ResetPassword(emailSender))
    v1.POST("/updatepassword", controllers.UpdatePassword)
    v1.POST("/send-verification-email",middlewares.IsAuthorized(), controllers.SendEmailVerificationLink(emailSender))
    v1.GET("/verify-email", controllers.VerifyEmail)
}
