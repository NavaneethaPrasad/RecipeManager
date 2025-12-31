package dto

type RecipeIngredientRequest struct {
	Name   string  `json:"name" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
	Unit   string  `json:"unit" binding:"required"`
}

type CreateRecipeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Servings    int    `json:"servings" binding:"required,gt=0"`
	PrepTime    int    `json:"prep_time"`
	CookTime    int    `json:"cook_time"`
	Category    string `json:"category" binding:"required"`

	Ingredients  []RecipeIngredientRequest `json:"ingredients"`
	Instructions []string                  `json:"instructions"`
}

type RecipeResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Servings    int    `json:"servings"`
	TotalTime   int    `json:"total_time"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type UpdateRecipeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Servings    int    `json:"servings"`
	PrepTime    int    `json:"prep_time"`
	CookTime    int    `json:"cook_time"`
	Category    string `json:"category" binding:"required"`

	Ingredients  []RecipeIngredientRequest `json:"ingredients"`
	Instructions []string                  `json:"instructions"`
}

type RecipeDetailResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Servings    int    `json:"servings"`
	Category    string `json:"category"`
	PrepTime    int    `json:"prep_time"`
	CookTime    int    `json:"cook_time"`
	TotalTime   int    `json:"total_time"`

	Ingredients  []IngredientResponse  `json:"ingredients"`
	Instructions []InstructionResponse `json:"instructions"`
}

type IngredientResponse struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type InstructionResponse struct {
	ID         uint   `json:"id"`
	StepNumber int    `json:"step_number"`
	Text       string `json:"text"`
}
