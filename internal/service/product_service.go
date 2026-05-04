package service

import (
	"github.com/user/go-commerce-api/internal/model"
	"github.com/user/go-commerce-api/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) CreateProduct(product *model.Product) error {
	return s.productRepo.Create(product)
}

func (s *ProductService) ListProducts(query string) ([]model.Product, error) {
	return s.productRepo.FindAll(query)
}
