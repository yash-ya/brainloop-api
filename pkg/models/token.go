package models

import "time"

type Token struct {
	Success   bool      `json:"success"`
	Token     string    `json:"token"`
	ExpiresIn int       `json:"expires_in"`
	ExpiresAt time.Time `json:"expires_at"`
}
