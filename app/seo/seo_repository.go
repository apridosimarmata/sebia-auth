package seo

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/seo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type seoRepository struct {
	seoCollection *mongo.Collection
}

func NewSEORepository(repositoryParam domain.RepositoryParam) seo.SEORepository {
	return &seoRepository{
		seoCollection: repositoryParam.Mongo.Collection("seo"),
	}
}

func (repository *seoRepository) UpsertFooterGroupByCategoryId(ctx context.Context, entity seo.FooterGroupByCategoryID) error {
	filter := bson.M{
		"category_id": entity.CategoryId,
	}

	update := bson.M{"$set": bson.M{
		"items": entity.Items,
	}}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	result := repository.seoCollection.FindOneAndUpdate(ctx, filter, update, opts)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
