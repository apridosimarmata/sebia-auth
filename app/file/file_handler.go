package file

import (
	"context"
	"fmt"
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/file"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type fileHandler struct {
	fileUsecase file.FileUsecase
}

func SetFileHandler(router *chi.Mux, usecase domain.Usecases) {
	fileHandler := fileHandler{
		fileUsecase: usecase.FileUsecase,
	}

	router.Route("/file", func(r chi.Router) {
		r.Post("/", fileHandler.UploadFile)
	})
}

func (handler *fileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	resp := &response.Response[string]{
		Writer: w,
	}

	// Retrieve the file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		resp.BadRequest(err.Error(), nil)
		return
	}
	defer file.Close()

	public := r.URL.Query().Get("public")
	publicBool, err := strconv.ParseBool(public)
	if err != nil {
		resp.BadRequest(fmt.Sprintf("invalid public param : %v", &public), nil)
		return
	}

	res := handler.fileUsecase.UploadFile(context.Background(), file, header, publicBool)
	res.Writer = w
	res.WriteResponse()
}
