package models

import "time"

type ShoppingListItem struct {
	ID             uint         `gorm:"primaryKey"`
	ShoppingListID uint         `gorm:"not null"`
	ShoppingList   ShoppingList `gorm:"foreignKey:ShoppingListID;constraint:OnDelete:CASCADE"`
	IngredientID   uint         `gorm:"not null"`
	Ingredient     Ingredient   `gorm:"foreignKey:IngredientID;constraint:OnDelete:RESTRICT"`
	Quantity       float64      `gorm:"not null"`
	Unit           string       `gorm:"not null"`
	Checked        bool         `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
