package router

import (
	"maos-cloud-project-api/controllers/aws"

	"github.com/gin-gonic/gin"
)

func ProjectRoutes( r *gin.Engine) {
	v1 := r.Group("/api/v1")

	v1.POST("/project", controllers.CreateProject)
	// v1.POST("/:project_name/stack", controllers.CreateStack)
	v1.DELETE("/:project_name/:stack_name", controllers.DeleteStack)

}
