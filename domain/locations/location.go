package locations

import (
	"context"
	"mini-wallet/domain/common/response"
)

type LocationRepository interface {
	GetProvinces(ctx context.Context) ([]Location, error)
	GetCitiesByProvinceID(ctx context.Context, provinceID int) ([]City, error)
	GetDistrictByCityID(ctx context.Context, provinceID int) ([]District, error)
}

type LocationUsecase interface {
	GetProvinces(ctx context.Context) response.Response[[]Location]
	GetCitiesByProvinceID(ctx context.Context, provinceID int) response.Response[[]City]
	GetDistrictByCityID(ctx context.Context, provinceID int) response.Response[[]District]
}
