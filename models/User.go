package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"unique" json:"email"`
	Password string `json:"password"`
	Roled    string  `json:"role"`
}

type Claims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}