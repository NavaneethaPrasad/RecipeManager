package handlers

import (
	"net/http"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Register(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set JWT in HTTP-only cookie
	c.SetCookie(
		"auth_token",
		token,
		3600*24, // 1 day
		"/",
		"",
		false, // true in production (HTTPS)
		true,  // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}
