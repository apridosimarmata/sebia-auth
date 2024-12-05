package affiliate

import (
	"mini-wallet/domain"
	"mini-wallet/domain/affiliate"
	"mini-wallet/domain/auth"
	_auth "mini-wallet/domain/auth"
	"mini-wallet/domain/common/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type affiliatesHandler struct {
	affiliateUsecase affiliate.AffiliateUsecase
}

func SetAffiliatesHandler(router *chi.Mux, usecases domain.Usecases, middleware _auth.AuthMiddleware) {
	handler := affiliatesHandler{
		affiliateUsecase: usecases.AffiliateUsecase,
	}

	router.Route("/affiliates/", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/", handler.ApplyForAffiliate)
		r.Get("/status", handler.GetUserAffiliatesStatus)
	})

}

func (handler *affiliatesHandler) ApplyForAffiliate(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := affiliate.AffiliateAppicationDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	userID := r.Context().Value(auth.UserIDContext{}).(*string)
	req.UserID = *userID
	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.affiliateUsecase.ApplyForAffiliate(r.Context(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *affiliatesHandler) GetUserAffiliatesStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDContext{}).(*string)
	res := handler.affiliateUsecase.GetUserAffiliateStatus(r.Context(), *userID)
	res.Writer = w
	res.WriteResponse()
}
