package main

import (

	"github.com/joho/godotenv"
	utils "maos-cloud-project-api/utils"
	"maos-cloud-project-api/router"
	"github.com/sirupsen/logrus"

)

func main() {
	utils.EnsurePlugins()
	// utils.CreateAwsSession()

	logrus.Info("Env file not loaded. Exiting..." )

	//load .env file
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatal("Env file not loaded. Exiting...", err )
	}
	r := utils.SetUpRouter()
	router.AuthRoutes(r)
	r.Run(":8080")

}
