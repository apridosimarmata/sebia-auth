package review

import (
	"context"
	"log"
	"mini-wallet/domain"
	"mini-wallet/domain/review"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type reviewRepository struct {
	reviewsCollection *mongo.Collection
}

func NewReviewRepository(repositoryParam domain.RepositoryParam) review.ReviewRepository {
	return &reviewRepository{
		reviewsCollection: repositoryParam.Mongo.Collection("reviews"),
	}
}

func (repo *reviewRepository) InsertReview(ctx context.Context, req review.ReviewEntity) (err error) {
	_, err = repo.reviewsCollection.InsertOne(ctx, req)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil

}

func (repo *reviewRepository) GetServiceTopReview(ctx context.Context, serviceId string) (res *review.ReviewEntity, err error) {
	filter := bson.M{
		"service_id": serviceId}
	opts := options.FindOne().SetSort(bson.D{
		{
			Key:   "score",
			Value: -1,
		},
	})

	result := repo.reviewsCollection.FindOne(ctx, filter, opts)
	if result.Err() != nil {
		return nil, err
	}

	result.Decode(&res)

	return res, nil
}
