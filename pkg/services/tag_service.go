package services

import (
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/repositories"
	"brainloop-api/pkg/utils"
	"net/http"
)

func CreateTag(tagName string) *models.ErrorResponse {
	tag := &models.Tag{
		Name: tagName,
	}
	err := repositories.FindOrCreateTag(tag)
	if err != nil {
		return utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to create tag")
	}
	return nil
}

func GetAllTags() ([]*models.Tag, *models.ErrorResponse) {
	tags, err := repositories.GetAllTags()
	if err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to get tags")
	}
	return tags, nil
}
