package services

import (
	"errors"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
)

var ErrIngredientUnauthorized = errors.New("not authorized to modify ingredient")

type IngredientService interface {
	CreateIngredient(req dto.CreateIngredientRequest) error
	GetIngredients() ([]dto.IngredientMasterResponse, error)

	AddIngredientToRecipe(recipeID uint, userID uint, req dto.AddRecipeIngredientRequest) error
	GetRecipeIngredients(recipeID uint, userID uint) ([]dto.IngredientResponse, error)
	RemoveRecipeIngredient(recipeIngredientID uint, userID uint) error
}

type ingredientService struct {
	IngredientRepo       repository.IngredientRepository
	RecipeIngredientRepo repository.RecipeIngredientRepository
	RecipeRepo           repository.RecipeRepository
}

func NewIngredientService(
	ingredientRepo repository.IngredientRepository,
	recipeIngredientRepo repository.RecipeIngredientRepository,
	recipeRepo repository.RecipeRepository,
) IngredientService {
	return &ingredientService{
		IngredientRepo:       ingredientRepo,
		RecipeIngredientRepo: recipeIngredientRepo,
		RecipeRepo:           recipeRepo,
	}
}

func (s *ingredientService) CreateIngredient(req dto.CreateIngredientRequest) error {
	ingredient := &models.Ingredient{
		Name: req.Name,
	}
	return s.IngredientRepo.Create(ingredient)
}

func (s *ingredientService) GetIngredients() ([]dto.IngredientMasterResponse, error) {
	ingredients, err := s.IngredientRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.IngredientMasterResponse
	for _, ing := range ingredients {
		response = append(response, dto.IngredientMasterResponse{
			ID:   ing.ID,
			Name: ing.Name,
		})
	}

	return response, nil
}

func (s *ingredientService) AddIngredientToRecipe(
	recipeID uint,
	userID uint,
	req dto.AddRecipeIngredientRequest,
) error {

	recipe, err := s.RecipeRepo.FindByID(recipeID)
	if err != nil {
		return err
	}
	if recipe.UserID != userID {
		return ErrIngredientUnauthorized
	}

	recipeIngredient := &models.RecipeIngredient{
		RecipeID:     recipeID,
		IngredientID: req.IngredientID,
		Quantity:     req.Quantity,
		Unit:         req.Unit,
	}

	return s.RecipeIngredientRepo.Create(recipeIngredient)
}

func (s *ingredientService) GetRecipeIngredients(
	recipeID uint,
	userID uint,
) ([]dto.IngredientResponse, error) {

	recipe, err := s.RecipeRepo.FindByID(recipeID)
	if err != nil {
		return nil, err
	}
	if recipe.UserID != userID {
		return nil, ErrIngredientUnauthorized
	}

	items, err := s.RecipeIngredientRepo.FindByRecipeID(recipeID)
	if err != nil {
		return nil, err
	}

	var response []dto.IngredientResponse
	for _, item := range items {
		response = append(response, dto.IngredientResponse{
			Name:     item.Ingredient.Name,
			Quantity: item.Quantity,
			Unit:     item.Unit,
		})
	}

	return response, nil
}

func (s *ingredientService) RemoveRecipeIngredient(
	recipeIngredientID uint,
	userID uint,
) error {

	ri, err := s.RecipeIngredientRepo.FindByID(recipeIngredientID)
	if err != nil {
		return err
	}

	recipe, err := s.RecipeRepo.FindByID(ri.RecipeID)
	if err != nil {
		return err
	}
	if recipe.UserID != userID {
		return ErrIngredientUnauthorized
	}

	return s.RecipeIngredientRepo.Delete(recipeIngredientID)
}
