package handlers

import (
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/services"
	"brainloop-api/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogRevision(ctx *gin.Context) {
	var revision models.RevisionHistory
	if err := ctx.ShouldBindJSON(&revision); err != nil {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format: "+err.Error())
		return
	}

	createdRevision, errResp := services.LogRevision(&revision)
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, errResp)
		return
	}

	ctx.JSON(http.StatusCreated, createdRevision)
}

func GetAllRevisionHistory(ctx *gin.Context) {

	questionID, ok := getQuestionIDFromParam(ctx)
	if !ok {
		return
	}

	revisionHistory, errResp := services.GetAllRevisionHistory(questionID)
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, errResp)
		return
	}

	ctx.JSON(http.StatusOK, revisionHistory)
}
