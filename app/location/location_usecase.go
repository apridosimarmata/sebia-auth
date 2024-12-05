package location

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/locations"
)

type locationUsecase struct {
	locationRepository locations.LocationRepository
}

func NewLocationUsecase(repositories domain.Repositories) locations.LocationUsecase {
	return &locationUsecase{
		locationRepository: repositories.LocationRepository,
	}
}

func (usecase *locationUsecase) GetProvinces(ctx context.Context) (res response.Response[[]locations.Location]) {
	provinces, err := usecase.locationRepository.GetProvinces(ctx)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.Success(provinces)
	return
}

func (usecase *locationUsecase) GetCitiesByProvinceID(ctx context.Context, provinceID int) (res response.Response[[]locations.City]) {
	cities, err := usecase.locationRepository.GetCitiesByProvinceID(ctx, provinceID)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.Success(cities)
	return
}

func (usecase *locationUsecase) GetDistrictByCityID(ctx context.Context, provinceID int) (res response.Response[[]locations.District]) {

	return
}
