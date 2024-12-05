package business

import (
	"mini-wallet/domain"
	_auth "mini-wallet/domain/auth"
	"mini-wallet/domain/business"
	"mini-wallet/domain/common/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type businessHandler struct {
	businessUsecase business.BusinessUsecase
}

func SetBusinessHandler(router *chi.Mux, usecases domain.Usecases, middleware _auth.AuthMiddleware) {
	businessHandler := businessHandler{
		businessUsecase: usecases.BusinessUsecase,
	}

	router.Route("/businesses/", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/", businessHandler.CreateBusiness)
		r.Get("/status", businessHandler.GetUserBusinessStatus)
	})

}

func (handler *businessHandler) CreateBusiness(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := business.BusinessCreationDTO{}
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
		resp.WriteResponse()
		return
	}

	res := handler.businessUsecase.CreateBusiness(r.Context(), req)
	res.Writer = w
	res.WriteResponse()
}

func (handler *businessHandler) GetUserBusinessStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(_auth.UserIDContext{}).(*string)

	res := handler.businessUsecase.GetUserBusinessStatus(r.Context(), *userID)
	res.Writer = w
	res.WriteResponse()
}
