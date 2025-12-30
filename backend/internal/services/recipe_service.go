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
	CreateRecipe(userID uint, req dto.CreateRecipeRequest) (uint, error)
	GetMyRecipes(userID uint) ([]dto.RecipeResponse, error)
	UpdateRecipe(recipeID uint, userID uint, req dto.UpdateRecipeRequest) error
	DeleteRecipe(recipeID uint, userID uint) error
	GetRecipeByID(recipeID uint, userID uint) (*dto.RecipeDetailResponse, error)
}

type recipeService struct {
	Repo repository.RecipeRepository
	DB   *gorm.DB
}

func NewRecipeService(repo repository.RecipeRepository, db *gorm.DB) RecipeService {
	return &recipeService{Repo: repo, DB: db}
}

func (s *recipeService) CreateRecipe(userID uint, req dto.CreateRecipeRequest) (uint, error) {

	recipe := models.Recipe{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		PrepTime:    req.PrepTime,
		CookTime:    req.CookTime,
		Servings:    req.Servings,
		Category:    req.Category,
	}

	for _, ingDTO := range req.Ingredients {
		var ingredient models.Ingredient
		if err := s.DB.FirstOrCreate(&ingredient, models.Ingredient{Name: ingDTO.Name}).Error; err != nil {
			return 0, err
		}

		recipe.Ingredients = append(recipe.Ingredients, models.RecipeIngredient{
			IngredientID: ingredient.ID,
			Quantity:     ingDTO.Amount,
			Unit:         ingDTO.Unit,
		})
	}

	for i, stepText := range req.Instructions {
		recipe.Instructions = append(recipe.Instructions, models.Instruction{
			StepNumber: i + 1,
			Text:       stepText,
		})
	}

	if err := s.Repo.Create(&recipe); err != nil {
		return 0, err
	}

	return recipe.ID, nil
}

func (s *recipeService) GetMyRecipes(userID uint) ([]dto.RecipeResponse, error) {
	recipes, err := s.Repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var response []dto.RecipeResponse
	for _, r := range recipes {
		response = append(response, dto.RecipeResponse{
			ID:          r.ID,
			Name:        r.Name,
			Servings:    r.Servings,
			TotalTime:   r.PrepTime + r.CookTime,
			Description: r.Description,
			Category:    r.Category,
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

	tx := s.DB.Begin()

	recipe.Name = req.Name
	recipe.Description = req.Description
	recipe.Servings = req.Servings
	recipe.PrepTime = req.PrepTime
	recipe.CookTime = req.CookTime
	recipe.Category = req.Category

	if err := tx.Save(recipe).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("recipe_id = ?", recipeID).Delete(&models.Instruction{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i, stepText := range req.Instructions {
		newInstruction := models.Instruction{
			RecipeID:   recipeID,
			StepNumber: i + 1,
			Text:       stepText,
		}
		if err := tx.Create(&newInstruction).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Where("recipe_id = ?", recipeID).Delete(&models.RecipeIngredient{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, ingDTO := range req.Ingredients {
		var ingredient models.Ingredient

		if err := tx.FirstOrCreate(&ingredient, models.Ingredient{Name: ingDTO.Name}).Error; err != nil {
			tx.Rollback()
			return err
		}

		ri := models.RecipeIngredient{
			RecipeID:     recipeID,
			IngredientID: ingredient.ID,
			Quantity:     ingDTO.Amount,
			Unit:         ingDTO.Unit,
		}
		if err := tx.Create(&ri).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
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
