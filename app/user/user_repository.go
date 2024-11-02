package user

import (
	"context"
	"fmt"
	"log"
	"mini-wallet/domain"
	"mini-wallet/domain/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	userCollection              *mongo.Collection
	temporaryUserCollection     *mongo.Collection
	userPasswordResetCollection *mongo.Collection
}

func NewUserRepository(repositoryParam domain.RepositoryParam) user.UserRepository {
	return &userRepository{
		userCollection:              repositoryParam.Mongo.Collection("user"),
		temporaryUserCollection:     repositoryParam.Mongo.Collection("user_temp"),
		userPasswordResetCollection: repositoryParam.Mongo.Collection("user_password_reset"),
	}
}

func (repository *userRepository) GetUserByIdentifier(ctx context.Context, identifier string) (user *user.UserEntity, err error) {
	filter := bson.M{
		"$or": bson.A{
			bson.M{"email": identifier},
			bson.M{"phone_number": identifier},
		},
	}

	res := repository.userCollection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, err
	}

	res.Decode(&user)

	return user, nil
}

func (repository *userRepository) UpsertUser(ctx context.Context, user user.UserEntity) (err error) {
	// Create the upsert option
	opts := options.Replace().SetUpsert(true)

	filter := bson.M{"email": user.Email}

	// Perform the upsert (update or insert)
	result, err := repository.userCollection.ReplaceOne(context.TODO(), filter, user, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		fmt.Println("No document matched the filter, a new document was inserted.")
	} else {
		fmt.Println("Existing document updated.")
	}

	return nil
}

func (repository *userRepository) InsertUserPasswordResetEntity(ctx context.Context, entity user.UserPasswordResetEntity) (err error) {
	_, err = repository.userPasswordResetCollection.InsertOne(ctx, entity)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func (repository *userRepository) DeleteUserPasswordResetEntity(ctx context.Context, email string) (err error) {
	filter := bson.M{"email": email}

	res := repository.userPasswordResetCollection.FindOneAndDelete(ctx, filter)
	if res.Err() != nil {
		return err
	}

	return nil

}

func (repository *userRepository) GetUserPasswordResetEntity(ctx context.Context, token string, now int64) (passwordReset *user.UserPasswordResetEntity, err error) {
	filter := bson.M{
		"password_reset_token": token,
		"expired_at": bson.M{
			"$gt": now,
		},
	}

	res := repository.userPasswordResetCollection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, err
	}

	res.Decode(&passwordReset)

	return passwordReset, nil
}

func (repository *userRepository) DeleteTemporaryUser(ctx context.Context, email string) (err error) {
	filter := bson.M{"email": email}

	res := repository.temporaryUserCollection.FindOneAndDelete(ctx, filter)
	if res.Err() != nil {
		return err
	}

	return nil
}

func (repository *userRepository) GetTemporaryUserByIdentifier(ctx context.Context, identifier string, now int64) (user *user.TemporaryUserEntity, err error) {
	filter := bson.M{
		"$or": bson.A{
			bson.M{"email": identifier},
			bson.M{"phone_number": identifier},
		},
	}

	res := repository.temporaryUserCollection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, err
	}

	res.Decode(&user)

	return user, nil
}

func (repository *userRepository) GetTemporaryUserByVerificationToken(ctx context.Context, token string, now int64) (user *user.TemporaryUserEntity, err error) {
	filter := bson.M{
		"verification_token": token,
		"expired_at": bson.M{
			"$gt": now,
		},
	}

	res := repository.temporaryUserCollection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, err
	}

	res.Decode(&user)

	return user, nil
}

func (repository *userRepository) InsertTemporaryUser(ctx context.Context, user user.TemporaryUserEntity) (err error) {
	_, err = repository.temporaryUserCollection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func (repository *userRepository) InsertUser(ctx context.Context, user user.UserEntity) (err error) {
	_, err = repository.userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func (repository *userRepository) GetUserByEmail(ctx context.Context, email string) (user *user.UserEntity, err error) {
	filter := bson.M{"email": email}

	res := repository.userCollection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, err
	}

	res.Decode(&user)

	return user, nil
}

func (repository *userRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user *user.UserEntity, err error) {
	filter := bson.M{"phone_number": phoneNumber}

	res := repository.userCollection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, err
	}

	res.Decode(&user)

	return user, nil
}
