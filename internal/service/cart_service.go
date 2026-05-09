package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/user/go-commerce-api/internal/model"
	"github.com/user/go-commerce-api/internal/repository"
	"gorm.io/gorm"
)

type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewCartService(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *CartService) AddToCart(userID uuid.UUID, productID uint, quantity int) error {
	product, err := s.productRepo.FindByID(fmt.Sprintf("%d", productID))
	if err != nil {
		return errors.New("product not found")
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	cart, err := s.cartRepo.FindActiveByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = &model.Cart{UserID: userID, Status: model.CartActive}
			if err := s.cartRepo.Create(cart); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	item, err := s.cartRepo.FindItem(cart.ID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			item = &model.CartItem{
				CartID:    cart.ID,
				ProductID: productID,
				Quantity:  quantity,
			}
			return s.cartRepo.AddItem(item)
		}
		return err
	}

	if product.Stock < item.Quantity+quantity {
		return errors.New("insufficient stock for updated quantity")
	}

	item.Quantity += quantity
	return s.cartRepo.UpdateItem(item)
}

func (s *CartService) GetCart(userID uuid.UUID) (*model.Cart, error) {
	return s.cartRepo.FindActiveWithItems(userID)
}

func (s *CartService) RemoveFromCart(userID uuid.UUID, itemID string) error {
	cart, err := s.cartRepo.FindActiveByUserID(userID)
	if err != nil {
		return errors.New("no active cart found")
	}

	item, err := s.cartRepo.FindItemByID(itemID)
	if err != nil {
		return errors.New("item not found in cart")
	}

	if item.CartID != cart.ID {
		return errors.New("item does not belong to user's cart")
	}

	return s.cartRepo.DeleteItem(itemID)
}

func (s *CartService) UpdateQuantity(userID uuid.UUID, itemID string, quantity int) error {
	cart, err := s.cartRepo.FindActiveByUserID(userID)
	if err != nil {
		return errors.New("no active cart found")
	}

	item, err := s.cartRepo.FindItemByID(itemID)
	if err != nil {
		return errors.New("item not found in cart")
	}

	if item.CartID != cart.ID {
		return errors.New("item does not belong to user's cart")
	}

	product, err := s.productRepo.FindByID(fmt.Sprintf("%d", item.ProductID))
	if err != nil {
		return errors.New("product not found")
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	item.Quantity = quantity
	return s.cartRepo.UpdateItem(item)
}
