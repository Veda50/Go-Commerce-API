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

func (s *ProductService) GetProduct(id string) (*model.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) UpdateProduct(id string, input *model.Product) (*model.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Stock = input.Stock
	product.CategoryID = input.CategoryID

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}
