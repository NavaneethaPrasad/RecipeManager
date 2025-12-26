package services

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
)

type InstructionService interface {
	AddInstruction(recipeID uint, userID uint, req dto.CreateInstructionRequest) error
	GetInstructions(recipeID uint, userID uint) ([]dto.InstructionResponse, error)
	UpdateInstruction(instructionID uint, userID uint, req dto.UpdateInstructionRequest) error
	DeleteInstruction(instructionID uint, userID uint) error
}

type instructionService struct {
	InstructionRepo repository.InstructionRepository
	RecipeRepo      repository.RecipeRepository
}

func NewInstructionService(
	instructionRepo repository.InstructionRepository,
	recipeRepo repository.RecipeRepository,
) InstructionService {
	return &instructionService{
		InstructionRepo: instructionRepo,
		RecipeRepo:      recipeRepo,
	}
}

func (s *instructionService) AddInstruction(recipeID uint, userID uint, req dto.CreateInstructionRequest) error {
	recipe, err := s.RecipeRepo.FindByID(recipeID)
	if err != nil {
		return err
	}

	if recipe.UserID != userID {
		return ErrUnauthorized
	}

	instruction := &models.Instruction{
		RecipeID:   recipeID,
		StepNumber: req.StepNumber,
		Text:       req.Text,
	}

	return s.InstructionRepo.Create(instruction)
}

func (s *instructionService) GetInstructions(recipeID uint, userID uint) ([]dto.InstructionResponse, error) {
	recipe, err := s.RecipeRepo.FindByID(recipeID)
	if err != nil {
		return nil, err
	}

	if recipe.UserID != userID {
		return nil, ErrUnauthorized
	}

	instructions, err := s.InstructionRepo.FindByRecipeID(recipeID)
	if err != nil {
		return nil, err
	}

	var response []dto.InstructionResponse
	for _, ins := range instructions {
		response = append(response, dto.InstructionResponse{
			ID:         ins.ID,
			StepNumber: ins.StepNumber,
			Text:       ins.Text,
		})
	}

	return response, nil
}

func (s *instructionService) UpdateInstruction(instructionID uint, userID uint, req dto.UpdateInstructionRequest) error {
	instruction, err := s.InstructionRepo.FindByID(instructionID)
	if err != nil {
		return err
	}

	recipe, err := s.RecipeRepo.FindByID(instruction.RecipeID)
	if err != nil {
		return err
	}

	if recipe.UserID != userID {
		return ErrUnauthorized
	}

	instruction.StepNumber = req.StepNumber
	instruction.Text = req.Text

	return s.InstructionRepo.Update(instruction)
}

func (s *instructionService) DeleteInstruction(instructionID uint, userID uint) error {
	instruction, err := s.InstructionRepo.FindByID(instructionID)
	if err != nil {
		return err
	}

	recipe, err := s.RecipeRepo.FindByID(instruction.RecipeID)
	if err != nil {
		return err
	}

	if recipe.UserID != userID {
		return ErrUnauthorized
	}

	return s.InstructionRepo.Delete(instruction)
}
