package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/middleware"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRecipeRoutes(r *gin.Engine, db *gorm.DB) {

	recipeRepo := repository.NewRecipeRepository(db)
	recipeService := services.NewRecipeService(recipeRepo)
	recipeHandler := handlers.NewRecipeHandler(recipeService)
	scaleService := services.NewRecipeScaleService(recipeRepo)
	scaleHandler := handlers.NewRecipeScaleHandler(scaleService)

	protected := r.Group("/api/recipes")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("", recipeHandler.CreateRecipe)
		protected.GET("", recipeHandler.GetMyRecipes)

		protected.GET("/:id", recipeHandler.GetRecipeByID)
		protected.PUT("/:id", recipeHandler.UpdateRecipe)
		protected.DELETE("/:id", recipeHandler.DeleteRecipe)
		protected.GET("/:id/scale", scaleHandler.ScaleRecipe)
	}
}
