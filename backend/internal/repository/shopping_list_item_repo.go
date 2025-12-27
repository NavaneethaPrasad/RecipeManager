package repository

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type ShoppingListItemRepository interface {
	Create(item *models.ShoppingListItem) error
	DeleteByShoppingListID(listID uint) error
}

type shoppingListItemRepository struct {
	DB *gorm.DB
}

func NewShoppingListItemRepository(db *gorm.DB) ShoppingListItemRepository {
	return &shoppingListItemRepository{DB: db}
}

func (r *shoppingListItemRepository) Create(item *models.ShoppingListItem) error {
	return r.DB.Create(item).Error
}

func (r *shoppingListItemRepository) DeleteByShoppingListID(listID uint) error {
	return r.DB.Where("shopping_list_id = ?", listID).
		Delete(&models.ShoppingListItem{}).Error
}
