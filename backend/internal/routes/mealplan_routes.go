package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMealPlanRoutes(r *gin.RouterGroup, db *gorm.DB) {

	mealRepo := repository.NewMealPlanRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)
	mealPlanService := services.NewMealPlanService(mealRepo, recipeRepo)
	mealPlanHandler := handlers.NewMealPlanHandler(mealPlanService)

	mealPlans := r.Group("/meal-plans")
	{
		mealPlans.POST("", mealPlanHandler.Create)
		mealPlans.GET("", mealPlanHandler.GetByDate)
		mealPlans.PUT("/:id", mealPlanHandler.Update)
		mealPlans.DELETE("/:id", mealPlanHandler.Delete)
	}
}
