package services

import (
	"errors"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(req dto.RegisterRequest) error
	Login(req dto.LoginRequest) (string, *dto.UserResponse, error)
}

type authService struct {
	UserRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{UserRepo: repo}
}

func (s *authService) Register(req dto.RegisterRequest) error {
	// Check if email already exists
	if _, err := s.UserRepo.FindByEmail(req.Email); err == nil {
		return errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	return s.UserRepo.Create(user)
}

func (s *authService) Login(req dto.LoginRequest) (string, *dto.UserResponse, error) {
	user, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}
	// Generate Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	// Prepare User Data for Frontend
	userResponse := &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	// Return Token AND User Data
	return token, userResponse, nil
}
