package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/middleware"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterIngredientRoutes(r *gin.Engine, db *gorm.DB) {

	ingredientRepo := repository.NewIngredientRepository(db)
	recipeIngredientRepo := repository.NewRecipeIngredientRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)

	ingredientService := services.NewIngredientService(
		ingredientRepo,
		recipeIngredientRepo,
		recipeRepo,
	)

	ingredientHandler := handlers.NewIngredientHandler(ingredientService)

	api := r.Group("/api")
	api.Use(middleware.JWTAuthMiddleware())
	{
		api.POST("/ingredients", ingredientHandler.CreateIngredient)
		api.GET("/ingredients", ingredientHandler.GetIngredients)

		api.POST("/recipes/:id/ingredients", ingredientHandler.AddIngredientToRecipe)
		api.GET("/recipes/:id/ingredients", ingredientHandler.GetRecipeIngredients)
		api.DELETE("/recipe-ingredients/:id", ingredientHandler.RemoveRecipeIngredient)
	}
}
