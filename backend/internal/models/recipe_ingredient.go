package models

import "time"

type RecipeIngredient struct {
	ID       uint   `gorm:"primaryKey"`
	RecipeID uint   `gorm:"not null"`
	Recipe   Recipe `gorm:"constraint:OnDelete:CASCADE;"`

	IngredientID uint       `gorm:"not null"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID"`

	Quantity float64 `gorm:"not null"`
	Unit     string  `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
