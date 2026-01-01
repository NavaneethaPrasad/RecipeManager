package services

import (
	"errors"
	"testing"
	"time"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

// =====================================================
// MOCK REPOSITORIES
// =====================================================

type MockMealPlanRepoForShoppingList struct {
	FindRangeFn func(uint, time.Time, time.Time) ([]models.MealPlan, error)
}

func (m *MockMealPlanRepoForShoppingList) FindByUserAndDateRange(u uint, s, e time.Time) ([]models.MealPlan, error) {
	if m.FindRangeFn != nil {
		return m.FindRangeFn(u, s, e)
	}
	return []models.MealPlan{}, nil
}

func (m *MockMealPlanRepoForShoppingList) Create(mp *models.MealPlan) error { return nil }
func (m *MockMealPlanRepoForShoppingList) FindByUserAndDate(u uint, d time.Time) ([]models.MealPlan, error) {
	return nil, nil
}
func (m *MockMealPlanRepoForShoppingList) FindByID(u uint) (*models.MealPlan, error) { return nil, nil }
func (m *MockMealPlanRepoForShoppingList) Update(mp *models.MealPlan) error          { return nil }
func (m *MockMealPlanRepoForShoppingList) Delete(mp *models.MealPlan) error          { return nil }
func (m *MockMealPlanRepoForShoppingList) FindDuplicate(u uint, d time.Time, t string) error {
	return gorm.ErrRecordNotFound
}

type MockRecipeIngredientRepo struct{} // Not used in the current Generate logic but required by interface

func (m *MockRecipeIngredientRepo) FindByRecipeID(id uint) ([]models.RecipeIngredient, error) {
	return nil, nil
}
func (m *MockRecipeIngredientRepo) Create(ri *models.RecipeIngredient) error { return nil }
func (m *MockRecipeIngredientRepo) FindByID(id uint) (*models.RecipeIngredient, error) {
	return nil, nil
}
func (m *MockRecipeIngredientRepo) Delete(id uint) error { return nil }

type MockShoppingListRepo struct {
	CreateFn     func(*models.ShoppingList) error
	CreateItemFn func(*models.ShoppingListItem) error
	FindByIDFn   func(uint) (*models.ShoppingList, error)
	FindItemsFn  func(uint) ([]models.ShoppingListItem, error)
	FindItemFn   func(uint) (*models.ShoppingListItem, error)
	UpdateItemFn func(*models.ShoppingListItem) error
}

func (m *MockShoppingListRepo) Create(sl *models.ShoppingList) error {
	if m.CreateFn != nil {
		sl.ID = 1
		return m.CreateFn(sl)
	}
	return nil
}
func (m *MockShoppingListRepo) CreateItem(i *models.ShoppingListItem) error {
	if m.CreateItemFn != nil {
		return m.CreateItemFn(i)
	}
	return nil
}
func (m *MockShoppingListRepo) FindByID(id uint) (*models.ShoppingList, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *MockShoppingListRepo) FindItemsByListID(id uint) ([]models.ShoppingListItem, error) {
	if m.FindItemsFn != nil {
		return m.FindItemsFn(id)
	}
	return []models.ShoppingListItem{}, nil
}
func (m *MockShoppingListRepo) FindItemByID(id uint) (*models.ShoppingListItem, error) {
	if m.FindItemFn != nil {
		return m.FindItemFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *MockShoppingListRepo) UpdateItem(i *models.ShoppingListItem) error {
	if m.UpdateItemFn != nil {
		return m.UpdateItemFn(i)
	}
	return nil
}

// =====================================================
// TESTS
// =====================================================

func TestGenerateShoppingList(t *testing.T) {
	t.Run("Success Path with Scaling", func(t *testing.T) {
		service := NewShoppingListService(
			&MockMealPlanRepoForShoppingList{
				FindRangeFn: func(u uint, s, e time.Time) ([]models.MealPlan, error) {
					return []models.MealPlan{{
						RecipeID:       1,
						TargetServings: 4,
						Recipe: models.Recipe{
							Servings: 2,
							Ingredients: []models.RecipeIngredient{
								{
									IngredientID: 10,
									Quantity:     100,
									Unit:         "g",
									Ingredient:   models.Ingredient{Name: "Flour"},
								},
							},
						},
					}}, nil
				},
			},
			&MockRecipeIngredientRepo{},
			&MockShoppingListRepo{
				CreateFn:     func(*models.ShoppingList) error { return nil },
				CreateItemFn: func(*models.ShoppingListItem) error { return nil },
			},
		)

		resp, err := service.Generate(1, "2025-01-01", "2025-01-07")
		if err != nil {
			t.Fatalf("Expected nil, got %v", err)
		}
		if len(resp.Items) != 1 {
			t.Fatalf("Expected 1 item, got %d", len(resp.Items))
		}
		if resp.Items[0].Quantity != 200 {
			t.Errorf("Scaling failed: expected 200, got %f", resp.Items[0].Quantity)
		}
	})

	t.Run("MealPlan Repo Error", func(t *testing.T) {
		service := NewShoppingListService(
			&MockMealPlanRepoForShoppingList{
				FindRangeFn: func(u uint, s, e time.Time) ([]models.MealPlan, error) {
					return nil, errors.New("db error")
				},
			},
			&MockRecipeIngredientRepo{},
			&MockShoppingListRepo{},
		)
		_, err := service.Generate(1, "2025-01-01", "2025-01-07")
		if err == nil || err.Error() != "db error" {
			t.Errorf("Expected db error, got %v", err)
		}
	})

	t.Run("Shopping List Create Error", func(t *testing.T) {
		service := NewShoppingListService(
			&MockMealPlanRepoForShoppingList{
				FindRangeFn: func(u uint, s, e time.Time) ([]models.MealPlan, error) {
					return []models.MealPlan{}, nil
				},
			},
			&MockRecipeIngredientRepo{},
			&MockShoppingListRepo{
				CreateFn: func(*models.ShoppingList) error { return errors.New("header fail") },
			},
		)
		_, err := service.Generate(1, "2025-01-01", "2025-01-07")
		if err == nil || err.Error() != "header fail" {
			t.Errorf("Expected header fail, got %v", err)
		}
	})

	t.Run("Shopping Item Create Error", func(t *testing.T) {
		service := NewShoppingListService(
			&MockMealPlanRepoForShoppingList{
				FindRangeFn: func(u uint, s, e time.Time) ([]models.MealPlan, error) {
					return []models.MealPlan{{
						Recipe: models.Recipe{
							Servings: 1,
							Ingredients: []models.RecipeIngredient{
								{IngredientID: 1, Ingredient: models.Ingredient{Name: "X"}},
							},
						},
					}}, nil
				},
			},
			&MockRecipeIngredientRepo{},
			&MockShoppingListRepo{
				CreateFn:     func(*models.ShoppingList) error { return nil },
				CreateItemFn: func(*models.ShoppingListItem) error { return errors.New("item fail") },
			},
		)
		_, err := service.Generate(1, "2025-01-01", "2025-01-07")
		if err == nil || err.Error() != "item fail" {
			t.Errorf("Expected item fail, got %v", err)
		}
	})

	t.Run("Date Errors", func(t *testing.T) {
		service := NewShoppingListService(&MockMealPlanRepoForShoppingList{}, &MockRecipeIngredientRepo{}, &MockShoppingListRepo{})
		_, err := service.Generate(1, "invalid", "2025-01-01")
		if err == nil {
			t.Error("Expected parsing error")
		}
		_, err = service.Generate(1, "2025-01-10", "2025-01-01")
		if err != ErrInvalidDateRange {
			t.Error("Expected ErrInvalidDateRange")
		}
	})
}

func TestGetShoppingListByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		service := NewShoppingListService(nil, &MockRecipeIngredientRepo{}, &MockShoppingListRepo{
			FindByIDFn: func(uint) (*models.ShoppingList, error) {
				return &models.ShoppingList{UserID: 1, StartDate: time.Now(), EndDate: time.Now()}, nil
			},
			FindItemsFn: func(uint) ([]models.ShoppingListItem, error) {
				return []models.ShoppingListItem{{IngredientID: 1, Ingredient: models.Ingredient{Name: "Salt"}, Quantity: 5}}, nil
			},
		})
		resp, err := service.GetShoppingListByID(1, 1)
		if err != nil || len(resp.Items) == 0 {
			t.Fatal("Failed to fetch list")
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		service := NewShoppingListService(nil, &MockRecipeIngredientRepo{}, &MockShoppingListRepo{
			FindByIDFn: func(uint) (*models.ShoppingList, error) {
				return &models.ShoppingList{UserID: 99}, nil
			},
		})
		_, err := service.GetShoppingListByID(1, 1)
		if err != ErrUnauthorized {
			t.Error("Expected unauthorized error")
		}
	})
}

func TestToggleItemChecked(t *testing.T) {
	service := NewShoppingListService(nil, &MockRecipeIngredientRepo{}, &MockShoppingListRepo{
		FindItemFn:   func(uint) (*models.ShoppingListItem, error) { return &models.ShoppingListItem{ShoppingListID: 1}, nil },
		FindByIDFn:   func(uint) (*models.ShoppingList, error) { return &models.ShoppingList{UserID: 1}, nil },
		UpdateItemFn: func(*models.ShoppingListItem) error { return nil },
	})
	err := service.ToggleItemChecked(1, 1)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
