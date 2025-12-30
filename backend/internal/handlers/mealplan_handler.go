package handlers

import (
	"net/http"
	"strconv"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type MealPlanHandler struct {
	Service services.MealPlanService
}

func NewMealPlanHandler(service services.MealPlanService) *MealPlanHandler {
	return &MealPlanHandler{Service: service}
}

func (h *MealPlanHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req dto.CreateMealPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.Create(userID, req); err != nil {
		if err == services.ErrMealExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *MealPlanHandler) GetByDate(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate != "" && endDate != "" {
		items, err := h.Service.GetByDateRange(userID, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
		return
	}

	date := c.Query("date")
	items, err := h.Service.GetByDate(userID, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *MealPlanHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req dto.UpdateMealPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.Update(uint(id), userID, req); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *MealPlanHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.Service.Delete(uint(id), userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
