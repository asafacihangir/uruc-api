package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User modeli tanımlama
type Book struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:char(36);primary_key"`
	Title           string    `gorm:"size:255;not null"`
	Author          string    `gorm:"size:255;not null"`
	Publisher       string    `gorm:"size:255;not null"`
	PublicationDate string    `gorm:"size:255;not null"`
	Isbn            string    `gorm:"size:255;not null"`
	Genre           string    `gorm:"size:255;not null"`
}
