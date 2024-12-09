package review

import (
	"context"
	"log"
	"mini-wallet/domain"
	"mini-wallet/domain/review"

	"go.mongodb.org/mongo-driver/mongo"
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
