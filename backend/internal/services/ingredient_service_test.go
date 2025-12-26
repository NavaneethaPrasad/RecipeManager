package services

import (
	"testing"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

//
// =======================
// MOCKS
// =======================
//

// ---------- Ingredient Repository Mock ----------
// Renamed to avoid collision

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

// ---------- Recipe Ingredient Repository Mock ----------
// Renamed clearly

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

// ---------- Recipe Repository Mock (for Ingredient Service) ----------
// Name changed to avoid collision with recipe tests

type MockRecipeRepoForIngredient struct {
	FindByIDFn func(uint) (*models.Recipe, error)
}

func (m *MockRecipeRepoForIngredient) FindByID(id uint) (*models.Recipe, error) {
	return m.FindByIDFn(id)
}

/* ---- Dummy methods to satisfy RecipeRepository interface ---- */

func (m *MockRecipeRepoForIngredient) Create(*models.Recipe) error {
	return nil
}

func (m *MockRecipeRepoForIngredient) FindByUserID(uint) ([]models.Recipe, error) {
	return nil, nil
}

func (m *MockRecipeRepoForIngredient) FindByIDWithDetails(uint) (*models.Recipe, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *MockRecipeRepoForIngredient) Update(*models.Recipe) error {
	return nil
}

func (m *MockRecipeRepoForIngredient) Delete(*models.Recipe) error {
	return nil
}

//
// =======================
// TESTS
// =======================
//

// ---------- Ingredient Master ----------

func TestCreateIngredient_Success(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{
			CreateFn: func(i *models.Ingredient) error {
				if i.Name == "" {
					t.Fatal("ingredient name should not be empty")
				}
				return nil
			},
		},
		&MockRecipeIngredientRepository{},
		&MockRecipeRepoForIngredient{},
	)

	err := service.CreateIngredient(dto.CreateIngredientRequest{
		Name: "Onion",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetIngredients_Success(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{
			FindAllFn: func() ([]models.Ingredient, error) {
				return []models.Ingredient{
					{ID: 1, Name: "Salt"},
					{ID: 2, Name: "Pepper"},
				}, nil
			},
		},
		&MockRecipeIngredientRepository{},
		&MockRecipeRepoForIngredient{},
	)

	items, err := service.GetIngredients()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 ingredients, got %d", len(items))
	}
}

// ---------- Recipe Ingredients ----------

func TestAddIngredientToRecipe_Unauthorized(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{},
		&MockRecipeIngredientRepository{},
		&MockRecipeRepoForIngredient{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 2, // not owner
				}, nil
			},
		},
	)

	err := service.AddIngredientToRecipe(
		1,
		1,
		dto.AddRecipeIngredientRequest{
			IngredientID: 1,
			Quantity:     2,
			Unit:         "pcs",
		},
	)

	if err != ErrIngredientUnauthorized {
		t.Fatalf("expected ErrIngredientUnauthorized, got %v", err)
	}
}

func TestAddIngredientToRecipe_Success(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{},
		&MockRecipeIngredientRepository{
			CreateFn: func(ri *models.RecipeIngredient) error {
				if ri.Quantity <= 0 {
					t.Fatal("quantity must be > 0")
				}
				return nil
			},
		},
		&MockRecipeRepoForIngredient{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 1,
				}, nil
			},
		},
	)

	err := service.AddIngredientToRecipe(
		1,
		1,
		dto.AddRecipeIngredientRequest{
			IngredientID: 1,
			Quantity:     2,
			Unit:         "pcs",
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRemoveRecipeIngredient_Unauthorized(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{},
		&MockRecipeIngredientRepository{
			FindByIDFn: func(id uint) (*models.RecipeIngredient, error) {
				return &models.RecipeIngredient{
					ID:       id,
					RecipeID: 1,
				}, nil
			},
		},
		&MockRecipeRepoForIngredient{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 2,
				}, nil
			},
		},
	)

	err := service.RemoveRecipeIngredient(1, 1)
	if err != ErrIngredientUnauthorized {
		t.Fatalf("expected ErrIngredientUnauthorized, got %v", err)
	}
}

func TestRemoveRecipeIngredient_Success(t *testing.T) {
	service := NewIngredientService(
		&MockIngredientRepository{},
		&MockRecipeIngredientRepository{
			FindByIDFn: func(id uint) (*models.RecipeIngredient, error) {
				return &models.RecipeIngredient{
					ID:       id,
					RecipeID: 1,
				}, nil
			},
			DeleteFn: func(id uint) error {
				return nil
			},
		},
		&MockRecipeRepoForIngredient{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 1,
				}, nil
			},
		},
	)

	err := service.RemoveRecipeIngredient(1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
