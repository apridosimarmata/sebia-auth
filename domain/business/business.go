package business

import (
	"context"
	"crypto/rand"
	"errors"
	"mini-wallet/domain/common/response"
	"mini-wallet/utils"
	"time"

	"github.com/oklog/ulid/v2"
)

type BusinessEntity struct {
	ID                 string `bson:"id"`
	Name               string `bson:"name"`
	Handle             string `bson:"handle"`
	PhoneNumber        string `bson:"phone_number"`
	Address            string `bson:"address"`
	RequirementFileUrl string `bson:"requirement_file_url"`
	CityID             int64  `bson:"city_id"`
	ProvinceID         int64  `bson:"province_id"`
	DistrictID         int64  `bson:"district_id"`

	UserID    string `bson:"user_id"`
	Status    int64  `bson:"status"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
}

func (p *BusinessCreationDTO) ToBusinessEntity() BusinessEntity {
	now, _ := utils.GetJktTime()
	// Create an entropy source for random number generation TODO
	entropy := ulid.Monotonic(rand.Reader, 0)

	// Generate a ULID
	t := time.Now().UTC()
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	handle := utils.GenerateSlug(p.Name)

	// Print the generated ULID
	// fmt.Printf("Generated ULID: %s\n", id)
	return BusinessEntity{
		ID:                 id.String(),
		Name:               p.Name,
		Handle:             handle,
		PhoneNumber:        p.PhoneNumber,
		Address:            p.Address,
		RequirementFileUrl: p.RequirementFileUrl,
		CityID:             p.CityID,
		ProvinceID:         p.ProvinceID,
		DistrictID:         p.DistrictID,
		UserID:             p.UserID,
		Status:             1,
		CreatedAt:          now.Unix(),
		UpdatedAt:          now.Unix(),
	}
}

func (p *BusinessCreationDTO) Validate() error {
	err := utils.ValidateRequired(p.Name)
	if err != nil {
		return err
	}

	err = utils.ValidateFullName(p.Name)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.Name)
	if err != nil {
		return err
	}

	newPhoneNumber, err := utils.ValidatePhoneNumber(p.PhoneNumber)
	if err != nil {
		return err
	}
	p.PhoneNumber = *newPhoneNumber

	err = utils.ValidateRequired(p.Address)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.UserID)
	if err != nil {
		return err
	}

	if p.CityID == 0 || p.ProvinceID == 0 {
		return errors.New("alamat tidak lengkap")
	}

	return nil
}

type BusinessCreationDTO struct {
	Name               string `json:"name"`
	PhoneNumber        string `json:"phone_number"`
	Address            string `json:"address"`
	RequirementFileUrl string `json:"requirement_file_url"`
	CityID             int64  `json:"city_id"`
	ProvinceID         int64  `json:"province_id"`
	DistrictID         int64  `json:"district_id"`
	UserID             string `json:"user_id"`
	// UserID derived from access_token
}

type BusinessDTO struct {
	Name               string `json:"name"`
	PhoneNumber        string `json:"phone_number"`
	Address            string `json:"address"`
	RequirementFileUrl string `json:"requirement_file_url"`
	CityID             int64  `json:"city_id"`
	ProvinceID         int64  `json:"province_id"`
	Status             string `json:"status"`
	DistrictID         int64  `json:"district_id"`
	// UserID derived from access_token
}

type PublicBusinessDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Handle string `json:"handle"`
	// PhoneNumber string `json:"phone_number"`
	// Address     string `json:"address"`
	// CityID     int64 `json:"city_id"`
	// ProvinceID int64 `json:"province_id"`
	// DistrictID int64 `json:"district_id"`

	City     string `json:"city"`
	Province string `json:"province"`
	// District string `json:"district"`

	// UserID derived from access_token
}

func (p *BusinessEntity) ToPublicBusinessDTO(province string, city string) PublicBusinessDTO {
	return PublicBusinessDTO{
		ID:       p.ID,
		Name:     p.Name,
		Handle:   p.Handle,
		Province: province,
		City:     city,
		// District: district,
	}
}

type BusinessUsecase interface {
	CreateBusiness(ctx context.Context, business BusinessCreationDTO) (res response.Response[string])
	GetUserBusinessStatus(ctx context.Context, userID string) (res response.Response[*string])
	GetBusinessByHandle(ctx context.Context, slug string) (res response.Response[*PublicBusinessDTO])
	GetBusinessByID(ctx context.Context, id string) (res response.Response[*PublicBusinessDTO])
}

type BusinessRepository interface {
	InsertBusiness(ctx context.Context, entity BusinessEntity) (err error)
	GetBusinessByUserId(ctx context.Context, userID string) (res *BusinessEntity, err error)
	GetBusinessById(ctx context.Context, id string) (res *BusinessEntity, err error)
	GetBusinessByHandle(ctx context.Context, slug string) (res *BusinessEntity, err error)
}
