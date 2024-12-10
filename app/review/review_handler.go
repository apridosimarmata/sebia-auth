package review

import (
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/review"
	"net/http"

	_auth "mini-wallet/domain/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type reviewHandler struct {
	reviewUsecase review.ReviewUsecase
}

func SetReviewHandler(router *chi.Mux, usecases domain.Usecases, middleware _auth.AuthMiddleware) {
	reviewHandler := reviewHandler{
		reviewUsecase: usecases.ReviewUsecase,
	}

	router.Route("/reviews", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/", reviewHandler.CreateReview)
	})

	router.Route("/public/reviews", func(r chi.Router) {
		r.Get("/top/{serviceId}", reviewHandler.GetServiceTopReview)
	})
}

func (handler *reviewHandler) GetServiceTopReview(w http.ResponseWriter, r *http.Request) {
	resp := response.Response[string]{
		Writer: w,
	}

	serviceId := chi.URLParam(r, "serviceId")
	if serviceId == "" {
		resp.BadRequest("serviceId can not be ampty", nil)
		return
	}

	res := handler.reviewUsecase.GetServiceTopReview(r.Context(), serviceId)
	res.Writer = w
	res.WriteResponse()
}

func (handler *reviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	resp := response.Response[string]{
		Writer: w,
	}

	req := review.ReviewDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	userID := r.Context().Value(_auth.UserIDContext{}).(*string)
	req.UserID = *userID
	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	res := handler.reviewUsecase.CreateReview(r.Context(), req)
	res.Writer = w
	res.WriteResponse()
}
