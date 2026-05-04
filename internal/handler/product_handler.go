package handler

import (
	"encoding/json"
	"net/http"

	"github.com/user/go-commerce-api/internal/service"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	products, err := h.productService.ListProducts(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
