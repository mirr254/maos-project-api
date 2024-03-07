package main

import (
	models "maos-cloud-project-api/models"
	"maos-cloud-project-api/router"
	utils "maos-cloud-project-api/utils"

	gin "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"os"

	"github.com/joho/godotenv"
)

func main() {
	utils.EnsurePlugins()
	// utils.CreateAwsSession()
	r := gin.Default()

	//load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	config := models.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSL_MODE"),
	}

	models.InitDB(config)

	//load routes
	router.AuthRoutes(r)
	r.Run(":8080")

}
