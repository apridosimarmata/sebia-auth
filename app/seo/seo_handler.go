package seo

import (
	"mini-wallet/domain"
	"mini-wallet/domain/seo"
	"net/http"

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
	})

}

func (handler *seoHandler) PopulateGroupsByCategoryID(w http.ResponseWriter, r *http.Request) {

	res := handler.seoUsecase.PopulateFooterGroupForEachCategoryId(r.Context())
	res.Writer = w
	res.WriteResponse()
}
