package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User modeli tanımlama
type User struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:char(36);primary_key"`
	FirstName string    `gorm:"size:100;not null"`
	LastName  string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;unique;not null"`
	Password  string    `gorm:"size:255;not null"`
}
