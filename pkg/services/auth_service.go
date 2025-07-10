package services

import (
	"brainloop-api/pkg/config"
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/repositories"
	"brainloop-api/pkg/utils"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	secretKey := []byte(config.AppConfig.JWTSecretKey)
	expirationInHours := config.AppConfig.JWTExpiration
	expirationTimeUTC := time.Now().UTC().Add(time.Duration(expirationInHours) * time.Hour)

	claims := &models.JWTClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTimeUTC),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Could not generate token")
	}

	tokenResponse := models.Token{
		Success:   true,
		Token:     tokenString,
		ExpiresIn: expirationInHours,
		ExpiresAt: expirationTimeUTC,
	}
	return &tokenResponse, nil
}
