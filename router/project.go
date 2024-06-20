package router

import (
	"maos-cloud-project-api/controllers"

	"github.com/gin-gonic/gin"
)

func ProjectRoutes( r *gin.Engine) {
	v1 := r.Group("/api/v1")

	v1.POST("/project", controllers.CreateProject)

}