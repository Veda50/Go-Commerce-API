package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/user/go-commerce-api/internal/middleware"
	"github.com/user/go-commerce-api/internal/repository"
	"github.com/user/go-commerce-api/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
	orderRepo    *repository.OrderRepository
}

func NewOrderHandler(orderService *service.OrderService, orderRepo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		orderRepo:    orderRepo,
	}
}

func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	order, invoiceURL, err := h.orderService.Checkout(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order":       order,
		"invoice_url": invoiceURL,
	})
}

func (h *OrderHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusUnauthorized)
		return
	}

	orders, err := h.orderRepo.FindByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
