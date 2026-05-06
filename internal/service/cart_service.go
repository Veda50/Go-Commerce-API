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
	// 1. Cek stok produk
	product, err := s.productRepo.FindByID(fmt.Sprintf("%d", productID))
	if err != nil {
		return errors.New("product not found")
	}

	if product.Stock < quantity {
		return errors.New("insufficient stock")
	}

	// 2. Ambil atau buat cart aktif
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

	// 3. Tambah atau update item
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

	// Update quantity jika sudah ada
	if product.Stock < item.Quantity+quantity {
		return errors.New("insufficient stock for updated quantity")
	}

	item.Quantity += quantity
	return s.cartRepo.UpdateItem(item)
}
