package services

import (
	"mini-wallet/domain"
	_auth "mini-wallet/domain/auth"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/services"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

type servicesHandler struct {
	servicesUsecase services.ServicesUsecase
	decoder         *schema.Decoder
}

func SetServicesHandler(router *chi.Mux, usecases domain.Usecases, middleware _auth.AuthMiddleware) {
	servicesHandler := servicesHandler{
		servicesUsecase: usecases.ServicesUsecase,
		decoder:         schema.NewDecoder(),
	}

	router.Route("/services", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Post("/", servicesHandler.CreateService)
		r.Put("/", servicesHandler.UpdateService)

		r.Get("/", servicesHandler.GetServices)
		r.Get("/{slug}", servicesHandler.GetServiceBySlug)
	})
	router.Route("/public/services", func(r chi.Router) {
		r.Get("/", servicesHandler.GetPublicServices)
		r.Get("/{slug}", servicesHandler.GetServiceBySlug)
		r.Get("/search", servicesHandler.SearchServicesByKeyword)
	})
}

func (handler *servicesHandler) SearchServicesByKeyword(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		resp.BadRequest("Kata kunci harus di isi", nil)
		return
	}

	res := handler.servicesUsecase.SearchServicesByKeyword(r.Context(), keyword)
	res.Writer = w
	res.WriteResponse()
}

func (handler *servicesHandler) GetPublicServices(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	var params services.GetPublicServicesRequest
	if err := handler.decoder.Decode(&params, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := params.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	res := handler.servicesUsecase.GetPublicServices(r.Context(), params)
	res.Writer = w
	res.WriteResponse()
}

func (handler *servicesHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	businessCookie, err := r.Cookie("business_id")
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	req := services.ServiceDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	req.BusinessID = businessCookie.Value
	userID := r.Context().Value(_auth.UserIDContext{}).(*string)
	err = req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.servicesUsecase.UpdateService(r.Context(), req, *userID)
	res.Writer = w
	res.WriteResponse()
}

func (handler *servicesHandler) GetServiceBySlug(w http.ResponseWriter, r *http.Request) {

	// var params services.GetServicesRequest
	// businessCookie, err := r.Cookie("business_id")
	slug := chi.URLParam(r, "slug")

	res := handler.servicesUsecase.GetServiceBySlug(r.Context(), slug)
	res.Writer = w
	res.WriteResponse()
}

func (handler *servicesHandler) GetServices(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	var params services.GetServicesRequest
	businessCookie, err := r.Cookie("business_id")
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	if err := handler.decoder.Decode(&params, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	params.BusinessID = &businessCookie.Value
	err = params.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}

	res := handler.servicesUsecase.GetServices(r.Context(), params)
	res.Writer = w
	res.WriteResponse()
}

func (handler *servicesHandler) CreateService(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	businessCookie, err := r.Cookie("business_id")
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	req := services.ServiceDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	req.BusinessID = businessCookie.Value
	userID := r.Context().Value(_auth.UserIDContext{}).(*string)
	err = req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	res := handler.servicesUsecase.CreateService(r.Context(), req, *userID)
	res.Writer = w
	res.WriteResponse()
}
