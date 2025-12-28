package services

import (
	"testing"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type MockRecipeRepository struct {
	CreateFn              func(recipe *models.Recipe) error
	FindByUserIDFn        func(userID uint) ([]models.Recipe, error)
	FindByIDFn            func(id uint) (*models.Recipe, error)
	FindByIDWithDetailsFn func(id uint) (*models.Recipe, error)
	UpdateFn              func(recipe *models.Recipe) error
	DeleteFn              func(recipe *models.Recipe) error
}

func (m *MockRecipeRepository) Create(r *models.Recipe) error {
	return m.CreateFn(r)
}

func (m *MockRecipeRepository) FindByUserID(userID uint) ([]models.Recipe, error) {
	return m.FindByUserIDFn(userID)
}

func (m *MockRecipeRepository) FindByID(id uint) (*models.Recipe, error) {
	return m.FindByIDFn(id)
}

func (m *MockRecipeRepository) FindByIDWithDetails(id uint) (*models.Recipe, error) {
	if m.FindByIDWithDetailsFn != nil {
		return m.FindByIDWithDetailsFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockRecipeRepository) Update(r *models.Recipe) error {
	return m.UpdateFn(r)
}

func (m *MockRecipeRepository) Delete(r *models.Recipe) error {
	return m.DeleteFn(r)
}

func TestCreateRecipe_Success(t *testing.T) {
	repo := &MockRecipeRepository{
		CreateFn: func(r *models.Recipe) error {
			return nil
		},
	}

	service := NewRecipeService(repo)

	req := dto.CreateRecipeRequest{
		Name:     "Test Recipe",
		Servings: 2,
	}

	err := service.CreateRecipe(1, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateRecipe_Unauthorized(t *testing.T) {
	repo := &MockRecipeRepository{
		FindByIDFn: func(id uint) (*models.Recipe, error) {
			return &models.Recipe{ID: id, UserID: 2}, nil
		},
	}

	service := NewRecipeService(repo)

	err := service.UpdateRecipe(1, 1, dto.UpdateRecipeRequest{})
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestDeleteRecipe_Success(t *testing.T) {
	repo := &MockRecipeRepository{
		FindByIDFn: func(id uint) (*models.Recipe, error) {
			return &models.Recipe{ID: id, UserID: 1}, nil
		},
		DeleteFn: func(r *models.Recipe) error {
			return nil
		},
	}

	service := NewRecipeService(repo)

	err := service.DeleteRecipe(1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetMyRecipes(t *testing.T) {
	repo := &MockRecipeRepository{
		FindByUserIDFn: func(userID uint) ([]models.Recipe, error) {
			return []models.Recipe{
				{ID: 1, Name: "A", PrepTime: 10, CookTime: 20},
			}, nil
		},
	}

	service := NewRecipeService(repo)

	res, err := service.GetMyRecipes(1)
	if err != nil || len(res) != 1 {
		t.Fatal("expected one recipe")
	}
}
