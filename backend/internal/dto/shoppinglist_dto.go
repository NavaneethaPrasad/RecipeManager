package dto

type GenerateShoppingListRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type ShoppingListResponse struct {
	ID        uint                       `json:"id"`
	StartDate string                     `json:"start_date"`
	EndDate   string                     `json:"end_date"`
	Items     []ShoppingListItemResponse `json:"items"`
}

type ShoppingListItemResponse struct {
	ID           uint    `json:"id"`
	IngredientID uint    `json:"ingredient_id"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	Checked      bool    `json:"checked"`
}
