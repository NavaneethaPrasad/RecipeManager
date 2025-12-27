package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/middleware"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMealPlanRoutes(r *gin.Engine, db *gorm.DB) {

	mealRepo := repository.NewMealPlanRepository(db)
	recipeRepo := repository.NewRecipeRepository(db)
	service := services.NewMealPlanService(mealRepo, recipeRepo)
	handler := handlers.NewMealPlanHandler(service)

	protected := r.Group("/api/meal-plans")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.POST("", handler.Create)
		protected.GET("", handler.GetByDate)
		protected.PUT("/:id", handler.Update)
		protected.DELETE("/:id", handler.Delete)
	}
}
