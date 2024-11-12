package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserProfile struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement;<-:create"`
	UserID    uint64         `gorm:"not null;uniqueIndex"`
	Birthday  time.Time      `gorm:"not null;type:date"`
	Quote     string         `gom:"size:256"`
	DeletedAt gorm.DeletedAt `gorm:"default:null"`
}

type UserProfiles []*UserProfile
