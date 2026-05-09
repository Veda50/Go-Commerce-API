package handler

import (
	"encoding/json"
	"net/http"

	"github.com/user/go-commerce-api/internal/config"
	"github.com/user/go-commerce-api/internal/model"
	"github.com/user/go-commerce-api/internal/service"
)

type WebhookHandler struct {
	orderService   *service.OrderService
	webhookToken string
}

func NewWebhookHandler(orderService *service.OrderService, cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{
		orderService:   orderService,
		webhookToken: cfg.XenditWebToken,
	}
}

func (h *WebhookHandler) XenditCallback(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-callback-token")
	if token != h.webhookToken {
		http.Error(w, "invalid callback token", http.StatusUnauthorized)
		return
	}

	var payload struct {
		ID         string  `json:"id"`
		ExternalID string  `json:"external_id"`
		Status     string  `json:"status"`
		Amount     float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if payload.Status == "PAID" {
		if err := h.orderService.UpdateOrderStatus(payload.ID, model.OrderPaid); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
