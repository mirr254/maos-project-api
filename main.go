package main

import (

	"github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
	
	utils "maos-cloud-project-api/utils"
	"maos-cloud-project-api/router"
	

)

func main() {
	//load .env file
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Env file not loaded. Exiting...", err )
	} 

	utils.EnsurePlugins()
	// utils.CreateAwsSession()

	r := utils.SetUpRouter()
	router.AuthRoutes(r)
	r.Run(":8080")

}
