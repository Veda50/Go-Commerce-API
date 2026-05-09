package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/user/go-commerce-api/internal/response"
	"github.com/user/go-commerce-api/internal/service"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	products, err := h.service.ListProducts(query)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	product, err := h.service.GetProduct(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "product not found")
		return
	}

	response.JSON(w, http.StatusOK, product)
}
