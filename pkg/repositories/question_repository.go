package repositories

import (
	"brainloop-api/pkg/database"
	"brainloop-api/pkg/models"
)

func CreateQuestion(question *models.Question) error {
	result := database.DB.Create(question)
	return result.Error
}

func GetQuestions(userID uint, status, difficulty string) (*[]models.Question, error) {
	db := database.GetDB()
	result := db.Where("user_id = ?", userID)

	if status != "" {
		result = result.Where("status = ?", status)
	}

	if difficulty != "" {
		result = result.Where("difficulty = ?", difficulty)
	}

	var questions []models.Question
	if err := result.Find(&questions).Error; err != nil {
		return nil, err
	}

	return &questions, nil
}

func GetQuestionByID(userID uint, questionID uint) (*models.Question, error) {
	db := database.GetDB()
	var question models.Question
	result := db.Where("user_id = ? AND id = ?", userID, questionID).First(&question)

	if result.Error != nil {
		return nil, result.Error
	}

	return &question, nil
}
