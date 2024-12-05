package location

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/locations"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type locationRepository struct {
	locationsCollection *mongo.Collection
}

func NewLocationRepository(repositoryParam domain.RepositoryParam) locations.LocationRepository {
	return &locationRepository{
		locationsCollection: repositoryParam.Mongo.Collection("locations"),
	}

}

func (repository *locationRepository) GetProvinces(ctx context.Context) (provinces []locations.Location, err error) {
	options := options.Find().SetProjection(
		bson.D{
			{"province_id", 1},
			{"name", 1},
		},
	).SetSort(bson.D{{"name", 1}})

	res, err := repository.locationsCollection.Find(ctx, bson.D{}, options)
	if res.Err() != nil {
		return nil, err
	}

	res.All(ctx, &provinces)
	return provinces, nil
}

func (repository *locationRepository) GetCitiesByProvinceID(ctx context.Context, provinceID int) (cities []locations.City, err error) {
	options := options.FindOne().SetProjection(
		bson.D{
			{Key: "cities", Value: 1},
		},
	)

	filter := bson.M{"province_id": provinceID}

	res := repository.locationsCollection.FindOne(ctx, filter, options)
	if res.Err() != nil {
		return nil, err
	}

	location := locations.Location{}
	res.Decode(&location)

	for _, p := range location.Cities {
		cities = append(cities, p)
	}

	return cities, nil
}

func (repository *locationRepository) GetDistrictByCityID(ctx context.Context, provinceID int) ([]locations.District, error) {

	return nil, nil
}
