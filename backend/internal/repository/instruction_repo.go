package repository

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type InstructionRepository interface {
	Create(instruction *models.Instruction) error
	FindByRecipeID(recipeID uint) ([]models.Instruction, error)
	FindByID(id uint) (*models.Instruction, error)
	Update(instruction *models.Instruction) error
	Delete(instruction *models.Instruction) error
}

type instructionRepository struct {
	db *gorm.DB
}

func NewInstructionRepository(db *gorm.DB) InstructionRepository {
	return &instructionRepository{db: db}
}

func (r *instructionRepository) Create(instruction *models.Instruction) error {
	return r.db.Create(instruction).Error
}

func (r *instructionRepository) FindByRecipeID(recipeID uint) ([]models.Instruction, error) {
	var instructions []models.Instruction
	err := r.db.
		Where("recipe_id = ?", recipeID).
		Order("step_number asc").
		Find(&instructions).Error
	return instructions, err
}

func (r *instructionRepository) FindByID(id uint) (*models.Instruction, error) {
	var instruction models.Instruction
	err := r.db.First(&instruction, id).Error
	return &instruction, err
}

func (r *instructionRepository) Update(instruction *models.Instruction) error {
	return r.db.Save(instruction).Error
}

func (r *instructionRepository) Delete(instruction *models.Instruction) error {
	return r.db.Delete(instruction).Error
}
