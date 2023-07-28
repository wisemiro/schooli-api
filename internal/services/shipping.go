package services

import (
	"context"
	"fmt"
	"schooli-api/internal/models"
	"schooli-api/internal/repository/postgresql/db"
)

type ShippingService interface {
	CreateShipping(ctx context.Context, sm models.Shipping) error
	UpdateShippingAddress(ctx context.Context, sm models.Shipping) error
	ListShipping(ctx context.Context) ([]*models.Shipping, error)
	ListUserShipping(ctx context.Context, userID int64) ([]*models.Shipping, error)
}

func (sq *SQLStore) CreateShipping(ctx context.Context, sm models.Shipping) error {
	k := fmt.Sprintf("POINT(%v %v)", sm.Geo.Latitude, sm.Geo.Longitude)
	if err := sq.store.CreateShipping(ctx, db.CreateShippingParams{
		StGeomfromtext: k,
		UserID:         sm.User.ID,
	}); err != nil {
		return err
	}
	return nil
}

func (sq *SQLStore) UpdateShippingAddress(ctx context.Context, sm models.Shipping) error {
	k := fmt.Sprintf("POINT(%v %v)", sm.Geo.Latitude, sm.Geo.Longitude)
	if err := sq.store.UpdateShipping(ctx, db.UpdateShippingParams{
		StGeomfromtext: k,
		ID:             sm.ID,
	}); err != nil {
		return err
	}
	return nil
}

func (sq *SQLStore) ListShipping(ctx context.Context) ([]*models.Shipping, error) {
	shipping, err := sq.store.ListShipping(ctx)
	if err != nil {
		return nil, err
	}
	shList := make([]*models.Shipping, len(shipping))
	for i, v := range shipping {
		shList[i] = &models.Shipping{
			ID:        v.ID,
			CreatedAt: v.CreatedAt.Time,
			UpdatedAt: v.UpdatedAt.Time,
			User: models.User{
				ID:          v.ID_2.Int64,
				Email:       v.Email.String,
				PhoneNumber: v.PhoneNumber.String,
			},
			Geo: models.Geo{
				Latitude:  v.Latitude.(float64),
				Longitude: v.Longitude.(float64),
			},
		}
	}
	return shList, nil
}

func (sq *SQLStore) ListUserShipping(ctx context.Context, userID int64) ([]*models.Shipping, error) {
	shipping, err := sq.store.ListUserShipping(ctx, userID)
	if err != nil {
		return nil, err
	}
	shList := make([]*models.Shipping, len(shipping))
	for i, v := range shipping {
		shList[i] = &models.Shipping{
			ID:        v.ID,
			CreatedAt: v.CreatedAt.Time,
			UpdatedAt: v.UpdatedAt.Time,
			User: models.User{
				ID:          v.ID_2.Int64,
				Email:       v.Email.String,
				PhoneNumber: v.PhoneNumber.String,
			},
			Geo: models.Geo{
				Latitude:  v.Latitude.(float64),
				Longitude: v.Longitude.(float64),
			},
		}
	}
	return shList, nil
}
