package dto

type ScaledIngredientResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type ScaledRecipeResponse struct {
	RecipeID         uint                       `json:"recipe_id"`
	Name             string                     `json:"name"`
	OriginalServings int                        `json:"original_servings"`
	ScaledServings   int                        `json:"scaled_servings"`
	Ingredients      []ScaledIngredientResponse `json:"ingredients"`
}
