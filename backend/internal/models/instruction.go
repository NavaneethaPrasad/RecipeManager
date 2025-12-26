package models

type Instruction struct {
	ID         uint   `gorm:"primaryKey"`
	RecipeID   uint   `gorm:"not null"`
	Recipe     Recipe `gorm:"constraint:OnDelete:CASCADE;"`
	StepNumber int    `gorm:"not null"`
	Text       string `gorm:"not null"`
}
