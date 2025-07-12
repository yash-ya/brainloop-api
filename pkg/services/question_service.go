package services

import (
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/repositories"
	"brainloop-api/pkg/utils"
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func CreateQuestion(question *models.Question, userID uint) (*models.Question, *models.ErrorResponse) {
	question.UserID = userID
	processedTags, err := repositories.FindOrCreateTags(question.Tags)
	if err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "TAG_PROCESSING_ERROR", "Failed to process tags.")
	}
	question.Tags = processedTags
	err = repositories.CreateQuestion(question)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, utils.SendError(http.StatusConflict, "DUPLICATE_QUESTION", "A question with this title already exists for this user.")
		}
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to save question to the database.")
	}
	fullQuestion, err := repositories.GetQuestionByID(question.UserID, question.ID)
	if err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to retrieve newly created question.")
	}

	return fullQuestion, nil
}

func GetQuestions(userID uint, status, difficulty string) ([]models.Question, *models.ErrorResponse) {
	questions, err := repositories.GetQuestions(userID, status, difficulty)
	if err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to retrieve questions from the database.")
	}
	return questions, nil
}

func GetQuestionByID(userID, questionID uint) (*models.Question, *models.ErrorResponse) {
	question, err := repositories.GetQuestionByID(userID, questionID)
	if err != nil {
		return nil, handleRepositoryError(err, "question")
	}
	return question, nil
}

func UpdateQuestion(questionUpdates *models.Question) (*models.Question, *models.ErrorResponse) {
	processedTags, err := repositories.FindOrCreateTags(questionUpdates.Tags)
	if err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "TAG_PROCESSING_ERROR", "Failed to process tags.")
	}
	questionUpdates.Tags = processedTags
	updatedQuestion, err := repositories.UpdateQuestion(questionUpdates)
	if err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to update question.")
	}

	return updatedQuestion, nil
}

func DeleteQuestion(userID, questionID uint) *models.ErrorResponse {
	err := repositories.DeleteQuestion(userID, questionID)
	if err != nil {
		return handleRepositoryError(err, "question")
	}
	return nil
}

func handleRepositoryError(err error, resourceName string) *models.ErrorResponse {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.SendError(http.StatusNotFound, "NOT_FOUND", "No "+resourceName+" found with this ID for the current user.")
	}
	return utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "A database error occurred on resource: "+resourceName)
}
