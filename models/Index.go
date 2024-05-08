package models

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	dsn := fmt.Sprintf("host=%s user=%s port=%s password=%s dbname=%s sslmode=%s", 
	         cfg.Host, cfg.User, cfg.Port ,cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		// panic(err)
		return nil, err
	}
	if err := db.AutoMigrate(&Users{}); err != nil {
		// panic(err)
		return nil, err
	}


	logrus.Info("Database Migration successful")

	// DB = db

	return db, nil
}
