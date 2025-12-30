package services

import (
	"errors"
	"time"

	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/dto"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/repository"
)

var ErrInvalidDateRange = errors.New("invalid date range")

type ShoppingListService interface {
	Generate(userID uint, startDate string, endDate string) (*dto.ShoppingListResponse, error)
	GetShoppingListByID(listID uint, userID uint) (*dto.ShoppingListResponse, error)
	ToggleItemChecked(itemID uint, userID uint) error
}

type shoppingListService struct {
	MealPlanRepo         repository.MealPlanRepository
	RecipeIngredientRepo repository.RecipeIngredientRepository
	ShoppingListRepo     repository.ShoppingListRepository
}

func NewShoppingListService(
	mealPlanRepo repository.MealPlanRepository,
	recipeIngredientRepo repository.RecipeIngredientRepository,
	shoppingListRepo repository.ShoppingListRepository,
) ShoppingListService {
	return &shoppingListService{
		MealPlanRepo:         mealPlanRepo,
		RecipeIngredientRepo: recipeIngredientRepo,
		ShoppingListRepo:     shoppingListRepo,
	}
}

func (s *shoppingListService) Generate(
	userID uint,
	startDateStr string,
	endDateStr string,
) (*dto.ShoppingListResponse, error) {

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return nil, err
	}

	if endDate.Before(startDate) {
		return nil, ErrInvalidDateRange
	}

	mealPlans, err := s.MealPlanRepo.FindByUserAndDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	type key struct {
		IngredientID uint
		Unit         string
	}

	type aggrItem struct {
		Name     string
		Quantity float64
	}

	aggregated := make(map[key]*aggrItem)

	for _, mp := range mealPlans {
		// Note: Ideally, use Preload("Recipe.Ingredients") in MealPlanRepo to avoid this N+1 loop
		items, err := s.RecipeIngredientRepo.FindByRecipeID(mp.RecipeID)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			k := key{item.IngredientID, item.Unit}

			if v, ok := aggregated[k]; ok {
				v.Quantity += item.Quantity
			} else {
				aggregated[k] = &aggrItem{
					Name:     item.Ingredient.Name,
					Quantity: item.Quantity,
				}
			}
		}
	}

	list := &models.ShoppingList{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := s.ShoppingListRepo.Create(list); err != nil {
		return nil, err
	}

	var responseItems []dto.ShoppingListItemResponse

	for k, v := range aggregated {
		slItem := &models.ShoppingListItem{
			ShoppingListID: list.ID,
			IngredientID:   k.IngredientID,
			Quantity:       v.Quantity,
			Unit:           k.Unit,
			Checked:        false, // Matches your Model
		}

		if err := s.ShoppingListRepo.CreateItem(slItem); err != nil {
			return nil, err
		}

		responseItems = append(responseItems, dto.ShoppingListItemResponse{
			ID:           slItem.ID,
			IngredientID: k.IngredientID,
			Name:         v.Name, // Use the name we aggregated
			Quantity:     v.Quantity,
			Unit:         k.Unit,
			Checked:      slItem.Checked,
		})
	}

	return &dto.ShoppingListResponse{
		ID:        list.ID,
		StartDate: startDateStr,
		EndDate:   endDateStr,
		Items:     responseItems,
	}, nil
}

func (s *shoppingListService) GetShoppingListByID(
	listID uint,
	userID uint,
) (*dto.ShoppingListResponse, error) {

	list, err := s.ShoppingListRepo.FindByID(listID)
	if err != nil {
		return nil, err
	}

	if list.UserID != userID {
		return nil, ErrUnauthorized
	}

	items, err := s.ShoppingListRepo.FindItemsByListID(listID)
	if err != nil {
		return nil, err
	}

	var responseItems []dto.ShoppingListItemResponse
	for _, item := range items {
		// FIX: Use item.Ingredient.Name because 'ShoppingListItem' has no 'Name' field
		itemName := "Unknown"
		if item.IngredientID != 0 {
			itemName = item.Ingredient.Name
		}

		responseItems = append(responseItems, dto.ShoppingListItemResponse{
			ID:           item.ID,
			IngredientID: item.IngredientID,
			Name:         itemName, // Corrected logic
			Quantity:     item.Quantity,
			Unit:         item.Unit,
			Checked:      item.Checked, // FIX: Use Checked (not IsPurchased)
		})
	}

	return &dto.ShoppingListResponse{
		ID:        list.ID,
		StartDate: list.StartDate.Format("2006-01-02"),
		EndDate:   list.EndDate.Format("2006-01-02"),
		Items:     responseItems,
	}, nil
}

func (s *shoppingListService) ToggleItemChecked(
	itemID uint,
	userID uint,
) error {

	item, err := s.ShoppingListRepo.FindItemByID(itemID)
	if err != nil {
		return err
	}

	list, err := s.ShoppingListRepo.FindByID(item.ShoppingListID)
	if err != nil {
		return err
	}

	if list.UserID != userID {
		return ErrUnauthorized
	}

	// FIX: Use Checked
	item.Checked = !item.Checked
	return s.ShoppingListRepo.UpdateItem(item)
}
