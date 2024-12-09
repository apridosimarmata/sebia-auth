package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepository interface {
	GetTransaction(ctx context.Context) (tx *mongo.SessionContext, err error)
	CommitTransaction(ctx context.Context, tx mongo.SessionContext) (err error)
	AbortTransaction(ctx context.Context, tx mongo.SessionContext) (err error)
}

type baseRepository struct {
	mongoClient mongo.Client
}

func NewBaseRepository(mongoClient mongo.Client) BaseRepository {
	return &baseRepository{
		mongoClient: mongoClient,
	}
}

func (baseRepo *baseRepository) GetTransaction(ctx context.Context) (tx *mongo.SessionContext, err error) {
	session, err := baseRepo.mongoClient.StartSession()
	if err != nil {
		return nil, err
	}

	session.StartTransaction()
	transactionCtx := mongo.NewSessionContext(ctx, session)

	return &transactionCtx, nil
}

func (baseRepo *baseRepository) CommitTransaction(ctx context.Context, tx mongo.SessionContext) (err error) {
	err = tx.CommitTransaction(tx)

	if err != nil {
		tx.AbortTransaction(ctx)
		return err
	}

	return nil
}

func (base *baseRepository) AbortTransaction(ctx context.Context, tx mongo.SessionContext) (err error) {
	err = tx.AbortTransaction(tx)

	if err != nil {
		tx.AbortTransaction(ctx)
		return err
	}

	return nil
}
