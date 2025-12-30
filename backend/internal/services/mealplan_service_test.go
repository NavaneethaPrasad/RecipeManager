package services

import (
	"testing"
	"time"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type MockMealPlanRepo struct {
	CreateFn                 func(*models.MealPlan) error
	FindByUserAndDateFn      func(uint, time.Time) ([]models.MealPlan, error)
	FindByIDFn               func(uint) (*models.MealPlan, error)
	UpdateFn                 func(*models.MealPlan) error
	DeleteFn                 func(*models.MealPlan) error
	FindDuplicateFn          func(uint, time.Time, string) error
	FindByUserAndDateRangeFn func(uint, time.Time, time.Time) ([]models.MealPlan, error)
}

func (m *MockMealPlanRepo) Create(mp *models.MealPlan) error {
	return m.CreateFn(mp)
}

func (m *MockMealPlanRepo) FindByUserAndDate(userID uint, date time.Time) ([]models.MealPlan, error) {
	return m.FindByUserAndDateFn(userID, date)
}

func (m *MockMealPlanRepo) FindByID(id uint) (*models.MealPlan, error) {
	return m.FindByIDFn(id)
}

func (m *MockMealPlanRepo) Update(mp *models.MealPlan) error {
	return m.UpdateFn(mp)
}

func (m *MockMealPlanRepo) Delete(mp *models.MealPlan) error {
	return m.DeleteFn(mp)
}

func (m *MockMealPlanRepo) FindDuplicate(userID uint, date time.Time, mealType string) error {
	if m.FindDuplicateFn != nil {
		return m.FindDuplicateFn(userID, date, mealType)
	}
	return gorm.ErrRecordNotFound
}

func (m *MockMealPlanRepo) FindByUserAndDateRange(userID uint, start, end time.Time) ([]models.MealPlan, error) {
	if m.FindByUserAndDateRangeFn != nil {
		return m.FindByUserAndDateRangeFn(userID, start, end)
	}
	return []models.MealPlan{}, nil
}

type MockRecipeRepoForMealPlan struct {
	FindByIDFn func(uint) (*models.Recipe, error)
}

func (m *MockRecipeRepoForMealPlan) FindByID(id uint) (*models.Recipe, error) {
	return m.FindByIDFn(id)
}

func (m *MockRecipeRepoForMealPlan) Create(*models.Recipe) error                { return nil }
func (m *MockRecipeRepoForMealPlan) FindByUserID(uint) ([]models.Recipe, error) { return nil, nil }
func (m *MockRecipeRepoForMealPlan) FindByIDWithDetails(uint) (*models.Recipe, error) {
	return nil, nil
}
func (m *MockRecipeRepoForMealPlan) Update(*models.Recipe) error { return nil }
func (m *MockRecipeRepoForMealPlan) Delete(*models.Recipe) error { return nil }

func TestCreateMealPlan_Success(t *testing.T) {
	service := NewMealPlanService(
		&MockMealPlanRepo{
			CreateFn: func(mp *models.MealPlan) error {
				if mp.TargetServings != 4 {
					t.Fatal("expected target servings to be saved")
				}
				return nil
			},
			FindDuplicateFn: func(uint, time.Time, string) error {
				return gorm.ErrRecordNotFound
			},
		},
		&MockRecipeRepoForMealPlan{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{ID: id, UserID: 1}, nil
			},
		},
	)

	err := service.Create(
		1,
		dto.CreateMealPlanRequest{
			RecipeID:       1,
			Date:           "2025-01-01",
			MealType:       "dinner",
			TargetServings: 4,
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreateMealPlan_Duplicate(t *testing.T) {
	service := NewMealPlanService(
		&MockMealPlanRepo{
			FindDuplicateFn: func(uint, time.Time, string) error {
				return nil
			},
		},
		&MockRecipeRepoForMealPlan{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{ID: id, UserID: 1}, nil
			},
		},
	)

	err := service.Create(
		1,
		dto.CreateMealPlanRequest{
			RecipeID:       1,
			Date:           "2025-01-01",
			MealType:       "dinner",
			TargetServings: 2,
		},
	)

	if err != ErrMealExists {
		t.Fatal("expected ErrMealExists")
	}
}

func TestGetMealPlansByDate_Success(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2025-01-01")

	service := NewMealPlanService(
		&MockMealPlanRepo{
			FindByUserAndDateFn: func(userID uint, d time.Time) ([]models.MealPlan, error) {
				if !d.Equal(date) {
					t.Fatal("date mismatch")
				}
				return []models.MealPlan{
					{
						ID:             1,
						UserID:         userID,
						TargetServings: 2,
						Recipe:         models.Recipe{Name: "Pasta"},
					},
				}, nil
			},
		},
		&MockRecipeRepoForMealPlan{},
	)

	plans, err := service.GetByDate(1, "2025-01-01")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(plans) != 1 {
		t.Fatalf("expected 1 meal plan, got %d", len(plans))
	}
	if plans[0].Recipe.Name != "Pasta" {
		t.Error("expected recipe name to be populated")
	}
}

func TestUpdateMealPlan_Unauthorized(t *testing.T) {
	service := NewMealPlanService(
		&MockMealPlanRepo{
			FindByIDFn: func(id uint) (*models.MealPlan, error) {
				return &models.MealPlan{ID: id, UserID: 2}, nil // Different User
			},
		},
		&MockRecipeRepoForMealPlan{},
	)

	err := service.Update(
		1,
		1,
		dto.UpdateMealPlanRequest{
			RecipeID:       1,
			MealType:       "breakfast",
			TargetServings: 2,
		},
	)

	if err != ErrMealUnauthorized {
		t.Fatalf("expected ErrMealUnauthorized, got %v", err)
	}
}

func TestUpdateMealPlan_Success(t *testing.T) {
	service := NewMealPlanService(
		&MockMealPlanRepo{
			FindByIDFn: func(id uint) (*models.MealPlan, error) {
				return &models.MealPlan{ID: id, UserID: 1}, nil
			},
			UpdateFn: func(mp *models.MealPlan) error {
				if mp.TargetServings != 5 {
					t.Fatal("target servings not updated")
				}
				return nil
			},
		},
		&MockRecipeRepoForMealPlan{},
	)

	err := service.Update(
		1,
		1,
		dto.UpdateMealPlanRequest{
			TargetServings: 5,
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteMealPlan_Success(t *testing.T) {
	service := NewMealPlanService(
		&MockMealPlanRepo{
			FindByIDFn: func(id uint) (*models.MealPlan, error) {
				return &models.MealPlan{ID: id, UserID: 1}, nil
			},
			DeleteFn: func(mp *models.MealPlan) error {
				return nil
			},
		},
		&MockRecipeRepoForMealPlan{},
	)

	err := service.Delete(1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
