package utils

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func GetSecretKey() (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
	    logrus.Fatal("SECRET_KEY environment variable is not set")
	    return secretKey, fmt.Errorf("Secret key not set")
    }

	return secretKey, nil

}