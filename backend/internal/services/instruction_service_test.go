package services

import (
	"testing"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type MockInstructionRepository struct {
	CreateFn         func(*models.Instruction) error
	FindByRecipeIDFn func(uint) ([]models.Instruction, error)
	FindByIDFn       func(uint) (*models.Instruction, error)
	UpdateFn         func(*models.Instruction) error
	DeleteFn         func(*models.Instruction) error
}

func (m *MockInstructionRepository) Create(i *models.Instruction) error {
	if m.CreateFn != nil {
		return m.CreateFn(i)
	}
	return nil
}

func (m *MockInstructionRepository) FindByRecipeID(recipeID uint) ([]models.Instruction, error) {
	if m.FindByRecipeIDFn != nil {
		return m.FindByRecipeIDFn(recipeID)
	}
	return []models.Instruction{}, nil
}

func (m *MockInstructionRepository) FindByID(id uint) (*models.Instruction, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockInstructionRepository) Update(i *models.Instruction) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(i)
	}
	return nil
}

func (m *MockInstructionRepository) Delete(i *models.Instruction) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(i)
	}
	return nil
}

type MockRecipeRepoForInstruction struct {
	FindByIDFn func(uint) (*models.Recipe, error)
}

func (m *MockRecipeRepoForInstruction) FindByID(id uint) (*models.Recipe, error) {
	return m.FindByIDFn(id)
}

func (m *MockRecipeRepoForInstruction) Create(*models.Recipe) error {
	return nil
}

func (m *MockRecipeRepoForInstruction) FindByUserID(uint) ([]models.Recipe, error) {
	return nil, nil
}

func (m *MockRecipeRepoForInstruction) FindByIDWithDetails(uint) (*models.Recipe, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *MockRecipeRepoForInstruction) Update(*models.Recipe) error {
	return nil
}

func (m *MockRecipeRepoForInstruction) Delete(*models.Recipe) error {
	return nil
}

func TestAddInstruction_Unauthorized(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{},
		&MockRecipeRepoForInstruction{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 2,
				}, nil
			},
		},
	)

	err := service.AddInstruction(
		1,
		1,
		dto.CreateInstructionRequest{
			StepNumber: 1,
			Text:       "Heat oil",
		},
	)

	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestAddInstruction_Success(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{
			CreateFn: func(i *models.Instruction) error {
				if i.Text == "" {
					t.Fatal("instruction text should not be empty")
				}
				return nil
			},
		},
		&MockRecipeRepoForInstruction{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 1,
				}, nil
			},
		},
	)

	err := service.AddInstruction(
		1,
		1,
		dto.CreateInstructionRequest{
			StepNumber: 1,
			Text:       "Heat oil",
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestGetInstructions_Success(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{
			FindByRecipeIDFn: func(recipeID uint) ([]models.Instruction, error) {
				return []models.Instruction{
					{ID: 1, StepNumber: 1, Text: "Heat oil"},
					{ID: 2, StepNumber: 2, Text: "Add onions"},
				}, nil
			},
		},
		&MockRecipeRepoForInstruction{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 1,
				}, nil
			},
		},
	)

	items, err := service.GetInstructions(1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 instructions, got %d", len(items))
	}
}

func TestUpdateInstruction_Unauthorized(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{
			FindByIDFn: func(id uint) (*models.Instruction, error) {
				return &models.Instruction{
					ID:       id,
					RecipeID: 1,
				}, nil
			},
		},
		&MockRecipeRepoForInstruction{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 2,
				}, nil
			},
		},
	)

	err := service.UpdateInstruction(
		1,
		1,
		dto.UpdateInstructionRequest{
			Text: "Updated step",
		},
	)

	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestUpdateInstruction_Success(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{
			FindByIDFn: func(id uint) (*models.Instruction, error) {
				return &models.Instruction{
					ID:       id,
					RecipeID: 1,
					Text:     "Old",
				}, nil
			},
			UpdateFn: func(i *models.Instruction) error {
				if i.Text != "Updated" {
					t.Fatal("instruction text not updated")
				}
				return nil
			},
		},
		&MockRecipeRepoForInstruction{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 1,
				}, nil
			},
		},
	)

	err := service.UpdateInstruction(
		1,
		1,
		dto.UpdateInstructionRequest{
			Text: "Updated",
		},
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteInstruction_Success(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{
			FindByIDFn: func(id uint) (*models.Instruction, error) {
				return &models.Instruction{
					ID:       id,
					RecipeID: 1,
				}, nil
			},
			DeleteFn: func(i *models.Instruction) error {
				return nil
			},
		},
		&MockRecipeRepoForInstruction{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 1,
				}, nil
			},
		},
	)

	err := service.DeleteInstruction(1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
