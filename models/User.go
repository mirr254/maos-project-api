package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name                 string    `json:"name"`
	Email                string    `json:"email" gorm:"unique"`
	Password             string    `json:"password"`
	Role                 string    `json:"role"`
	IsEmailVerified      bool      `json:"is_email_verified"`
	EmailVerificationToken string    `json:"email_verification_token"`
	EmailTokenGeneratedAt time.Time `json:"email_token_generated_at"`
}