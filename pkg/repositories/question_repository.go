package repositories

import (
	"brainloop-api/pkg/database"
	"brainloop-api/pkg/models"

	"gorm.io/gorm"
)

func CreateQuestion(question *models.Question) error {
	db := database.GetDB()
	result := db.Create(question)
	return result.Error
}

func GetQuestions(userID uint, status, difficulty string) ([]models.Question, error) {
	db := database.GetDB()
	result := db.Where("user_id = ?", userID)

	if status != "" {
		result = result.Where("status = ?", status)
	}

	if difficulty != "" {
		result = result.Where("difficulty = ?", difficulty)
	}

	var questions []models.Question
	if err := result.Preload("Tags").Preload("Revisions").Find(&questions).Error; err != nil {
		return nil, err
	}

	return questions, nil
}

func GetQuestionByID(userID, questionID uint) (*models.Question, error) {
	db := database.GetDB()
	var question models.Question

	result := db.Preload("Tags").Preload("Revisions").Where("user_id = ? AND id = ?", userID, questionID).First(&question)

	if result.Error != nil {
		return nil, result.Error
	}

	return &question, nil
}

func UpdateQuestion(userID, questionID uint, input *models.UpdateQuestion) error {
	db := database.GetDB()
	result := db.Model(&models.Question{}).
		Where("user_id = ? AND id = ?", userID, questionID).
		Updates(input)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func DeleteQuestion(userID, questionID uint) error {
	db := database.GetDB()
	result := db.Where("user_id = ? AND id = ?", userID, questionID).Delete(&models.Question{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
