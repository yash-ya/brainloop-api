package main

import (
	"brainloop-api/pkg/config"
	"brainloop-api/pkg/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadConfig()
}

func main() {
	router := gin.Default()
	routes.SetupRoutes(router)

	log.Println("Server running on port", config.AppConfig.Port)
	router.Run(":" + config.AppConfig.Port)
}
