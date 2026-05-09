package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/user/go-commerce-api/internal/middleware"
	"github.com/user/go-commerce-api/internal/repository"
	"github.com/user/go-commerce-api/internal/response"
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
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, _ := uuid.Parse(userIDStr)
	order, invoiceURL, err := h.orderService.Checkout(userID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"order":       order,
		"invoice_url": invoiceURL,
	})
}

func (h *OrderHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userIDStr, _ := r.Context().Value(middleware.UserIDKey).(string)
	userID, _ := uuid.Parse(userIDStr)

	orders, err := h.orderRepo.FindByUserID(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, orders)
}
