package payment

import (
	"context"
	"mini-wallet/domain/common/response"
)

type PaymentCallbackDTO struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	SignatureKey      string `json:"signature_key"`
}

type PaymentUsecase interface {
	HandlePaymentCallback(ctx context.Context, req PaymentCallbackDTO) (res response.Response[string])
}
