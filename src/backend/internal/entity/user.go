package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        int
	Name      string
	Surname   string
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string
	Role      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
