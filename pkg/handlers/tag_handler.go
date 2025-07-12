package handlers

import (
	"brainloop-api/pkg/services"
	"brainloop-api/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTag(ctx *gin.Context) {
	tagName := ctx.Param("name")
	if tagName == "" {
		utils.SendContextError(ctx, http.StatusBadRequest, "INVALID_REQUEST", "Question ID parameter is missing")
		return
	}

	errResp := services.CreateTag(tagName)
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, errResp)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Tag created",
	})
}

func GetAllTags(ctx *gin.Context) {
	tags, errResp := services.GetAllTags()
	if errResp != nil {
		ctx.JSON(errResp.StatusCode, errResp)
		return
	}

	ctx.JSON(http.StatusOK, tags)
}
