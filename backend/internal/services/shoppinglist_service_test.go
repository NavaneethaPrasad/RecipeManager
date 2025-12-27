package services

import (
	"testing"
	"time"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type MockMealPlanRepoForShoppingList struct {
	FindRangeFn func(uint, time.Time, time.Time) ([]models.MealPlan, error)
}

func (m *MockMealPlanRepoForShoppingList) FindByUserAndDateRange(
	userID uint,
	start, end time.Time,
) ([]models.MealPlan, error) {
	return m.FindRangeFn(userID, start, end)
}

func (m *MockMealPlanRepoForShoppingList) Create(*models.MealPlan) error { return nil }
func (m *MockMealPlanRepoForShoppingList) FindByUserAndDate(uint, time.Time) ([]models.MealPlan, error) {
	return nil, nil
}
func (m *MockMealPlanRepoForShoppingList) FindByID(uint) (*models.MealPlan, error) { return nil, nil }
func (m *MockMealPlanRepoForShoppingList) Update(*models.MealPlan) error           { return nil }
func (m *MockMealPlanRepoForShoppingList) Delete(*models.MealPlan) error           { return nil }
func (m *MockMealPlanRepoForShoppingList) FindDuplicate(uint, time.Time, string) error {
	return gorm.ErrRecordNotFound
}

type MockRecipeIngredientRepo struct {
	FindFn     func(uint) ([]models.RecipeIngredient, error)
	FindByIDFn func(uint) (*models.RecipeIngredient, error)
	CreateFn   func(*models.RecipeIngredient) error
	DeleteFn   func(uint) error
}

func (m *MockRecipeIngredientRepo) FindByRecipeID(recipeID uint) ([]models.RecipeIngredient, error) {
	if m.FindFn != nil {
		return m.FindFn(recipeID)
	}
	return nil, nil
}

func (m *MockRecipeIngredientRepo) FindByID(id uint) (*models.RecipeIngredient, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockRecipeIngredientRepo) Create(ri *models.RecipeIngredient) error {
	if m.CreateFn != nil {
		return m.CreateFn(ri)
	}
	return nil
}

func (m *MockRecipeIngredientRepo) Delete(id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

type MockShoppingListRepo struct {
	CreateFn     func(*models.ShoppingList) error
	CreateItemFn func(*models.ShoppingListItem) error
	FindByIDFn   func(uint) (*models.ShoppingList, error)
	FindItemsFn  func(uint) ([]models.ShoppingListItem, error)
	FindItemFn   func(uint) (*models.ShoppingListItem, error)
	UpdateItemFn func(*models.ShoppingListItem) error
}

func (m *MockShoppingListRepo) Create(sl *models.ShoppingList) error {
	return m.CreateFn(sl)
}
func (m *MockShoppingListRepo) CreateItem(item *models.ShoppingListItem) error {
	return m.CreateItemFn(item)
}
func (m *MockShoppingListRepo) FindByID(id uint) (*models.ShoppingList, error) {
	return m.FindByIDFn(id)
}
func (m *MockShoppingListRepo) FindItemsByListID(id uint) ([]models.ShoppingListItem, error) {
	return m.FindItemsFn(id)
}
func (m *MockShoppingListRepo) FindItemByID(id uint) (*models.ShoppingListItem, error) {
	return m.FindItemFn(id)
}
func (m *MockShoppingListRepo) UpdateItem(item *models.ShoppingListItem) error {
	return m.UpdateItemFn(item)
}

func TestGenerateShoppingList_Success(t *testing.T) {
	service := NewShoppingListService(
		&MockMealPlanRepoForShoppingList{
			FindRangeFn: func(uint, time.Time, time.Time) ([]models.MealPlan, error) {
				return []models.MealPlan{
					{RecipeID: 1},
				}, nil
			},
		},
		&MockRecipeIngredientRepo{
			FindFn: func(uint) ([]models.RecipeIngredient, error) {
				return []models.RecipeIngredient{
					{
						IngredientID: 1,
						Quantity:     2,
						Unit:         "pcs",
						Ingredient:   models.Ingredient{Name: "Onion"},
					},
				}, nil
			},
		},
		&MockShoppingListRepo{
			CreateFn:     func(*models.ShoppingList) error { return nil },
			CreateItemFn: func(*models.ShoppingListItem) error { return nil },
		},
	)

	resp, err := service.Generate(1, "2025-01-01", "2025-01-07")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp.Items))
	}
}

func TestGenerateShoppingList_InvalidDate(t *testing.T) {
	service := NewShoppingListService(
		&MockMealPlanRepoForShoppingList{},
		&MockRecipeIngredientRepo{},
		&MockShoppingListRepo{},
	)

	_, err := service.Generate(1, "2025-01-10", "2025-01-01")
	if err != ErrInvalidDateRange {
		t.Fatalf("expected ErrInvalidDateRange")
	}
}

func TestToggleItemChecked_Success(t *testing.T) {
	service := NewShoppingListService(
		&MockMealPlanRepoForShoppingList{},
		&MockRecipeIngredientRepo{},
		&MockShoppingListRepo{
			FindItemFn: func(uint) (*models.ShoppingListItem, error) {
				return &models.ShoppingListItem{
					ID:             1,
					ShoppingListID: 1,
					Checked:        false,
				}, nil
			},
			FindByIDFn: func(uint) (*models.ShoppingList, error) {
				return &models.ShoppingList{
					ID:     1,
					UserID: 1,
				}, nil
			},
			UpdateItemFn: func(item *models.ShoppingListItem) error {
				if !item.Checked {
					t.Fatal("item should be checked")
				}
				return nil
			},
		},
	)

	err := service.ToggleItemChecked(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
