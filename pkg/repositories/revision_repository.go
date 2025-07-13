package repositories

import (
	"brainloop-api/pkg/database"
	"brainloop-api/pkg/models"
)

func LogRevision(revision *models.RevisionHistory) error {
	db := database.GetDB()
	result := db.Create(revision)
	return result.Error
}

func GetAllRevisionHistory(questionID uint) ([]models.RevisionHistory, error) {
	db := database.GetDB()
	var history []models.RevisionHistory
	if err := db.Where("question_id=?", questionID).Order("created_at asc").Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func CountRevisionsForQuestion(questionID uint) (int64, error) {
	db := database.GetDB()
	var count int64
	result := db.Model(&models.RevisionHistory{}).Where("question_id = ?", questionID).Count(&count)
	return count, result.Error
}
