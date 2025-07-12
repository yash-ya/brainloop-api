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

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "API is healthy"})
		})

		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/register", handlers.Register)
			authRoutes.POST("/login", handlers.Login)
		}

		protectedRoutes := apiV1.Group("/")
		protectedRoutes.Use(middleware.AuthMiddleware())
		{
			questionRoutes := protectedRoutes.Group("/questions")
			{
				questionRoutes.POST("", handlers.CreateQuestion)                     // POST /api/v1/questions
				questionRoutes.GET("", handlers.GetQuestions)                        // GET /api/v1/questions
				questionRoutes.GET("/:id", handlers.GetQuestionByID)                 // GET /api/v1/questions/:id
				questionRoutes.PUT("/:id", handlers.UpdateQuestion)                  // PUT /api/v1/questions/:id
				questionRoutes.DELETE("/:id", handlers.DeleteQuestion)               // DELETE /api/v1/questions/:id
				questionRoutes.GET("/:id/revisions", handlers.GetAllRevisionHistory) // GET /api/v1/questions/:id/revisions
			}

			revisionRoutes := protectedRoutes.Group("/revisions")
			{
				revisionRoutes.POST("", handlers.LogRevision) // POST /api/v1/revisions
			}

			tagRoutes := protectedRoutes.Group("/tags")
			{
				tagRoutes.GET("", handlers.GetAllTags)
				tagRoutes.POST("/:name", handlers.CreateTag)
			}
		}
	}
}
