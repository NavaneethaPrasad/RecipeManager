package services

import (
	"os"
	"testing"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 1. Define the Mock Repo
type MockUserRepo struct {
	FindByEmailFn func(email string) (*models.User, error)
	CreateFn      func(user *models.User) error
}

// Implement the interface methods so MockUserRepo satisfies repository.UserRepository
func (m *MockUserRepo) FindByEmail(email string) (*models.User, error) {
	return m.FindByEmailFn(email)
}

func (m *MockUserRepo) Create(user *models.User) error {
	return m.CreateFn(user)
}

func TestRegister_Success(t *testing.T) {
	// Setup Mock
	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			// Return "Record Not Found" so registration can proceed
			return nil, gorm.ErrRecordNotFound
		},
		CreateFn: func(user *models.User) error {
			// Simulate successful creation
			user.ID = 1
			return nil
		},
	}

	// FIX 1: NewAuthService only takes 1 argument (repo)
	authService := NewAuthService(mockRepo)

	req := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// FIX 2: Method is named 'Register', not 'RegisterUser'
	err := authService.Register(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			// Return a user to simulate "Email Exists"
			return &models.User{}, nil
		},
		CreateFn: func(user *models.User) error { return nil },
	}

	authService := NewAuthService(mockRepo)

	req := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	err := authService.Register(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	// Setup Env for Utils (if your utils.GenerateToken uses env vars)
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	// Create a real hashed password for comparison
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

	// FIX 3: Capture 3 return values (token, user, error)
	token, user, err := authService.Login(req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected JWT token, got empty string")
	}
	if user == nil {
		t.Fatal("expected user data, got nil")
	}
	if user.Email != req.Email {
		t.Errorf("expected email %s, got %s", req.Email, user.Email)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	hashedPassword := "$2a$10$SomeWrongHashValue..."

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

	// Expect an error
	_, _, err := authService.Login(req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
