package utils

import (
	// "maos-cloud-project-api/models"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {

	r := gin.Default()
	// config := GetEnvVars()

	// models.InitDB(config)

	return r
}
