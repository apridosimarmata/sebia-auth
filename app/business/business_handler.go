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

	router.Route("/public/businesses", func(r chi.Router) {
		r.Get("/{handle}", businessHandler.GetBusinessByHandle)
		r.Get("/id/{id}", businessHandler.GetBusinessById)

	})

}

func (handler *businessHandler) GetBusinessById(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		resp.BadRequest("invalid business id", nil)
		return
	}

	res := handler.businessUsecase.GetBusinessByID(r.Context(), id)
	res.Writer = w
	res.WriteResponse()
}

func (handler *businessHandler) GetBusinessByHandle(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	handle := chi.URLParam(r, "handle")
	if handle == "" {
		resp.BadRequest("invalid business handle", nil)
		return
	}

	res := handler.businessUsecase.GetBusinessByHandle(r.Context(), handle)
	res.Writer = w
	res.WriteResponse()
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
