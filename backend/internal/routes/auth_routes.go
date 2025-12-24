package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Public routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// 🔒 Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{
				"user_id": userID,
			})
		})
	}

	return r
}
