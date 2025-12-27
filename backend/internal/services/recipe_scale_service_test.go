package services

import (
	"testing"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type MockRecipeRepoForScale struct {
	FindByIDWithDetailsFn func(uint) (*models.Recipe, error)
}

func (m *MockRecipeRepoForScale) FindByIDWithDetails(id uint) (*models.Recipe, error) {
	return m.FindByIDWithDetailsFn(id)
}

func (m *MockRecipeRepoForScale) Create(*models.Recipe) error {
	return nil
}

func (m *MockRecipeRepoForScale) FindByUserID(uint) ([]models.Recipe, error) {
	return nil, nil
}

func (m *MockRecipeRepoForScale) FindByID(uint) (*models.Recipe, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *MockRecipeRepoForScale) Update(*models.Recipe) error {
	return nil
}

func (m *MockRecipeRepoForScale) Delete(*models.Recipe) error {
	return nil
}

func TestScaleRecipe_Success(t *testing.T) {
	service := NewRecipeScaleService(
		&MockRecipeRepoForScale{
			FindByIDWithDetailsFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:       id,
					UserID:   1,
					Name:     "Chicken Curry",
					Servings: 2,
					Ingredients: []models.RecipeIngredient{
						{
							IngredientID: 1,
							Quantity:     2,
							Unit:         "pcs",
							Ingredient: models.Ingredient{
								Name: "Onion",
							},
						},
					},
				}, nil
			},
		},
	)

	resp, err := service.ScaleRecipe(1, 1, 4)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.ScaledServings != 4 {
		t.Fatalf("expected servings 4, got %d", resp.ScaledServings)
	}

	if resp.Ingredients[0].Quantity != 4 {
		t.Fatalf("expected scaled quantity 4, got %f", resp.Ingredients[0].Quantity)
	}
}

func TestScaleRecipe_Unauthorized(t *testing.T) {
	service := NewRecipeScaleService(
		&MockRecipeRepoForScale{
			FindByIDWithDetailsFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:       id,
					UserID:   2,
					Servings: 2,
				}, nil
			},
		},
	)

	_, err := service.ScaleRecipe(1, 1, 4)
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestScaleRecipe_InvalidServings(t *testing.T) {
	service := NewRecipeScaleService(&MockRecipeRepoForScale{})

	_, err := service.ScaleRecipe(1, 1, 0)
	if err != ErrInvalidServings {
		t.Fatalf("expected ErrInvalidServings, got %v", err)
	}
}
