package file

import (
	"context"
	"mime/multipart"
	"mini-wallet/domain/common/response"
)

type FileUsecase interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, public bool) response.Response[string]
}
