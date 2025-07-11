package services

import (
	"brainloop-api/pkg/models"
	"brainloop-api/pkg/repositories"
	"brainloop-api/pkg/utils"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func CreateQuestion(question *models.Question, userID uint) *models.ErrorResponse {
	question.UserID = userID
	err := repositories.CreateQuestion(question)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return utils.SendError(http.StatusConflict, "DUPLICATE_QUESTION", "A question with this title already exists.")
		}
		return utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to save question to the database.")
	}
	return nil
}

func GetQuestions(userId uint, status, difficulty string) (*[]models.Question, *models.ErrorResponse) {
	questions, err := repositories.GetQuestions(userId, status, difficulty)
	if err != nil {
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to retrieve questions from the database.")
	}
	return questions, nil
}

func GetQuestionByID(userId uint, questionID uint) (*models.Question, *models.ErrorResponse) {
	questions, err := repositories.GetQuestionByID(userId, questionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.SendError(http.StatusNotFound, "QUESTION_NOT_FOUND", "No question found with this ID for the current user.")
		}
		return nil, utils.SendError(http.StatusInternalServerError, "DATABASE_ERROR", "Failed to retrieve question from the database.")
	}
	return questions, nil
}
