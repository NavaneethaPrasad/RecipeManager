package dto

type CreateIngredientRequest struct {
	Name string `json:"name" binding:"required"`
}

type IngredientMasterResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AddRecipeIngredientRequest struct {
	IngredientID uint    `json:"ingredient_id" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required"`
	Unit         string  `json:"unit" binding:"required"`
}
