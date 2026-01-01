package services

import (
	"errors"
	"testing"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

// =====================================================
// MOCK REPOSITORIES
// =====================================================

type MockIngredientRepository struct {
	CreateFn   func(*models.Ingredient) error
	FindAllFn  func() ([]models.Ingredient, error)
	FindByIDFn func(uint) (*models.Ingredient, error)
}

func (m *MockIngredientRepository) Create(i *models.Ingredient) error {
	if m.CreateFn != nil {
		return m.CreateFn(i)
	}
	return nil
}
func (m *MockIngredientRepository) FindAll() ([]models.Ingredient, error) {
	if m.FindAllFn != nil {
		return m.FindAllFn()
	}
	return []models.Ingredient{}, nil
}
func (m *MockIngredientRepository) FindByID(id uint) (*models.Ingredient, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

type MockRecipeIngredientRepository struct {
	CreateFn         func(*models.RecipeIngredient) error
	FindByRecipeIDFn func(uint) ([]models.RecipeIngredient, error)
	FindByIDFn       func(uint) (*models.RecipeIngredient, error)
	DeleteFn         func(uint) error
}

func (m *MockRecipeIngredientRepository) Create(ri *models.RecipeIngredient) error {
	if m.CreateFn != nil {
		return m.CreateFn(ri)
	}
	return nil
}
func (m *MockRecipeIngredientRepository) FindByRecipeID(recipeID uint) ([]models.RecipeIngredient, error) {
	if m.FindByRecipeIDFn != nil {
		return m.FindByRecipeIDFn(recipeID)
	}
	return []models.RecipeIngredient{}, nil
}
func (m *MockRecipeIngredientRepository) FindByID(id uint) (*models.RecipeIngredient, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *MockRecipeIngredientRepository) Delete(id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

type MockRecipeRepoForIngredient struct {
	FindByIDFn func(uint) (*models.Recipe, error)
}

func (m *MockRecipeRepoForIngredient) FindByID(id uint) (*models.Recipe, error) {
	return m.FindByIDFn(id)
}
func (m *MockRecipeRepoForIngredient) Create(*models.Recipe) error                { return nil }
func (m *MockRecipeRepoForIngredient) FindByUserID(uint) ([]models.Recipe, error) { return nil, nil }
func (m *MockRecipeRepoForIngredient) FindByIDWithDetails(uint) (*models.Recipe, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *MockRecipeRepoForIngredient) Update(*models.Recipe) error { return nil }
func (m *MockRecipeRepoForIngredient) Delete(*models.Recipe) error { return nil }

// =====================================================
// TESTS
// =====================================================

func TestCreateIngredient_Success(t *testing.T) {
	service := NewIngredientService(&MockIngredientRepository{}, &MockRecipeIngredientRepository{}, &MockRecipeRepoForIngredient{})
	err := service.CreateIngredient(dto.CreateIngredientRequest{Name: "Onion"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetIngredients_Success(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{
			FindAllFn: func() ([]models.Ingredient, error) {
				return []models.Ingredient{{ID: 1, Name: "Salt"}}, nil
			},
		},
		&MockRecipeIngredientRepository{}, &MockRecipeRepoForIngredient{},
	)
	items, err := service.GetIngredients()
	if err != nil || len(items) != 1 {
		t.Fatal("failed to get ingredients")
	}
}

func TestGetIngredients_Error(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{
			FindAllFn: func() ([]models.Ingredient, error) {
				return nil, errors.New("db error")
			},
		},
		&MockRecipeIngredientRepository{}, &MockRecipeRepoForIngredient{},
	)
	_, err := service.GetIngredients()
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

func TestAddIngredientToRecipe_Success(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{},
		&MockRecipeIngredientRepository{},
		&MockRecipeRepoForIngredient{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{ID: id, UserID: 1}, nil
			},
		},
	)
	err := service.AddIngredientToRecipe(1, 1, dto.AddRecipeIngredientRequest{Quantity: 2})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestAddIngredientToRecipe_Unauthorized(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{}, &MockRecipeIngredientRepository{},
		&MockRecipeRepoForIngredient{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{ID: id, UserID: 2}, nil // Owner is 2, user is 1
			},
		},
	)
	err := service.AddIngredientToRecipe(1, 1, dto.AddRecipeIngredientRequest{})
	if err != ErrIngredientUnauthorized {
		t.Fatal("expected unauthorized error")
	}
}

func TestGetRecipeIngredients_Success(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{},
		&MockRecipeIngredientRepository{
			FindByRecipeIDFn: func(uint) ([]models.RecipeIngredient, error) {
				return []models.RecipeIngredient{
					{
						Quantity:   10,
						Unit:       "g",
						Ingredient: models.Ingredient{Name: "Sugar"}, // Preloaded data
					},
				}, nil
			},
		},
		&MockRecipeRepoForIngredient{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{ID: id, UserID: 1}, nil
			},
		},
	)

	items, err := service.GetRecipeIngredients(1, 1)
	if err != nil || len(items) != 1 || items[0].Name != "Sugar" {
		t.Fatalf("failed to fetch ingredients for recipe correctly")
	}
}

func TestGetRecipeIngredients_Unauthorized(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{}, &MockRecipeIngredientRepository{},
		&MockRecipeRepoForIngredient{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{UserID: 2}, nil
			},
		},
	)
	_, err := service.GetRecipeIngredients(1, 1)
	if err != ErrIngredientUnauthorized {
		t.Fatal("expected unauthorized error")
	}
}

func TestGetRecipeIngredients_RepoError(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{},
		&MockRecipeIngredientRepository{
			FindByRecipeIDFn: func(uint) ([]models.RecipeIngredient, error) {
				return nil, errors.New("query failed")
			},
		},
		&MockRecipeRepoForIngredient{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{UserID: 1}, nil
			},
		},
	)
	_, err := service.GetRecipeIngredients(1, 1)
	if err == nil || err.Error() != "query failed" {
		t.Fatal("expected repo error")
	}
}

func TestRemoveRecipeIngredient_Success(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{},
		&MockRecipeIngredientRepository{
			FindByIDFn: func(id uint) (*models.RecipeIngredient, error) {
				return &models.RecipeIngredient{ID: id, RecipeID: 10}, nil
			},
		},
		&MockRecipeRepoForIngredient{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{ID: 10, UserID: 1}, nil
			},
		},
	)
	err := service.RemoveRecipeIngredient(1, 1)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestRemoveRecipeIngredient_RecipeError(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{},
		&MockRecipeIngredientRepository{
			FindByIDFn: func(id uint) (*models.RecipeIngredient, error) {
				return &models.RecipeIngredient{RecipeID: 10}, nil
			},
		},
		&MockRecipeRepoForIngredient{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return nil, errors.New("recipe not found")
			},
		},
	)
	err := service.RemoveRecipeIngredient(1, 1)
	if err == nil {
		t.Fatal("expected error when recipe lookup fails")
	}
}
