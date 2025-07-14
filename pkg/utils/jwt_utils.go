package utils

import (
	"brainloop-api/pkg/config"
	"brainloop-api/pkg/models"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *models.User) (*models.Token, *models.ErrorResponse) {
	secretKey := []byte(config.AppConfig.JWTSecretKey)
	expirationInHours := config.AppConfig.JWTExpiration
	expirationTimeUTC := time.Now().UTC().Add(time.Duration(expirationInHours) * time.Hour)

	claims := &models.JWTClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserEmail: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTimeUTC),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, SendError(http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Could not generate token")
	}

	tokenResponse := models.Token{
		Success:   true,
		Token:     tokenString,
		ExpiresIn: expirationInHours,
		ExpiresAt: expirationTimeUTC,
	}
	return &tokenResponse, nil
}

func GenerateSecurePassword() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
