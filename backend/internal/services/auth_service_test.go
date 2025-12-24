package services

import (
	"os"
	"testing"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type MockUserRepo struct {
	FindByEmailFn func(email string) (*models.User, error)
	CreateFn      func(user *models.User) error
}

func (m *MockUserRepo) FindByEmail(email string) (*models.User, error) {
	return m.FindByEmailFn(email)
}

func (m *MockUserRepo) Create(user *models.User) error {
	return m.CreateFn(user)
}

func TestRegister_Success(t *testing.T) {

	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		CreateFn: func(user *models.User) error {
			return nil
		},
	}

	authService := NewAuthService(mockRepo)

	req := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	err := authService.Register(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {

	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			return &models.User{}, nil
		},
		CreateFn: func(user *models.User) error {
			return nil
		},
	}

	authService := NewAuthService(mockRepo)

	req := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	err := authService.Register(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogin_Success(t *testing.T) {

	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	hashedPasswordBytes, _ := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)

	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			return &models.User{
				ID:       1,
				Email:    email,
				Password: string(hashedPasswordBytes),
			}, nil
		},
	}

	authService := NewAuthService(mockRepo)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	token, err := authService.Login(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("expected JWT token, got empty string")
	}
}

func TestLogin_InvalidPassword(t *testing.T) {

	hashedPassword := "$2a$10$e0MYzXyjpJS7Pd0RVvHwHeFx8LxH0VZ3yN1YxkQ0Y5qZ5bM7bP5d6"

	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			return &models.User{
				ID:       1,
				Email:    email,
				Password: hashedPassword,
			}, nil
		},
	}

	authService := NewAuthService(mockRepo)

	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err := authService.Login(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
