package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement;<-:create"`
	Name      string         `gorm:"not null;size:256"`
	Email     string         `gorm:"not null;size:256;uique"`
	Password  string         `gorm:"not null;size:256"`
	CreatedAt time.Time      `gorm:"not null;<-:create"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"default:null"`
	Profile   UserProfile    `gorm:"constraint:OnUpdate:cascade,OnDelete:cascade"` // one-to-one relationship, this is included as default in User data.
	// Posts     []Post         `gorm:"foreignKey:AuthorID"` // FIXME: not required, enable if needed.
}

type Users []*User
