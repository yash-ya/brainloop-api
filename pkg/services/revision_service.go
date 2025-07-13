package services

import (
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/repositories"
	"brainloop-api/pkg/srs"
	"brainloop-api/pkg/utils"
	"net/http"
	"time"
)

func LogRevision(revision *models.RevisionHistory) (*models.RevisionHistory, *models.ErrorResponse) {
	question, err := repositories.GetQuestionForRevision(revision.QuestionID)
	if err != nil {
		return nil, utils.SendError(http.StatusNotFound, "NOT_FOUND", "Question not found.")
	}

	if err := repositories.LogRevision(revision); err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to log revision.")
	}

	isDue := question.NextRevisionDate == nil || !time.Now().Before(*question.NextRevisionDate)

	if isDue {
		newSRSLevel := question.SrsLevel + 1
		nextRevisionDate := srs.CalculateNextRevisionDate(newSRSLevel)
		err := repositories.UpdateQuestionSchedule(question.ID, newSRSLevel, &nextRevisionDate)
		if err != nil {
			return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to update question schedule.")
		}
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
