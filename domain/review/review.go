package review

import (
	"context"
	"errors"
	"mini-wallet/domain/common/response"
	"mini-wallet/utils"
	"time"
)

type ReviewEntity struct {
	ID        string `bson:"id"`
	ServiceID string `bson:"service_id"`
	InquiryID string `bson:"inquiry_id"`
	UserID    string `bson:"user_id"`
	Content   string `bson:"content"`
	Score     int    `bson:"score"`
	CreatedAt string `bson:"created_at"`
	Status    int    `bson:"status"`
}

type ReviewDTO struct {
	InquiryID string `json:"inquiry_id"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	Content   string `json:"content"`
	Score     int    `json:"score"`
	CreatedAt string `json:"created_at"`
}

func (p *ReviewEntity) ToReviewDTO(userName string) ReviewDTO {
	return ReviewDTO{
		UserName:  userName,
		Content:   p.Content,
		Score:     p.Score,
		CreatedAt: p.CreatedAt,
	}
}

type ReviewUsecase interface {
	CreateReview(ctx context.Context, req ReviewDTO) (res response.Response[string])
	GetServiceTopReview(ctx context.Context, serviceId string) (res response.Response[*ReviewDTO])
}

type ReviewRepository interface {
	GetServiceTopReview(ctx context.Context, serviceId string) (res *ReviewEntity, err error)
	InsertReview(ctx context.Context, review ReviewEntity) (err error)
}

func (p *ReviewDTO) ToReviewEntity() ReviewEntity {
	now, _ := utils.GetJktTime()
	return ReviewEntity{
		ID: utils.GenerateUniqueId(),
		// ServiceID: p.ServiceID,
		InquiryID: p.InquiryID,
		Score:     p.Score,
		Content:   p.Content,
		Status:    0,
		UserID:    p.UserID,
		CreatedAt: now.Format(time.RFC3339),
	}
}

func (p *ReviewDTO) Validate() error {
	err := utils.ValidateRequiredInt(p.Score)
	if err != nil {
		return err
	}

	if p.Score < 1 || p.Score > 5 {
		return errors.New("skor tidak valid")
	}

	err = utils.ValidateRequired(p.InquiryID)
	if err != nil {
		return err
	}

	// err = utils.ValidateRequired(p.ServiceID)
	// if err != nil {
	// 	return err
	// }

	return nil
}
