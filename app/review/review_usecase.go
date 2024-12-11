package review

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/inquiry"
	"mini-wallet/domain/review"
	"mini-wallet/domain/services"
	"mini-wallet/domain/user"
	"mini-wallet/utils"
	"time"
)

type reviewUsecase struct {
	reviewRepository  review.ReviewRepository
	serviceRepository services.ServicesRepository
	inquiryRepository inquiry.InquiryRepository
	userRepository    user.UserRepository
}

func NewReviewUsecase(repositories domain.Repositories) review.ReviewUsecase {

	return &reviewUsecase{
		reviewRepository:  repositories.ReviewRepository,
		serviceRepository: repositories.ServicesRepository,
		inquiryRepository: repositories.InquiryRepository,
		userRepository:    repositories.UserRepository,
	}
}

func (uc *reviewUsecase) GetServiceTopReview(ctx context.Context, serviceId string) (res response.Response[*review.ReviewDTO]) {
	review, err := uc.reviewRepository.GetServiceTopReview(ctx, serviceId)
	if err != nil {
		res.InternalServerError("no such review")
		return
	}

	if review == nil {
		res.Success(nil)
		return
	}

	user, err := uc.userRepository.GetUserByUserID(ctx, review.UserID)
	if err != nil {
		res.InternalServerError("no such review")
		return
	}

	if user == nil {
		res.Success(nil)
		return
	}

	result := review.ToReviewDTO(user.Name)

	res.Success(&result)
	return
}

func (uc *reviewUsecase) CreateReview(ctx context.Context, req review.ReviewDTO) (res response.Response[string]) {
	reviewEntity := req.ToReviewEntity()

	inquiryEntity, err := uc.inquiryRepository.GetInquiryById(ctx, req.InquiryID)
	if err != nil {
		res.InternalServerError(err.Error())
		return

	}

	reviewEntity.ServiceID = inquiryEntity.ServiceID

	if req.UserID != *inquiryEntity.UserID {
		res.Unauthorized("pesanan milik pengguna lain")
		return
	}

	if inquiryEntity == nil {
		res.NotFound("pesanan tidak ditemukan", nil)
		return
	}

	serviceEntity, err := uc.serviceRepository.GetServiceByID(ctx, inquiryEntity.ServiceID)
	if err != nil {
		res.InternalServerError(err.Error())
		return

	}

	if serviceEntity == nil {
		res.NotFound("layanan tidak ditemukan", nil)
		return
	}

	now, _ := utils.GetJktTime()
	lastSelectedDate := inquiryEntity.SelectedDates[len(inquiryEntity.SelectedDates)-1]
	// Define the layout that matches the input string
	layout := "2006/1/2"

	// Parse the string into a time.Time object
	parsedLastSelectedDate, _ := time.Parse(layout, lastSelectedDate)

	if inquiryEntity.Status != 3 || !now.After(parsedLastSelectedDate) || inquiryEntity.ReviewMade {
		res.BadRequest("Belum bisa membuat review", nil)
		return
	}

	err = uc.reviewRepository.InsertReview(ctx, reviewEntity)
	if err != nil {
		res.InternalServerError(err.Error())
		return
	}

	newScore := serviceEntity.TotalScore + req.Score
	newReviewCount := serviceEntity.ReviewCount + 1
	serviceUpdated := serviceEntity.ToServiceEntity(serviceEntity.ID)
	serviceUpdated.TotalScore += newScore
	serviceUpdated.ReviewCount += newReviewCount
	err = uc.serviceRepository.UpdateService(ctx, serviceUpdated)
	if err != nil {
		res.InternalServerError(err.Error())
		return

	}

	inquiryEntity.ReviewMade = true
	err = uc.inquiryRepository.UpdateInquiry(ctx, *inquiryEntity)
	if err != nil {
		res.InternalServerError(err.Error())
		return

	}

	res.Success("review berhasil dibuat")
	return
}
