package handlers

import (
	"net/http"
	"strconv"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type InstructionHandler struct {
	Service services.InstructionService
}

func NewInstructionHandler(service services.InstructionService) *InstructionHandler {
	return &InstructionHandler{Service: service}
}

func (h *InstructionHandler) AddInstruction(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	recipeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req dto.CreateInstructionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.AddInstruction(uint(recipeID), userID, req); err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *InstructionHandler) GetInstructions(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	recipeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	instructions, err := h.Service.GetInstructions(uint(recipeID), userID)
	if err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, instructions)
}

func (h *InstructionHandler) UpdateInstruction(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	instructionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req dto.UpdateInstructionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.UpdateInstruction(uint(instructionID), userID, req); err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *InstructionHandler) DeleteInstruction(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	instructionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.Service.DeleteInstruction(uint(instructionID), userID); err != nil {
		if err == services.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
