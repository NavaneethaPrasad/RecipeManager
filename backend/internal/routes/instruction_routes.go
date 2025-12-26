package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/middleware"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterInstructionRoutes(r *gin.Engine, db *gorm.DB) {

	instructionRepo := repository.NewInstructionRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)

	instructionService := services.NewInstructionService(instructionRepo, recipeRepo)
	instructionHandler := handlers.NewInstructionHandler(instructionService)

	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("/recipes/:id/instructions", instructionHandler.AddInstruction)
		protected.GET("/recipes/:id/instructions", instructionHandler.GetInstructions)
		protected.PUT("/recipes/:id/instructions/:instructionId", instructionHandler.UpdateInstruction)
		protected.DELETE("/recipes/:id/instructions/:instructionId", instructionHandler.DeleteInstruction)

	}
}
