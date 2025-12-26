package dto

type CreateRecipeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Servings    int    `json:"servings" binding:"required"`
	PrepTime    int    `json:"prep_time"`
	CookTime    int    `json:"cook_time"`
	Category    string `json:"category" binding:"required"`
}

type RecipeResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Servings  int    `json:"servings"`
	TotalTime int    `json:"total_time"`
	Category  string `json:"category"`
}

type UpdateRecipeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Servings    int    `json:"servings"`
	PrepTime    int    `json:"prep_time"`
	CookTime    int    `json:"cook_time"`
	Category    string `json:"category" binding:"required"`
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
