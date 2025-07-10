package routes

import (
	"brainloop-api/pkg/handlers"
	"brainloop-api/pkg/middleware"
	"brainloop-api/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "BrainLoop API is running",
		})
	})

	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Healthy"})
		})
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", func(c *gin.Context) {
				userID, exists := c.Get("userID")
				if !exists {
					utils.SendContextError(c, http.StatusInternalServerError, "CONTEXT_ERROR", "User ID not found in context")
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"message": "This is a protected route",
					"userID":  userID,
				})
			})
		}
	}
}
