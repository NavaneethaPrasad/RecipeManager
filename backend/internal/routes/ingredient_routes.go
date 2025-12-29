package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterIngredientRoutes(r *gin.RouterGroup, db *gorm.DB) {

	ingredientRepo := repository.NewIngredientRepository(db)
	recipeIngredientRepo := repository.NewRecipeIngredientRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)

	ingredientService := services.NewIngredientService(
		ingredientRepo,
		recipeIngredientRepo,
		recipeRepo,
	)

	ingredientHandler := handlers.NewIngredientHandler(ingredientService)

	ingredients := r.Group("/ingredients")
	{
		ingredients.POST("", ingredientHandler.CreateIngredient)
		ingredients.GET("", ingredientHandler.GetIngredients)

		ingredients.POST("/recipes/:id/ingredients", ingredientHandler.AddIngredientToRecipe)
		ingredients.GET("/recipes/:id/ingredients", ingredientHandler.GetRecipeIngredients)
		ingredients.DELETE("/recipe-ingredients/:id", ingredientHandler.RemoveRecipeIngredient)
	}
}
