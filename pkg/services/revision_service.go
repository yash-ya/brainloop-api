package services

import (
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/repositories"
	"brainloop-api/pkg/utils"
	"net/http"
)

func LogRevision(revision *models.RevisionHistory) (*models.RevisionHistory, *models.ErrorResponse) {
	err := repositories.LogRevision(revision)
	if err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to log revision to the database.")
	}
	return revision, nil
}

func GetAllRevisionHistory(questionID uint) ([]models.RevisionHistory, *models.ErrorResponse) {
	revisionHistory, err := repositories.GetAllRevisionHistory(questionID)
	if err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to retrieve questions from the database.")
	}
	return revisionHistory, nil
}
