package repository

import (
	"github.com/user/go-commerce-api/internal/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) FindAll(query string) ([]model.Product, error) {
	var products []model.Product
	db := r.db.Preload("Category")
	if query != "" {
		db = db.Where("name ILIKE ?", "%"+query+"%")
	}
	err := db.Find(&products).Error
	return products, err
}

func (r *ProductRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id string) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *ProductRepository) DeductStock(tx *gorm.DB, productID uint, quantity int) error {
	return tx.Model(&model.Product{}).Where("id = ?", productID).
		Update("stock", gorm.Expr("stock - ?", quantity)).Error
}

func (r *ProductRepository) UpdateStock(id string, stock int) error {
	return r.db.Model(&model.Product{}).Where("id = ?", id).Update("stock", stock).Error
}
