package dto

type CreateMealPlanRequest struct {
	RecipeID uint   `json:"recipe_id" binding:"required"`
	Date     string `json:"date" binding:"required"` // YYYY-MM-DD
	MealType string `json:"meal_type" binding:"required"`
}

type UpdateMealPlanRequest struct {
	RecipeID uint   `json:"recipe_id" binding:"required"`
	MealType string `json:"meal_type" binding:"required"`
}

type MealPlanResponse struct {
	ID         uint   `json:"id"`
	Date       string `json:"date"`
	MealType   string `json:"meal_type"`
	RecipeID   uint   `json:"recipe_id"`
	RecipeName string `json:"recipe_name"`
}
