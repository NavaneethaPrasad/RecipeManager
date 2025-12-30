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
	GetByDateRange(userID uint, startDateStr, endDateStr string) ([]dto.MealPlanResponse, error)
}

type mealPlanService struct {
	Repo       repository.MealPlanRepository
	RecipeRepo repository.RecipeRepository
}

func NewMealPlanService(repo repository.MealPlanRepository, recipeRepo repository.RecipeRepository) MealPlanService {
	return &mealPlanService{Repo: repo, RecipeRepo: recipeRepo}
}

func (s *mealPlanService) Create(userID uint, req dto.CreateMealPlanRequest) error {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return err
	}

	recipe, err := s.RecipeRepo.FindByID(req.RecipeID)
	if err != nil {
		return err
	}
	if recipe.UserID != userID {
		return ErrMealUnauthorized
	}

	err = s.Repo.FindDuplicate(userID, date, req.MealType)
	if err == nil {
		return ErrMealExists
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	mp := &models.MealPlan{
		UserID:         userID,
		RecipeID:       req.RecipeID,
		Date:           date,
		MealType:       req.MealType,
		TargetServings: req.TargetServings,
	}

	return s.Repo.Create(mp)
}

func (s *mealPlanService) GetByDateRange(userID uint, startDateStr, endDateStr string) ([]dto.MealPlanResponse, error) {
	layout := "2006-01-02"
	start, err := time.Parse(layout, startDateStr)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(layout, endDateStr)
	if err != nil {
		return nil, err
	}

	plans, err := s.Repo.FindByUserAndDateRange(userID, start, end)
	if err != nil {
		return nil, err
	}

	var response []dto.MealPlanResponse
	for _, p := range plans {
		response = append(response, dto.MealPlanResponse{
			ID:             p.ID,
			Date:           p.Date.Format(layout),
			MealType:       p.MealType,
			TargetServings: p.TargetServings,

			Recipe: dto.RecipeResponse{
				ID:   p.Recipe.ID,
				Name: p.Recipe.Name,
			},
		})
	}

	return response, nil
}

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
			ID:             p.ID,
			Date:           dateStr,
			MealType:       p.MealType,
			TargetServings: p.TargetServings,
			Recipe: dto.RecipeResponse{
				ID:   p.Recipe.ID,
				Name: p.Recipe.Name,
			},
		})
	}

	return resp, nil
}

func (s *mealPlanService) Update(id uint, userID uint, req dto.UpdateMealPlanRequest) error {
	mp, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	if mp.UserID != userID {
		return ErrMealUnauthorized
	}

	if req.RecipeID != 0 {
		mp.RecipeID = req.RecipeID
	}
	if req.MealType != "" {
		mp.MealType = req.MealType
	}
	if req.TargetServings > 0 {
		mp.TargetServings = req.TargetServings
	}

	return s.Repo.Update(mp)
}

func (s *mealPlanService) Delete(id uint, userID uint) error {
	mp, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	if mp.UserID != userID {
		return ErrMealUnauthorized
	}

	return s.Repo.Delete(mp)
}
