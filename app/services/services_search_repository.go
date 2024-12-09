package services

import (
	"context"
	"log"
	"mini-wallet/domain"
	"mini-wallet/domain/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type servicesSearchRepository struct {
	servicesCollection *mongo.Collection
}

func NewServicesSearchRepository(repositoryParam domain.RepositoryParam) services.ServicesSearchRepository {
	return &servicesSearchRepository{
		servicesCollection: repositoryParam.Mongo.Collection("services"),
	}
}

func (repo *servicesSearchRepository) SearchServices(ctx context.Context, keyword string) (res []services.ServiceSearchResultDTO, err error) {
	searchStage := bson.D{
		{"$search", bson.D{
			{"index", "default"}, // Replace "default" with your Atlas Search index name
			{"autocomplete", bson.D{
				{"query", keyword}, // The search keyword
				{"path", "title"},  // The field to search
				{"fuzzy", bson.D{
					{"maxEdits", 1},
				},
				},
			}},
		}},
	}
	// Project only the title and slug fields
	projectStage := bson.D{
		{"$project", bson.D{
			{"title", 1},         // Include title
			{"slug", 1},          // Include slug
			{"category_path", 1}, // Include slug
			{"_id", 0},           // Exclude the _id field
		}},
	}

	limitStage := bson.D{{"$limit", 10}} // Limit results to 10 documents

	// Execute the aggregation pipeline
	cursor, err := repo.servicesCollection.Aggregate(ctx, mongo.Pipeline{searchStage, projectStage, limitStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	cursor.All(ctx, &res)
	if err := cursor.Err(); err != nil {
		log.Fatalf("Error while iterating results: %v", err)
	}

	return res, err
}
