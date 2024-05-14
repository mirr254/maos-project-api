package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	DB_HOST     string
	DB_USER     string
	DB_PORT     string
	DB_PASSWORD string
	DB_NAME     string
	DB_SSL_MODE string

	SECRET_KEY     string
	BASE_URL       string
	SMTP_HOST      string
	SMTP_PORT      string
	FROM_EMAIL     string
	EMAIL_PASSWORD string 
	ENV            string
}

func LoadConfig(configPaths ...string) *Config {
    viper.SetConfigName("config") // name of config file (without extension)
    viper.SetConfigType("yaml")   
    viper.AutomaticEnv()          

    viper.BindEnv("DB_HOST")
    viper.BindEnv("DB_USER")
    viper.BindEnv("DB_PORT")
    viper.BindEnv("DB_PASSWORD")
    viper.BindEnv("DB_NAME")
    viper.BindEnv("DB_SSL_MODE")
    viper.BindEnv("SECRET_KEY")
    viper.BindEnv("BASE_URL")
    viper.BindEnv("SMTP_HOST")
    viper.BindEnv("SMTP_PORT")
    viper.BindEnv("FROM_EMAIL")
    viper.BindEnv("EMAIL_PASSWORD")
    viper.BindEnv("ENV")

    for _, path := range configPaths {
        viper.AddConfigPath(path) // path to look for the config file in
    }

    if err := viper.ReadInConfig(); err != nil {
        logrus.Fatalf("Error reading config file, %s", err)
    }

    var config Config
    err := viper.Unmarshal(&config)
    if err != nil {
        logrus.Fatalf("Unable to decode into struct, %v", err)
    }

	// Add this check to make sure your configuration is correctly loaded
    if config.DB_HOST == "" || config.DB_USER == "" || config.DB_PORT == "" || config.DB_PASSWORD == "" || config.DB_NAME == "" || config.DB_SSL_MODE == "" {
        logrus.Fatalf("Configuration not correctly loaded: %+v", config)
    }

    return &config
}
