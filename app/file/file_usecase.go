package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/file"
	"mini-wallet/utils"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type fileUsecase struct {
	s3Service s3.S3
}

func NewFileUsecase(infra domain.Infrastructure) file.FileUsecase {
	return &fileUsecase{
		s3Service: infra.S3,
	}
}

func (usecase *fileUsecase) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, public bool) (res response.Response[string]) {
	// Read the contents of the file into a buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		res.InternalServerError(err.Error())
		return
	}

	bucketName := "sebia"
	if public {
		bucketName = "sebia-public"
	}

	ext := filepath.Ext(header.Filename)
	random, err := utils.GenerateRandomString(10)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	fileName := random + ext

	// This uploads the contents of the buffer to S3
	_, err = usecase.s3Service.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		fmt.Println(err.Error())
		res.InternalServerError(err.Error())
		return
	}

	res.Success(fileName)
	return
}
