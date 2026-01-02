package services

import (
	"errors"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"gorm.io/gorm"
)

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
	recipe, err := s.RecipeRepo.FindByID(recipeID)
	if err != nil {
		return err
	}

	if recipe.UserID != userID {
		return ErrUnauthorizedInstruction
	}

	instruction := &models.Instruction{
		RecipeID:   recipeID,
		StepNumber: req.StepNumber,
		Text:       req.Text,
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
	ins, err := s.InstructionRepo.FindByID(instructionID)
	if err != nil {
		return err
	}

	recipe, err := s.RecipeRepo.FindByID(ins.RecipeID)
	if err != nil {
		return err
	}
	if recipe.UserID != userID {
		return ErrUnauthorizedInstruction
	}

	ins.StepNumber = req.StepNumber
	ins.Text = req.Text

	return s.InstructionRepo.Update(ins)
}

func (s *instructionService) DeleteInstruction(instructionID uint, userID uint) error {
	ins, err := s.InstructionRepo.FindByID(instructionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	recipe, err := s.RecipeRepo.FindByID(ins.RecipeID)
	if err != nil {
		return err
	}
	if recipe.UserID != userID {
		return ErrUnauthorizedInstruction
	}

	return s.InstructionRepo.Delete(ins.ID)
}
