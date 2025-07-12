package repositories

import (
	"brainloop-api/pkg/database"
	"brainloop-api/pkg/models"
)

func FindOrCreateTags(tags []*models.Tag) ([]*models.Tag, error) {
	db := database.GetDB()
	var processedTags []*models.Tag
	if len(tags) == 0 {
		return processedTags, nil
	}

	for _, tag := range tags {
		result := db.Where("name = ?", tag.Name).FirstOrCreate(&tag)
		if result.Error != nil {
			return nil, result.Error
		}
		processedTags = append(processedTags, tag)
	}

	return processedTags, nil
}
