package services

import (
	"context"
	"errors"
	"mini-wallet/domain/common/response"
	"mini-wallet/utils"
)

var typePathMap = map[int]string{
	1: "activities",
	2: "events",
	3: "rentals",
	4: "opentrips",
	5: "trips",
}

var categoryStringMap = map[int]string{
	1: "Camping",
	2: "Olahraga",
	3: "Musik",
	4: "Seni",
	5: "Kegiatan di Alam",
}

var measurementStringMap = map[int]string{
	1: "s.d. selesai",
	2: "menit",
	3: "jam",
	4: "hari",
	5: "malam",
}

// export this so can be used outside this package
var MeasurementUnitRegex = map[int]string{
	1: "^$",
	2: "^(0[0-9]|1[0-9]|2[0-3]):(00|15|30|45)$",
	3: "^(0[0-9]|1[0-9]|2[0-3]):00$",
	4: "^$",
}

type ItineraryEntity struct {
	ServiceID string `bson:"service_id" json:"itinerary"`
	Items     string `json:"items" bson:"items"`
}

type ItineraryActivity struct {
	Start       string `json:"start" bson:"start"`
	End         string `json:"end" bson:"end"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

type ItineraryItem struct {
	Day      int                 `json:"day" bson:"day"`
	Activity []ItineraryActivity `json:"activities" bson:"activities"`
}

type ServiceEntity struct {
	ID              string           `bson:"id"`
	Title           string           `bson:"title"`
	Variants        []ServiceVariant `bson:"variants"`
	Description     string           `bson:"description"`
	WhatAreIncluded string           `bson:"what_are_included"`
	Photos          []string         `bson:"photos"`
	IsEvent         bool             `json:"is_event" bson:"is_event"`

	Slug string `bson:"slug"`

	CategoryID        int `bson:"category_id"`
	TypeID            int `bson:"type_id"`
	MeasurementUnitID int `bson:"measurement_unit_id"`

	TypePath          string `bson:"type_path"`
	MeasurementString string `bson:"measurement_string"`
	TypeString        string `bson:"type_string"`
	CategoryString    string `bson:"category_string"`

	BusinessID string `bson:"business_id"`
	CreatedAt  int64  `bson:"created_at"`
	UpdatedAt  int64  `bson:"updated_at"`

	EventDetails *EventDetails `json:"event_details,omitempty" bson:"event_details"`

	OpenForAffiliate     int  `json:"open_for_affiliate" bson:"open_for_affiliate"`
	OpenForAffiliateBool bool `json:"open_for_affiliate_bool" bson:"open_for_affiliate_bool"`
	AffiliateComission   int  `json:"affiliate_comission" bson:"affiliate_comission"`

	TotalScore  int `json:"-" bson:"total_score"` // never allow this value modified by client
	ReviewCount int `json:"-" bson:"review_count"`
}

type ServiceVariant struct {
	Price                int `bson:"price" json:"price"`
	Duration             int `bson:"duration" json:"duration"`
	Pax                  int `bson:"pax" json:"pax"`
	MaxReservationPerDay int `bson:"max_reservation_per_day" json:"max_reservation_per_day"`
}

type ServiceDTO struct {
	ID                *string          `json:"id" bson:"id"`
	Title             string           `json:"title" bson:"title"`
	Variants          []ServiceVariant `json:"variants" bson:"variants"`
	TypeID            int              `json:"type_id" bson:"type_id"`
	CategoryID        int              `json:"category_id" bson:"category_id"`
	Description       string           `json:"description" bson:"description"`
	WhatAreIncluded   string           `json:"what_are_included" bson:"what_are_included"`
	Photos            []string         `json:"photos" bson:"photos"`
	MeasurementUnitID int              `json:"measurement_unit_id" bson:"measurement_unit_id"`
	BusinessID        string           `json:"business_id" bson:"business_id"`
	Slug              string           `json:"slug,omitempty" bson:"slug"`
	EventDetails      *EventDetails    `json:"event_details,omitempty" bson:"event_details"`

	TypePath          string `json:"type_path" bson:"type_path"`
	TypeString        string `json:"type_string" bson:"type_string"`
	CategoryString    string `json:"category_string" bson:"category_string"`
	MeasurementString string `json:"measurement_string" bson:"measurement_string"`

	//
	OpenForAffiliate   int `json:"open_for_affiliate" bson:"open_for_affiliate"`
	AffiliateComission int `json:"affiliate_comission" bson:"affiliate_comission"`
	TotalScore         int `json:"total_score" bson:"total_score"` // never allow this value modified by client
	ReviewCount        int `json:"review_count" bson:"review_count"`
}

type EventDetails struct {
	StartDate           string `json:"event_start_date,omitempty" bson:"event_start_date"`
	EventDurationInDays int    `json:"event_duration_in_days,omitempty" bson:"event_duration_in_days"`
	EventVenueAddress   string `json:"event_venue_address" bson:"event_venue_address"`
}

type MiniServiceDTO struct {
	Title              string         `json:"title" bson:"title"`
	MoreThanOneVariant bool           `json:"more_than_one_variant" bson:"more_than_one_variant"`
	FirstVariant       ServiceVariant `json:"first_variant" bson:"first_variant"`
	Slug               string         `json:"slug" bson:"slug"`
	TypeID             int            `json:"type_id" bson:"type_id"`
	CategoryID         int            `json:"category_id" bson:"category_id"`
	Description        string         `json:"description" bson:"description"`
	Photos             []string       `json:"photos" bson:"photos"`
	TypePath           string         `json:"type_path" bson:"type_path"`
	IsEvent            bool           `json:"is_event" bson:"is_event"`
	MeasurementString  string         `json:"measurement_string" bson:"measurement_string"`
	TypeString         string         `json:"type_string" bson:"type_string"`
	CategoryString     string         `json:"category_string" bson:"category_string"`
}

func (p *ServiceDTO) Validate() error {
	err := utils.ValidateRequired(p.Title)
	if err != nil {
		return err
	}

	err = utils.ValidateRequiredInt(p.TypeID)
	if err != nil {
		return err
	}

	err = utils.ValidateRequiredInt(p.CategoryID)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.Description)
	if err != nil {
		return err
	}

	err = utils.ValidateRequiredInt(p.MeasurementUnitID)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.BusinessID)
	if err != nil {
		return err
	}

	err = utils.ValidateRequiredInt(p.OpenForAffiliate)
	if err != nil {
		return err
	}

	err = utils.ValidateRequiredIntAllowsZero(p.AffiliateComission)
	if err != nil {
		return err
	}

	if p.OpenForAffiliate == 1 && (p.AffiliateComission < 5 || p.AffiliateComission > 10) {
		return errors.New("komisi minimal 5 persen dan maksimal 10 persen")
	}

	if len(p.Photos) < 1 {
		return errors.New("tambahkan minimal 1 foto")
	}

	if len(p.Variants) < 1 {
		return errors.New("tambahkan minimal 1 varian layanan")
	}

	return nil
}

func (p *ServiceDTO) ToServiceEntity(serviceId *string) ServiceEntity {
	now, _ := utils.GetJktTime()
	slug := p.Slug
	if p.Slug == "" {
		slug = utils.GenerateSlug(p.Title)
	}

	var eventDetails *EventDetails

	if p.TypeID == 2 && p.EventDetails != nil {
		eventDetails = p.EventDetails
	}

	if serviceId == nil {
		newId := utils.GenerateUniqueId()
		serviceId = &newId
	}

	openForAffiliateBool := true
	if p.OpenForAffiliate == 2 {
		openForAffiliateBool = false
	}

	comission := 0
	if openForAffiliateBool {
		comission = p.AffiliateComission
	}

	return ServiceEntity{
		ID:                *serviceId,
		Title:             p.Title,
		Variants:          p.Variants,
		TypeID:            p.TypeID,
		CategoryID:        p.CategoryID,
		Description:       p.Description,
		WhatAreIncluded:   p.WhatAreIncluded,
		Photos:            p.Photos,
		MeasurementUnitID: p.MeasurementUnitID,
		Slug:              slug,
		BusinessID:        p.BusinessID,
		CreatedAt:         now.Unix(),
		UpdatedAt:         now.Unix(),
		TypePath:          typePathMap[p.TypeID],
		EventDetails:      eventDetails,
		IsEvent:           eventDetails != nil,
		// strings
		MeasurementString:    measurementStringMap[p.MeasurementUnitID],
		CategoryString:       categoryStringMap[p.CategoryID],
		OpenForAffiliate:     p.OpenForAffiliate,
		OpenForAffiliateBool: openForAffiliateBool,
		AffiliateComission:   comission,
	}
}

type ServicesRepository interface {
	InsertService(ctx context.Context, entity ServiceEntity) (err error)
	UpdateService(ctx context.Context, entity ServiceEntity) (err error)
	GetPublicServices(ctx context.Context, req GetPublicServicesRequest) ([]MiniServiceDTO, error)
	GetBusinessPublicServices(ctx context.Context, req GetPublicServicesRequest) ([]MiniServiceDTO, error)

	GetServices(ctx context.Context, req GetServicesRequest) ([]MiniServiceDTO, error)
	GetServiceBySlug(ctx context.Context, slug string) (res *ServiceDTO, err error)
	GetServiceByID(ctx context.Context, id string) (res *ServiceDTO, err error)
	GetServicesByCategoryID(ctx context.Context, id int) (res []ServiceEntity, err error)
}

type ServicesSearchRepository interface {
	SearchServices(ctx context.Context, keyword string) (res []ServiceSearchResultDTO, err error)
}

type ServicesUsecase interface {
	CreateService(ctx context.Context, req ServiceDTO, userID string) (res response.Response[string])
	UpdateService(ctx context.Context, req ServiceDTO, userID string) (res response.Response[string])
	GetServices(ctx context.Context, req GetServicesRequest) (res response.Response[[]MiniServiceDTO])
	GetServiceBySlug(ctx context.Context, slug string) (res response.Response[*ServiceDTO])
	GetPublicServices(ctx context.Context, req GetPublicServicesRequest) (res response.Response[[]MiniServiceDTO])
	GetBusinessPublicServices(ctx context.Context, req GetPublicServicesRequest) (res response.Response[[]MiniServiceDTO])
	SearchServicesByKeyword(ctx context.Context, keyword string) (res response.Response[[]ServiceSearchResultDTO])
}

type ServiceSearchResultDTO struct {
	Title    string `bson:"title" json:"title"`
	Slug     string `bson:"slug" json:"slug"`
	TypePath string `bson:"type_path" json:"type_path"`
}

type GetPublicServicesRequest struct {
	Page       int    `json:"page"`
	Size       int    `json:"size"`
	CategoryID int    `json:"categoryId"`
	BusinessID string `json:"businessId"`
}

func (p *GetPublicServicesRequest) Validate() error {
	err := utils.ValidateRequiredInt(p.Page)
	if err != nil {
		return err
	}

	err = utils.ValidateRequiredInt(p.Size)
	if err != nil {
		return err
	}

	// err = utils.ValidateRequiredInt(p.CategoryID)
	// if err != nil {
	// 	return err
	// }

	return nil
}

type GetServicesRequest struct {
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	BusinessID *string `json:"business_id,omitempty"`
}

func (p *GetServicesRequest) Validate() error {
	err := utils.ValidateRequiredInt(p.Page)
	if err != nil {
		return err
	}

	err = utils.ValidateRequiredInt(p.Size)
	if err != nil {
		return err
	}

	return nil
}

func (p *GetServicesRequest) ToMapInterface() map[string]interface{} {
	result := make(map[string]interface{})

	result["page"] = p.Page
	result["size"] = p.Size

	if p.BusinessID != nil {
		result["business_id"] = *p.BusinessID
	}

	return result
}

func (p *GetPublicServicesRequest) ToMapInterface() map[string]interface{} {
	result := make(map[string]interface{})

	result["page"] = p.Page
	result["size"] = p.Size
	result["category_id"] = p.CategoryID
	result["business_id"] = p.BusinessID

	return result
}
