package booking

import (
	"context"
	"mini-wallet/utils"
)

// duration type -> key
// half hour -> hh:15 30 45 00
// hour -> hh:00
// day -> date

type BookingCreationRequest struct {
	InquiryID string `json:"inquiry_id"`
}

type ServiceBookings struct {
	ID             string                     `json:"id" bson:"id"`
	ServiceID      string                     `json:"service_id" bson:"service_id"`
	VariantPax     int                        `json:"variant_pax" bson:"variant_pax"`
	YearMonth      string                     `json:"year_month" bson:"year_month"` // yy:mm of the inquiry
	BookingsByDate map[string]BookingsPerDate `json:"bookings_by_date" bson:"bookings_by_date"`
}

type BookingsPerDate map[string]BookingsByHour // key is the date 1 .. 31

type BookingsByHour map[string][]Booking // key is the hour -> 06:00, 10:15, 21:30, etc

type Booking struct {
	ConfirmationCode string `json:"confirmation_code"`
	InquiryID        string `json:"inquiry_id"`
}

type BookingUsecase interface {
	CreateBooking(ctx context.Context, inquiryID string) error
}

type BookingRepository interface {
	GetBookings(ctx context.Context, serviceId string, variantPax int, yearMonths []string) (res []ServiceBookings, err error)
	UpsertBookingsDocument(ctx context.Context, documents []ServiceBookings) (err error)
}

func (p *ServiceBookings) Init(serviceId string, variantPax int, yearMonth string) {
	id, _ := utils.GenerateRandomString(16)
	p.ID = id
	p.ServiceID = serviceId
	p.VariantPax = variantPax
	p.YearMonth = yearMonth
	p.BookingsByDate = map[string]BookingsPerDate{}
}
