package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	DBUrl               string
	Port                string
	JWTSecretKey        string
	JWTExpiration       int
	FrontendCallbackURL string
	FrontendURL         string
	SMTPHost            string
	SMTPPort            int
	SMTPUsername        string
	SMTPPassword        string
	GoogleLoginConfig   oauth2.Config
}

var AppConfig Config

func LoadConfig() {
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, reading from environment variables")
		}
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
		port = "8080"
	}

	jwtSecretKey := os.Getenv("JWT_SECRET")
	if jwtSecretKey == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	jwtExpirationStr := os.Getenv("JWT_EXPIRATION")
	if jwtExpirationStr == "" {
		log.Fatal("JWT_EXPIRATION environment variable is not set")
	}

	jwtExpiration, err := strconv.Atoi(jwtExpirationStr)
	if err != nil {
		log.Fatal("Invalid JWT_EXPIRATION_HOURS value. Must be an integer.")
	}

	googleOauthClientId := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	if googleOauthClientId == "" {
		log.Fatal("GOOGLE_OAUTH_CLIENT_ID environment variable is not set")
	}

	googleOauthClientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	if googleOauthClientSecret == "" {
		log.Fatal("GOOGLE_OAUTH_CLIENT_SECRET environment variable is not set")
	}

	googleOauthRedirectURL := os.Getenv("GOOGLE_OAUTH_REDIRECT_URL")
	if googleOauthRedirectURL == "" {
		log.Fatal("GOOGLE_OAUTH_REDIRECT_URL environment variable is not set")
	}

	frontendCallbackURL := os.Getenv("FRONTEND_CALLBACK_URL")
	if frontendCallbackURL == "" {
		log.Fatal("FRONTEND_CALLBACK_URL environment variable is not set")
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Fatal("FRONTEND_URL environment variable is not set")
	}

	smtpPortStr := os.Getenv("SMTP_PORT")
	if smtpPortStr == "" {
		log.Fatal("SMTP_PORT environment variable is not set")
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatal("Invalid SMTP_PORT value. Must be an integer.")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		log.Fatal("SMTP_HOST environment variable is not set")
	}

	smtpUsername := os.Getenv("SMTP_USERNAME")
	if smtpUsername == "" {
		log.Fatal("SMTP_USERNAME environment variable is not set")
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		log.Fatal("SMTP_PASSWORD environment variable is not set")
	}

	AppConfig = Config{
		DBUrl:               dbURL,
		Port:                port,
		JWTSecretKey:        jwtSecretKey,
		JWTExpiration:       jwtExpiration,
		FrontendCallbackURL: frontendCallbackURL,
		SMTPHost:            smtpHost,
		SMTPPort:            smtpPort,
		SMTPUsername:        smtpUsername,
		SMTPPassword:        smtpPassword,
		FrontendURL:         frontendURL,
	}
	AppConfig.GoogleLoginConfig = oauth2.Config{
		ClientID:     googleOauthClientId,
		ClientSecret: googleOauthClientSecret,
		RedirectURL:  googleOauthRedirectURL,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
	log.Println("Configuration loaded successfully")
}
