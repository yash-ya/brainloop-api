package handlers

import (
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/services"
	"brainloop-api/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateQuestion(ctx *gin.Context) {
	var question models.Question
	if err := ctx.ShouldBind(&question); err != nil {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format: "+err.Error())
		return
	}
	userIDContext, exist := ctx.Get("userID")
	if !exist || userIDContext == "" {
		utils.SendContextError(ctx, http.StatusInternalServerError, "CONTEXT_ERROR", "User ID not found in context")
		return
	}

	userIDFloat, ok := userIDContext.(float64)
	if !ok {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID format")
		return
	}

	userID := uint(userIDFloat)
	err := services.CreateQuestion(&question, userID)
	if err != nil {
		ctx.JSON(err.StatusCode, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Question created successfully!",
	})
}

func GetQuestions(ctx *gin.Context) {
	userIDContext, exist := ctx.Get("userID")
	if !exist || userIDContext == "" {
		utils.SendContextError(ctx, http.StatusInternalServerError, "CONTEXT_ERROR", "User ID not found in context")
		return
	}

	userIDFloat, ok := userIDContext.(float64)
	if !ok {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID format")
		return
	}

	userID := uint(userIDFloat)
	var query struct {
		Status     string `form:"status"`
		Difficulty string `form:"difficulty"`
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format: "+err.Error())
		return
	}

	questions, err := services.GetQuestions(userID, query.Status, query.Difficulty)
	if err != nil {
		ctx.JSON(err.StatusCode, err)
		return
	}

	ctx.JSON(http.StatusOK, questions)
}

func GetQuestionByID(ctx *gin.Context) {
	userIDContext, exist := ctx.Get("userID")
	if !exist || userIDContext == "" {
		utils.SendContextError(ctx, http.StatusInternalServerError, "CONTEXT_ERROR", "User ID not found in context")
		return
	}

	userIDFloat, ok := userIDContext.(float64)
	if !ok {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID format")
		return
	}

	userID := uint(userIDFloat)
	questionIDStr := ctx.Param("id")
	if questionIDStr == "" {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Empty question id")
		return
	}

	questionID, err := strconv.ParseUint(questionIDStr, 10, 64)
	if err != nil {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_QUESTION_ID", "Question ID must be a valid number")
		return
	}

	question, errResp := services.GetQuestionByID(userID, uint(questionID))
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, err)
		return
	}

	ctx.JSON(http.StatusOK, question)
}
