package handlers

import (
	"net/http"
	"strconv"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type RecipeScaleHandler struct {
	Service services.RecipeScaleService
}

func NewRecipeScaleHandler(service services.RecipeScaleService) *RecipeScaleHandler {
	return &RecipeScaleHandler{Service: service}
}

func (h *RecipeScaleHandler) ScaleRecipe(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	recipeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recipe id"})
		return
	}

	servings, err := strconv.Atoi(c.Query("servings"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid servings"})
		return
	}

	resp, err := h.Service.ScaleRecipe(uint(recipeID), userID, servings)
	if err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
