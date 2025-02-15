package models

import (
	"time"

	"gorm.io/gorm"
)

type Verification struct {
    ID                uint           `gorm:"primaryKey"`
    VerifiedUserID    uint           `gorm:"not null"`
    VerifierID        uint           `gorm:"not null"`
    SubmissionWindow  uint           `gorm:"not null;index"`
    CreatedAt         time.Time
    UpdatedAt         time.Time
    DeletedAt         gorm.DeletedAt `gorm:"index"`
}
