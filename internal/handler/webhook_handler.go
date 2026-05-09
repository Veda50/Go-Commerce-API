package handler

import (
	"encoding/json"
	"net/http"

	"github.com/user/go-commerce-api/internal/config"
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
	// 1. Verifikasi Webhook Token
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

	// Logic update order status akan ada di commit berikutnya
	// Untuk sekarang kita log saja
	// fmt.Printf("Received webhook for invoice %s with status %s\n", payload.ID, payload.Status)

	w.WriteHeader(http.StatusOK)
}
