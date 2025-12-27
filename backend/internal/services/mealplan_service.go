package services

import (
	"errors"
	"time"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"gorm.io/gorm"
)

var ErrMealExists = errors.New("meal already planned for this date and meal type")
var ErrMealUnauthorized = errors.New("not authorized")

type MealPlanService interface {
	Create(userID uint, req dto.CreateMealPlanRequest) error
	GetByDate(userID uint, date string) ([]dto.MealPlanResponse, error)
	Update(id uint, userID uint, req dto.UpdateMealPlanRequest) error
	Delete(id uint, userID uint) error
}

type mealPlanService struct {
	Repo       repository.MealPlanRepository
	RecipeRepo repository.RecipeRepository
}

func NewMealPlanService(
	repo repository.MealPlanRepository,
	recipeRepo repository.RecipeRepository,
) MealPlanService {
	return &mealPlanService{Repo: repo, RecipeRepo: recipeRepo}
}

/* ---------- Create ---------- */

func (s *mealPlanService) Create(userID uint, req dto.CreateMealPlanRequest) error {

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return err
	}

	// 🔐 recipe ownership
	recipe, err := s.RecipeRepo.FindByID(req.RecipeID)
	if err != nil {
		return err
	}
	if recipe.UserID != userID {
		return ErrUnauthorized
	}

	// 🚫 duplicate check
	err = s.Repo.FindDuplicate(userID, date, req.MealType)
	if err == nil {
		return errors.New("meal plan already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	mp := &models.MealPlan{
		UserID:   userID,
		RecipeID: req.RecipeID,
		Date:     date,
		MealType: req.MealType,
	}

	return s.Repo.Create(mp)
}

/* ---------- Read ---------- */

func (s *mealPlanService) GetByDate(userID uint, dateStr string) ([]dto.MealPlanResponse, error) {

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	plans, err := s.Repo.FindByUserAndDate(userID, date)
	if err != nil {
		return nil, err
	}

	var resp []dto.MealPlanResponse
	for _, p := range plans {
		resp = append(resp, dto.MealPlanResponse{
			ID:         p.ID,
			Date:       dateStr,
			MealType:   p.MealType,
			RecipeID:   p.RecipeID,
			RecipeName: p.Recipe.Name,
		})
	}

	return resp, nil
}

/* ---------- Update ---------- */

func (s *mealPlanService) Update(id uint, userID uint, req dto.UpdateMealPlanRequest) error {
	mp, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	if mp.UserID != userID {
		return ErrUnauthorized
	}

	mp.RecipeID = req.RecipeID
	mp.MealType = req.MealType

	return s.Repo.Update(mp)
}

/* ---------- Delete ---------- */
func (s *mealPlanService) Delete(id uint, userID uint) error {
	mp, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	if mp.UserID != userID {
		return ErrUnauthorized
	}

	return s.Repo.Delete(mp)
}
