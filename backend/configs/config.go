package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() string {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, relying on System Env Vars")
	}
	DBHost := getEnv("DB_HOST", "localhost")
	DBUser := getEnv("DB_USER", "postgres")
	DBPassword := getEnv("DB_PASSWORD", "password")
	DBName := getEnv("DB_NAME", "recipe_db")
	DBPort := getEnv("DB_PORT", "5432")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DBHost, DBUser, DBPassword, DBName, DBPort)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
