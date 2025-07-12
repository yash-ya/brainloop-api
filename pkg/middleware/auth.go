package middleware

import (
	"brainloop-api/pkg/config"
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendContextError(ctx, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.SendContextError(ctx, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Authorization header format must be Bearer {token}")
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		claims := &models.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.AppConfig.JWTSecretKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				utils.SendContextError(ctx, http.StatusUnauthorized, "INVALID_SIGNATURE", "Invalid token signature")
			} else if err == jwt.ErrTokenExpired {
				utils.SendContextError(ctx, http.StatusUnauthorized, "TOKEN_EXPIRED", "Token has expired")
			} else {
				utils.SendContextError(ctx, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid token: "+err.Error())
			}
			ctx.Abort()
			return
		}

		if !token.Valid {
			utils.SendContextError(ctx, http.StatusUnauthorized, "INVALID_TOKEN", "Token is not valid")
			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Next()
	}
}
