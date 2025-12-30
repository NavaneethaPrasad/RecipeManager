package repository

import (
	"github.com/NavaneethaPrasad/RecipeManager/backend/internal/models"
	"gorm.io/gorm"
)

type ShoppingListRepository interface {
	Create(list *models.ShoppingList) error
	CreateItem(item *models.ShoppingListItem) error
	FindByID(id uint) (*models.ShoppingList, error)
	FindItemsByListID(listID uint) ([]models.ShoppingListItem, error)
	FindItemByID(id uint) (*models.ShoppingListItem, error)
	UpdateItem(item *models.ShoppingListItem) error
}

type shoppingListRepository struct {
	DB *gorm.DB
}

func NewShoppingListRepository(db *gorm.DB) ShoppingListRepository {
	return &shoppingListRepository{DB: db}
}

func (r *shoppingListRepository) Create(list *models.ShoppingList) error {
	return r.DB.Create(list).Error
}

func (r *shoppingListRepository) CreateItem(item *models.ShoppingListItem) error {
	return r.DB.Create(item).Error
}

func (r *shoppingListRepository) FindByID(id uint) (*models.ShoppingList, error) {
	var list models.ShoppingList
	err := r.DB.First(&list, id).Error
	return &list, err
}

func (r *shoppingListRepository) FindItemsByListID(listID uint) ([]models.ShoppingListItem, error) {
	var items []models.ShoppingListItem
	err := r.DB.Preload("Ingredient").Where("shopping_list_id = ?", listID).Find(&items).Error
	return items, err
}

func (r *shoppingListRepository) FindItemByID(id uint) (*models.ShoppingListItem, error) {
	var item models.ShoppingListItem
	err := r.DB.First(&item, id).Error
	return &item, err
}

func (r *shoppingListRepository) UpdateItem(item *models.ShoppingListItem) error {
	return r.DB.Save(item).Error
}
