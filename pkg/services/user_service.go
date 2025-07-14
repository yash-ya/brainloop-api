package services

import (
	"brainloop-api/pkg/email"
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/repositories"
	"brainloop-api/pkg/utils"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func FindOrCreateUserByGoogle(userInfo *models.GoogleUserInfo) (*models.User, *models.ErrorResponse) {
	existingUser, err := repositories.FindUserByEmail(userInfo.Email)
	if err == nil {
		return existingUser, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Error checking for existing user.")
	}

	placeholderPassword, err := utils.GenerateSecurePassword()
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

func VerifyUserEmail(token string) *models.ErrorResponse {
	user, err := repositories.FindUserByVerificationToken(token)
	if err != nil {
		return utils.SendError(http.StatusBadRequest, "INVALID_TOKEN", "This verification link is invalid.")
	}

	if time.Now().UTC().After(user.VerificationTokenExpiresAt) {
		return utils.SendError(http.StatusBadRequest, "TOKEN_EXPIRED", "This verification link has expired. Please register again.")
	}

	if err := repositories.ActivateUser(user); err != nil {
		return utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to activate user account.")
	}

	return nil
}

func ResendVerificationEmail(userEmail string) *models.ErrorResponse {
	existingUser, err := repositories.FindUserByEmail(userEmail)
	if err != nil {
		return utils.SendError(http.StatusNotFound, "USER_NOT_FOUND", "No account found with that email address.")
	}

	if existingUser.IsEmailVerified {
		return utils.SendError(http.StatusBadRequest, "ALREADY_VERIFIED", "This email address has already been verified.")
	}

	existingUser.IsEmailVerified = false
	existingUser.VerificationToken, _ = utils.GenerateSecurePassword()
	existingUser.VerificationTokenExpiresAt = time.Now().UTC().Add(30 * time.Minute)

	if err := repositories.UpdateUserVerification(existingUser); err != nil {
		return utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to update user verification details.")
	}

	if err := email.SendVerificationEmail(existingUser.Email, existingUser.VerificationToken); err != nil {
		return utils.SendError(http.StatusInternalServerError, "EMAIL_SEND_FAILED", "Failed to send new verification email.")
	}

	return nil
}
