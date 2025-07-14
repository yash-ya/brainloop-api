package services

import (
	"brainloop-api/pkg/email"
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/repositories"
	"brainloop-api/pkg/utils"
	"log"
	"net/http"
	"strings"
	"time"

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

	user.IsEmailVerified = false
	user.VerificationToken, _ = utils.GenerateSecurePassword()
	user.VerificationTokenExpiresAt = time.Now().UTC().Add(30 * time.Minute)
	user.Password = string(hashedPassword)

	err = repositories.CreateUser(user)
	if err != nil {
		return utils.SendError(http.StatusInternalServerError, "REGISTRATION_FAILED", "Failed to register user")
	}

	go func() {
		if err := email.SendVerificationEmail(user.Email, user.VerificationToken); err != nil {
			log.Printf("ERROR: Failed to send verification email to %s: %v\n", user.Email, err)
		}
	}()

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

	if !user.IsEmailVerified {
		return nil, utils.SendError(http.StatusForbidden, "EMAIL_VERIFICATION_PENDING", "Please verify your email address before logging in.")
	}

	token, errResp := utils.GenerateToken(user)
	if errResp != nil {
		return nil, errResp
	}

	return token, nil
}
