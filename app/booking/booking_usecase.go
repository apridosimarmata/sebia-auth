package booking

import (
	"context"
	"fmt"
	"mini-wallet/domain"
	"mini-wallet/domain/booking"
	"mini-wallet/domain/inquiry"
	"mini-wallet/domain/services"
	"mini-wallet/integration"
	"mini-wallet/utils"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type bookingUsecase struct {
	baseRepository      domain.BaseRepository
	inquiryRepository   inquiry.InquiryRepository
	bookingRepository   booking.BookingRepository
	serviceRepository   services.ServicesRepository
	notificationService integration.NotificationService
}

func NewBookingUsecase(repositories domain.Repositories, integrations domain.Infrastructure) booking.BookingUsecase {
	return &bookingUsecase{
		baseRepository:      repositories.BaseRepository,
		inquiryRepository:   repositories.InquiryRepository,
		bookingRepository:   repositories.BookingRepository,
		notificationService: integrations.NotificationService,
		serviceRepository:   repositories.ServicesRepository,
	}
}

func (usecase *bookingUsecase) CreateBooking(ctx context.Context, inquiryID string) error {
	var bookingsDocument []booking.ServiceBookings

	inquiryEntity, err := usecase.inquiryRepository.GetInquiryById(ctx, inquiryID)
	if err != nil {
		return err
	}

	confirmationCode, _ := utils.GenerateRandomString(7)

	yearMonths := make(map[string]struct{})

	for _, date := range inquiryEntity.SelectedDates {
		dateSplitted := strings.Split(date, "/")
		yearMonth := fmt.Sprintf("%s/%s", dateSplitted[0], dateSplitted[1])
		yearMonths[yearMonth] = struct{}{}
	}

	yearMonthsSlice := []string{}
	for yearMonth := range yearMonths {
		yearMonthsSlice = append(yearMonthsSlice, yearMonth)
	}

	bookingsDocument, err = usecase.bookingRepository.GetBookings(ctx, inquiryEntity.ServiceID, inquiryEntity.SelectedVariant.Pax, yearMonthsSlice)
	if err != nil {
		return err
	}

	bookingsDocumentMapped := make(map[string]booking.ServiceBookings)

	for _, document := range bookingsDocument {
		bookingsDocumentMapped[document.YearMonth] = document
	}

	var newBookingDocuments []booking.ServiceBookings
	if len(bookingsDocumentMapped) != len(yearMonths) {
		for _, yearMonth := range yearMonthsSlice {
			if _, found := bookingsDocumentMapped[yearMonth]; !found {
				newDocument := booking.ServiceBookings{}
				newDocument.Init(inquiryEntity.ServiceID, inquiryEntity.SelectedVariant.Pax, yearMonth)
				newBookingDocuments = append(newBookingDocuments, newDocument)
			}
		}
	}

	for _, newDocument := range newBookingDocuments {
		bookingsDocumentMapped[newDocument.YearMonth] = newDocument
	}

	for _, date := range inquiryEntity.SelectedDates {
		dateSplitted := strings.Split(date, "/")
		yearMonth := fmt.Sprintf("%s/%s", dateSplitted[0], dateSplitted[1])
		bookingDocument := bookingsDocumentMapped[yearMonth]

		date := dateSplitted[2]
		if bookingPerDate, found := bookingDocument.BookingsByDate[date]; !found {
			bookingsByHour := booking.BookingsByHour{}
			bookingsByHour[inquiryEntity.SelectedHour] = []booking.Booking{
				{
					ConfirmationCode: confirmationCode,
					InquiryID:        inquiryEntity.ID,
				},
			}
			_bookingPerDate := booking.BookingsPerDate{}
			_bookingPerDate[inquiryEntity.SelectedHour] = bookingsByHour

			bookingDocument.BookingsByDate[date] = _bookingPerDate
		} else {
			if bookingsByHour, found := bookingPerDate[inquiryEntity.SelectedHour]; !found {
				bookingsByHour := booking.BookingsByHour{}
				bookingsByHour[inquiryEntity.SelectedHour] = []booking.Booking{
					{
						ConfirmationCode: confirmationCode,
						InquiryID:        inquiryEntity.ID,
					},
				}
				bookingPerDate[inquiryEntity.SelectedHour] = bookingsByHour
				bookingDocument.BookingsByDate[date] = bookingPerDate

			} else {
				bookingsOnThisHour := bookingsByHour[inquiryEntity.SelectedHour]
				bookingsOnThisHour = append(bookingsOnThisHour, booking.Booking{
					ConfirmationCode: confirmationCode,
					InquiryID:        inquiryEntity.ID,
				})
				bookingsByHour[inquiryEntity.SelectedHour] = bookingsOnThisHour
				bookingPerDate[date] = bookingsByHour
				bookingDocument.BookingsByDate[date] = bookingPerDate
			}
		}

		bookingsDocumentMapped[yearMonth] = bookingDocument
	}

	for _, document := range bookingsDocumentMapped {
		bookingsDocument = append(bookingsDocument, document)
	}

	// update stage
	tx, err := usecase.baseRepository.GetTransaction(ctx)
	if err != nil {
		return err
	}

	defer usecase.CleanUpTransaction(ctx, *tx, err)

	err = usecase.bookingRepository.UpsertBookingsDocument(ctx, tx, bookingsDocument)
	if err != nil {
		return err
	}

	confirmationCodeUppercase := strings.ToUpper(confirmationCode)
	service, err := usecase.serviceRepository.GetServiceByID(ctx, inquiryEntity.ServiceID)
	if err != nil {
		return err
	}

	inquiryEntity.Status = 3
	inquiryEntity.ConfirmationCode = &confirmationCodeUppercase
	err = usecase.inquiryRepository.UpdateInquiryWithTx(ctx, tx, *inquiryEntity)
	if err != nil {
		return err
	}

	// commit stage
	err = usecase.baseRepository.CommitTransaction(ctx, *tx)
	if err != nil {
		return err
	}

	_ = usecase.notificationService.
		SendWhatsAppMessage(ctx, fmt.Sprintf("Hi %s, Berikut adalah kode konfirmasimu [%s] untuk pemesanan di %s!", inquiryEntity.FullName, confirmationCodeUppercase, service.Title), inquiryEntity.PhoneNumber)

	return nil
}

func (usecase *bookingUsecase) CleanUpTransaction(ctx context.Context, tx mongo.SessionContext, err error) error {
	if err != nil {
		err = usecase.baseRepository.AbortTransaction(ctx, tx)
		if err != nil {
			return err
		}
		return err

	}

	return nil
}
