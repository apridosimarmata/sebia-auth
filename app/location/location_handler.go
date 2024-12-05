package location

import (
	"fmt"
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/locations"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type locationHandler struct {
	locationUsecase locations.LocationUsecase
}

func SetLocationHandler(router *chi.Mux, usecase domain.Usecases) {
	locationHandler := locationHandler{
		locationUsecase: usecase.LocationUsecase,
	}

	router.Route("/locations", func(r chi.Router) {
		r.Get("/provinces", locationHandler.GetProvinces)
		r.Get("/cities", locationHandler.GetCitiesByProvinceID)
		r.Get("/districts", locationHandler.GetProvinces)

	})
}

func (handler *locationHandler) GetProvinces(w http.ResponseWriter, r *http.Request) {
	res := handler.locationUsecase.GetProvinces(r.Context())
	res.Writer = w
	res.WriteResponse()
}

func (handler *locationHandler) GetCitiesByProvinceID(w http.ResponseWriter, r *http.Request) {
	errRes := response.Response[interface{}]{}
	errRes.Writer = w
	provinceId, err := strconv.Atoi(r.URL.Query().Get("provinceID"))
	if err != nil {
		errRes.BadRequest(fmt.Sprintf("invalid provinceID %s", r.URL.Query().Get("provinceID")), nil)
		errRes.WriteResponse()
		return
	}

	res := handler.locationUsecase.GetCitiesByProvinceID(r.Context(), provinceId)
	res.Writer = w
	res.WriteResponse()
}

func (handler *locationHandler) GetDistrictsByCityID(w http.ResponseWriter, r *http.Request) {
	errRes := response.Response[interface{}]{}
	errRes.Writer = w

	cityId, err := strconv.Atoi(r.URL.Query().Get("cityID"))
	if err != nil {
		errRes.BadRequest(fmt.Sprintf("invalid cityID %s", r.URL.Query().Get("cityID")), nil)
		errRes.WriteResponse()
		return
	}

	res := handler.locationUsecase.GetDistrictByCityID(r.Context(), cityId)
	res.Writer = w
	res.WriteResponse()
}
