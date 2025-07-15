package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username                    string `gorm:"uniqueIndex;not null"`
	Email                       string `gorm:"uniqueIndex;not null"`
	Password                    string `gorm:"not null"`
	IsEmailVerified             bool
	VerificationToken           string
	VerificationTokenExpiresAt  time.Time
	PasswordResetToken          string
	PasswordResetTokenExpiresAt time.Time
	Questions                   []Question `gorm:"foreignKey:UserID"`
}

type GoogleUserInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
