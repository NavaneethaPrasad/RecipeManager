package services

import (
	"errors"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"gorm.io/gorm"
)

// Define error here (or import from recipe_service if shared)
var ErrUnauthorizedInstruction = errors.New("not authorized")

type InstructionService interface {
	AddInstruction(recipeID uint, userID uint, req dto.CreateInstructionRequest) error
	GetInstructions(recipeID uint, userID uint) ([]models.Instruction, error)
	UpdateInstruction(instructionID uint, userID uint, req dto.UpdateInstructionRequest) error
	DeleteInstruction(instructionID uint, userID uint) error
}

type instructionService struct {
	InstructionRepo repository.InstructionRepository
	RecipeRepo      repository.RecipeRepository
}

func NewInstructionService(instRepo repository.InstructionRepository, recipeRepo repository.RecipeRepository) InstructionService {
	return &instructionService{
		InstructionRepo: instRepo,
		RecipeRepo:      recipeRepo,
	}
}

func (s *instructionService) AddInstruction(recipeID uint, userID uint, req dto.CreateInstructionRequest) error {
	// 1. Check if recipe exists
	recipe, err := s.RecipeRepo.FindByID(recipeID)
	if err != nil {
		return err
	}

	// 2. Check Ownership
	if recipe.UserID != userID {
		return ErrUnauthorizedInstruction
	}

	// 3. Create Instruction using data from DTO
	instruction := &models.Instruction{
		RecipeID:   recipeID,
		StepNumber: req.StepNumber, // Using the StepNumber from your DTO
		Text:       req.Text,       // Using the Text from your DTO
	}

	return s.InstructionRepo.Create(instruction)
}

func (s *instructionService) GetInstructions(recipeID uint, userID uint) ([]models.Instruction, error) {
	recipe, err := s.RecipeRepo.FindByID(recipeID)
	if err != nil {
		return nil, err
	}

	if recipe.UserID != userID {
		return nil, ErrUnauthorizedInstruction
	}

	return s.InstructionRepo.FindByRecipeID(recipeID)
}

func (s *instructionService) UpdateInstruction(instructionID uint, userID uint, req dto.UpdateInstructionRequest) error {
	// 1. Get Instruction
	ins, err := s.InstructionRepo.FindByID(instructionID)
	if err != nil {
		return err
	}

	// 2. Check Ownership via Recipe
	recipe, err := s.RecipeRepo.FindByID(ins.RecipeID)
	if err != nil {
		return err
	}
	if recipe.UserID != userID {
		return ErrUnauthorizedInstruction
	}

	// 3. Update fields from DTO
	ins.StepNumber = req.StepNumber
	ins.Text = req.Text

	return s.InstructionRepo.Update(ins)
}

func (s *instructionService) DeleteInstruction(instructionID uint, userID uint) error {
	// 1. Get Instruction
	ins, err := s.InstructionRepo.FindByID(instructionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	// 2. Check Ownership
	recipe, err := s.RecipeRepo.FindByID(ins.RecipeID)
	if err != nil {
		return err
	}
	if recipe.UserID != userID {
		return ErrUnauthorizedInstruction
	}

	return s.InstructionRepo.Delete(ins.ID)
}
