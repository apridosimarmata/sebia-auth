package presentation

import (
	"context"
	"fmt"
	"mini-wallet/app/auth"
	"mini-wallet/app/user"

	"mini-wallet/domain"
	"mini-wallet/infrastructure"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitServer() chi.Router {
	ctx := context.Background()
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// config := infrastructure.GetConfig()
	//

	mongoDb, err := infrastructure.GetMongoDatabase(ctx)
	if err != nil {
		panic(err)
	}

	repositoryParam := domain.RepositoryParam{
		Mongo: mongoDb,
	}

	repositories := domain.Repositories{
		UserRepository: user.NewUserRepository(repositoryParam),
	}

	usecases := domain.Usecases{
		AuthUsecase: auth.NewAuthUsecase(repositories),
	}

	// in terms of authorization, a token should not be a forever-lived value
	// provided a /refresh endpoint to get fresh token
	auth.SetAuthHandler(router, usecases)

	fmt.Println("server listening on port 3000")

	return router
}

func StopServer() {

}
