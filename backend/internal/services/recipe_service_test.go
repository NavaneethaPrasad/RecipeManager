package services

import (
	"errors"
	"testing"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Ingredient{}, &models.Recipe{}, &models.Instruction{}, &models.RecipeIngredient{})
	return db
}

type MockRecipeRepository struct {
	CreateFn              func(recipe *models.Recipe) error
	FindByUserIDFn        func(userID uint) ([]models.Recipe, error)
	FindByIDFn            func(id uint) (*models.Recipe, error)
	FindByIDWithDetailsFn func(id uint) (*models.Recipe, error)
	UpdateFn              func(recipe *models.Recipe) error
	DeleteFn              func(recipe *models.Recipe) error
}

func (m *MockRecipeRepository) Create(r *models.Recipe) error { return m.CreateFn(r) }
func (m *MockRecipeRepository) FindByUserID(u uint) ([]models.Recipe, error) {
	return m.FindByUserIDFn(u)
}
func (m *MockRecipeRepository) FindByID(id uint) (*models.Recipe, error) { return m.FindByIDFn(id) }
func (m *MockRecipeRepository) FindByIDWithDetails(id uint) (*models.Recipe, error) {
	if m.FindByIDWithDetailsFn != nil {
		return m.FindByIDWithDetailsFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *MockRecipeRepository) Update(r *models.Recipe) error { return m.UpdateFn(r) }
func (m *MockRecipeRepository) Delete(r *models.Recipe) error { return m.DeleteFn(r) }

func TestCreateRecipe_Success(t *testing.T) {
	db := setupTestDB()
	repo := &MockRecipeRepository{
		CreateFn: func(r *models.Recipe) error { r.ID = 1; return nil },
	}
	service := NewRecipeService(repo, db)

	req := dto.CreateRecipeRequest{
		Name: "Pasta", Servings: 2,
		Ingredients:  []dto.RecipeIngredientRequest{{Name: "Flour", Amount: 200, Unit: "g"}},
		Instructions: []string{"Boil water", "Cook pasta"},
	}

	id, err := service.CreateRecipe(1, req)
	if err != nil || id != 1 {
		t.Fatalf("expected success, got err: %v", err)
	}
}

func TestCreateRecipe_NoIngredients(t *testing.T) {
	service := NewRecipeService(&MockRecipeRepository{}, nil)
	_, err := service.CreateRecipe(1, dto.CreateRecipeRequest{Ingredients: []dto.RecipeIngredientRequest{}})
	if err == nil || err.Error() != "at least one ingredient is required" {
		t.Fatal("expected error for no ingredients")
	}
}

func TestCreateRecipe_RepoError(t *testing.T) {
	db := setupTestDB()
	repo := &MockRecipeRepository{
		CreateFn: func(r *models.Recipe) error { return errors.New("db error") },
	}
	service := NewRecipeService(repo, db)
	req := dto.CreateRecipeRequest{Ingredients: []dto.RecipeIngredientRequest{{Name: "A", Amount: 1}}}

	_, err := service.CreateRecipe(1, req)
	if err == nil || err.Error() != "db error" {
		t.Fatal("expected repo error to propagate")
	}
}

func TestGetMyRecipes_Success(t *testing.T) {
	repo := &MockRecipeRepository{
		FindByUserIDFn: func(u uint) ([]models.Recipe, error) {
			return []models.Recipe{{ID: 1, Name: "A", PrepTime: 5, CookTime: 5}}, nil
		},
	}
	service := NewRecipeService(repo, nil)
	res, err := service.GetMyRecipes(1)
	if err != nil || len(res) != 1 || res[0].TotalTime != 10 {
		t.Fatal("failed to get recipes or calculate total time")
	}
}

func TestGetRecipeByID_Success(t *testing.T) {
	repo := &MockRecipeRepository{
		FindByIDWithDetailsFn: func(id uint) (*models.Recipe, error) {
			return &models.Recipe{
				ID: 1, UserID: 1, Name: "A",
				Ingredients:  []models.RecipeIngredient{{Quantity: 1, Ingredient: models.Ingredient{Name: "X"}}},
				Instructions: []models.Instruction{{Text: "Step 1"}},
			}, nil
		},
	}
	service := NewRecipeService(repo, nil)
	res, err := service.GetRecipeByID(1, 1)
	if err != nil || res.Name != "A" || len(res.Ingredients) != 1 {
		t.Fatal("failed to get recipe details")
	}
}

func TestUpdateRecipe_Success(t *testing.T) {
	db := setupTestDB()
	repo := &MockRecipeRepository{
		FindByIDFn: func(id uint) (*models.Recipe, error) {
			return &models.Recipe{ID: 1, UserID: 1}, nil
		},
	}
	service := NewRecipeService(repo, db)

	req := dto.UpdateRecipeRequest{
		Name:         "New Name",
		Instructions: []string{"New Step"},
		Ingredients:  []dto.RecipeIngredientRequest{{Name: "New Ing", Amount: 10}},
	}

	err := service.UpdateRecipe(1, 1, req)
	if err != nil {
		t.Fatalf("expected successful update, got %v", err)
	}
}

func TestUpdateRecipe_NotFound(t *testing.T) {
	repo := &MockRecipeRepository{
		FindByIDFn: func(id uint) (*models.Recipe, error) { return nil, gorm.ErrRecordNotFound },
	}
	service := NewRecipeService(repo, nil)
	err := service.UpdateRecipe(1, 1, dto.UpdateRecipeRequest{})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatal("expected record not found error")
	}
}

func TestUpdateRecipe_Unauthorized(t *testing.T) {
	repo := &MockRecipeRepository{
		FindByIDFn: func(id uint) (*models.Recipe, error) { return &models.Recipe{UserID: 2}, nil },
	}
	service := NewRecipeService(repo, nil)
	err := service.UpdateRecipe(1, 1, dto.UpdateRecipeRequest{})
	if err != ErrUnauthorized {
		t.Fatal("expected unauthorized error")
	}
}

func TestDeleteRecipe_Success(t *testing.T) {
	repo := &MockRecipeRepository{
		FindByIDFn: func(id uint) (*models.Recipe, error) { return &models.Recipe{ID: 1, UserID: 1}, nil },
		DeleteFn:   func(r *models.Recipe) error { return nil },
	}
	service := NewRecipeService(repo, nil)
	err := service.DeleteRecipe(1, 1)
	if err != nil {
		t.Fatal("expected successful delete")
	}
}

func TestDeleteRecipe_NotFound(t *testing.T) {
	repo := &MockRecipeRepository{
		FindByIDFn: func(id uint) (*models.Recipe, error) { return nil, errors.New("not found") },
	}
	service := NewRecipeService(repo, nil)
	err := service.DeleteRecipe(1, 1)
	if err == nil {
		t.Fatal("expected error for non-existent recipe")
	}
}
