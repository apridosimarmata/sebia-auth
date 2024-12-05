package inquiry

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/inquiry"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type inquiryRepository struct {
	inquiryCollection *mongo.Collection
}

func NewInquiryRepository(repositoryParam domain.RepositoryParam) inquiry.InquiryRepository {
	return &inquiryRepository{
		inquiryCollection: repositoryParam.Mongo.Collection("inquiries"),
	}
}

func (repo *inquiryRepository) InsertInquiry(ctx context.Context, req inquiry.InquiryEntity) (err error) {
	_, err = repo.inquiryCollection.InsertOne(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
func (repo *inquiryRepository) UpdateInquiry(ctx context.Context, req inquiry.InquiryEntity) (err error) {
	filter := bson.M{"id": req.ID}

	result := repo.inquiryCollection.FindOneAndReplace(ctx, filter, req)
	if result.Err() != nil {
		return err
	}

	return nil

}
func (repo *inquiryRepository) GetInquiryById(ctx context.Context, id string) (res *inquiry.InquiryEntity, err error) {
	filter := &bson.M{
		"id": id,
	}

	result := repo.inquiryCollection.FindOne(ctx, filter, nil)
	if result.Err() != nil {
		return nil, err
	}

	result.Decode(&res)

	return res, nil
}
