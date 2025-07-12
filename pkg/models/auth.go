package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	UserEmail string `json:"user_email"`
	jwt.RegisteredClaims
}
