package handlers

import (
	"net/http"
	"strconv"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	Service services.IngredientService
}

func NewIngredientHandler(service services.IngredientService) *IngredientHandler {
	return &IngredientHandler{Service: service}
}

func (h *IngredientHandler) CreateIngredient(c *gin.Context) {
	var req dto.CreateIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.CreateIngredient(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *IngredientHandler) GetIngredients(c *gin.Context) {
	items, err := h.Service.GetIngredients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *IngredientHandler) AddIngredientToRecipe(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	recipeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req dto.AddRecipeIngredientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.AddIngredientToRecipe(uint(recipeID), userID, req); err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *IngredientHandler) GetRecipeIngredients(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	recipeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe id"})
		return
	}

	items, err := h.Service.GetRecipeIngredients(uint(recipeID), userID)
	if err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *IngredientHandler) RemoveRecipeIngredient(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ingredient id"})
		return
	}

	if err := h.Service.RemoveRecipeIngredient(uint(id), userID); err != nil {
		if err == services.ErrIngredientUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
