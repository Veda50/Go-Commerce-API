package repository

import (
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

func (r *OrderRepository) FindByID(id string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items.Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) UpdateStatus(id string, status model.OrderStatus) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) FindAll(status string) ([]model.Order, error) {
	var orders []model.Order
	db := r.db.Preload("Items.Product")
	if status != "" {
		db = db.Where("status = ?", status)
	}
	err := db.Order("created_at desc").Find(&orders).Error
	return orders, err
}
