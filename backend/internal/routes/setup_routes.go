package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	RegisterAuthRoutes(r, db)
	RegisterRecipeRoutes(r, db)
	RegisterIngredientRoutes(r, db)
	RegisterInstructionRoutes(r, db)
	RegisterMealPlanRoutes(r, db)

	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{"user_id": userID})
		})
	}

	return r
}
