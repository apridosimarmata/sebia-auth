package domain

import (
	"mini-wallet/domain/auth"
	"mini-wallet/domain/user"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	UserRepository user.UserRepository
}

type Usecases struct {
	AuthUsecase auth.AuthUsecase
}

type RepositoryParam struct {
	Mongo *mongo.Database
}
