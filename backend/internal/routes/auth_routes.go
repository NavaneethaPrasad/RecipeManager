package routes

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/handlers"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB) {

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
