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
)

type bookingUsecase struct {
	inquiryRepository   inquiry.InquiryRepository
	bookingRepository   booking.BookingRepository
	serviceRepository   services.ServicesRepository
	notificationService integration.NotificationService
}

func NewBookingUsecase(repositories domain.Repositories, integrations domain.Infrastructure) booking.BookingUsecase {
	return &bookingUsecase{
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

	err = usecase.bookingRepository.UpsertBookingsDocument(ctx, bookingsDocument)

	service, err := usecase.serviceRepository.GetServiceByID(ctx, inquiryEntity.ServiceID)
	err = usecase.notificationService.
		SendWhatsAppMessage(ctx, fmt.Sprintf("Hi %s, Berikut adalah kode konfirmasimu [%s] untuk pemesanan di %s! :D", inquiryEntity.FullName, strings.ToUpper(confirmationCode), service.Title), inquiryEntity.PhoneNumber)

	inquiryEntity.Status = 3
	err = usecase.inquiryRepository.UpdateInquiry(ctx, *inquiryEntity)
	if err != nil {
		return err
	}
	
	return nil
}
