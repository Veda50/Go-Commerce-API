package repository

import (
	"github.com/google/uuid"
	"github.com/user/go-commerce-api/internal/model"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) FindActiveByUserID(userID uuid.UUID) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Where("user_id = ? AND status = ?", userID, model.CartActive).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) FindActiveWithItems(userID uuid.UUID) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Preload("Items.Product").Where("user_id = ? AND status = ?", userID, model.CartActive).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) Create(cart *model.Cart) error {
	return r.db.Create(cart).Error
}

func (r *CartRepository) AddItem(item *model.CartItem) error {
	return r.db.Create(item).Error
}

func (r *CartRepository) FindItem(cartID uint, productID uint) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CartRepository) FindItemByID(itemID string) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.First(&item, itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CartRepository) UpdateItem(item *model.CartItem) error {
	return r.db.Save(item).Error
}

func (r *CartRepository) DeleteItem(itemID string) error {
	return r.db.Delete(&model.CartItem{}, itemID).Error
}
