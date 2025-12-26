package repository

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type IngredientRepository interface {
	Create(ingredient *models.Ingredient) error
	FindAll() ([]models.Ingredient, error)
	FindByID(id uint) (*models.Ingredient, error)
}

type ingredientRepository struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) IngredientRepository {
	return &ingredientRepository{db: db}
}

func (r *ingredientRepository) Create(ingredient *models.Ingredient) error {
	return r.db.Create(ingredient).Error
}

func (r *ingredientRepository) FindAll() ([]models.Ingredient, error) {
	var ingredients []models.Ingredient
	err := r.db.Find(&ingredients).Error
	return ingredients, err
}

func (r *ingredientRepository) FindByID(id uint) (*models.Ingredient, error) {
	var ingredient models.Ingredient
	err := r.db.First(&ingredient, id).Error
	return &ingredient, err
}
