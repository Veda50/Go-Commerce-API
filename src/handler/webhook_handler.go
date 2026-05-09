package handler

import (
	"encoding/json"
	"net/http"

	"github.com/user/go-commerce-api/src/config"
	"github.com/user/go-commerce-api/src/model"
	"github.com/user/go-commerce-api/src/response"
	"github.com/user/go-commerce-api/src/service"
)

type WebhookHandler struct {
	orderService *service.OrderService
	webhookToken string
}

func NewWebhookHandler(orderService *service.OrderService, cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{
		orderService: orderService,
		webhookToken: cfg.XenditWebToken,
	}
}

func (h *WebhookHandler) XenditCallback(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-callback-token")
	if token != h.webhookToken {
		response.Error(w, http.StatusUnauthorized, "invalid callback token")
		return
	}

	var payload struct {
		ID         string  `json:"id"`
		ExternalID string  `json:"external_id"`
		Status     string  `json:"status"`
		Amount     float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if payload.Status == "PAID" {
		if err := h.orderService.UpdateOrderStatus(payload.ID, model.OrderPaid); err != nil {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
