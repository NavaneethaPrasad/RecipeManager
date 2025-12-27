package models

import "time"

type ShoppingList struct {
	ID        uint               `gorm:"primaryKey"`
	UserID    uint               `gorm:"not null"`
	User      User               `gorm:"foreignKey:UserID"`
	StartDate time.Time          `gorm:"not null"`
	EndDate   time.Time          `gorm:"not null"`
	Items     []ShoppingListItem `gorm:"foreignKey:ShoppingListID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
