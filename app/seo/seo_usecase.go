package seo

import (
	"context"
	"fmt"
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/seo"
	"mini-wallet/domain/services"
)

type seoUsecase struct {
	SEORepository     seo.SEORepository
	serviceRepository services.ServicesRepository
}

func NewSEOUsecase(repositories domain.Repositories) seo.SEOUsecase {
	return &seoUsecase{
		SEORepository:     repositories.SEORepository,
		serviceRepository: repositories.ServicesRepository,
	}
}

func (usecase *seoUsecase) GetItemsByCategoryId(ctx context.Context, id int) (res response.Response[[]seo.FooterServiceItem]) {

	group, err := usecase.SEORepository.GetGroupByCategoryId(ctx, id)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	if group == nil {
		res.Success([]seo.FooterServiceItem{})
		return
	}

	res.Success(group.Items)
	return
}

func (usecase *seoUsecase) PopulateFooterGroupForEachCategoryId(ctx context.Context) (res response.Response[string]) {

	for _, categoryId := range []int{1, 2, 3, 4, 5} {
		services, err := usecase.serviceRepository.GetServicesByCategoryID(ctx, categoryId)
		if err != nil {
			res.InternalServerError(err.Error())
			return
		}

		group := seo.FooterGroupByCategoryID{
			CategoryId: categoryId,
		}
		for _, service := range services {
			group.Items = append(group.Items, seo.FooterServiceItem{
				Title: service.Title,
				Url:   fmt.Sprintf("/%s/%s", service.TypePath, service.Slug),
			})
		}

		err = usecase.SEORepository.UpsertFooterGroupByCategoryId(ctx, group)
		if err != nil {
			res.InternalServerError(err.Error())
			return
		}
	}

	res.Success("groups populated!")

	return
}
