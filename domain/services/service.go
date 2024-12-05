package services

import (
	"context"
	"errors"
	"mini-wallet/domain/common/response"
	"mini-wallet/utils"
)

var categoryPathMap = map[int]string{
	1: "activities",
	2: "events",
	3: "rentals",
	4: "opentrips",
}

var categoryStringMap = map[int]string{
	1: "Aktivitas Wisata",
	2: "Event",
	3: "Rental",
	4: "Open trip",
}

var measurementStringMap = map[int]string{
	1: "s.d. selesai",
	2: "menit",
	3: "jam",
	4: "hari",
}

// export this so can be used outside this package
var MeasurementUnitRegex = map[int]string{
	1: "^$",
	2: "^(0[0-9]|1[0-9]|2[0-3]):(00|15|30|45)$",
	3: "^(0[0-9]|1[0-9]|2[0-3]):00$",
	4: "^$",
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

	CategoryPath      string `bson:"category_path"`
	MeasurementString string `bson:"measurement_string"`
	TypeString        string `bson:"type_string"`
	CategoryString    string `bson:"category_string"`

	BusinessID string `bson:"business_id"`
	CreatedAt  int64  `bson:"created_at"`
	UpdatedAt  int64  `bson:"updated_at"`

	EventDetails *EventDetails `json:"event_details,omitempty" bson:"event_details"`
}

type ServiceVariant struct {
	Price                int `bson:"price" json:"price"`
	Duration             int `bson:"duration" json:"duration"`
	Pax                  int `bson:"pax" json:"pax"`
	MaxReservationPerDay int `bson:"max_reservation_per_day" json:"max_reservation_per_day"`
}

type ServiceDTO struct {
	ID                *string          ` json:"string" bson:"id"`
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
	CategoryPath      string           `json:"category_path" bson:"category_path"`
	EventDetails      *EventDetails    `json:"event_details,omitempty" bson:"event_details"`

	TypeString        string `json:"type_string" bson:"type_string"`
	CategoryString    string `json:"category_string" bson:"category_string"`
	MeasurementString string `json:"measurement_string" bson:"measurement_string"`
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
	CategoryPath       string         `json:"category_path" bson:"category_path"`
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
		CategoryPath:      categoryPathMap[p.CategoryID],
		EventDetails:      eventDetails,
		IsEvent:           eventDetails != nil,
		// strings
		MeasurementString: measurementStringMap[p.MeasurementUnitID],
		CategoryString:    categoryStringMap[p.CategoryID],
	}
}

type ServicesRepository interface {
	InsertService(ctx context.Context, entity ServiceEntity) (err error)
	UpdateService(ctx context.Context, entity ServiceEntity) (err error)
	GetPublicServices(ctx context.Context, req GetPublicServicesRequest) ([]MiniServiceDTO, error)
	GetServices(ctx context.Context, req GetServicesRequest) ([]MiniServiceDTO, error)
	GetServiceBySlug(ctx context.Context, slug string) (res *ServiceDTO, err error)
	GetServiceByID(ctx context.Context, id string) (res *ServiceDTO, err error)
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
	SearchServicesByKeyword(ctx context.Context, keyword string) (res response.Response[[]ServiceSearchResultDTO])
}

type ServiceSearchResultDTO struct {
	Title string `bson:"title"`
	Slug  string `bson:"slug"`
}

type GetPublicServicesRequest struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	CategoryID int `json:"categoryId"`
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

	err = utils.ValidateRequiredInt(p.CategoryID)
	if err != nil {
		return err
	}

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

	return result
}
