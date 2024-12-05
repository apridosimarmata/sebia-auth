package services

import (
	"context"
	"log"
	"mini-wallet/domain"
	"mini-wallet/domain/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type servicesRepository struct {
	servicesCollection *mongo.Collection
}

func NewServicesRepository(repositoryParam domain.RepositoryParam) services.ServicesRepository {
	return &servicesRepository{
		servicesCollection: repositoryParam.Mongo.Collection("services"),
	}
}

func (repository *servicesRepository) GetPublicServices(ctx context.Context, req services.GetPublicServicesRequest) ([]services.MiniServiceDTO, error) {
	options := options.Find().SetProjection(
		bson.D{
			{
				Key:   "measurement_unit",
				Value: 1,
			},
			{
				Key:   "type_id",
				Value: 1,
			},
			{
				Key:   "category_id",
				Value: 1,
			},
			{
				Key:   "measurement_string",
				Value: 1,
			},
			{
				Key:   "category_string",
				Value: 1,
			},
			{
				Key:   "category_path",
				Value: 1,
			},
			{
				Key:   "description",
				Value: 1,
			},
			{
				Key:   "photos",
				Value: 1,
			},
			{"title", 1},
			{"slug", 1},
			{"more_than_one_variant", bson.D{
				{"$gt", bson.A{
					bson.D{{"$size", "$variants"}}, // Get the size of the variants array
					1,                              // Compare if greater than 1
				}},
			}},
			{"is_event", 1},
			{"first_variant", bson.D{
				{"$arrayElemAt", bson.A{"$variants", 0}}, // Get the first element of the variants array
			}},
		},
	).SetSort(bson.D{{"name", 1}})

	filter := &bson.M{
		"category_id": req.ToMapInterface()["category_id"],
	}

	res, err := repository.servicesCollection.Find(ctx, filter, options)
	if res.Err() != nil {
		return nil, err
	}

	result := []services.MiniServiceDTO{}

	res.All(ctx, &result)

	return result, nil
}

func (repository *servicesRepository) UpdateService(ctx context.Context, entity services.ServiceEntity) (err error) {
	filter := bson.M{"slug": entity.Slug}

	result := repository.servicesCollection.FindOneAndReplace(ctx, filter, entity)
	if result.Err() != nil {
		return err
	}

	return nil
}

func (repository *servicesRepository) GetServiceBySlug(ctx context.Context, slug string) (res *services.ServiceDTO, err error) {
	filter := bson.M{"slug": slug}

	result := repository.servicesCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, err
	}

	result.Decode(&res)

	return res, nil
}

func (repository *servicesRepository) GetServiceByID(ctx context.Context, id string) (res *services.ServiceDTO, err error) {
	filter := bson.M{"id": id}

	result := repository.servicesCollection.FindOne(ctx, filter)
	if result.Err() != nil {
		return nil, err
	}

	result.Decode(&res)

	return res, nil
}

func (repository *servicesRepository) GetServices(ctx context.Context, req services.GetServicesRequest) ([]services.MiniServiceDTO, error) {
	options := options.Find().SetProjection(
		bson.D{
			{"measurement_unit", 1},
			{"measurement_string", 1},
			{"category_string", 1},
			{"type_id", 1},
			{"category_id", 1},
			{"description", 1},
			{"photos", 1},
			{"title", 1},
			{"slug", 1},
			{"more_than_one_variant", bson.D{
				{"$gt", bson.A{
					bson.D{{"$size", "$variants"}}, // Get the size of the variants array
					1,                              // Compare if greater than 1
				}},
			}},
			{"is_event", 1},
			{"first_variant", bson.D{
				{"$arrayElemAt", bson.A{"$variants", 0}}, // Get the first element of the variants array
			}},
		},
	).SetSort(bson.D{{"name", 1}})

	var filter *bson.M
	if req.BusinessID != nil {
		filter = &bson.M{
			"business_id": req.ToMapInterface()["business_id"],
		}
	}

	res, err := repository.servicesCollection.Find(ctx, filter, options)
	if res.Err() != nil {
		return nil, err
	}

	result := []services.MiniServiceDTO{}

	res.All(ctx, &result)

	return result, nil
}

func (repository *servicesRepository) InsertService(ctx context.Context, entity services.ServiceEntity) (err error) {
	_, err = repository.servicesCollection.InsertOne(ctx, entity)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

// func (repository *servicesRepository) GetServicesByUserId(ctx context.Context, userID string) (res *services.ServiceEntity, err error) {
// 	filter := bson.M{"user_id": userID}

// 	result := repository.servicesCollection.FindOne(ctx, filter)
// 	if result.Err() != nil {
// 		return nil, err
// 	}

// 	result.Decode(&res)

// 	return res, nil
// }
