package repository

import (
	"github.com/google/uuid"
	"github.com/user/go-commerce-api/internal/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) UpdateStatusByPaymentID(paymentID string, status model.OrderStatus) error {
	return r.db.Model(&model.Order{}).Where("payment_intent_id = ?", paymentID).Update("status", status).Error
}

func (r *OrderRepository) FindAll(status string) ([]model.Order, error) {
	var orders []model.Order
	db := r.db.Preload("Items.Product")
	if status != "" {
		db = db.Where("status = ?", status)
	}
	err := db.Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) FindByUserID(userID uuid.UUID) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) FindByPaymentID(paymentID string) (*model.Order, error) {
	var order model.Order
	err := r.db.Where("payment_intent_id = ?", paymentID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
