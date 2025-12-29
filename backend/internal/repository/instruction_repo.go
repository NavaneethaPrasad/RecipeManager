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
	Delete(id uint) error
}

type instructionRepository struct {
	DB *gorm.DB
}

func NewInstructionRepository(db *gorm.DB) InstructionRepository {
	return &instructionRepository{DB: db}
}

func (r *instructionRepository) Create(instruction *models.Instruction) error {
	return r.DB.Create(instruction).Error
}

func (r *instructionRepository) FindByRecipeID(recipeID uint) ([]models.Instruction, error) {
	var instructions []models.Instruction
	err := r.DB.Where("recipe_id = ?", recipeID).Order("step_number asc").Find(&instructions).Error
	return instructions, err
}

func (r *instructionRepository) FindByID(id uint) (*models.Instruction, error) {
	var instruction models.Instruction
	err := r.DB.First(&instruction, id).Error
	return &instruction, err
}

func (r *instructionRepository) Update(instruction *models.Instruction) error {
	return r.DB.Save(instruction).Error
}

func (r *instructionRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Instruction{}, id).Error
}
