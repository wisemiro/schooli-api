package models

import "time"

type Category struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
}

type Product struct {
	ID                   int64                  `json:"id"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	DeletedAt            time.Time              `json:"-"`
	Name                 string                 `json:"name"`
	Price                int64                  `json:"price"`
	AverageRating        int                    `json:"average_rating"`
	DiscountPrice        int64                  `json:"discount_price"`
	Sku                  string                 `json:"sku"`
	Description          string                 `json:"description"`
	StockCount           int64                  `json:"stock_count"`
	MinStockCount        int64                  `json:"min_stock_count"`
	Category             Category               `json:"category"`
	DefaultImage         string                 `json:"default_image"`
	TotalRatings         int                    `json:"total_ratings"`
	TotalViews           int                    `json:"total_views"`
	ProductVariant       []ProductVariant       `json:"product_variants"`
	ProductSpecification []ProductSpecification `json:"product_specifications"`
}

type ProductVariant struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	ProductID int64     `json:"product_id"`
	Type      int       `json:"type"`
}

type OrderProduct struct {
	ID              int64          `json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       time.Time      `json:"-"`
	Product         Product        `json:"product"`
	ProductVariants []int32        `json:"product_variants"`
	Variants        []ProductVariant `json:"variants"`
	Quantity        int64          `json:"quantity"`
	TotalPrice      int64          `json:"total_price"`
	Order           Order          `json:"order_id"`
}

type Order struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"-"`
	SerialNumber string    `json:"serial_number"`
	GrandTotal   int64     `json:"grand_total"`
	User         User      `json:"user"`
	Shipping     Shipping  `json:"shipping_address"`
	Confirmed    bool      `json:"confirmed"`
}

type ProductSpecification struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ProductID   int64     `json:"product_id"`
}
