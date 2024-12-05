package payment

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/payment"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type paymentHandler struct {
	paymentUsecase payment.PaymentUsecase
}

func SetPaymentHandler(router *chi.Mux, usecases domain.Usecases) {
	paymentHandler := paymentHandler{
		paymentUsecase: usecases.PaymentUsecase,
	}

	router.Route("/payments", func(r chi.Router) {
		r.Post("/callback", paymentHandler.HandlePaymentCallback)
	})

}

func (handler *paymentHandler) HandlePaymentCallback(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := payment.PaymentCallbackDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.paymentUsecase.HandlePaymentCallback(context.Background(), req)
	res.Writer = w
	res.WriteResponse()
}
