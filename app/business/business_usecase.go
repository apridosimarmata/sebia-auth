package business

import (
	"context"
	"errors"
	"mini-wallet/domain"
	"mini-wallet/domain/business"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/locations"
)

type businessUsecase struct {
	businessRepository business.BusinessRepository
	locationRepository locations.LocationRepository
}

func NewBusinessUsecase(repositories domain.Repositories) business.BusinessUsecase {
	return &businessUsecase{
		businessRepository: repositories.BusinessRepository,
		locationRepository: repositories.LocationRepository,
	}
}

func (uc *businessUsecase) GetBusinessByID(ctx context.Context, id string) (res response.Response[*business.PublicBusinessDTO]) {
	businessEntity, err := uc.businessRepository.GetBusinessById(ctx, id)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if businessEntity == nil {
		res.NotFound("business not found", nil)
		return
	}

	provinces, err := uc.locationRepository.GetProvinces(ctx)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	province := locations.Location{}
	for _, _province := range provinces {
		if _province.ProvinceID == int(businessEntity.ProvinceID) {
			province = _province
			break
		}
	}

	if province.ProvinceID == 0 {
		res.InternalServerError(errors.New("province not found").Error())
		return
	}

	cities, err := uc.locationRepository.GetCitiesByProvinceID(ctx, int(businessEntity.ProvinceID))
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	city := ""
	// todo, optimize
	for _, _city := range cities {
		if _city.ID == int(businessEntity.CityID) {
			city = _city.Name
			break
		}
	}

	if city == "" {
		res.NotFound("city not found", nil)
		return
	}

	// districts, err := uc.locationRepository.GetDistrictByCityID(ctx, int(businessEntity.CityID))
	// if err != nil {
	// 	res.InternalServerError(err.Error())
	// 	return
	// }

	// district := ""
	// for _, _district := range districts {
	// 	if _district.ID == int(businessEntity.DistrictID) {
	// 		district = _district.Name
	// 		break
	// 	}
	// }

	// if district == "" {
	// 	res.NotFound("district not found", nil)
	// 	return
	// }

	result := businessEntity.ToPublicBusinessDTO(province.Name, city)

	res.Success(&result)
	return

}

func (uc *businessUsecase) GetBusinessByHandle(ctx context.Context, slug string) (res response.Response[*business.PublicBusinessDTO]) {

	businessEntity, err := uc.businessRepository.GetBusinessByHandle(ctx, slug)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if businessEntity == nil {
		res.NotFound("business not found", nil)
		return
	}

	provinces, err := uc.locationRepository.GetProvinces(ctx)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	province := locations.Location{}
	for _, _province := range provinces {
		if _province.ProvinceID == int(businessEntity.ProvinceID) {
			province = _province
			break
		}
	}

	if province.ProvinceID == 0 {
		res.InternalServerError(errors.New("province not found").Error())
		return
	}

	cities, err := uc.locationRepository.GetCitiesByProvinceID(ctx, int(businessEntity.ProvinceID))
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	city := ""
	// todo, optimize
	for _, _city := range cities {
		if _city.ID == int(businessEntity.CityID) {
			city = _city.Name
			break
		}
	}

	if city == "" {
		res.NotFound("city not found", nil)
		return
	}

	// districts, err := uc.locationRepository.GetDistrictByCityID(ctx, int(businessEntity.CityID))
	// if err != nil {
	// 	res.InternalServerError(err.Error())
	// 	return
	// }

	// district := ""
	// for _, _district := range districts {
	// 	if _district.ID == int(businessEntity.DistrictID) {
	// 		district = _district.Name
	// 		break
	// 	}
	// }

	// if district == "" {
	// 	res.NotFound("district not found", nil)
	// 	return
	// }

	result := businessEntity.ToPublicBusinessDTO(province.Name, city)

	res.Success(&result)
	return
}

func (uc *businessUsecase) CreateBusiness(ctx context.Context, req business.BusinessCreationDTO) (res response.Response[string]) {
	err := req.Validate()
	if err != nil {
		res.BadRequest(err.Error(), nil)
		return
	}

	businessEntity := req.ToBusinessEntity()
	err = uc.businessRepository.InsertBusiness(ctx, businessEntity)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.SuccessWithMessage("Permintaan terkirim! 🎉")
	return
}

func (uc *businessUsecase) GetUserBusinessStatus(ctx context.Context, userID string) (res response.Response[*string]) {
	userBusiness, err := uc.businessRepository.GetBusinessByUserId(ctx, userID)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if userBusiness == nil {
		res.Success(nil)
		return
	}

	if userBusiness.Status == 1 {
		pending := "pending"
		res.Success(&pending)
		return
	}

	res.Success(&userBusiness.ID)
	return
}
