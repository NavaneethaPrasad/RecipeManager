package utils

import (
	"os"
	"testing"
)

func TestGenerateAndValidateToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	token, err := GenerateToken(10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}

	if claims.UserID != 10 {
		t.Fatalf("expected userID 10, got %d", claims.UserID)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	_, err := ValidateToken("invalid.token.value")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestValidateToken_EmptySecret(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	_, err := ValidateToken("some.token.value")
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is not set")
	}
}
