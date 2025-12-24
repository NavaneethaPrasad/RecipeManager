package database

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/configs"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// creates a database connection
func Connect() (*gorm.DB, error) {

	dsn := configs.LoadConfig()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// database migrations
func CreateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Recipe{},
		&models.Ingredient{},
		&models.RecipeIngredient{},
		&models.MealPlan{},
		&models.ShoppingList{},
		&models.ShoppingListItem{},
	)
}
