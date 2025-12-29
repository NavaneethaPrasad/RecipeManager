package handlers

import (
	"net/http"
	"strconv"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type RecipeHandler struct {
	Service services.RecipeService
}

func NewRecipeHandler(service services.RecipeService) *RecipeHandler {
	return &RecipeHandler{Service: service}
}

func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req dto.CreateRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipeID, err := h.Service.CreateRecipe(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": recipeID,
	})
}

func (h *RecipeHandler) GetMyRecipes(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	recipes, err := h.Service.GetMyRecipes(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func (h *RecipeHandler) UpdateRecipe(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	recipeIDParam := c.Param("id")
	recipeID, err := strconv.ParseUint(recipeIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe id"})
		return
	}

	var req dto.UpdateRecipeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.UpdateRecipe(uint(recipeID), userID, req)
	if err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	recipeIDParam := c.Param("id")
	recipeID, err := strconv.ParseUint(recipeIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe id"})
		return
	}

	err = h.Service.DeleteRecipe(uint(recipeID), userID)
	if err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *RecipeHandler) GetRecipeByID(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	recipeIDParam := c.Param("id")
	recipeID, err := strconv.ParseUint(recipeIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe id"})
		return
	}

	recipe, err := h.Service.GetRecipeByID(uint(recipeID), userID)
	if err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}
