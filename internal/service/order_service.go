package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/user/go-commerce-api/internal/model"
	"github.com/user/go-commerce-api/internal/repository"
	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
	db          *gorm.DB
}

func NewOrderService(orderRepo *repository.OrderRepository, cartRepo *repository.CartRepository, productRepo *repository.ProductRepository, db *gorm.DB) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		db:          db,
	}
}

func (s *OrderService) Checkout(userID uuid.UUID) (*model.Order, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var cart model.Cart
	if err := tx.Preload("Items.Product").Where("user_id = ? AND status = ?", userID, model.CartActive).First(&cart).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("cart not found or empty")
	}

	if len(cart.Items) == 0 {
		tx.Rollback()
		return nil, errors.New("cart is empty")
	}

	var total float64
	for _, item := range cart.Items {
		total += item.Product.Price * float64(item.Quantity)
	}

	order := &model.Order{
		UserID:      userID,
		Status:      model.OrderPending,
		TotalAmount: total,
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, item := range cart.Items {
		orderItem := model.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Product.Price,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := s.productRepo.DeductStock(tx, item.ProductID, item.Quantity); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Model(&cart).Update("status", model.CartCompleted).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrder(id string) (*model.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s *OrderService) UpdateOrderStatus(id string, status model.OrderStatus) error {
	return s.orderRepo.UpdateStatus(id, status)
}

func (s *OrderService) ListOrders(status string) ([]model.Order, error) {
	return s.orderRepo.FindAll(status)
}
