package inquiry

import (
	"context"
	"errors"
	"mini-wallet/domain/business"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/services"
	"mini-wallet/utils"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var inquiryStatusMap = map[int]string{
	0: "Menunggu Pembayaran",
	2: "Sudah Dibayar",
	3: "Kode Booking Terbit",
}

type MaskedInquiryContactDTO struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

type InquiryDTO struct {
	ID          *string `json:"id,omitempty"`
	ServiceSlug string  `json:"service_slug"`

	SelectedDates     []string                `json:"selected_dates"`
	SelectedVariantID string                  `json:"selected_variant_id"`
	SelectedVariant   services.ServiceVariant `json:"selected_variant_details"`

	SelectedHour string `json:"selected_hour,omitempty"`

	FullName         string  `json:"full_name,omitempty"`
	PhoneNumber      string  `json:"phone_number,omitempty"`
	Email            string  `json:"email,omitempty"`
	ConfirmationCode *string `json:"confirmation_code,omitempty"`

	UserID *string `json:"user_id"`
	Status int     `json:"status"`

	// for response details in booking page
	StatusString             string                `json:"status_string"`
	ServiceDetails           InquiryServiceDetails `json:"service_details"`
	CreatedAt                string                `json:"created_at"`
	ServiceMeasurementUnitID int                   `json:"service_measurement_unit_id"`
	ServiceMeasurementUnit   string                `json:"service_measurement_unit"`
	ReviewAvailable          bool                  `json:"review_available"`
	ReviewMade               bool                  `json:"review_made"`

	// If status == 3
}
type InquiryServiceDetails struct {
	Photo    string `json:"photo"`
	Slug     string `json:"slug"`
	HostName string `json:"host_name"`
	HostSlug string `json:"host_slug"`
	Title    string `json:"title"`
	Path     string `json:"path"`
}

type InquiryEntity struct {
	ID string `bson:"id"`

	ServiceID         string                  `bson:"service_id"`
	SelectedDates     []string                `bson:"selected_dates"`
	SelectedVariantID string                  `bson:"selected_variant_id"`
	SelectedVariant   services.ServiceVariant `bson:"selected_variant_details"`
	SelectedHour      string                  `bson:"selected_hour"`

	FullName    string `bson:"full_name"`
	PhoneNumber string `bson:"phone_number"`
	Email       string `bson:"email"`

	UserID *string `bson:"user_id"`
	Status int     `bson:"status"`

	CreatedDate              string `bson:"created_date"`
	UpdatedDate              string `bson:"updated_date"`
	TotalPayment             int    `bson:"total_payment"`
	ServiceMeasurementUnitID int    `bson:"service_measurement_unit_id"`
	ServiceMeasurementUnit   string `bson:"service_measurement_unit"`

	ReviewMade       bool    `bson:"review_made"`
	ConfirmationCode *string `bson:"confirmation_code"`
}

func (p *InquiryEntity) ToInquiryDetailsResponse(service services.ServiceEntity, host business.BusinessEntity) InquiryDTO {
	statusString := inquiryStatusMap[p.Status]

	reviewAvailable := false
	lastSelectedDate := p.SelectedDates[len(p.SelectedDates)-1]
	now, _ := utils.GetJktTime()
	// Define the layout that matches the input string
	layout := "2006/1/2"

	// Parse the string into a time.Time object
	parsedLastSelectedDate, _ := time.Parse(layout, lastSelectedDate)
	parsedLastSelectedDate = parsedLastSelectedDate.Add(time.Hour * 24)

	if p.Status == 3 && now.After(parsedLastSelectedDate) && !p.ReviewMade {
		reviewAvailable = true
	}

	return InquiryDTO{
		ID:              &p.ID,
		ServiceSlug:     service.Slug,
		SelectedDates:   p.SelectedDates,
		SelectedVariant: p.SelectedVariant,

		FullName:    p.FullName,
		PhoneNumber: p.PhoneNumber,
		Email:       p.Email,

		Status:       p.Status,
		StatusString: statusString,
		ServiceDetails: InquiryServiceDetails{
			Photo:    service.Photos[0],
			HostName: host.Name,
			Slug:     service.Slug,
			Title:    service.Title,
			Path:     service.TypePath,
		},
		CreatedAt:                p.CreatedDate,
		ServiceMeasurementUnit:   p.ServiceMeasurementUnit,
		SelectedHour:             p.SelectedHour,
		ServiceMeasurementUnitID: p.ServiceMeasurementUnitID,
		UserID:                   p.UserID,
		ReviewAvailable:          reviewAvailable,
		ReviewMade:               p.ReviewMade,
		ConfirmationCode:         p.ConfirmationCode,
	}
}

func (p *InquiryDTO) Validate() error {
	err := utils.ValidateRequired(p.ServiceSlug)
	if err != nil {
		return err
	}

	selectedDates := []interface{}{}
	selectedDates = append(selectedDates, selectedDates...)

	err = utils.ValidateRequiredSlice(selectedDates)
	if err != nil {
		return err
	}

	newPhoneNumber, err := utils.ValidatePhoneNumber(p.PhoneNumber)
	if err != nil {
		return err
	}
	p.PhoneNumber = *newPhoneNumber

	return nil
}

func (p *InquiryDTO) ToInquiryEntity(service services.ServiceEntity, selectedVariant services.ServiceVariant, total int) (res *InquiryEntity, err error) {
	now, _ := utils.GetJktTime()

	id := p.ID
	if id == nil || *id == "" {
		newId := utils.GenerateUniqueId()
		id = &newId
	}

	// validate supplied hour, if any
	re := regexp.MustCompile(services.MeasurementUnitRegex[service.MeasurementUnitID])
	match := re.MatchString(p.SelectedHour)
	if !match {
		return nil, errors.New("Jam yang dipilih tidak tersedia")
	}

	if p.SelectedHour == "" {
		p.SelectedHour = "00:00"
	}

	return &InquiryEntity{
		ID:                     *id,
		ServiceID:              service.ID,
		SelectedDates:          p.SelectedDates,
		SelectedVariantID:      p.SelectedVariantID,
		UserID:                 p.UserID,
		Status:                 p.Status, //todo
		FullName:               p.FullName,
		PhoneNumber:            p.PhoneNumber,
		Email:                  p.Email,
		TotalPayment:           total,
		CreatedDate:            now.Format(time.RFC3339),
		UpdatedDate:            now.Format(time.RFC3339),
		SelectedVariant:        selectedVariant,
		ServiceMeasurementUnit: service.MeasurementString,
		SelectedHour:           p.SelectedHour,
	}, nil
}

type InquiryUsecase interface {
	CreateInquiry(ctx context.Context, req InquiryDTO) (res response.Response[string])
	GetInquiry(ctx context.Context, id string) (res response.Response[InquiryDTO])
	GetInquiryMaskedContact(ctx context.Context, id string) (res response.Response[MaskedInquiryContactDTO])
}

type InquiryRepository interface {
	InsertInquiry(ctx context.Context, req InquiryEntity) (err error)
	UpdateInquiry(ctx context.Context, req InquiryEntity) (err error)
	UpdateInquiryWithTx(ctx context.Context, txSession *mongo.SessionContext, req InquiryEntity) (err error)

	GetInquiryById(ctx context.Context, id string) (res *InquiryEntity, err error)
}
