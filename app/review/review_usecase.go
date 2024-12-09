package review

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/common/response"
	"mini-wallet/domain/inquiry"
	"mini-wallet/domain/review"
	"mini-wallet/domain/services"
	"mini-wallet/utils"
	"time"
)

type reviewUsecase struct {
	reviewRepository  review.ReviewRepository
	serviceRepository services.ServicesRepository
	inquiryRepository inquiry.InquiryRepository
}

func NewReviewUsecase(repositories domain.Repositories) review.ReviewUsecase {

	return &reviewUsecase{
		reviewRepository:  repositories.ReviewRepository,
		serviceRepository: repositories.ServicesRepository,
		inquiryRepository: repositories.InquiryRepository,
	}
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

	avgScore := float32(1.0)
	if serviceEntity.AverageScore == 0 {
		avgScore = float32(req.Score)
	} else {
		sum := serviceEntity.AverageScore + float32(req.Score)
		avgScore = sum / 2
	}

	serviceUpdated := serviceEntity.ToServiceEntity(serviceEntity.ID)
	serviceUpdated.AverageScore = avgScore
	serviceUpdated.ReviewCount += 1
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
