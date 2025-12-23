package models

import "time"

type Ingredient struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
