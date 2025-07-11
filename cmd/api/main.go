package main

import (
	"brainloop-api/pkg/config"
	"brainloop-api/pkg/database"
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/routes"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadConfig()
	database.ConnectDB()
}

func main() {
	db := database.GetDB()
	migrationErr := db.AutoMigrate(&models.User{}, &models.Tag{}, &models.Question{}, &models.RevisionHistory{})
	if migrationErr != nil {
		log.Fatal("Failed to migrate database: ", migrationErr)
	}
	log.Println("Database migrated successfully.")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		// We allow our local frontend to make requests
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	err := router.SetTrustedProxies([]string{"127.0.0.1", "::1", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})
	if err != nil {
		log.Fatal("Failed to set trusted proxies: ", err)
	}
	routes.SetupRoutes(router)
	router.Run(":" + config.AppConfig.Port)
}
