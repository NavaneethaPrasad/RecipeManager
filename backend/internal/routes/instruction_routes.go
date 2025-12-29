package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterInstructionRoutes(r *gin.RouterGroup, db *gorm.DB) {
	instructionRepo := repository.NewInstructionRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)

	service := services.NewInstructionService(instructionRepo, recipeRepo)
	handler := handlers.NewInstructionHandler(service)

	instructions := r.Group("/recipes/:id/instructions")
	{
		instructions.POST("", handler.AddInstruction)
		instructions.GET("", handler.GetInstructions)
		instructions.PUT("/:instructionId", handler.UpdateInstruction)
		instructions.DELETE("/:instructionId", handler.DeleteInstruction)
	}
}
