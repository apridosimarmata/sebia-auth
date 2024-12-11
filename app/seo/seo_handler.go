package seo

import (
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/seo"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type seoHandler struct {
	seoUsecase seo.SEOUsecase
}

func SetSeoHandler(router *chi.Mux, usecases domain.Usecases) {
	seoHandler := seoHandler{
		seoUsecase: usecases.SEOUsecase,
	}

	router.Route("/seo", func(r chi.Router) {
		r.Get("/populate/category", seoHandler.PopulateGroupsByCategoryID)
		r.Get("/category/{categoryId}", seoHandler.GetItemsByCategoryID)
	})

}

func (handler *seoHandler) PopulateGroupsByCategoryID(w http.ResponseWriter, r *http.Request) {

	res := handler.seoUsecase.PopulateFooterGroupForEachCategoryId(r.Context())
	res.Writer = w
	res.WriteResponse()
}

func (handler *seoHandler) GetItemsByCategoryID(w http.ResponseWriter, r *http.Request) {
	errResp := response.Response[string]{
		Writer: w,
	}
	categoryId := chi.URLParam(r, "categoryId")
	categoryIdInt, err := strconv.Atoi(categoryId)

	if err != nil {
		errResp.BadRequest("invalid category id", nil)
		return
	}
	res := handler.seoUsecase.GetItemsByCategoryId(r.Context(), categoryIdInt)
	res.Writer = w
	res.WriteResponse()
}
