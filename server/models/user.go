package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	// Email        string         `gorm:"unique;not null"`
	Privilege    int  `gorm:"not null;default:0"` // 0: Normal, 1: Admin, 2: Super Admin
	IsEliminated bool `gorm:"not null;default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
