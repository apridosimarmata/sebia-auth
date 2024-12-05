package affiliate

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/affiliate"
	"mini-wallet/domain/common/response"
)

type affilatesUsecase struct {
	affiliateRepository affiliate.AffiliateRepository
}

func NewAffiliatesUsecase(repositories domain.Repositories) affiliate.AffiliateUsecase {
	return &affilatesUsecase{
		affiliateRepository: repositories.AffiliateRepository,
	}
}

func (usecase *affilatesUsecase) ApplyForAffiliate(ctx context.Context, req affiliate.AffiliateAppicationDTO) (res response.Response[string]) {
	affiliateEntity := req.ToAffiliateEntity()
	err := usecase.affiliateRepository.InsertAffiliate(ctx, affiliateEntity)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	res.SuccessWithMessage("Permintaan terkirim! 🎉")
	return
}

func (usecase *affilatesUsecase) GetUserAffiliateStatus(ctx context.Context, userID string) (res response.Response[int]) {
	userAffiliates, err := usecase.affiliateRepository.GetAffiliateByUserId(ctx, userID)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if userAffiliates == nil {
		res.Success(0)
		return
	}

	res.Success(int(userAffiliates.Status))
	return
}
