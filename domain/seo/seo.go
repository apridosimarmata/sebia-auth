package seo

import (
	"context"
	"mini-wallet/domain/common/response"
)

// entity + dto
type FooterServiceItem struct {
	Title string `json:"title" bson:"title"`
	Url   string `json:"url" bson:"url"`
}

type FooterGroupByCategoryID struct {
	CategoryId int                 `json:"category_id" bson:"category_id"`
	Items      []FooterServiceItem `json:"items"`
}

type SEOUsecase interface {
	PopulateFooterGroupForEachCategoryId(ctx context.Context) (res response.Response[string])
}

type SEORepository interface {
	UpsertFooterGroupByCategoryId(ctx context.Context, entity FooterGroupByCategoryID) error
}
