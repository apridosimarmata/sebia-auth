package services

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/business"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/services"
)

type servicesUsecase struct {
	servicesRepository       services.ServicesRepository
	servicesSearchRepository services.ServicesSearchRepository

	BusinessRepository business.BusinessRepository
}

func NewServicesUsecase(repositories domain.Repositories) services.ServicesUsecase {
	return &servicesUsecase{
		servicesRepository:       repositories.ServicesRepository,
		BusinessRepository:       repositories.BusinessRepository,
		servicesSearchRepository: repositories.ServicesSearchRepository,
	}
}

func (usecase *servicesUsecase) GetBusinessPublicServices(ctx context.Context, req services.GetPublicServicesRequest) (res response.Response[[]services.MiniServiceDTO]) {
	result, err := usecase.servicesRepository.GetBusinessPublicServices(ctx, req)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.Success(result)
	return
}

func (usecase *servicesUsecase) SearchServicesByKeyword(ctx context.Context, keyword string) (res response.Response[[]services.ServiceSearchResultDTO]) {
	result, err := usecase.servicesSearchRepository.SearchServices(ctx, keyword)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.Success(result)
	return
}

func (usecase *servicesUsecase) UpdateService(ctx context.Context, req services.ServiceDTO, userID string) (res response.Response[string]) {
	serviceEntity, err := usecase.servicesRepository.GetServiceBySlug(ctx, req.Slug)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if serviceEntity == nil {
		res.Forbidden("service not found", nil)
		return
	}

	if serviceEntity.BusinessID != req.BusinessID {
		res.Unauthorized("has no right on this service")
		return
	}

	updatedServiceEntity := req.ToServiceEntity(serviceEntity.ID)

	err = usecase.servicesRepository.UpdateService(ctx, updatedServiceEntity)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.Success(updatedServiceEntity.Slug)
	return
}

func (usecase *servicesUsecase) GetServiceBySlug(ctx context.Context, slug string) (res response.Response[*services.ServiceDTO]) {
	result, err := usecase.servicesRepository.GetServiceBySlug(ctx, slug)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if result == nil {
		res.NotFound("service not found", nil)
		return
	}

	res.Success(result)
	return
}

func (usecase *servicesUsecase) GetPublicServices(ctx context.Context, req services.GetPublicServicesRequest) (res response.Response[[]services.MiniServiceDTO]) {
	result, err := usecase.servicesRepository.GetPublicServices(ctx, req)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.Success(result)
	return
}

func (usecase *servicesUsecase) GetServices(ctx context.Context, req services.GetServicesRequest) (res response.Response[[]services.MiniServiceDTO]) {
	result, err := usecase.servicesRepository.GetServices(ctx, req)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.Success(result)
	return
}

func (usecase *servicesUsecase) CreateService(ctx context.Context, req services.ServiceDTO, userID string) (res response.Response[string]) {
	businessEntity, err := usecase.BusinessRepository.GetBusinessByUserId(ctx, userID)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if businessEntity == nil {
		res.Forbidden("user has no business", nil)
		return
	}

	if businessEntity.ID != req.BusinessID {
		res.Unauthorized("has no right on this business")
		return
	}

	serviceEntity := req.ToServiceEntity(nil)

	err = usecase.servicesRepository.InsertService(ctx, serviceEntity)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.Success(serviceEntity.Slug)
	return
}
