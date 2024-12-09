package business

import (
	"context"
	"log"
	"mini-wallet/domain"
	"mini-wallet/domain/business"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type businessRepository struct {
	businessCollection *mongo.Collection
}

func NewBusinessRepository(repositoryParam domain.RepositoryParam) business.BusinessRepository {
	return &businessRepository{
		businessCollection: repositoryParam.Mongo.Collection("business"),
	}
}

func (repository *businessRepository) GetBusinessByHandle(ctx context.Context, handle string) (res *business.BusinessEntity, err error) {
	filter := bson.M{"handle": handle}

	result := repository.businessCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, err
	}

	result.Decode(&res)

	return res, nil
}

func (repository *businessRepository) InsertBusiness(ctx context.Context, entity business.BusinessEntity) (err error) {
	_, err = repository.businessCollection.InsertOne(ctx, entity)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func (repository *businessRepository) GetBusinessByUserId(ctx context.Context, userID string) (res *business.BusinessEntity, err error) {
	filter := bson.M{"user_id": userID}

	result := repository.businessCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, err
	}

	result.Decode(&res)

	return res, nil
}

func (repository *businessRepository) GetBusinessById(ctx context.Context, id string) (res *business.BusinessEntity, err error) {
	filter := bson.M{"id": id}

	result := repository.businessCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, err
	}

	result.Decode(&res)

	return res, nil
}
