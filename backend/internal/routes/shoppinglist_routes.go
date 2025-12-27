package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/middleware"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterShoppingListRoutes(r *gin.Engine, db *gorm.DB) {

	shoppingListRepo := repository.NewShoppingListRepository(db)
	mealPlanRepo := repository.NewMealPlanRepository(db)
	recipeIngredientRepo := repository.NewRecipeIngredientRepository(db)

	shoppingListService := services.NewShoppingListService(
		mealPlanRepo,
		recipeIngredientRepo,
		shoppingListRepo,
	)

	shoppingListHandler := handlers.NewShoppingListHandler(shoppingListService)

	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("/shopping-lists", shoppingListHandler.GenerateShoppingList)
		protected.GET("/shopping-lists/:id", shoppingListHandler.GetShoppingList)
		protected.PATCH("/shopping-lists/items/:id/toggle", shoppingListHandler.ToggleItem)
	}
}
