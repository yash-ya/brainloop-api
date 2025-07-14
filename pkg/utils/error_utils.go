package utils

import (
	"brainloop-api/pkg/models"

	"github.com/gin-gonic/gin"
)

func SendContextError(ctx *gin.Context, statusCode int, errorCode, errorMessage string) {
	response := models.ErrorResponse{
		Success: false,
	}
	response.Error.Code = errorCode
	response.Error.Message = errorMessage

	ctx.JSON(statusCode, response)
}

func SendError(statusCode int, errorCode, errorMessage string) *models.ErrorResponse {
	response := models.ErrorResponse{
		Success:    false,
		StatusCode: statusCode,
	}
	response.Error.Code = errorCode
	response.Error.Message = errorMessage

	return &response
}
