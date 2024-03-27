package utils

import (
	"maos-cloud-project-api/models"

	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {

	r := gin.Default()
	config := GetEnvVars()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	models.InitDB(config)

	return r
}
