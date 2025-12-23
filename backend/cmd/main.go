package main

import (
	"log"
	"os"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Connect to Database
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	_ = db // db will be used later in repositories

	// 2. Setup Router
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "Database Connected",
		})
	})

	// 3. Get App Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 4. Run Server
	r.Run(":" + port)
}
