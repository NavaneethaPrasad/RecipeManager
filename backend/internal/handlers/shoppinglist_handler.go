package handlers

import (
	"net/http"
	"strconv"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ShoppingListHandler struct {
	Service services.ShoppingListService
}

func NewShoppingListHandler(service services.ShoppingListService) *ShoppingListHandler {
	return &ShoppingListHandler{Service: service}
}

func (h *ShoppingListHandler) GenerateShoppingList(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req dto.GenerateShoppingListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list, err := h.Service.Generate(
		userID,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, list)
}

func (h *ShoppingListHandler) GetShoppingList(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	listID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shopping list id"})
		return
	}

	list, err := h.Service.GetShoppingListByID(uint(listID), userID)

	if err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *ShoppingListHandler) ToggleItem(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	if err := h.Service.ToggleItemChecked(uint(itemID), userID); err != nil {

		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
