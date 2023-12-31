// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Categories struct {
	ID        int64              `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	Name      string             `json:"name"`
	Image     string             `json:"image"`
}

type Devices struct {
	ID        int64              `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UserID    int64              `json:"user_id"`
	Device    string             `json:"device"`
}

type OrderProducts struct {
	ID              int64              `json:"id"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
	DeletedAt       pgtype.Timestamptz `json:"deleted_at"`
	Quantity        int32              `json:"quantity"`
	TotalPrice      int32              `json:"total_price"`
	ProductVariants []int32            `json:"product_variants"`
	ProductID       int64              `json:"product_id"`
	OrderID         pgtype.Int8        `json:"order_id"`
}

type Orders struct {
	ID              int64              `json:"id"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
	UpdatedAt       pgtype.Timestamptz `json:"updated_at"`
	GrandTotal      int32              `json:"grand_total"`
	SerialNumber    string             `json:"serial_number"`
	UserID          pgtype.Int8        `json:"user_id"`
	ShippingAddress int64              `json:"shipping_address"`
	Confirmed       pgtype.Bool        `json:"confirmed"`
}

type ProductRatings struct {
	ID        int64              `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	UserID    int64              `json:"user_id"`
	Stars     int32              `json:"stars"`
	Feedback  pgtype.Text        `json:"feedback"`
	ProductID int64              `json:"product_id"`
}

type ProductSpecifications struct {
	ID          int64              `json:"id"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
	DeletedAt   pgtype.Timestamptz `json:"deleted_at"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	ProductID   int64              `json:"product_id"`
}

type ProductVariants struct {
	ID        int64              `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	Type      int32              `json:"type"`
	Name      string             `json:"name"`
	ProductID int64              `json:"product_id"`
}

type Products struct {
	ID            int64              `json:"id"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	UpdatedAt     pgtype.Timestamptz `json:"updated_at"`
	DeletedAt     pgtype.Timestamptz `json:"deleted_at"`
	AverageRating pgtype.Int4        `json:"average_rating"`
	Name          string             `json:"name"`
	Price         int32              `json:"price"`
	DiscountPrice pgtype.Int4        `json:"discount_price"`
	Sku           string             `json:"sku"`
	Description   string             `json:"description"`
	StockCount    int32              `json:"stock_count"`
	MinStockCount int32              `json:"min_stock_count"`
	CategoryID    int64              `json:"category_id"`
	TotalRatings  int32              `json:"total_ratings"`
	TotalView     pgtype.Int4        `json:"total_view"`
	DefaultImage  string             `json:"default_image"`
}

type SchemaMigrations struct {
	Version int64 `json:"version"`
	Dirty   bool  `json:"dirty"`
}

type Shipping struct {
	ID        int64              `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	Location  interface{}        `json:"location"`
	UserID    int64              `json:"user_id"`
}

type Users struct {
	ID           int64              `json:"id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
	DeletedAt    pgtype.Timestamptz `json:"deleted_at"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"password_hash"`
	PhoneNumber  string             `json:"phone_number"`
	IsVerified   pgtype.Bool        `json:"is_verified"`
}

type Wishlists struct {
	ID        int64              `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
	ProductID int64              `json:"product_id"`
	UserID    int64              `json:"user_id"`
}
