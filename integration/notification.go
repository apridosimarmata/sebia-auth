package integration

import (
	"context"
	"errors"
	"mini-wallet/infrastructure/proto/generated/notifications"

	grpc "google.golang.org/grpc"
)

type notificationService struct {
	client *grpc.ClientConn
}

type NotificationService interface {
	SendWhatsAppMessage(ctx context.Context, message string, destination string) (err error)
}

func NewNotificationService(client *grpc.ClientConn) NotificationService {
	return &notificationService{
		client: client,
	}
}

func (service *notificationService) SendWhatsAppMessage(ctx context.Context, message string, destination string) (err error) {
	x := notifications.NewNotificationServiceClient(
		service.client,
	)

	res, err := x.SendWhatsAppMessage(ctx, &notifications.WhatsAppMessage{
		Content: message,
		Target:  destination,
	})

	if err != nil || res.Error {
		return errors.New("an error occured")
	}

	return nil
}
