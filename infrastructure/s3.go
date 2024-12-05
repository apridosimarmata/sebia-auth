package infrastructure

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewS3Service() (*s3.S3, error) {
	region := "ap-southeast-2"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			"AKIAZDZTBVPPF3YJBB4F",
			"Cf9UB9Awyad8wRiIFe8ExJusZw2jQ0faAE0fz2N+",
			"",
		),
	})
	if err != nil {
		fmt.Println("Error creating session:", err)
		return nil, err
	}

	svc := s3.New(sess)

	return svc, nil
}
