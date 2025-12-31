package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/middleware"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	RegisterAuthRoutes(r, db)

	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		RegisterRecipeRoutes(protected, db)
		RegisterIngredientRoutes(protected, db)
		// RegisterInstructionRoutes(protected, db)
		RegisterMealPlanRoutes(protected, db)
		RegisterShoppingListRoutes(protected, db)

		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")

			var user models.User
			if err := db.First(&user, userID).Error; err != nil {
				c.JSON(404, gin.H{"error": "User not found"})
				return
			}

			c.JSON(200, gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			})
		})
	}

	return r
}
