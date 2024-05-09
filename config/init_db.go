package config

import (
	"fmt"
	"maos-cloud-project-api/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



func InitDB(cfg *Config) ( *gorm.DB, error){

	dsn := fmt.Sprintf("host=%s user=%s port=%s password=%s dbname=%s sslmode=%s", 
	         cfg.DB_HOST, cfg.DB_USER, cfg.DB_PORT ,cfg.DB_PASSWORD, cfg.DB_NAME, cfg.DB_SSL_MODE)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}
	if err := db.AutoMigrate(&models.Users{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}
	logrus.Info("Database Migration successful")

	return db, nil
}
