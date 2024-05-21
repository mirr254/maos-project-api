package router

import (
	"maos-cloud-project-api/controllers"
	"maos-cloud-project-api/middlewares"
	"maos-cloud-project-api/utils"
	"maos-cloud-project-api/config"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, cfg *config.Config) {
    v1 := r.Group("/api/v1")

    emailSender := &utils.SMTPSender{}


    v1.POST("/login", controllers.Login)
    v1.POST("/signup", controllers.Signup(emailSender))
    v1.GET("/dashboard",middlewares.IsAuthorized(), controllers.Dashboard)
    v1.GET("/logout", controllers.Logout)
    v1.POST("/resetpassword", controllers.ResetPassword(emailSender))
    v1.POST("/updatepassword", controllers.UpdatePassword)
    v1.POST("/send-verification-email",middlewares.IsAuthorized(), controllers.SendEmailVerificationLink(emailSender))
    v1.GET("/verify-email", controllers.VerifyEmail)
}
