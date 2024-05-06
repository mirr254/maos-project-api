package router

import (
	"maos-cloud-project-api/controllers"
	"github.com/gin-gonic/gin"
)

func HealthCheck (r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.GET("/health", controllers.HealthCheck)
	
}