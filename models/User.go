package models

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Name                   string    `json:"name"`
	Email                  string    `json:"email" gorm:"unique"`
	Password               string    `json:"password"`
	Role                   string    `json:"role"`
	IsEmailVerified        bool      `json:"is_email_verified" gorm:"default:false"`
	EmailVerificationToken string    `json:"email_verification_token"`
	ResetPasswordToken     string    `json:"reset_password_token"`
}