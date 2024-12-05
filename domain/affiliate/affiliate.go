package affiliate

import (
	"context"
	"mini-wallet/domain/common/response"
	"mini-wallet/utils"
)

type AffiliateAppicationDTO struct {
	InstagramUsername *string `json:"instagram_username"`
	TiktokUsername    *string `json:"tiktok_username"`
	Age               int     `json:"age"`
	GenderID          int     `json:"gender_id"`
	Address           string  `json:"address"`
	ProvinceID        int     `json:"province_id"`
	CityID            int     `json:"city_id"`
	DistrictID        int     `json:"district_id"`
	UserID            string  `json:"user_id"`
}

func (p *AffiliateAppicationDTO) Validate() error {
	instagramUsername := ""
	if p.InstagramUsername != nil {
		instagramUsername = *p.InstagramUsername
	}

	tiktokUsername := ""
	if p.TiktokUsername != nil {
		tiktokUsername = *p.TiktokUsername
	}

	err := utils.ValidateRequired(instagramUsername)
	if err != nil {
		err := utils.ValidateRequired(tiktokUsername)
		if err != nil {
			return err
		}
	}

	err = utils.ValidateRequiredInt(p.Age)
	if err != nil {
		return err
	}

	err = utils.ValidateRequiredInt(p.GenderID)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.Address)
	if err != nil {
		return err
	}

	err = utils.ValidateRequiredInt(p.ProvinceID)
	if err != nil {
		return err
	}

	err = utils.ValidateRequiredInt(p.CityID)
	if err != nil {
		return err
	}

	err = utils.ValidateRequiredInt(p.DistrictID)
	if err != nil {
		return err
	}

	err = utils.ValidateRequired(p.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (p *AffiliateAppicationDTO) ToAffiliateEntity() AffiliateEntity {
	now, _ := utils.GetJktTime()

	return AffiliateEntity{
		InstagramUsername: p.InstagramUsername,
		TiktokUsername:    p.TiktokUsername,
		Age:               p.Age,
		GenderID:          p.GenderID,
		Address:           p.Address,
		ProvinceID:        p.ProvinceID,
		CityID:            p.CityID,
		DistrictID:        p.DistrictID,
		UserID:            p.UserID,
		CreatedAt:         now.Unix(),
		UpdatedAt:         now.Unix(),
		Status:            1,
	}
}

type AffiliateEntity struct {
	InstagramUsername *string `bson:"instagram_username"`
	TiktokUsername    *string `bson:"tiktok_username"`
	Age               int     `bson:"age"`
	GenderID          int     `bson:"gender_id"`
	Address           string  `bson:"address"`
	ProvinceID        int     `bson:"province_id"`
	CityID            int     `bson:"city_id"`
	DistrictID        int     `bson:"district_id"`
	UserID            string  `bson:"user_id"`

	Status    int   `bson:"status"`
	CreatedAt int64 `bson:"created_at"`
	UpdatedAt int64 `bson:"updated_at"`
}

type AffiliateUsecase interface {
	ApplyForAffiliate(ctx context.Context, req AffiliateAppicationDTO) response.Response[string]
	GetUserAffiliateStatus(ctx context.Context, userID string) (res response.Response[int])
}

type AffiliateRepository interface {
	InsertAffiliate(ctx context.Context, entity AffiliateEntity) (err error)
	GetAffiliateByUserId(ctx context.Context, userID string) (res *AffiliateEntity, err error)
}
