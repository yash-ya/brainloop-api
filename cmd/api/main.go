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
	err := router.SetTrustedProxies([]string{"127.0.0.1", "::1", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})
	if err != nil {
		log.Fatal("Failed to set trusted proxies: ", err)
	}
	routes.SetupRoutes(router)

	log.Println("Server running on port", config.AppConfig.Port)
	router.Run(":" + config.AppConfig.Port)
}
