package models

import (
	"time"

	"gorm.io/gorm"
)

type SubmissionWindow struct {
    ID        uint           `gorm:"primaryKey"`
    OpenTime  time.Time      `gorm:"not null;index"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
