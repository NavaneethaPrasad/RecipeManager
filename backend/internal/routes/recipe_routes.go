package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRecipeRoutes(r *gin.RouterGroup, db *gorm.DB) {

	recipeRepo := repository.NewRecipeRepository(db)
	recipeService := services.NewRecipeService(recipeRepo, db)
	recipeHandler := handlers.NewRecipeHandler(recipeService)

	scaleService := services.NewRecipeScaleService(recipeRepo)
	scaleHandler := handlers.NewRecipeScaleHandler(scaleService)

	instRepo := repository.NewInstructionRepository(db)
	instService := services.NewInstructionService(instRepo, recipeRepo)
	instHandler := handlers.NewInstructionHandler(instService)

	recipes := r.Group("/recipes")
	{
		recipes.POST("", recipeHandler.CreateRecipe)
		recipes.GET("", recipeHandler.GetMyRecipes)
		recipes.GET("/:id", recipeHandler.GetRecipeByID)
		recipes.PUT("/:id", recipeHandler.UpdateRecipe)
		recipes.DELETE("/:id", recipeHandler.DeleteRecipe)

		recipes.GET("/:id/scale", scaleHandler.ScaleRecipe)

		recipes.POST("/:id/instructions", instHandler.AddInstruction)
		recipes.GET("/:id/instructions", instHandler.GetInstructions)
		recipes.PUT("/instructions/:id", instHandler.UpdateInstruction)
		recipes.DELETE("/instructions/:id", instHandler.DeleteInstruction)
	}
}
