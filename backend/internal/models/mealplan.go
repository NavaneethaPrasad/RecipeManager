package models

import "time"

type MealPlan struct {
	ID             uint      `gorm:"primaryKey"`
	UserID         uint      `gorm:"not null"`
	User           User      `gorm:"foreignKey:UserID"`
	RecipeID       uint      `gorm:"not null"`
	Recipe         Recipe    `gorm:"foreignKey:RecipeID"`
	Date           time.Time `gorm:"not null"`
	MealType       string    `gorm:"not null"` // breakfast/lunch/dinner
	TargetServings int       `json:"target_servings"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
