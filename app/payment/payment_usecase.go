package payment

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/booking"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/inquiry"
	"mini-wallet/domain/payment"
	"mini-wallet/infrastructure"
	"mini-wallet/utils"
	"time"
)

type paymentUsecase struct {
	inquiryRepository inquiry.InquiryRepository
	paymentService    infrastructure.Payment
	messageProducer   infrastructure.MessagingProducer
	config            *utils.AppConfig
}

func NewPaymentUsecase(repositories domain.Repositories, integrations domain.Infrastructure, config *utils.AppConfig) payment.PaymentUsecase {
	return &paymentUsecase{
		inquiryRepository: repositories.InquiryRepository,
		paymentService:    integrations.PaymentService,
		messageProducer:   integrations.MesageProducer,
		config:            config,
	}
}

func (usecase *paymentUsecase) HandlePaymentCallback(ctx context.Context, req payment.PaymentCallbackDTO) (res response.Response[string]) {
	err := usecase.paymentService.VerifyCallback(ctx, req.SignatureKey, req.OrderID, req.GrossAmount, req.StatusCode)
	if err != nil {
		res.BadRequest(err.Error(), nil)
		return
	}

	inquiryEntity, err := usecase.inquiryRepository.GetInquiryById(ctx, req.OrderID)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	// this should be wrapped by transaction
	if req.TransactionStatus == "capture" || req.TransactionStatus == "settlement" {
		inquiryEntity.Status = 2
		now, _ := utils.GetJktTime()
		inquiryEntity.UpdatedDate = now.Format(time.RFC3339)
		err = usecase.inquiryRepository.UpdateInquiry(ctx, *inquiryEntity)
		if err != nil {
			res.InternalServerError(err.Error())
			return
		}

		err = usecase.messageProducer.PublishMessage(ctx, usecase.config.BookingTopic, "", booking.BookingCreationRequest{
			InquiryID: req.OrderID,
		})
		if err != nil {
			res.InternalServerError(err.Error())
			return
		}

		// commit here TODO
	}

	res.Success("thanks! <3 callback received")
	return
}
