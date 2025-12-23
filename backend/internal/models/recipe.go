package models

import "time"

type Recipe struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"not null"`
	User         User   `gorm:"foreignKey:UserID"`
	Name         string `gorm:"not null"`
	Description  string
	Servings     int
	PrepTime     int    // minutes
	CookTime     int    // minutes
	Instructions string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
