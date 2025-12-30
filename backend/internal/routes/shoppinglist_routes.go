package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterShoppingListRoutes(r *gin.RouterGroup, db *gorm.DB) {

	mealPlanRepo := repository.NewMealPlanRepository(db)
	shoppingRepo := repository.NewShoppingListRepository(db)

	recipeIngRepo := repository.NewRecipeIngredientRepository(db)

	service := services.NewShoppingListService(mealPlanRepo, recipeIngRepo, shoppingRepo)

	handler := handlers.NewShoppingListHandler(service)

	shopping := r.Group("/shopping-lists")
	{
		shopping.POST("/generate", handler.GenerateShoppingList)
		shopping.GET("/:id", handler.GetShoppingList)
		shopping.PATCH("/items/:id/toggle", handler.ToggleItem)
	}
}
