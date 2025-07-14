package repositories

import (
	"brainloop-api/pkg/database"
	"brainloop-api/pkg/models"
	"time"
)

func CreateUser(user *models.User) error {
	result := database.DB.Create(user)
	return result.Error
}

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func FindUserByVerificationToken(token string) (*models.User, error) {
	db := database.GetDB()
	var user models.User
	if err := db.Where("verification_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func ActivateUser(user *models.User) error {
	db := database.GetDB()

	updates := map[string]interface{}{
		"is_email_verified":             true,
		"verification_token":            "",
		"verification_token_expires_at": time.Time{},
	}

	result := db.Model(&user).Updates(updates)
	return result.Error
}
