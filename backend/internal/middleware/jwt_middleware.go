package middleware

import (
	"net/http"
	"strings"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		var tokenStr string

		// Try cookie first (browser apps)
		if cookieToken, err := c.Cookie("auth_token"); err == nil {
			tokenStr = cookieToken
		} else {
			// Fallback to Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				c.Abort()
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Authorization header must be Bearer <token>",
				})
				c.Abort()
				return
			}

			tokenStr = parts[1]
		}

		// Validate token via utils
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store user_id in context
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
