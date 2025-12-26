package services

import (
	"errors"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"gorm.io/gorm"
)

var ErrUnauthorized = errors.New("not authorized")

type RecipeService interface {
	CreateRecipe(userID uint, req dto.CreateRecipeRequest) error
	GetMyRecipes(userID uint) ([]dto.RecipeResponse, error)

	UpdateRecipe(recipeID uint, userID uint, req dto.UpdateRecipeRequest) error
	DeleteRecipe(recipeID uint, userID uint) error

	GetRecipeByID(recipeID uint, userID uint) (*dto.RecipeDetailResponse, error)
}

type recipeService struct {
	Repo repository.RecipeRepository
}

func NewRecipeService(repo repository.RecipeRepository) RecipeService {
	return &recipeService{Repo: repo}
}

func (s *recipeService) CreateRecipe(userID uint, req dto.CreateRecipeRequest) error {
	recipe := &models.Recipe{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Servings:    req.Servings,
		PrepTime:    req.PrepTime,
		CookTime:    req.CookTime,
		Category:    req.Category,
	}

	return s.Repo.Create(recipe)
}

func (s *recipeService) GetMyRecipes(userID uint) ([]dto.RecipeResponse, error) {
	recipes, err := s.Repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var response []dto.RecipeResponse
	for _, r := range recipes {
		response = append(response, dto.RecipeResponse{
			ID:        r.ID,
			Name:      r.Name,
			Servings:  r.Servings,
			TotalTime: r.PrepTime + r.CookTime,
		})
	}

	return response, nil
}

func (s *recipeService) UpdateRecipe(recipeID uint, userID uint, req dto.UpdateRecipeRequest) error {
	recipe, err := s.Repo.FindByID(recipeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return err
	}

	if recipe.UserID != userID {
		return ErrUnauthorized
	}

	recipe.Name = req.Name
	recipe.Description = req.Description
	recipe.Servings = req.Servings
	recipe.PrepTime = req.PrepTime
	recipe.CookTime = req.CookTime
	recipe.Category = req.Category

	return s.Repo.Update(recipe)
}

func (s *recipeService) DeleteRecipe(recipeID uint, userID uint) error {
	recipe, err := s.Repo.FindByID(recipeID)
	if err != nil {
		return err
	}

	if recipe.UserID != userID {
		return ErrUnauthorized
	}

	return s.Repo.Delete(recipe)
}

func (s *recipeService) GetRecipeByID(recipeID uint, userID uint) (*dto.RecipeDetailResponse, error) {

	recipe, err := s.Repo.FindByIDWithDetails(recipeID)
	if err != nil {
		return nil, err
	}

	if recipe.UserID != userID {
		return nil, ErrUnauthorized
	}

	var ingredients []dto.IngredientResponse
	for _, ri := range recipe.Ingredients {
		ingredients = append(ingredients, dto.IngredientResponse{
			Name:     ri.Ingredient.Name,
			Quantity: ri.Quantity,
			Unit:     ri.Unit,
		})
	}

	var instructions []dto.InstructionResponse
	for _, ins := range recipe.Instructions {
		instructions = append(instructions, dto.InstructionResponse{
			StepNumber: ins.StepNumber,
			Text:       ins.Text,
		})
	}

	response := &dto.RecipeDetailResponse{
		ID:           recipe.ID,
		Name:         recipe.Name,
		Description:  recipe.Description,
		Servings:     recipe.Servings,
		Category:     recipe.Category,
		PrepTime:     recipe.PrepTime,
		CookTime:     recipe.CookTime,
		TotalTime:    recipe.PrepTime + recipe.CookTime,
		Ingredients:  ingredients,
		Instructions: instructions,
	}

	return response, nil
}
