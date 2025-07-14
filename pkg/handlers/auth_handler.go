package handlers

import (
	"brainloop-api/pkg/config"
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/services"
	"brainloop-api/pkg/utils"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const GoogleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

func Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBind(&user); err != nil {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format: "+err.Error())
		return
	}

	err := services.CreateUser(&user)
	if err != nil {
		ctx.JSON(err.StatusCode, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully!",
	})
}

func Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBind(&req); err != nil {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format: "+err.Error())
		return
	}

	token, err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		ctx.JSON(err.StatusCode, err)
		return
	}

	ctx.JSON(http.StatusOK, token)
}

func GoogleLogin(ctx *gin.Context) {
	state, err := generateRandomState()
	if err != nil {
		utils.SendContextError(ctx, http.StatusInternalServerError, "SERVER_ERROR", "Failed to generate state for authentication: "+err.Error())
		return
	}
	ctx.SetCookie("oauthstate", state, int(10*time.Minute.Seconds()), "/api/v1/auth/google", "", true, true)
	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(ctx *gin.Context) {
	stateFromURL := ctx.Query("state")
	stateFromCookie, err := ctx.Cookie("oauthstate")
	if err != nil {
		utils.SendContextError(ctx, http.StatusBadRequest, "STATE_COOKIE_MISSING", "Session state is missing or expired. Please try logging in again.")
		return
	}
	if stateFromURL != stateFromCookie {
		utils.SendContextError(ctx, http.StatusUnauthorized, "STATE_MISMATCH", "Invalid state parameter. CSRF attempt suspected.")
		return
	}

	code := ctx.Query("code")
	oauthToken, err := config.AppConfig.GoogleLoginConfig.Exchange(context.Background(), code)
	if err != nil {
		utils.SendContextError(ctx, http.StatusUnauthorized, "TOKEN_EXCHANGE_FAILED", "Failed to exchange authorization code for token: "+err.Error())
		return
	}

	client := config.AppConfig.GoogleLoginConfig.Client(ctx, oauthToken)
	response, err := client.Get(GoogleUserInfoURL)
	if err != nil {
		utils.SendContextError(ctx, http.StatusBadGateway, "GOOGLE_API_FAILED", "Failed to contact Google's services: "+err.Error())
		return
	}
	defer response.Body.Close()

	userInfoBytes, err := io.ReadAll(response.Body)
	if err != nil {
		utils.SendContextError(ctx, http.StatusInternalServerError, "RESPONSE_READ_FAILED", "Failed to read user info response from Google.")
		return
	}

	var userInfo models.GoogleUserInfo
	if err := json.Unmarshal(userInfoBytes, &userInfo); err != nil {
		utils.SendContextError(ctx, http.StatusInternalServerError, "JSON_UNMARSHAL_FAILED", "Failed to parse user info from Google.")
		return
	}

	user, errResp := services.FindOrCreateUserByGoogle(&userInfo)
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, errResp)
		return
	}

	token, errResp := utils.GenerateToken(user)
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, errResp)
		return
	}

	frontendCallbackURL := "http://localhost:3000/auth/callback"
	redirectURL := fmt.Sprintf("%s?token=%s", frontendCallbackURL, token.Token)

	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
