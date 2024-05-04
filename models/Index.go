package models

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"os"
)

type Config struct {
	Host     string
	User     string
	Port     string
	Password string
	DBName   string
	SSLMode  string
}

func InitDB(cfg Config) ( *gorm.DB, error){

	if err := godotenv.Load("../.env"); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s port=%s password=%s dbname=%s sslmode=%s", 
	         os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PORT") ,os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_SSL_MODE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		// panic(err)
		return nil, err
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		// panic(err)
		return nil, err
	}


	logrus.Info("Database Migration successful")

	// DB = db

	return db, nil
}
