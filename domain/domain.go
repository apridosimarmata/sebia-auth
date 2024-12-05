package domain

import (
	"mini-wallet/domain/affiliate"
	"mini-wallet/domain/auth"
	"mini-wallet/domain/booking"
	"mini-wallet/domain/business"
	"mini-wallet/domain/file"
	"mini-wallet/domain/inquiry"
	"mini-wallet/domain/locations"
	"mini-wallet/domain/payment"
	"mini-wallet/domain/services"
	"mini-wallet/domain/user"
	"mini-wallet/infrastructure"
	"mini-wallet/integration"

	"github.com/aws/aws-sdk-go/service/s3"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	UserRepository           user.UserRepository
	LocationRepository       locations.LocationRepository
	BusinessRepository       business.BusinessRepository
	AffiliateRepository      affiliate.AffiliateRepository
	ServicesRepository       services.ServicesRepository
	ServicesSearchRepository services.ServicesSearchRepository

	InquiryRepository inquiry.InquiryRepository
	BookingRepository booking.BookingRepository
}

type Usecases struct {
	AuthUsecase      auth.AuthUsecase
	BusinessUsecase  business.BusinessUsecase
	AffiliateUsecase affiliate.AffiliateUsecase
	FileUsecase      file.FileUsecase
	LocationUsecase  locations.LocationUsecase
	ServicesUsecase  services.ServicesUsecase
	InquiryUsecase   inquiry.InquiryUsecase
	PaymentUsecase   payment.PaymentUsecase
	BookingUsecase   booking.BookingUsecase
}

type Infrastructure struct {
	S3                  s3.S3
	NotificationService integration.NotificationService
	PaymentService      infrastructure.Payment
	MesageProducer      infrastructure.MessagingProducer
}

type RepositoryParam struct {
	Mongo *mongo.Database
}
