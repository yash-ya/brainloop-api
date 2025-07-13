package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;not null"`
}

type Question struct {
	gorm.Model
	Title            string `gorm:"not null"`
	Problem          string `gorm:"type:text"`
	Examples         string `gorm:"type:text"`
	TimeTaken        string
	Status           string            `gorm:"not null;default:'To Do'"`
	Difficulty       string            `gorm:"not null;default:'Medium'"`
	Notes            string            `gorm:"type:text"`
	UserID           uint              `gorm:"not null"`
	NextRevisionDate *time.Time        `gorm:"index"`
	SrsLevel         int               `gorm:"not null;default:0"`
	Tags             []*Tag            `gorm:"many2many:question_tags;"`
	Revisions        []RevisionHistory `gorm:"foreignKey:QuestionID"`
}

type RevisionHistory struct {
	gorm.Model
	QuestionID uint `gorm:"not null"`
	TimeTaken  string
	RevisedAt  time.Time `gorm:"not null"`
}
