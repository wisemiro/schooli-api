// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateCategory(ctx context.Context, arg CreateCategoryParams) error
	//
	//
	CreateDevice(ctx context.Context, arg CreateDeviceParams) error
	//
	//
	CreateOrder(ctx context.Context, arg CreateOrderParams) error
	//
	//
	CreateOrderProduct(ctx context.Context, arg CreateOrderProductParams) error
	//
	//
	CreateProduct(ctx context.Context, arg CreateProductParams) error
	//
	//
	CreateProductSpec(ctx context.Context, arg CreateProductSpecParams) error
	CreateRating(ctx context.Context, arg CreateRatingParams) error
	CreateShipping(ctx context.Context, arg CreateShippingParams) error
	//
	//
	CreateUser(ctx context.Context, arg CreateUserParams) error
	//
	//
	CreateVariant(ctx context.Context, arg CreateVariantParams) error
	//
	//
	CreateWishlist(ctx context.Context, arg CreateWishlistParams) error
	//
	//
	DeleteCategory(ctx context.Context, id int64) error
	//
	//
	DeleteOrder(ctx context.Context, id int64) error
	//
	//
	DeleteOrderProduct(ctx context.Context, id int64) error
	//
	//
	DeleteProduct(ctx context.Context, id int64) error
	//
	//
	DeleteProductSpec(ctx context.Context, id int64) error
	//
	//
	DeleteRatings(ctx context.Context, id int64) error
	//
	//
	DeleteUser(ctx context.Context, id int64) error
	//
	//
	DeleteVariant(ctx context.Context, id int64) error
	//
	//
	DeleteWishlist(ctx context.Context, id int64) error
	//
	//
	DiscountedProducts(ctx context.Context) ([]*Products, error)
	//
	//
	GetCategory(ctx context.Context, id int64) (*Categories, error)
	//
	//
	GetOneDevice(ctx context.Context, userID int64) (*Devices, error)
	//
	//
	GetOrderProduct(ctx context.Context, id int64) (*GetOrderProductRow, error)
	//
	//
	GetRating(ctx context.Context, id int64) (*GetRatingRow, error)
	//
	//
	GetUser(ctx context.Context, id int64) (*Users, error)
	//
	//
	GetVariant(ctx context.Context, id int64) (*ProductVariants, error)
	//
	//
	GetWishlist(ctx context.Context, arg GetWishlistParams) (*Wishlists, error)
	//
	//
	ListByVerification(ctx context.Context, isVerified pgtype.Bool) ([]*Users, error)
	//
	//
	ListCategories(ctx context.Context) ([]*Categories, error)
	ListOrderProducts(ctx context.Context) ([]*ListOrderProductsRow, error)
	//
	//
	ListOrders(ctx context.Context) ([]*ListOrdersRow, error)
	//
	//
	ListProductRatings(ctx context.Context, productID int64) ([]*ListProductRatingsRow, error)
	//
	//
	ListProducts(ctx context.Context) ([]*ListProductsRow, error)
	//
	//
	ListShipping(ctx context.Context) ([]*ListShippingRow, error)
	//
	//
	ListUserOrders(ctx context.Context, userID pgtype.Int8) ([]*ListUserOrdersRow, error)
	//
	//
	ListUserShipping(ctx context.Context, userID int64) ([]*ListUserShippingRow, error)
	ListUsers(ctx context.Context) ([]*Users, error)
	//
	//
	ListVariants(ctx context.Context, productID int64) ([]*ProductVariants, error)
	//
	//
	ListWishlist(ctx context.Context, userID int64) ([]*ListWishlistRow, error)
	//
	//
	OneProduct(ctx context.Context, id int64) (*OneProductRow, error)
	//
	//
	ProductSpecifications(ctx context.Context, productID int64) ([]*ProductSpecifications, error)
	ProductsByCategory(ctx context.Context, arg ProductsByCategoryParams) ([]*Products, error)
	//
	// TODO: Delete stale order_products
	//
	SearchProducts(ctx context.Context, similarity string) ([]*SearchProductsRow, error)
	//
	//
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) error
	//
	//
	UpdateDevice(ctx context.Context, arg UpdateDeviceParams) error
	//
	//
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (pgtype.Int8, error)
	//
	//
	UpdateOrderProduct(ctx context.Context, arg UpdateOrderProductParams) error
	//
	//
	UpdateProduct(ctx context.Context, arg UpdateProductParams) error
	//
	//
	UpdateProductSpec(ctx context.Context, arg UpdateProductSpecParams) error
	//
	//
	UpdateRating(ctx context.Context, arg UpdateRatingParams) error
	//
	//
	UpdateShipping(ctx context.Context, arg UpdateShippingParams) error
	//
	//
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	//
	//
	UpdateVariant(ctx context.Context, arg UpdateVariantParams) error
	//
	//
	UserByEmail(ctx context.Context, email string) (*UserByEmailRow, error)
}

var _ Querier = (*Queries)(nil)
