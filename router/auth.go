package router

import (
    "maos-cloud-project-api/controllers"
    "maos-cloud-project-api/middlewares"

    "github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
    r.POST("/login", controllers.Login)
    r.POST("/signup", controllers.Signup)
    r.GET("/dashboard",middlewares.IsAuthorized(), controllers.Dashboard)
    r.GET("/logout", controllers.Logout)
    r.POST("/resetpassword", controllers.ResetPassword)
}
