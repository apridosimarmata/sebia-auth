package booking

import (
	"context"
	"mini-wallet/domain"
	"mini-wallet/domain/booking"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type bookingRepository struct {
	bookingCollection *mongo.Collection
}

func NewBookingRepository(repositoryParam domain.RepositoryParam) booking.BookingRepository {
	return &bookingRepository{
		bookingCollection: repositoryParam.Mongo.Collection("bookings"),
	}
}

func (repo *bookingRepository) GetBookings(ctx context.Context, serviceId string, variantPax int, yearMonths []string) (res []booking.ServiceBookings, err error) {
	filter := bson.M{
		"service_id":  serviceId,
		"variant_pax": variantPax,
		"year_month": bson.M{
			"$in": yearMonths,
		},
	}

	result, err := repo.bookingCollection.Find(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	result.All(ctx, &res)

	return res, nil
}

// CreateServiceBookingDocument(ctx context.Context, serviceID string, yearMonth string) (err error)
func (repo *bookingRepository) UpsertBookingsDocument(ctx context.Context, tx *mongo.SessionContext, documents []booking.ServiceBookings) (err error) {
	for _, document := range documents {
		filter := bson.M{"id": document.ID}
		update := bson.M{"$set": bson.M{
			"service_id":       document.ServiceID,
			"variant_pax":      document.VariantPax,
			"year_month":       document.YearMonth,
			"bookings_by_date": document.BookingsByDate,
		}}

		opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
		result := repo.bookingCollection.FindOneAndUpdate(*tx, filter, update, opts)

		if result.Err() != nil {
			return result.Err()
		}
	}

	return nil
}
