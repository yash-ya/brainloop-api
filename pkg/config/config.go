package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl         string
	Port          string
	JWTSecretKey  string
	JWTExpiration int
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

	AppConfig = Config{DBUrl: dbURL, Port: port, JWTSecretKey: jwtSecretKey, JWTExpiration: jwtExpiration}
	log.Println("Configuration loaded successfully")
}
