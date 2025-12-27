package services

import (
	"errors"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
)

var ErrInvalidServings = errors.New("invalid servings")

type RecipeScaleService interface {
	ScaleRecipe(recipeID uint, userID uint, newServings int) (*dto.ScaledRecipeResponse, error)
}

type recipeScaleService struct {
	RecipeRepo repository.RecipeRepository
}

func NewRecipeScaleService(repo repository.RecipeRepository) RecipeScaleService {
	return &recipeScaleService{RecipeRepo: repo}
}

func (s *recipeScaleService) ScaleRecipe(
	recipeID uint,
	userID uint,
	newServings int,
) (*dto.ScaledRecipeResponse, error) {

	if newServings <= 0 {
		return nil, ErrInvalidServings
	}

	recipe, err := s.RecipeRepo.FindByIDWithDetails(recipeID)
	if err != nil {
		return nil, err
	}

	// Ownership check
	if recipe.UserID != userID {
		return nil, ErrUnauthorized
	}

	scaleFactor := float64(newServings) / float64(recipe.Servings)

	var ingredients []dto.ScaledIngredientResponse

	for _, ri := range recipe.Ingredients {
		ingredients = append(ingredients, dto.ScaledIngredientResponse{
			ID:       ri.IngredientID,
			Name:     ri.Ingredient.Name,
			Quantity: ri.Quantity * scaleFactor,
			Unit:     ri.Unit,
		})
	}

	return &dto.ScaledRecipeResponse{
		RecipeID:         recipe.ID,
		Name:             recipe.Name,
		OriginalServings: recipe.Servings,
		ScaledServings:   newServings,
		Ingredients:      ingredients,
	}, nil
}
