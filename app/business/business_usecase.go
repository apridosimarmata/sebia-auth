package business

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/business"
	"mini-wallet/domain/common/response"
)

type businessUsecase struct {
	businessRepository business.BusinessRepository
}

func NewBusinessUsecase(repositories domain.Repositories) business.BusinessUsecase {
	return &businessUsecase{
		businessRepository: repositories.BusinessRepository,
	}
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

	res.Success(&userBusiness.ID)
	return
}
