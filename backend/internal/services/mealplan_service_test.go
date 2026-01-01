package services

import (
	"errors"
	"testing"
	"time"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

// =====================================================
// MOCK REPOSITORIES
// =====================================================

type MockMealPlanRepo struct {
	CreateFn                 func(*models.MealPlan) error
	FindByUserAndDateFn      func(uint, time.Time) ([]models.MealPlan, error)
	FindByIDFn               func(uint) (*models.MealPlan, error)
	UpdateFn                 func(*models.MealPlan) error
	DeleteFn                 func(*models.MealPlan) error
	FindDuplicateFn          func(uint, time.Time, string) error
	FindByUserAndDateRangeFn func(uint, time.Time, time.Time) ([]models.MealPlan, error)
}

func (m *MockMealPlanRepo) Create(mp *models.MealPlan) error { return m.CreateFn(mp) }
func (m *MockMealPlanRepo) FindByUserAndDate(u uint, d time.Time) ([]models.MealPlan, error) {
	return m.FindByUserAndDateFn(u, d)
}
func (m *MockMealPlanRepo) FindByID(id uint) (*models.MealPlan, error) { return m.FindByIDFn(id) }
func (m *MockMealPlanRepo) Update(mp *models.MealPlan) error           { return m.UpdateFn(mp) }
func (m *MockMealPlanRepo) Delete(mp *models.MealPlan) error           { return m.DeleteFn(mp) }
func (m *MockMealPlanRepo) FindDuplicate(u uint, d time.Time, mt string) error {
	if m.FindDuplicateFn != nil {
		return m.FindDuplicateFn(u, d, mt)
	}
	return gorm.ErrRecordNotFound
}
func (m *MockMealPlanRepo) FindByUserAndDateRange(u uint, s, e time.Time) ([]models.MealPlan, error) {
	return m.FindByUserAndDateRangeFn(u, s, e)
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

// =====================================================
// TESTS
// =====================================================

func TestCreateMealPlan_Success(t *testing.T) {
	service := NewMealPlanService(
		&MockMealPlanRepo{
			CreateFn: func(mp *models.MealPlan) error { return nil },
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

	err := service.Create(1, dto.CreateMealPlanRequest{
		RecipeID: 1, Date: "2025-01-01", MealType: "dinner", TargetServings: 4,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreateMealPlan_Errors(t *testing.T) {
	service := NewMealPlanService(&MockMealPlanRepo{}, &MockRecipeRepoForMealPlan{})

	t.Run("Invalid Date Format", func(t *testing.T) {
		err := service.Create(1, dto.CreateMealPlanRequest{Date: "01-01-2025"})
		if err == nil {
			t.Fatal("expected error for invalid date format")
		}
	})

	t.Run("Recipe Lookup Error", func(t *testing.T) {
		service.(*mealPlanService).RecipeRepo = &MockRecipeRepoForMealPlan{
			FindByIDFn: func(id uint) (*models.Recipe, error) { return nil, errors.New("db error") },
		}
		err := service.Create(1, dto.CreateMealPlanRequest{Date: "2025-01-01", RecipeID: 1})
		if err == nil {
			t.Fatal("expected error when recipe lookup fails")
		}
	})

	t.Run("Unauthorized Recipe Usage", func(t *testing.T) {
		service.(*mealPlanService).RecipeRepo = &MockRecipeRepoForMealPlan{
			FindByIDFn: func(id uint) (*models.Recipe, error) { return &models.Recipe{UserID: 2}, nil },
		}
		err := service.Create(1, dto.CreateMealPlanRequest{Date: "2025-01-01", RecipeID: 1})
		if err != ErrMealUnauthorized {
			t.Fatal("expected unauthorized error")
		}
	})
}

func TestGetByDateRange_Success(t *testing.T) {
	service := NewMealPlanService(
		&MockMealPlanRepo{
			FindByUserAndDateRangeFn: func(u uint, s, e time.Time) ([]models.MealPlan, error) {
				return []models.MealPlan{
					{ID: 1, Date: s, Recipe: models.Recipe{Name: "Pizza"}},
				}, nil
			},
		},
		&MockRecipeRepoForMealPlan{},
	)

	res, err := service.GetByDateRange(1, "2025-01-01", "2025-01-07")
	if err != nil || len(res) != 1 || res[0].Recipe.Name != "Pizza" {
		t.Fatal("failed to get meal plan range")
	}
}

func TestGetByDateRange_DateError(t *testing.T) {
	service := NewMealPlanService(&MockMealPlanRepo{}, &MockRecipeRepoForMealPlan{})
	_, err := service.GetByDateRange(1, "invalid", "2025-01-07")
	if err == nil {
		t.Fatal("expected error for invalid start date")
	}
	_, err = service.GetByDateRange(1, "2025-01-01", "invalid")
	if err == nil {
		t.Fatal("expected error for invalid end date")
	}
}

func TestGetByDate_Errors(t *testing.T) {
	service := NewMealPlanService(&MockMealPlanRepo{}, &MockRecipeRepoForMealPlan{})

	t.Run("Invalid Date", func(t *testing.T) {
		_, err := service.GetByDate(1, "invalid")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("Repo Error", func(t *testing.T) {
		service.(*mealPlanService).Repo = &MockMealPlanRepo{
			FindByUserAndDateFn: func(u uint, d time.Time) ([]models.MealPlan, error) {
				return nil, errors.New("db error")
			},
		}
		_, err := service.GetByDate(1, "2025-01-01")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestUpdateMealPlan_Failures(t *testing.T) {
	t.Run("Record Not Found", func(t *testing.T) {
		service := NewMealPlanService(&MockMealPlanRepo{
			FindByIDFn: func(id uint) (*models.MealPlan, error) { return nil, gorm.ErrRecordNotFound },
		}, &MockRecipeRepoForMealPlan{})
		err := service.Update(1, 1, dto.UpdateMealPlanRequest{})
		if err != gorm.ErrRecordNotFound {
			t.Fatal("expected record not found")
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		service := NewMealPlanService(&MockMealPlanRepo{
			FindByIDFn: func(id uint) (*models.MealPlan, error) { return &models.MealPlan{UserID: 2}, nil },
		}, &MockRecipeRepoForMealPlan{})
		err := service.Update(1, 1, dto.UpdateMealPlanRequest{})
		if err != ErrMealUnauthorized {
			t.Fatal("expected unauthorized")
		}
	})
}

func TestDeleteMealPlan_Unauthorized(t *testing.T) {
	service := NewMealPlanService(&MockMealPlanRepo{
		FindByIDFn: func(id uint) (*models.MealPlan, error) { return &models.MealPlan{UserID: 2}, nil },
	}, &MockRecipeRepoForMealPlan{})

	err := service.Delete(1, 1)
	if err != ErrMealUnauthorized {
		t.Fatal("expected unauthorized error")
	}
}

func TestDeleteMealPlan_FindError(t *testing.T) {
	service := NewMealPlanService(&MockMealPlanRepo{
		FindByIDFn: func(id uint) (*models.MealPlan, error) { return nil, errors.New("find error") },
	}, &MockRecipeRepoForMealPlan{})

	err := service.Delete(1, 1)
	if err == nil {
		t.Fatal("expected find error")
	}
}
