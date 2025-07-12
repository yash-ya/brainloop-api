package database

import (
	"brainloop-api/pkg/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := config.AppConfig.DBUrl
	if dsn == "" {
		log.Fatal("Environment variable DATABASE_URL is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	log.Println("Database connection established")
}

func GetDB() *gorm.DB {
	return DB
}
