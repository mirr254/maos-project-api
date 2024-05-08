package utils

import (
	"maos-cloud-project-api/models"
	
	"os"

)

func GetEnvVars() models.Config {

	return models.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSL_MODE"),
	}
}

func LoadTestConfig() models.Config {
    // Load configuration from a file
    return models.Config{
        Host:     "database",
        User:     "postgres",
        Port:     "5432",
        Password: "postgres",
        DBName:   "maosproject",
        SSLMode:  "disable",
    }
}