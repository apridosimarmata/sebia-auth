package affiliate

import (
	"context"
	"log"
	"mini-wallet/domain"
	"mini-wallet/domain/affiliate"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type affiliatesRepository struct {
	affiliatesCollection *mongo.Collection
}

func NewAffiliatesRepository(repositoryParam domain.RepositoryParam) affiliate.AffiliateRepository {
	return &affiliatesRepository{
		affiliatesCollection: repositoryParam.Mongo.Collection("affiliate"),
	}
}

func (repository *affiliatesRepository) InsertAffiliate(ctx context.Context, entity affiliate.AffiliateEntity) (err error) {
	_, err = repository.affiliatesCollection.InsertOne(ctx, entity)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func (repository *affiliatesRepository) GetAffiliateByUserId(ctx context.Context, userID string) (res *affiliate.AffiliateEntity, err error) {
	filter := bson.M{"user_id": userID}

	result := repository.affiliatesCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, err
	}

	result.Decode(&res)

	return res, nil
}
