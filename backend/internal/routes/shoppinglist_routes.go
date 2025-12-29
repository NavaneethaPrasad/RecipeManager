package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterShoppingListRoutes(r *gin.RouterGroup, db *gorm.DB) {

	shoppingListRepo := repository.NewShoppingListRepository(db)
	mealPlanRepo := repository.NewMealPlanRepository(db)
	recipeIngredientRepo := repository.NewRecipeIngredientRepository(db)

	shoppingListService := services.NewShoppingListService(
		mealPlanRepo,
		recipeIngredientRepo,
		shoppingListRepo,
	)

	shoppingListHandler := handlers.NewShoppingListHandler(shoppingListService)

	shopping := r.Group("/shopping-lists")
	{
		shopping.POST("", shoppingListHandler.GenerateShoppingList)
		shopping.GET("/:id", shoppingListHandler.GetShoppingList)
		shopping.PATCH("/items/:id/toggle", shoppingListHandler.ToggleItem)
	}
}
