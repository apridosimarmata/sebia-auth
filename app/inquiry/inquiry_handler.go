package inquiry

import (
	"mini-wallet/domain"
	_auth "mini-wallet/domain/auth"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/inquiry"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

type inquiryHandler struct {
	inquiryUsecase inquiry.InquiryUsecase
	decoder        *schema.Decoder
}

func SetInquiryHandler(router *chi.Mux, usecases domain.Usecases, middleware _auth.AuthMiddleware) {
	inquiryHandler := inquiryHandler{
		inquiryUsecase: usecases.InquiryUsecase,
	}

	router.Route("/inquiries", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
	})

	router.Route("/public/inquiries", func(r chi.Router) {
		// r.Use(middleware.PublicMiddleware)
		r.Post("/", inquiryHandler.CreateInquiry)
		r.Get("/{inquiryId}", inquiryHandler.GetInquiryDetails)

	})
}

func (handler *inquiryHandler) GetInquiryDetails(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	inquiryId := chi.URLParam(r, "inquiryId")
	if inquiryId == "" {
		resp.BadRequest("invalid inquiry id", nil)
		return
	}

	res := handler.inquiryUsecase.GetInquiry(r.Context(), inquiryId)
	res.Writer = w
	res.WriteResponse()

}

func (handler *inquiryHandler) CreateInquiry(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	req := inquiry.InquiryDTO{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}

	err := req.Validate()
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		resp.WriteResponse()
		return
	}
	// userID := r.Context().Value(_auth.UserIDContext{}).(*string)
	// err := req.Validate()
	// if err != nil {
	// 	resp.BadRequest(err.Error(), nil)
	// 	resp.WriteResponse()
	// 	return
	// }

	// req.UserID = userID

	res := handler.inquiryUsecase.CreateInquiry(r.Context(), req)
	res.Writer = w
	res.WriteResponse()
}
