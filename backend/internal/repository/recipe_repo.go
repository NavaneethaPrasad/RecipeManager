package repository

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type RecipeRepository interface {
	Create(recipe *models.Recipe) error
	FindByUserID(userID uint) ([]models.Recipe, error)
	FindByID(id uint) (*models.Recipe, error)
	FindByIDWithDetails(id uint) (*models.Recipe, error)
	Update(recipe *models.Recipe) error
	Delete(recipe *models.Recipe) error
}

type recipeRepository struct {
	DB *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) RecipeRepository {
	return &recipeRepository{DB: db}
}

func (r *recipeRepository) Create(recipe *models.Recipe) error {
	return r.DB.Create(recipe).Error
}

func (r *recipeRepository) FindByUserID(userID uint) ([]models.Recipe, error) {
	var recipes []models.Recipe
	// Preload minimal info if needed for list view, or just fetch basic
	err := r.DB.Where("user_id = ?", userID).Find(&recipes).Error
	return recipes, err
}

func (r *recipeRepository) FindByID(id uint) (*models.Recipe, error) {
	var recipe models.Recipe
	// It is safer to always preload basic ingredients to avoid nil pointer issues
	err := r.DB.Preload("Ingredients.Ingredient").First(&recipe, id).Error
	return &recipe, err
}

func (r *recipeRepository) FindByIDWithDetails(id uint) (*models.Recipe, error) {
	var recipe models.Recipe
	err := r.DB.
		Preload("Ingredients.Ingredient"). // Loads "Sugar", "Milk" names
		Preload("Instructions").           // Loads the instruction steps
		First(&recipe, id).Error
	return &recipe, err
}

func (r *recipeRepository) Update(recipe *models.Recipe) error {
	return r.DB.Save(recipe).Error
}

func (r *recipeRepository) Delete(recipe *models.Recipe) error {
	return r.DB.Delete(recipe).Error
}
