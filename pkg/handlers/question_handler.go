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
	if err := ctx.ShouldBindJSON(&question); err != nil {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format: "+err.Error())
		return
	}

	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	errResp := services.CreateQuestion(&question, userID)
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, errResp)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Question created successfully!",
	})
}

func GetQuestions(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	var query struct {
		Status     string `form:"status"`
		Difficulty string `form:"difficulty"`
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Invalid query format: "+err.Error())
		return
	}

	questions, errResp := services.GetQuestions(userID, query.Status, query.Difficulty)
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, errResp)
		return
	}

	ctx.JSON(http.StatusOK, questions)
}

func GetQuestionByID(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	questionID, ok := getQuestionIDFromParam(ctx)
	if !ok {
		return
	}

	question, errResp := services.GetQuestionByID(userID, questionID)
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, errResp)
		return
	}

	ctx.JSON(http.StatusOK, question)
}

func UpdateQuestion(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	questionID, ok := getQuestionIDFromParam(ctx)
	if !ok {
		return
	}

	var input models.UpdateQuestion
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format: "+err.Error())
		return
	}

	errResp := services.UpdateQuestion(userID, questionID, &input)
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, errResp)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Question updated successfully!",
	})
}

func DeleteQuestion(ctx *gin.Context) {
	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		return
	}

	questionID, ok := getQuestionIDFromParam(ctx)
	if !ok {
		return
	}

	errResp := services.DeleteQuestion(userID, questionID)
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, errResp)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func getUserIDFromContext(ctx *gin.Context) (uint, bool) {
	userIDContext, exists := ctx.Get("userID")
	if !exists {
		utils.SendContextError(ctx, http.StatusInternalServerError, "CONTEXT_ERROR", "User ID not found in context")
		return 0, false
	}

	userID, ok := userIDContext.(uint)
	if !ok {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_USER_ID", "User ID in context is not of type uint")
		return 0, false
	}

	return userID, true
}

func getQuestionIDFromParam(ctx *gin.Context) (uint, bool) {
	questionIDStr := ctx.Param("id")
	if questionIDStr == "" {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Question ID parameter is missing")
		return 0, false
	}

	id, err := strconv.ParseUint(questionIDStr, 10, 64)
	if err != nil {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_QUESTION_ID", "Question ID must be a valid number")
		return 0, false
	}
	return uint(id), true
}
