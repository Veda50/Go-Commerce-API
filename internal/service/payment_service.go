package service

import (
	"context"
	"fmt"

	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/invoice"
)

type PaymentService struct {
	client *xendit.APIClient
}

func NewPaymentService(secretKey string) *PaymentService {
	client := xendit.NewClient(secretKey)
	return &PaymentService{client: client}
}

func (s *PaymentService) CreateInvoice(orderID uint, amount float64, email string) (string, string, error) {
	createInvoiceRequest := invoice.NewCreateInvoiceRequest(
		fmt.Sprintf("ORDER-%d", orderID),
		amount,
	)
	createInvoiceRequest.SetPayerEmail(email)

	inv, _, err := s.client.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(*createInvoiceRequest).
		Execute()

	if err != nil {
		return "", "", err
	}

	return inv.GetId(), inv.GetExternalId(), nil
}

func (s *PaymentService) CreateInvoiceURL(orderID uint, amount float64, email string) (string, string, error) {
	createInvoiceRequest := invoice.NewCreateInvoiceRequest(
		fmt.Sprintf("ORDER-%d", orderID),
		amount,
	)
	createInvoiceRequest.SetPayerEmail(email)

	inv, _, err := s.client.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(*createInvoiceRequest).
		Execute()

	if err != nil {
		return "", "", err
	}

	return inv.GetInvoiceUrl(), inv.GetId(), nil
}
