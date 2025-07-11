package routes

import (
	"brainloop-api/pkg/handlers"
	"brainloop-api/pkg/middleware"
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
			questions := protected.Group("/questions")
			{
				questions.POST("", handlers.CreateQuestion)
				questions.GET("", handlers.GetQuestions)
				questions.GET("/:id", handlers.GetQuestionByID)
			}
		}
	}
}
