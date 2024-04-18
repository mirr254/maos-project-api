package router

import (
	"maos-cloud-project-api/controllers"
	"maos-cloud-project-api/middlewares"
	"maos-cloud-project-api/mocks"
    "maos-cloud-project-api/utils"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
    v1 := r.Group("/api/v1")
    env := os.Getenv("ENV")

    var emailSender utils.EmailSender
    if env == "test" {
        emailSender = &mocks.MockEmailSender{}
    } else {
        emailSender = &utils.SMTPSender{}
    }

    v1.POST("/login", controllers.Login)
    v1.POST("/signup", controllers.Signup)
    v1.GET("/dashboard",middlewares.IsAuthorized(), controllers.Dashboard)
    v1.GET("/logout", controllers.Logout)
    v1.POST("/resetpassword", controllers.ResetPassword)
    v1.GET("/verify-email", controllers.VerifyEmail)
    v1.POST("/send-verification-email",middlewares.IsAuthorized(), controllers.SendEmailVerification(emailSender))
}
