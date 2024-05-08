package main

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
	
	utils "maos-cloud-project-api/utils"
	"maos-cloud-project-api/router"
	

)

func main() {
	err := godotenv.Load()
	if err != nil {
		// Ignore error if it's a file not found error
		if !strings.Contains(err.Error(), "no such file or directory") {
			logrus.Errorf("Failed to load .env file: %v", err)
		} else {
			logrus.Info("No .env file found, using environment variables")
		}
	}

	utils.EnsurePlugins()
	// utils.CreateAwsSession()

	r := utils.SetUpRouter()
	router.AuthRoutes(r)
	router.HealthCheck(r)
	r.Run(":8080")

}
