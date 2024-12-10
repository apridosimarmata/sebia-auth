package infrastructure

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//

func GetMongoDatabase(ctx context.Context) (db *mongo.Database, err error) {
	uri := "mongodb+srv://isimarmata09:k2Bzkk90Ym1fbGhK@gerbangtobamaster.lvmfl.mongodb.net/?retryWrites=true&w=majority&appName=gerbangTobaMaster"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client.Database("sebia-dev"), nil
}
