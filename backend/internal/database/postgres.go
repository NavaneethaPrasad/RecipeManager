package database

import (
	"log"

	"github.com/NavaneethaPrasad/RecipeManager/backend/configs"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	//DSN (Data Source Name)

	dsn := configs.LoadConfig()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to Dockerized PostgreSQL")

	// Auto Migrate all models
	log.Println("Running Migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Recipe{},
		&models.Ingredient{},
		&models.RecipeIngredient{},
		&models.MealPlan{},
		&models.ShoppingList{},
		&models.ShoppingListItem{},
	)
	if err != nil {
		log.Fatal("Migration Failed:", err)
	}

	return db, nil
}
