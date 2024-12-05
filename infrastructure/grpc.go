package infrastructure

import (
	"log"

	grpc "google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
)

type GrpcConn struct {
	NotificationService grpc.ClientConn
}

func NewGrpcConn() (*GrpcConn, error) {
	conn, err := grpc.Dial("localhost:6005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}

	return &GrpcConn{
		NotificationService: *conn,
	}, nil
}
