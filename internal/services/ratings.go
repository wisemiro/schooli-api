package services

import (
	"context"
	"schooli-api/internal/models"
	"schooli-api/internal/repository/postgresql/db"
	"schooli-api/pkg/resterrors"

	"github.com/jackc/pgx/v5/pgtype"
)

type RatingsService interface {
	CreateRatings(ctx context.Context, pr models.ProductRatings) error
	UpdateRatings(ctx context.Context, pr models.ProductRatings) error
	DeleteRating(ctx context.Context, ratingID int64) error
	ProductRatings(ctx context.Context, productID int64) ([]*models.ProductRatings, error)
	GetRating(ctx context.Context, ratingID int64) (*models.ProductRatings, error)
}

func (sq *SQLStore) CreateRatings(ctx context.Context, pr models.ProductRatings) error {
	if err := sq.store.CreateRating(ctx, db.CreateRatingParams{
		UserID:    pr.User.ID,
		Stars:     int32(pr.Stars),
		Feedback:  pgtype.Text{String: pr.Feeedback, Valid: true},
		ProductID: pr.Product.ID,
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "RatingsService.CreateRatings")
	}
	return nil
}
func (sq *SQLStore) UpdateRatings(ctx context.Context, pr models.ProductRatings) error {
	if err := sq.store.UpdateRating(ctx, db.UpdateRatingParams{
		ID:       pr.ID,
		Stars:    int32(pr.Stars),
		Feedback: pgtype.Text{String: pr.Feeedback, Valid: true},
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "RatingsService.UpdateRatings")
	}
	return nil
}

func (sq *SQLStore) DeleteRating(ctx context.Context, ratingID int64) error {
	if err := sq.store.DeleteRatings(ctx, ratingID); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "RatingsService.DeleteRatings")
	}
	return nil
}

func (sq *SQLStore) ProductRatings(ctx context.Context, productID int64) ([]*models.ProductRatings, error) {
	ratings, err := sq.store.ListProductRatings(ctx, productID)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "RatingsService.ProductRatings")
	}
	ratingsList := make([]*models.ProductRatings, len(ratings))
	for i, v := range ratings {
		ratingsList[i] = &models.ProductRatings{
			ID:        v.ID,
			CreatedAt: v.CreatedAt.Time,
			User: models.User{
				ID:    v.UserID,
				Email: v.Email.String,
			},
			Product: models.Product{
				ID:   v.ProductID,
				Name: v.Name.String,
			},
			Stars:     int(v.Stars),
			Feeedback: v.Feedback.String,
		}
	}
	return ratingsList, nil
}

func (sq *SQLStore) GetRating(ctx context.Context, ratingID int64) (*models.ProductRatings, error) {
	rating, err := sq.store.GetRating(ctx, ratingID)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "RatingsService.GetRating")
	}
	return &models.ProductRatings{
		ID:        rating.ID,
		CreatedAt: rating.CreatedAt.Time,
		User: models.User{
			ID:    rating.UserID,
			Email: rating.Email.String,
		},
		Product: models.Product{
			ID:   rating.ProductID,
			Name: rating.Name.String,
		},
		Stars:     int(rating.Stars),
		Feeedback: rating.Feedback.String,
	}, nil
}
