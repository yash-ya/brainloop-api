package services

import (
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/repositories"
	"brainloop-api/pkg/utils"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(user *models.User) *models.ErrorResponse {
	user.Email = strings.ToLower(user.Email)
	_, err := repositories.FindUserByEmail(user.Email)
	if err == nil {
		return utils.SendError(http.StatusConflict, "USER_ALREADY_EXISTS", "User already registered with this email")
	} else if err != gorm.ErrRecordNotFound {
		return utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Error checking for existing user: "+err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.SendError(http.StatusInternalServerError, "PASSWORD_HASH_FAILED", "Error hashing password")
	}
	user.Password = string(hashedPassword)
	err = repositories.CreateUser(user)
	if err != nil {
		return utils.SendError(http.StatusInternalServerError, "REGISTRATION_FAILED", "Failed to register user")
	}
	return nil
}

func LoginUser(email, password string) (*models.Token, *models.ErrorResponse) {
	email = strings.ToLower(email)
	user, err := repositories.FindUserByEmail(email)
	if err != nil {
		return nil, utils.SendError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, utils.SendError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
	}

	token, errResp := utils.GenerateToken(user)
	if errResp != nil {
		return nil, errResp
	}

	return token, nil
}

func FindOrCreateUserByGoogle(userInfo *models.GoogleUserInfo) (*models.User, *models.ErrorResponse) {
	existingUser, err := repositories.FindUserByEmail(userInfo.Email)
	if err == nil {
		return existingUser, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Error checking for existing user.")
	}

	placeholderPassword, err := generateSecurePassword()
	if err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "SERVER_ERROR", "Could not generate secure password.")
	}

	newUser := &models.User{
		Username: userInfo.Name,
		Email:    userInfo.Email,
		Password: placeholderPassword,
	}

	if err := repositories.CreateUser(newUser); err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to create new user.")
	}

	return newUser, nil
}

func generateSecurePassword() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
