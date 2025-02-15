package models

import (
	"time"

	"gorm.io/gorm"
)

type Phrase struct {
    ID                uint           `gorm:"primaryKey"`
    Content          string         `gorm:"not null"`
    SubmittedBy      uint           `gorm:"not null"`
    SubmissionWindow uint           `gorm:"not null;index"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
    DeletedAt        gorm.DeletedAt `gorm:"index"`
}
