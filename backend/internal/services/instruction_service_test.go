package services

import (
	"testing"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

// --- MOCK INSTRUCTION REPO ---
type MockInstructionRepository struct {
	CreateFn         func(*models.Instruction) error
	FindByRecipeIDFn func(uint) ([]models.Instruction, error)
	FindByIDFn       func(uint) (*models.Instruction, error)
	UpdateFn         func(*models.Instruction) error
	DeleteFn         func(uint) error // FIX: Changed from *models.Instruction to uint
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

// FIX: Match the interface signature (Delete by ID)
func (m *MockInstructionRepository) Delete(id uint) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}

// --- MOCK RECIPE REPO ---
type MockRecipeRepoForInstruction struct {
	FindByIDFn func(uint) (*models.Recipe, error)
}

func (m *MockRecipeRepoForInstruction) FindByID(id uint) (*models.Recipe, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(id)
	}
	return nil, gorm.ErrRecordNotFound
}

// Stub other methods required by the interface but not used in these tests
func (m *MockRecipeRepoForInstruction) Create(*models.Recipe) error { return nil }
func (m *MockRecipeRepoForInstruction) FindByUserID(uint) ([]models.Recipe, error) {
	return nil, nil
}
func (m *MockRecipeRepoForInstruction) FindByIDWithDetails(uint) (*models.Recipe, error) {
	return nil, nil
}
func (m *MockRecipeRepoForInstruction) Update(*models.Recipe) error { return nil }
func (m *MockRecipeRepoForInstruction) Delete(*models.Recipe) error { return nil }

// --- TESTS ---

func TestAddInstruction_Unauthorized(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{},
		&MockRecipeRepoForInstruction{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 2, // Different user
				}, nil
			},
		},
	)

	err := service.AddInstruction(
		1,
		1, // Current User
		dto.CreateInstructionRequest{
			StepNumber: 1,
			Text:       "Heat oil",
		},
	)

	// Note: Use ErrUnauthorizedInstruction if you defined that specific variable,
	// otherwise use ErrUnauthorized from recipe_service.go
	if err != ErrUnauthorizedInstruction && err != ErrUnauthorized {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}

func TestAddInstruction_Success(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{
			CreateFn: func(i *models.Instruction) error {
				if i.Text == "" {
					t.Fatal("instruction text should not be empty")
				}
				// Verify DTO data mapped to Model
				if i.StepNumber != 1 {
					t.Fatalf("expected step 1, got %d", i.StepNumber)
				}
				return nil
			},
		},
		&MockRecipeRepoForInstruction{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 1, // Same User
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
					RecipeID: 100,
				}, nil
			},
		},
		&MockRecipeRepoForInstruction{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 2, // Different Owner
				}, nil
			},
		},
	)

	err := service.UpdateInstruction(
		1,
		1, // Current User
		dto.UpdateInstructionRequest{
			Text: "Updated step",
		},
	)

	if err != ErrUnauthorizedInstruction && err != ErrUnauthorized {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}

func TestDeleteInstruction_Success(t *testing.T) {
	service := NewInstructionService(
		&MockInstructionRepository{
			FindByIDFn: func(id uint) (*models.Instruction, error) {
				return &models.Instruction{
					ID:       id,
					RecipeID: 100,
				}, nil
			},
			// FIX: Changed signature to accept ID (uint)
			DeleteFn: func(id uint) error {
				if id != 1 {
					t.Fatalf("expected delete ID 1, got %d", id)
				}
				return nil
			},
		},
		&MockRecipeRepoForInstruction{
			FindByIDFn: func(id uint) (*models.Recipe, error) {
				return &models.Recipe{
					ID:     id,
					UserID: 1, // Correct Owner
				}, nil
			},
		},
	)

	err := service.DeleteInstruction(1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
