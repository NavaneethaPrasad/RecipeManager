package models

import "time"

type ShoppingList struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`
	Date   time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
