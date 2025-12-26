package repository

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type RecipeIngredientRepository interface {
	Create(ri *models.RecipeIngredient) error
	FindByRecipeID(recipeID uint) ([]models.RecipeIngredient, error)
	FindByID(id uint) (*models.RecipeIngredient, error)
	Delete(id uint) error
}

type recipeIngredientRepository struct {
	db *gorm.DB
}

func NewRecipeIngredientRepository(db *gorm.DB) RecipeIngredientRepository {
	return &recipeIngredientRepository{db: db}
}

func (r *recipeIngredientRepository) Create(ri *models.RecipeIngredient) error {
	return r.db.Create(ri).Error
}

func (r *recipeIngredientRepository) FindByRecipeID(recipeID uint) ([]models.RecipeIngredient, error) {
	var items []models.RecipeIngredient
	err := r.db.
		Preload("Ingredient").
		Where("recipe_id = ?", recipeID).
		Find(&items).Error
	return items, err
}

func (r *recipeIngredientRepository) Delete(id uint) error {
	return r.db.Delete(&models.RecipeIngredient{}, id).Error
}

func (r *recipeIngredientRepository) FindByID(id uint) (*models.RecipeIngredient, error) {
	var ri models.RecipeIngredient
	err := r.db.First(&ri, id).Error
	return &ri, err
}
