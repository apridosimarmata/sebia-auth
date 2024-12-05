package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"mini-wallet/domain"
	"mini-wallet/domain/booking"
	"mini-wallet/infrastructure"

	"github.com/nsqio/go-nsq"
)

type bookingMessageConsumer struct {
	bookingUsecase booking.BookingUsecase
}

func NewBookingMessageConsumer(usecases domain.Usecases) infrastructure.MessagingConsumerInterface {
	return &bookingMessageConsumer{
		bookingUsecase: usecases.BookingUsecase,
	}
}

func (consumer *bookingMessageConsumer) ConsumeMessage(message *nsq.Message) (err error) {
	req := booking.BookingCreationRequest{}

	err = json.Unmarshal(message.Body, &req)
	if err != nil {
		return err
	}

	err = consumer.bookingUsecase.CreateBooking(context.Background(), req.InquiryID)
	if err != nil {
		return err
	}

	fmt.Println("done processing booking creation request, error:", err)
	return nil
}
