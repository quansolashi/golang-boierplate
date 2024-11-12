package entity

import "gorm.io/gorm"

type Post struct {
	ID               uint64         `gorm:"primaryKey;autoIncrement;<-:create"`
	AuthorID         uint64         `gorm:"not null"`
	Title            string         `gorm:"not null;size:256"`
	ShortDescription string         `gomr:"size:256"`
	Content          string         `gorm:"not null;type:text"`
	DeletedAt        gorm.DeletedAt `gorm:"default:null"`
	User             User           `gorm:"foreignKey:AuthorID;constraint:OnUpdate:cascade,OnDelete:cascade"` // one-to-many relationship.
}
