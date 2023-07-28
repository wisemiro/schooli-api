package services

import (
	"context"
	"encoding/json"
	"schooli-api/internal/models"
	"schooli-api/internal/repository/postgresql/db"
	"schooli-api/pkg/resterrors"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrderService interface {
	CreateOrderItem(ctx context.Context, om models.OrderProduct) error
	UpdateOrderItem(ctx context.Context, om models.OrderProduct) error
	DeleteOrderItem(ctx context.Context, orderProductID int64) error
	ListOrderItems(ctx context.Context) ([]*models.OrderProduct, error)
	GetOrderProduct(ctx context.Context, orderProductID int64) (*models.OrderProduct, error)
	CreateOrder(ctx context.Context, om models.Order) error
	UpdateOrder(ctx context.Context, om models.Order) error
	DeleteOrder(ctx context.Context, orderID int64) error
	ListOrders(ctx context.Context) ([]*models.OrderProduct, error)
	ListUserOrders(ctx context.Context, userID int64) ([]*models.OrderProduct, error)
}

func (sq *SQLStore) CreateOrderItem(ctx context.Context, om models.OrderProduct) error {
	if err := sq.store.CreateOrderProduct(ctx, db.CreateOrderProductParams{
		Quantity:        int32(om.Quantity),
		TotalPrice:      int32(om.TotalPrice),
		ProductVariants: om.ProductVariants,
		ProductID:       om.Product.ID,
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.CreateOrderItem")
	}
	return nil
}

func (sq *SQLStore) UpdateOrderItem(ctx context.Context, om models.OrderProduct) error {
	if err := sq.store.UpdateOrderProduct(ctx, db.UpdateOrderProductParams{
		ID:              om.ID,
		Quantity:        int32(om.Quantity),
		TotalPrice:      int32(om.TotalPrice),
		ProductVariants: om.ProductVariants,
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.UpdateOrderItem")
	}
	return nil
}

func (sq *SQLStore) DeleteOrderItem(ctx context.Context, orderProductID int64) error {
	if err := sq.store.DeleteOrderProduct(ctx, orderProductID); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.DeleteOrderItem")
	}
	return nil
}

func (sq *SQLStore) ListOrderItems(ctx context.Context) ([]*models.OrderProduct, error) {
	orderProdcts, err := sq.store.ListOrderProducts(ctx)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.ListOrderItems")
	}
	orders := make([]*models.OrderProduct, len(orderProdcts))
	for i, v := range orderProdcts {
		var pv []models.ProductVariant
		err = json.Unmarshal(v.ProductVariants, &pv)
		if err != nil {
			return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.ListOrderItems")
		}
		orders[i] = &models.OrderProduct{
			ID:        v.ID_2,
			CreatedAt: v.CreatedAt.Time,
			UpdatedAt: v.UpdatedAt.Time,
			Product: models.Product{
				ID:            v.ID.Int64,
				Name:          v.Name.String,
				Price:         int64(v.Price.Int32),
				DiscountPrice: int64(v.DiscountPrice.Int32),
				DefaultImage:  sq.fileStore.BuildFilePath(v.DefaultImage.String),
			},
			Quantity:   int64(v.Quantity),
			TotalPrice: int64(v.TotalPrice),
			Variants:   pv,
		}
	}
	return orders, nil
}

func (sq *SQLStore) GetOrderProduct(ctx context.Context, orderProductID int64) (*models.OrderProduct, error) {
	v, err := sq.store.GetOrderProduct(ctx, orderProductID)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.GetOrderProduct")
	}
	var pv []models.ProductVariant
	err = json.Unmarshal(v.ProductVariants, &pv)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.GetOrderProduct")
	}
	return &models.OrderProduct{
		ID:        v.ID_2,
		CreatedAt: v.CreatedAt.Time,
		UpdatedAt: v.UpdatedAt.Time,
		Product: models.Product{
			ID:            v.ID.Int64,
			Name:          v.Name.String,
			Price:         int64(v.Price.Int32),
			DiscountPrice: int64(v.DiscountPrice.Int32),
			DefaultImage:  sq.fileStore.BuildFilePath(v.DefaultImage.String),
		},

		Quantity:   int64(v.Quantity),
		TotalPrice: int64(v.TotalPrice),
		Variants:   pv,
	}, nil
}

func (sq *SQLStore) CreateOrder(ctx context.Context, om models.Order) error {
	if err := sq.store.CreateOrder(ctx, db.CreateOrderParams{
		GrandTotal:      int32(om.GrandTotal),
		SerialNumber:    om.SerialNumber,
		ShippingAddress: om.Shipping.ID,
		UserID:          pgtype.Int8{Int64: om.User.ID, Valid: true},
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.CreateOrder")
	}
	return nil
}
func (sq *SQLStore) UpdateOrder(ctx context.Context, om models.Order) error {
	if err := sq.store.UpdateOrder(ctx, db.UpdateOrderParams{
		GrandTotal: int32(om.GrandTotal),
		Confirmed:  pgtype.Bool{Bool: om.Confirmed, Valid: true},
	}); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.CreateOrder")
	}
	return nil
}

func (sq *SQLStore) DeleteOrder(ctx context.Context, orderID int64) error {
	if err := sq.store.DeleteOrder(ctx, orderID); err != nil {
		return resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.CreateOrder")
	}
	return nil
}

func (sq *SQLStore) ListOrders(ctx context.Context) ([]*models.OrderProduct, error) {
	orders, err := sq.store.ListOrders(ctx)
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.ListOrders")

	}
	ordersList := make([]*models.OrderProduct, len(orders))

	for i, v := range orders {
		var pv []models.ProductVariant
		err = json.Unmarshal(v.ProductVariants, &pv)
		if err != nil {
			return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.GetOrderProduct")
		}
		ordersList[i] = &models.OrderProduct{
			CreatedAt:  v.CreatedAt.Time,
			Quantity:   int64(v.Quantity),
			TotalPrice: int64(v.TotalPrice),
			Variants:   pv,
			Product: models.Product{
				ID:            v.ID.Int64,
				Name:          v.Name.String,
				Price:         int64(v.Price.Int32),
				DiscountPrice: int64(v.DiscountPrice.Int32),
				Sku:           v.Sku.String,
				DefaultImage:  sq.fileStore.BuildFilePath(v.DefaultImage.String),
			},
			Order: models.Order{
				ID:           v.ID.Int64,
				SerialNumber: v.SerialNumber.String,
				GrandTotal:   int64(v.GrandTotal.Int32),
				User: models.User{
					ID:          v.ID_4.Int64,
					Email:       v.Email.String,
					PhoneNumber: v.PhoneNumber.String,
				},
				Shipping: models.Shipping{
					ID: v.ID_5.Int64,
					Geo: models.Geo{
						Latitude:  v.Latitude.(float64),
						Longitude: v.Longitude.(float64),
					},
				},
			},
		}

	}
	return ordersList, nil
}

func (sq *SQLStore) ListUserOrders(ctx context.Context, userID int64) ([]*models.OrderProduct, error) {
	orders, err := sq.store.ListUserOrders(ctx, pgtype.Int8{Int64: userID, Valid: true})
	if err != nil {
		return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.ListUserOrders")
	}
	ordersList := make([]*models.OrderProduct, len(orders))

	for i, v := range orders {
		var pv []models.ProductVariant
		err = json.Unmarshal(v.ProductVariants, &pv)
		if err != nil {
			return nil, resterrors.WrapErrorf(err, resterrors.ECodeUnknown, "OrdersService.GetOrderProduct")
		}
		ordersList[i] = &models.OrderProduct{
			CreatedAt:  v.CreatedAt.Time,
			Quantity:   int64(v.Quantity),
			TotalPrice: int64(v.TotalPrice),
			Variants:   pv,
			Product: models.Product{
				ID:            v.ID.Int64,
				Name:          v.Name.String,
				Price:         int64(v.Price.Int32),
				DiscountPrice: int64(v.DiscountPrice.Int32),
				Sku:           v.Sku.String,
				DefaultImage:  sq.fileStore.BuildFilePath(v.DefaultImage.String),
			},
			Order: models.Order{
				ID:           v.ID.Int64,
				SerialNumber: v.SerialNumber.String,
				GrandTotal:   int64(v.GrandTotal.Int32),
				User: models.User{
					ID:          v.ID_4.Int64,
					Email:       v.Email.String,
					PhoneNumber: v.PhoneNumber.String,
				},
				Shipping: models.Shipping{
					ID: v.ID_5.Int64,
					Geo: models.Geo{
						Latitude:  v.Latitude.(float64),
						Longitude: v.Longitude.(float64),
					},
				},
			},
		}

	}
	return ordersList, nil
}
