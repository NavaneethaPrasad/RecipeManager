package main

import (
	"log"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/database"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/routes"
)

func main() {

	// Connect to Postgres
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate
	err = database.CreateDB(db)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	//Setup routes and pass DB
	r := routes.SetupRoutes(db)

	log.Println("Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
