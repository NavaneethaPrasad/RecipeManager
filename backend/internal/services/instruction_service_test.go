package services

import (
	"errors"
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
	DeleteFn         func(uint) error
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
func (m *MockInstructionRepository) Delete(id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

type MockRecipeRepoForInstruction struct {
	FindByIDFn func(uint) (*models.Recipe, error)
}

func (m *MockRecipeRepoForInstruction) FindByID(id uint) (*models.Recipe, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *MockRecipeRepoForInstruction) Create(*models.Recipe) error                { return nil }
func (m *MockRecipeRepoForInstruction) FindByUserID(uint) ([]models.Recipe, error) { return nil, nil }
func (m *MockRecipeRepoForInstruction) FindByIDWithDetails(uint) (*models.Recipe, error) {
	return nil, nil
}
func (m *MockRecipeRepoForInstruction) Update(*models.Recipe) error { return nil }
func (m *MockRecipeRepoForInstruction) Delete(*models.Recipe) error { return nil }

func TestAddInstruction_Success(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{CreateFn: func(i *models.Instruction) error { return nil }},
		&MockRecipeRepoForInstruction{FindByIDFn: func(id uint) (*models.Recipe, error) {
			return &models.Recipe{ID: id, UserID: 1}, nil
		}},
	)
	err := service.AddInstruction(1, 1, dto.CreateInstructionRequest{StepNumber: 1, Text: "Test"})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAddInstruction_RecipeNotFound(t *testing.T) {
	service := NewInstructionService(&MockInstructionRepository{}, &MockRecipeRepoForInstruction{
		FindByIDFn: func(id uint) (*models.Recipe, error) { return nil, errors.New("db error") },
	})
	err := service.AddInstruction(1, 1, dto.CreateInstructionRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAddInstruction_Unauthorized(t *testing.T) {
	service := NewInstructionService(&MockInstructionRepository{}, &MockRecipeRepoForInstruction{
		FindByIDFn: func(id uint) (*models.Recipe, error) { return &models.Recipe{UserID: 2}, nil },
	})
	err := service.AddInstruction(1, 1, dto.CreateInstructionRequest{})
	if err != ErrUnauthorizedInstruction && err != ErrUnauthorized {
		t.Fatal("expected unauthorized")
	}
}

func TestGetInstructions_Success(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{FindByRecipeIDFn: func(u uint) ([]models.Instruction, error) {
			return []models.Instruction{{Text: "Step 1"}}, nil
		}},
		&MockRecipeRepoForInstruction{FindByIDFn: func(id uint) (*models.Recipe, error) {
			return &models.Recipe{UserID: 1}, nil
		}},
	)
	res, err := service.GetInstructions(1, 1)
	if err != nil || len(res) != 1 {
		t.Fatal("failed to get instructions")
	}
}

func TestGetInstructions_RecipeError(t *testing.T) {
	service := NewInstructionService(&MockInstructionRepository{}, &MockRecipeRepoForInstruction{
		FindByIDFn: func(id uint) (*models.Recipe, error) { return nil, errors.New("err") },
	})
	_, err := service.GetInstructions(1, 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetInstructions_Unauthorized(t *testing.T) {
	service := NewInstructionService(&MockInstructionRepository{}, &MockRecipeRepoForInstruction{
		FindByIDFn: func(id uint) (*models.Recipe, error) { return &models.Recipe{UserID: 2}, nil },
	})
	_, err := service.GetInstructions(1, 1)
	if err != ErrUnauthorizedInstruction && err != ErrUnauthorized {
		t.Fatal("expected unauthorized")
	}
}

func TestUpdateInstruction_Success(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{
			FindByIDFn: func(u uint) (*models.Instruction, error) { return &models.Instruction{RecipeID: 1}, nil },
			UpdateFn:   func(i *models.Instruction) error { return nil },
		},
		&MockRecipeRepoForInstruction{FindByIDFn: func(id uint) (*models.Recipe, error) {
			return &models.Recipe{UserID: 1}, nil
		}},
	)
	err := service.UpdateInstruction(1, 1, dto.UpdateInstructionRequest{Text: "Update"})
	if err != nil {
		t.Fatal("failed update")
	}
}

func TestUpdateInstruction_InstructionNotFound(t *testing.T) {
	service := NewInstructionService(&MockInstructionRepository{
		FindByIDFn: func(u uint) (*models.Instruction, error) { return nil, errors.New("not found") },
	}, &MockRecipeRepoForInstruction{})
	err := service.UpdateInstruction(1, 1, dto.UpdateInstructionRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateInstruction_RecipeError(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{FindByIDFn: func(u uint) (*models.Instruction, error) { return &models.Instruction{RecipeID: 1}, nil }},
		&MockRecipeRepoForInstruction{FindByIDFn: func(id uint) (*models.Recipe, error) { return nil, errors.New("err") }},
	)
	err := service.UpdateInstruction(1, 1, dto.UpdateInstructionRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDeleteInstruction_Success(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{
			FindByIDFn: func(u uint) (*models.Instruction, error) { return &models.Instruction{RecipeID: 1}, nil },
			DeleteFn:   func(u uint) error { return nil },
		},
		&MockRecipeRepoForInstruction{FindByIDFn: func(id uint) (*models.Recipe, error) {
			return &models.Recipe{UserID: 1}, nil
		}},
	)
	err := service.DeleteInstruction(1, 1)
	if err != nil {
		t.Fatal("failed delete")
	}
}

func TestDeleteInstruction_RecordNotFound(t *testing.T) {
	service := NewInstructionService(&MockInstructionRepository{
		FindByIDFn: func(u uint) (*models.Instruction, error) { return nil, gorm.ErrRecordNotFound },
	}, &MockRecipeRepoForInstruction{})
	err := service.DeleteInstruction(1, 1)
	if err != nil {
		t.Fatal("expected nil for record not found")
	}
}

func TestDeleteInstruction_GeneralError(t *testing.T) {
	service := NewInstructionService(&MockInstructionRepository{
		FindByIDFn: func(u uint) (*models.Instruction, error) { return nil, errors.New("db error") },
	}, &MockRecipeRepoForInstruction{})
	err := service.DeleteInstruction(1, 1)
	if err == nil || err.Error() != "db error" {
		t.Fatal("expected error to propagate")
	}
}

func TestDeleteInstruction_RecipeError(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{FindByIDFn: func(u uint) (*models.Instruction, error) { return &models.Instruction{RecipeID: 1}, nil }},
		&MockRecipeRepoForInstruction{FindByIDFn: func(id uint) (*models.Recipe, error) { return nil, errors.New("err") }},
	)
	err := service.DeleteInstruction(1, 1)
	if err == nil {
		t.Fatal("expected error")
	}
}
