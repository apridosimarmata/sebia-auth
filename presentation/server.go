package presentation

import (
	"context"
	"fmt"
	"mini-wallet/app/affiliate"
	"mini-wallet/app/auth"
	"mini-wallet/app/booking"
	"mini-wallet/app/business"
	"mini-wallet/app/file"
	"mini-wallet/app/inquiry"
	"mini-wallet/app/review"
	"mini-wallet/utils"

	"mini-wallet/app/location"
	"mini-wallet/app/payment"
	"mini-wallet/app/services"
	"mini-wallet/integration"

	"mini-wallet/app/user"

	"mini-wallet/domain"

	"mini-wallet/infrastructure"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Temporary struct {
	Message string `json:"message"`
}

func InitServer() (chi.Router, string) {
	ctx := context.Background()
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	fmt.Println("reading config")

	config, err := utils.GetConfig()
	if err != nil {
		panic("error getting config")
	}

	// config := infrastructure.GetConfig()
	//

	grpcConn, err := infrastructure.NewGrpcConn()
	notificationService := integration.NewNotificationService(&grpcConn.NotificationService)

	mongoDb, err := infrastructure.GetMongoDatabase(ctx, config.DatabaseName)
	if err != nil {
		panic(err)
	}

	repositoryParam := domain.RepositoryParam{
		Mongo: mongoDb,
	}

	// locations.PopulateData(mongoDb)
	// return router

	repositories := domain.Repositories{
		BaseRepository:           domain.NewBaseRepository(*mongoDb.Client()),
		UserRepository:           user.NewUserRepository(repositoryParam),
		LocationRepository:       location.NewLocationRepository(repositoryParam),
		BusinessRepository:       business.NewBusinessRepository(repositoryParam),
		AffiliateRepository:      affiliate.NewAffiliatesRepository(repositoryParam),
		ServicesRepository:       services.NewServicesRepository(repositoryParam),
		InquiryRepository:        inquiry.NewInquiryRepository(repositoryParam),
		BookingRepository:        booking.NewBookingRepository(repositoryParam),
		ServicesSearchRepository: services.NewServicesSearchRepository(repositoryParam),
		ReviewRepository:         review.NewReviewRepository(repositoryParam),
	}

	s3, err := infrastructure.NewS3Service()
	if err != nil {
		panic(err.Error())
	}

	snapClient := integration.NewSnapClient()
	messagingProducer := infrastructure.NewMessagingProducer()

	infra := domain.Infrastructure{
		S3:                  *s3,
		NotificationService: notificationService,
		PaymentService:      infrastructure.NewPayment(snapClient),
		MesageProducer:      messagingProducer,
	}

	usecases := domain.Usecases{
		AuthUsecase:      auth.NewAuthUsecase(repositories, infra, config),
		FileUsecase:      file.NewFileUsecase(infra),
		LocationUsecase:  location.NewLocationUsecase(repositories),
		BusinessUsecase:  business.NewBusinessUsecase(repositories),
		AffiliateUsecase: affiliate.NewAffiliatesUsecase(repositories),
		ServicesUsecase:  services.NewServicesUsecase(repositories),
		InquiryUsecase:   inquiry.NewInquiryUsecase(repositories, infra),
		PaymentUsecase:   payment.NewPaymentUsecase(repositories, infra, config),
		BookingUsecase:   booking.NewBookingUsecase(repositories, infra),
		ReviewUsecase:    review.NewReviewUsecase(repositories),
	}

	middlewares := auth.NewAuthMiddleware(repositories, config)

	// messaging
	go infrastructure.RegisterConsumers([]infrastructure.RegisterListenersParam{
		{
			Topic:    config.BookingTopic,
			Channel:  "creation",
			Listener: booking.NewBookingMessageConsumer(usecases),
		},
	})

	// in terms of authorization, a token should not be a forever-lived value
	// provided a /refresh endpoint to get fresh token
	auth.SetAuthHandler(router, usecases, middlewares, config)
	file.SetFileHandler(router, usecases)
	location.SetLocationHandler(router, usecases)
	business.SetBusinessHandler(router, usecases, middlewares)
	affiliate.SetAffiliatesHandler(router, usecases, middlewares)
	services.SetServicesHandler(router, usecases, middlewares)
	inquiry.SetInquiryHandler(router, usecases, middlewares)
	payment.SetPaymentHandler(router, usecases)
	review.SetReviewHandler(router, usecases, middlewares)

	fmt.Println("[" + config.AppEnvironment + "] server listening on port " + config.AppPort)

	return router, config.AppPort
}

func StopServer() {

}
