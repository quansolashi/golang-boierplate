package entity

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;<-:create"`
	Name      string    `gorm:"not null;size:256"`
	Email     string    `gorm:"not null;size:256;uique"`
	Password  string    `gorm:"not null;size:256"`
	CreatedAt time.Time `gorm:"not null;<-:create"`
	UpdatedAt time.Time `gorm:"not null"`
}

type Users []*User
