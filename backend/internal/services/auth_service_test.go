package services

import (
	"errors"
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
			user.ID = 1
			return nil
		},
	}
	authService := NewAuthService(mockRepo)
	req := dto.RegisterRequest{Name: "Test", Email: "test@ex.com", Password: "password123"}

	err := authService.Register(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			return &models.User{Email: email}, nil
		},
	}
	authService := NewAuthService(mockRepo)
	req := dto.RegisterRequest{Email: "exists@ex.com", Password: "password123"}

	err := authService.Register(req)
	if err == nil || err.Error() != "email already exists" {
		t.Fatal("expected 'email already exists' error")
	}
}

func TestRegister_DatabaseLookupError(t *testing.T) {
	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			return nil, errors.New("database connection failed")
		},
	}
	authService := NewAuthService(mockRepo)
	err := authService.Register(dto.RegisterRequest{Email: "test@ex.com"})

	if err == nil || err.Error() != "database connection failed" {
		t.Fatal("expected database error to propagate")
	}
}

func TestRegister_CreateError(t *testing.T) {
	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
		CreateFn: func(user *models.User) error {
			return errors.New("failed to save user")
		},
	}
	authService := NewAuthService(mockRepo)
	err := authService.Register(dto.RegisterRequest{Name: "T", Email: "t@e.com", Password: "pass"})

	if err == nil || err.Error() != "failed to save user" {
		t.Fatal("expected create error to propagate")
	}
}

func TestLogin_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			return &models.User{ID: 1, Email: email, Password: string(hash)}, nil
		},
	}
	authService := NewAuthService(mockRepo)

	token, user, err := authService.Login(dto.LoginRequest{Email: "test@e.com", Password: "password123"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" || user == nil {
		t.Fatal("expected token and user object")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	authService := NewAuthService(mockRepo)
	_, _, err := authService.Login(dto.LoginRequest{Email: "notfound@e.com", Password: "any"})

	if err == nil || err.Error() != "invalid email or password" {
		t.Fatal("expected 'invalid email or password' error")
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct_pass"), bcrypt.DefaultCost)
	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			return &models.User{Password: string(hash)}, nil
		},
	}
	authService := NewAuthService(mockRepo)
	_, _, err := authService.Login(dto.LoginRequest{Email: "a@b.com", Password: "wrong_password"})

	if err == nil || err.Error() != "invalid email or password" {
		t.Fatal("expected error for wrong password")
	}
}

func TestLogin_JWTSecretMissing(t *testing.T) {
	os.Setenv("JWT_SECRET", "")
	defer os.Setenv("JWT_SECRET", "test_secret")

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	mockRepo := &MockUserRepo{
		FindByEmailFn: func(email string) (*models.User, error) {
			return &models.User{ID: 1, Password: string(hash)}, nil
		},
	}
	authService := NewAuthService(mockRepo)

	_, _, err := authService.Login(dto.LoginRequest{Email: "a@b.com", Password: "pass"})

	if err == nil || err.Error() != "failed to generate token" {
		t.Fatalf("expected 'failed to generate token' error, got: %v", err)
	}
}
