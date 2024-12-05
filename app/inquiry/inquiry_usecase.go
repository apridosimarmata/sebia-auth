package inquiry

import (
	"context"
	"errors"
	"fmt"
	"mini-wallet/domain"
	"mini-wallet/domain/business"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/inquiry"
	"mini-wallet/domain/services"
	"mini-wallet/infrastructure"
	"mini-wallet/integration"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type inquiryUsecase struct {
	servicesRepository  services.ServicesRepository
	inquiryRepository   inquiry.InquiryRepository
	businessRepository  business.BusinessRepository
	notificationService integration.NotificationService
	paymentService      infrastructure.Payment
}

func NewInquiryUsecase(repositories domain.Repositories, integrations domain.Infrastructure) inquiry.InquiryUsecase {
	return &inquiryUsecase{
		servicesRepository:  repositories.ServicesRepository,
		inquiryRepository:   repositories.InquiryRepository,
		notificationService: integrations.NotificationService,
		paymentService:      integrations.PaymentService,
		businessRepository:  repositories.BusinessRepository,
	}
}

func (usecase *inquiryUsecase) CreateInquiry(ctx context.Context, req inquiry.InquiryDTO) (res response.Response[string]) {
	serviceEntity, err := usecase.servicesRepository.GetServiceBySlug(ctx, req.ServiceSlug)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if serviceEntity == nil {
		res.NotFound("layanan tidak ditemukan", nil)
		return
	}

	selectedVariant, err := usecase.validateSelectedVariant(req, serviceEntity.Variants)
	if err != nil {
		res.BadRequest(err.Error(), nil)
		return
	}

	total := len(req.SelectedDates) * selectedVariant.Price
	entity, err := req.ToInquiryEntity(serviceEntity.ToServiceEntity(serviceEntity.ID), *selectedVariant, total)
	if err != nil {
		res.BadRequest(err.Error(), nil)
		return
	}

	err = usecase.inquiryRepository.InsertInquiry(ctx, *entity)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	url, err := usecase.paymentService.CreatePaymentLink(ctx, *entity)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	p := message.NewPrinter(language.English)
	totalString := p.Sprintf("%d", total)
	err = usecase.notificationService.SendWhatsAppMessage(ctx,
		fmt.Sprintf("Halo %s,\nBerikut adalah link pembayaranmu untuk pemesanan %s sebesar Rp%s\n\n%s\n\nLakukan pembayaran sebelum 24 jam.\n\nCek status pembayaranmu di sini:\nhttps://tobacamping.id/bookings/%s", entity.FullName, serviceEntity.Title, totalString, url, entity.ID), req.PhoneNumber)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.Success(entity.ID)
	return
}

func (usecase *inquiryUsecase) GetInquiry(ctx context.Context, id string) (res response.Response[inquiry.InquiryDTO]) {

	inquiryEntity, err := usecase.inquiryRepository.GetInquiryById(ctx, id)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if inquiryEntity == nil {
		res.NotFound("inquiry not found", nil)
		return
	}

	// should store the service id instead of slug
	service, err := usecase.servicesRepository.GetServiceByID(ctx, inquiryEntity.ServiceID)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if service == nil {
		res.NotFound("service not found", nil)
		return
	}

	business, err := usecase.businessRepository.GetBusinessById(ctx, service.BusinessID)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if business == nil {
		res.NotFound("business not found", nil)
		return
	}

	result := inquiryEntity.ToInquiryDetailsResponse(service.ToServiceEntity(service.ID), *business)

	res.Success(result)
	return
}

func (usecase *inquiryUsecase) validateSelectedVariant(req inquiry.InquiryDTO, variants []services.ServiceVariant) (selectedVariant *services.ServiceVariant, err error) {
	for pos, variant := range variants {
		if fmt.Sprintf("%d", pos+1) == req.SelectedVariantID {
			selectedVariant = &variant
			if req.SelectedVariant.Duration != variant.Duration || req.SelectedVariant.Pax != variant.Pax || req.SelectedVariant.Price != variant.Price {
				return nil, errors.New("Terjadi perubahan harga, mohon coba lagi")
			}

			break
		}
	}

	return selectedVariant, nil
}
