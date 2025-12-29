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
	DB   *gorm.DB // Needed for transactions/FirstOrCreate logic
}

// UPDATE: Pass *gorm.DB here to handle Ingredient creation logic
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

	// 1. ADD INGREDIENTS (This was missing in your code)
	for _, ingDTO := range req.Ingredients {
		var ingredient models.Ingredient
		// Check if ingredient name exists (e.g. "Salt"), if not create it
		if err := s.DB.FirstOrCreate(&ingredient, models.Ingredient{Name: ingDTO.Name}).Error; err != nil {
			return 0, err
		}

		// Link it to the recipe
		recipe.Ingredients = append(recipe.Ingredients, models.RecipeIngredient{
			IngredientID: ingredient.ID,
			Quantity:     ingDTO.Amount, // Matches your frontend JSON
			Unit:         ingDTO.Unit,
		})
	}

	// 2. ADD INSTRUCTIONS (This was missing)
	for i, stepText := range req.Instructions {
		recipe.Instructions = append(recipe.Instructions, models.Instruction{
			StepNumber: i + 1,
			Text:       stepText, // Assuming req.Instructions is []string
		})
	}

	// 3. Save to DB
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
			Description: r.Description, // Added description so frontend card shows it
			Category:    r.Category,
		})
	}

	return response, nil
}

func (s *recipeService) UpdateRecipe(recipeID uint, userID uint, req dto.UpdateRecipeRequest) error {
	// 1. Find existing recipe
	recipe, err := s.Repo.FindByID(recipeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return err
		}
		return err
	}

	// 2. Check Authorization
	if recipe.UserID != userID {
		return ErrUnauthorized
	}

	// 3. Start a Transaction (Safety first!)
	tx := s.DB.Begin()

	// --- A. Update Basic Fields ---
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

	// --- B. Update Instructions ---
	// Strategy: Delete old instructions -> Add new ones

	// 1. Delete old instructions for this recipe
	if err := tx.Where("recipe_id = ?", recipeID).Delete(&models.Instruction{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Add new instructions
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

	// --- C. Update Ingredients ---
	// Strategy: Delete old recipe-ingredients links -> Add new ones

	// 1. Delete old links
	if err := tx.Where("recipe_id = ?", recipeID).Delete(&models.RecipeIngredient{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Add new ingredients
	for _, ingDTO := range req.Ingredients {
		var ingredient models.Ingredient

		// Find Ingredient ID by Name (or create it if it's new, e.g. "Saffron")
		if err := tx.FirstOrCreate(&ingredient, models.Ingredient{Name: ingDTO.Name}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Create the link
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

	// 4. Commit Transaction
	return tx.Commit().Error
}

func (s *recipeService) DeleteRecipe(recipeID uint, userID uint) error {
	// 1. Find the recipe to ensure it exists and belongs to the user
	recipe, err := s.Repo.FindByID(recipeID)
	if err != nil {
		return err
	}

	if recipe.UserID != userID {
		return ErrUnauthorized
	}

	// 2. Start a Transaction (Delete everything or nothing)
	tx := s.DB.Begin()

	// 3. Delete Instructions linked to this recipe
	if err := tx.Where("recipe_id = ?", recipeID).Delete(&models.Instruction{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 4. Delete Recipe-Ingredient links
	if err := tx.Where("recipe_id = ?", recipeID).Delete(&models.RecipeIngredient{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 5. Delete Meal Plans associated with this recipe (Critical for FK constraints)
	if err := tx.Where("recipe_id = ?", recipeID).Delete(&models.MealPlan{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 6. Finally, Delete the Recipe itself
	// We use the tx (transaction) to delete, ignoring the Repo for this step to keep it atomic
	if err := tx.Delete(&models.Recipe{}, recipeID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 7. Commit changes
	return tx.Commit().Error
}

func (s *recipeService) GetRecipeByID(recipeID uint, userID uint) (*dto.RecipeDetailResponse, error) {

	recipe, err := s.Repo.FindByIDWithDetails(recipeID)
	if err != nil {
		return nil, err
	}

	// Allow viewing if owner OR if you want to allow public viewing, remove this check
	if recipe.UserID != userID {
		// return nil, ErrUnauthorized // Uncomment to restrict viewing to owner only
	}

	var ingredients []dto.IngredientResponse
	for _, ri := range recipe.Ingredients {
		ingredients = append(ingredients, dto.IngredientResponse{
			Name:     ri.Ingredient.Name,
			Quantity: ri.Quantity, // Renamed from Amount to Quantity to match your previous code
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
