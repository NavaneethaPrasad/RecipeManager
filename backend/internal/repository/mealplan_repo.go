package repository

import (
	"time"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type MealPlanRepository interface {
	Create(mp *models.MealPlan) error
	FindByUserAndDate(userID uint, date time.Time) ([]models.MealPlan, error)
	FindByID(id uint) (*models.MealPlan, error)
	FindDuplicate(userID uint, date time.Time, mealType string) error
	Update(mp *models.MealPlan) error
	Delete(mp *models.MealPlan) error
}

type mealPlanRepository struct {
	DB *gorm.DB
}

func NewMealPlanRepository(db *gorm.DB) MealPlanRepository {
	return &mealPlanRepository{DB: db}
}

func (r *mealPlanRepository) Create(mp *models.MealPlan) error {
	return r.DB.Create(mp).Error
}

func (r *mealPlanRepository) FindByUserAndDate(userID uint, date time.Time) ([]models.MealPlan, error) {
	var plans []models.MealPlan
	err := r.DB.Preload("Recipe").
		Where("user_id = ? AND date = ?", userID, date).
		Find(&plans).Error
	return plans, err
}

func (r *mealPlanRepository) FindByID(id uint) (*models.MealPlan, error) {
	var mp models.MealPlan
	err := r.DB.First(&mp, id).Error
	return &mp, err
}

func (r *mealPlanRepository) FindDuplicate(userID uint, date time.Time, mealType string) error {
	return r.DB.
		Where("user_id = ? AND date = ? AND meal_type = ?", userID, date, mealType).
		First(&models.MealPlan{}).Error
}

func (r *mealPlanRepository) Update(mp *models.MealPlan) error {
	return r.DB.Save(mp).Error
}

func (r *mealPlanRepository) Delete(mp *models.MealPlan) error {
	return r.DB.Delete(mp).Error
}
