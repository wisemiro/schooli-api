// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: products.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createCategory = `-- name: CreateCategory :exec
insert into categories(created_at, name, image)
values(current_timestamp, $1, $2)
`

type CreateCategoryParams struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) error {
	_, err := q.db.Exec(ctx, createCategory, arg.Name, arg.Image)
	return err
}

const createProduct = `-- name: CreateProduct :exec
insert into products(
        created_at,
        name,
        price,
        discount_price,
        sku,
        description,
        stock_count,
        category_id,
        default_image
    )
values(
        current_timestamp,
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    )
`

type CreateProductParams struct {
	Name          string      `json:"name"`
	Price         int32       `json:"price"`
	DiscountPrice pgtype.Int4 `json:"discount_price"`
	Sku           string      `json:"sku"`
	Description   string      `json:"description"`
	StockCount    int32       `json:"stock_count"`
	CategoryID    int64       `json:"category_id"`
	DefaultImage  string      `json:"default_image"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) error {
	_, err := q.db.Exec(ctx, createProduct,
		arg.Name,
		arg.Price,
		arg.DiscountPrice,
		arg.Sku,
		arg.Description,
		arg.StockCount,
		arg.CategoryID,
		arg.DefaultImage,
	)
	return err
}

const createProductSpec = `-- name: CreateProductSpec :exec
insert into product_specifications(
        created_at,
        name,
        description,
        product_id
    )
values(
        current_timestamp,
        $1,
        $2,
        $3
    )
`

type CreateProductSpecParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProductID   int64  `json:"product_id"`
}

func (q *Queries) CreateProductSpec(ctx context.Context, arg CreateProductSpecParams) error {
	_, err := q.db.Exec(ctx, createProductSpec, arg.Name, arg.Description, arg.ProductID)
	return err
}

const createVariant = `-- name: CreateVariant :exec
insert into product_variants(created_at, name, product_id, type)
values(
        current_timestamp,
        $1,
        $2,
        $3
    )
`

type CreateVariantParams struct {
	Name      string `json:"name"`
	ProductID int64  `json:"product_id"`
	Type      int32  `json:"type"`
}

func (q *Queries) CreateVariant(ctx context.Context, arg CreateVariantParams) error {
	_, err := q.db.Exec(ctx, createVariant, arg.Name, arg.ProductID, arg.Type)
	return err
}

const createWishlist = `-- name: CreateWishlist :exec
insert into wishlists(created_at, product_id, user_id)
values(
        current_timestamp,
        $1,
        $2
    )
`

type CreateWishlistParams struct {
	ProductID int64 `json:"product_id"`
	UserID    int64 `json:"user_id"`
}

func (q *Queries) CreateWishlist(ctx context.Context, arg CreateWishlistParams) error {
	_, err := q.db.Exec(ctx, createWishlist, arg.ProductID, arg.UserID)
	return err
}

const deleteCategory = `-- name: DeleteCategory :exec
delete from categories
where id = $1
`

func (q *Queries) DeleteCategory(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteCategory, id)
	return err
}

const deleteProduct = `-- name: DeleteProduct :exec
update products
set deleted_at = current_timestamp
where id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteProduct, id)
	return err
}

const deleteProductSpec = `-- name: DeleteProductSpec :exec
delete from product_specifications
where id = $1
`

func (q *Queries) DeleteProductSpec(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteProductSpec, id)
	return err
}

const deleteVariant = `-- name: DeleteVariant :exec
delete from product_variants
where id = $1
`

func (q *Queries) DeleteVariant(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteVariant, id)
	return err
}

const deleteWishlist = `-- name: DeleteWishlist :exec
delete from wishlists
where id = $1
`

func (q *Queries) DeleteWishlist(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteWishlist, id)
	return err
}

const discountedProducts = `-- name: DiscountedProducts :many
SELECT id, created_at, updated_at, deleted_at, average_rating, name, price, discount_price, sku, description, stock_count, min_stock_count, category_id, total_ratings, total_view, default_image
from products
where products.discount_price > 0
    and products.deleted_at is null
`

func (q *Queries) DiscountedProducts(ctx context.Context) ([]*Products, error) {
	rows, err := q.db.Query(ctx, discountedProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Products{}
	for rows.Next() {
		var i Products
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.AverageRating,
			&i.Name,
			&i.Price,
			&i.DiscountPrice,
			&i.Sku,
			&i.Description,
			&i.StockCount,
			&i.MinStockCount,
			&i.CategoryID,
			&i.TotalRatings,
			&i.TotalView,
			&i.DefaultImage,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCategory = `-- name: GetCategory :one
select id, created_at, updated_at, deleted_at, name, image
from categories
where id = $1
`

func (q *Queries) GetCategory(ctx context.Context, id int64) (*Categories, error) {
	row := q.db.QueryRow(ctx, getCategory, id)
	var i Categories
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Name,
		&i.Image,
	)
	return &i, err
}

const getVariant = `-- name: GetVariant :one
select id, created_at, updated_at, deleted_at, type, name, product_id
from product_variants
where id = $1
`

func (q *Queries) GetVariant(ctx context.Context, id int64) (*ProductVariants, error) {
	row := q.db.QueryRow(ctx, getVariant, id)
	var i ProductVariants
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Type,
		&i.Name,
		&i.ProductID,
	)
	return &i, err
}

const getWishlist = `-- name: GetWishlist :one
select id, created_at, updated_at, deleted_at, product_id, user_id
from wishlists
where user_id = $1
    and product_id = $2
`

type GetWishlistParams struct {
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
}

func (q *Queries) GetWishlist(ctx context.Context, arg GetWishlistParams) (*Wishlists, error) {
	row := q.db.QueryRow(ctx, getWishlist, arg.UserID, arg.ProductID)
	var i Wishlists
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.ProductID,
		&i.UserID,
	)
	return &i, err
}

const listCategories = `-- name: ListCategories :many
select id, created_at, updated_at, deleted_at, name, image
from categories
`

func (q *Queries) ListCategories(ctx context.Context) ([]*Categories, error) {
	rows, err := q.db.Query(ctx, listCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Categories{}
	for rows.Next() {
		var i Categories
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Name,
			&i.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProducts = `-- name: ListProducts :many
select products.id, products.created_at, products.updated_at, products.deleted_at, average_rating, products.name, price, discount_price, sku, description, stock_count, min_stock_count, category_id, total_ratings, total_view, default_image, c.id, c.created_at, c.updated_at, c.deleted_at, c.name, image
from products
    left join categories c on c.id = products.id
where products.deleted_at is null
`

type ListProductsRow struct {
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
	ID_2          pgtype.Int8        `json:"id_2"`
	CreatedAt_2   pgtype.Timestamptz `json:"created_at_2"`
	UpdatedAt_2   pgtype.Timestamptz `json:"updated_at_2"`
	DeletedAt_2   pgtype.Timestamptz `json:"deleted_at_2"`
	Name_2        pgtype.Text        `json:"name_2"`
	Image         pgtype.Text        `json:"image"`
}

func (q *Queries) ListProducts(ctx context.Context) ([]*ListProductsRow, error) {
	rows, err := q.db.Query(ctx, listProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*ListProductsRow{}
	for rows.Next() {
		var i ListProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.AverageRating,
			&i.Name,
			&i.Price,
			&i.DiscountPrice,
			&i.Sku,
			&i.Description,
			&i.StockCount,
			&i.MinStockCount,
			&i.CategoryID,
			&i.TotalRatings,
			&i.TotalView,
			&i.DefaultImage,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.DeletedAt_2,
			&i.Name_2,
			&i.Image,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listVariants = `-- name: ListVariants :many
select id, created_at, updated_at, deleted_at, type, name, product_id
from product_variants
where product_id = $1
`

func (q *Queries) ListVariants(ctx context.Context, productID int64) ([]*ProductVariants, error) {
	rows, err := q.db.Query(ctx, listVariants, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*ProductVariants{}
	for rows.Next() {
		var i ProductVariants
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Type,
			&i.Name,
			&i.ProductID,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listWishlist = `-- name: ListWishlist :many
select wishlists.id, wishlists.created_at, wishlists.updated_at, wishlists.deleted_at, product_id, user_id, p.id, p.created_at, p.updated_at, p.deleted_at, average_rating, name, price, discount_price, sku, description, stock_count, min_stock_count, category_id, total_ratings, total_view, default_image
from wishlists
    left join products p on p.id = product_id
where user_id = $1
`

type ListWishlistRow struct {
	ID            int64              `json:"id"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	UpdatedAt     pgtype.Timestamptz `json:"updated_at"`
	DeletedAt     pgtype.Timestamptz `json:"deleted_at"`
	ProductID     int64              `json:"product_id"`
	UserID        int64              `json:"user_id"`
	ID_2          pgtype.Int8        `json:"id_2"`
	CreatedAt_2   pgtype.Timestamptz `json:"created_at_2"`
	UpdatedAt_2   pgtype.Timestamptz `json:"updated_at_2"`
	DeletedAt_2   pgtype.Timestamptz `json:"deleted_at_2"`
	AverageRating pgtype.Int4        `json:"average_rating"`
	Name          pgtype.Text        `json:"name"`
	Price         pgtype.Int4        `json:"price"`
	DiscountPrice pgtype.Int4        `json:"discount_price"`
	Sku           pgtype.Text        `json:"sku"`
	Description   pgtype.Text        `json:"description"`
	StockCount    pgtype.Int4        `json:"stock_count"`
	MinStockCount pgtype.Int4        `json:"min_stock_count"`
	CategoryID    pgtype.Int8        `json:"category_id"`
	TotalRatings  pgtype.Int4        `json:"total_ratings"`
	TotalView     pgtype.Int4        `json:"total_view"`
	DefaultImage  pgtype.Text        `json:"default_image"`
}

func (q *Queries) ListWishlist(ctx context.Context, userID int64) ([]*ListWishlistRow, error) {
	rows, err := q.db.Query(ctx, listWishlist, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*ListWishlistRow{}
	for rows.Next() {
		var i ListWishlistRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.ProductID,
			&i.UserID,
			&i.ID_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
			&i.DeletedAt_2,
			&i.AverageRating,
			&i.Name,
			&i.Price,
			&i.DiscountPrice,
			&i.Sku,
			&i.Description,
			&i.StockCount,
			&i.MinStockCount,
			&i.CategoryID,
			&i.TotalRatings,
			&i.TotalView,
			&i.DefaultImage,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const oneProduct = `-- name: OneProduct :one
SELECT products.id,
    products.created_at,
    products.updated_at,
    products.name,
    products.price,
    products.discount_price,
    products.sku,
    products.description,
    products.stock_count,
    products.category_id,
    products.default_image,
    products.average_rating,
    products.total_ratings,
    products.total_view,
    products.min_stock_count,
    COALESCE(
        (
            SELECT JSON_AGG(
                    JSON_BUILD_OBJECT(
                        'id',
                        ps.id,
                        'name',
                        ps.name,
                        'description',
                        ps.description
                    )
                )
            FROM product_specifications ps
            WHERE ps.product_id = products.id
        ),
        '[]'::JSON
    ) AS product_specifications,
    COALESCE(
        (
            SELECT JSON_AGG(
                    JSON_BUILD_OBJECT(
                        'id',
                        pr.id,
                        'name',
                        pr.name,
                        'type',
                        pr.type
                    )
                )
            FROM product_variants pr
            WHERE pr.product_id = products.id
        ),
        '[]'::JSON
    ) AS product_variants
FROM products
WHERE products.id = $1
`

type OneProductRow struct {
	ID                    int64              `json:"id"`
	CreatedAt             pgtype.Timestamptz `json:"created_at"`
	UpdatedAt             pgtype.Timestamptz `json:"updated_at"`
	Name                  string             `json:"name"`
	Price                 int32              `json:"price"`
	DiscountPrice         pgtype.Int4        `json:"discount_price"`
	Sku                   string             `json:"sku"`
	Description           string             `json:"description"`
	StockCount            int32              `json:"stock_count"`
	CategoryID            int64              `json:"category_id"`
	DefaultImage          string             `json:"default_image"`
	AverageRating         pgtype.Int4        `json:"average_rating"`
	TotalRatings          int32              `json:"total_ratings"`
	TotalView             pgtype.Int4        `json:"total_view"`
	MinStockCount         int32              `json:"min_stock_count"`
	ProductSpecifications interface{}        `json:"product_specifications"`
	ProductVariants       interface{}        `json:"product_variants"`
}

func (q *Queries) OneProduct(ctx context.Context, id int64) (*OneProductRow, error) {
	row := q.db.QueryRow(ctx, oneProduct, id)
	var i OneProductRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Price,
		&i.DiscountPrice,
		&i.Sku,
		&i.Description,
		&i.StockCount,
		&i.CategoryID,
		&i.DefaultImage,
		&i.AverageRating,
		&i.TotalRatings,
		&i.TotalView,
		&i.MinStockCount,
		&i.ProductSpecifications,
		&i.ProductVariants,
	)
	return &i, err
}

const productSpecifications = `-- name: ProductSpecifications :many
select id, created_at, updated_at, deleted_at, name, description, product_id
from product_specifications
where product_id = $1
`

func (q *Queries) ProductSpecifications(ctx context.Context, productID int64) ([]*ProductSpecifications, error) {
	rows, err := q.db.Query(ctx, productSpecifications, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*ProductSpecifications{}
	for rows.Next() {
		var i ProductSpecifications
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Name,
			&i.Description,
			&i.ProductID,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const productsByCategory = `-- name: ProductsByCategory :many
select id, created_at, updated_at, deleted_at, average_rating, name, price, discount_price, sku, description, stock_count, min_stock_count, category_id, total_ratings, total_view, default_image
from products
where category_id = $1
    and products.id > $2
    and products.deleted_at is null
order by products.created_at desc
limit 50
`

type ProductsByCategoryParams struct {
	CategoryID int64 `json:"category_id"`
	ID         int64 `json:"id"`
}

func (q *Queries) ProductsByCategory(ctx context.Context, arg ProductsByCategoryParams) ([]*Products, error) {
	rows, err := q.db.Query(ctx, productsByCategory, arg.CategoryID, arg.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Products{}
	for rows.Next() {
		var i Products
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.AverageRating,
			&i.Name,
			&i.Price,
			&i.DiscountPrice,
			&i.Sku,
			&i.Description,
			&i.StockCount,
			&i.MinStockCount,
			&i.CategoryID,
			&i.TotalRatings,
			&i.TotalView,
			&i.DefaultImage,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchProducts = `-- name: SearchProducts :many
select id, created_at, updated_at, deleted_at, average_rating, name, price, discount_price, sku, description, stock_count, min_stock_count, category_id, total_ratings, total_view, default_image,
    similarity (name, $1) as score
from products
where similarity (name, $1) > 0.4
order by score desc
`

type SearchProductsRow struct {
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
	Score         float32            `json:"score"`
}

// TODO: Delete stale order_products
func (q *Queries) SearchProducts(ctx context.Context, similarity string) ([]*SearchProductsRow, error) {
	rows, err := q.db.Query(ctx, searchProducts, similarity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchProductsRow{}
	for rows.Next() {
		var i SearchProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.AverageRating,
			&i.Name,
			&i.Price,
			&i.DiscountPrice,
			&i.Sku,
			&i.Description,
			&i.StockCount,
			&i.MinStockCount,
			&i.CategoryID,
			&i.TotalRatings,
			&i.TotalView,
			&i.DefaultImage,
			&i.Score,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCategory = `-- name: UpdateCategory :exec
update categories
SET updated_at = current_timestamp,
    name = COALESCE($1::text, name),
    image = CASE
        WHEN $2::text IS NOT NULL THEN $2::text
        ELSE image
    END
WHERE id = $3
`

type UpdateCategoryParams struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	ID    int64  `json:"id"`
}

func (q *Queries) UpdateCategory(ctx context.Context, arg UpdateCategoryParams) error {
	_, err := q.db.Exec(ctx, updateCategory, arg.Name, arg.Image, arg.ID)
	return err
}

const updateProduct = `-- name: UpdateProduct :exec
update products
set updated_at = current_timestamp,
    name = $2,
    price = $3,
    discount_price = $4,
    sku = $5,
    description = $6,
    stock_count = $7,
    category_id = $8,
    default_image = $9
where id = $1
    and deleted_at is null
`

type UpdateProductParams struct {
	ID            int64       `json:"id"`
	Name          string      `json:"name"`
	Price         int32       `json:"price"`
	DiscountPrice pgtype.Int4 `json:"discount_price"`
	Sku           string      `json:"sku"`
	Description   string      `json:"description"`
	StockCount    int32       `json:"stock_count"`
	CategoryID    int64       `json:"category_id"`
	DefaultImage  string      `json:"default_image"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.Exec(ctx, updateProduct,
		arg.ID,
		arg.Name,
		arg.Price,
		arg.DiscountPrice,
		arg.Sku,
		arg.Description,
		arg.StockCount,
		arg.CategoryID,
		arg.DefaultImage,
	)
	return err
}

const updateProductSpec = `-- name: UpdateProductSpec :exec
update product_specifications
set updated_at = current_timestamp,
    name = $2,
    description = $3,
    product_id = $4
where id = $1
`

type UpdateProductSpecParams struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProductID   int64  `json:"product_id"`
}

func (q *Queries) UpdateProductSpec(ctx context.Context, arg UpdateProductSpecParams) error {
	_, err := q.db.Exec(ctx, updateProductSpec,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.ProductID,
	)
	return err
}

const updateVariant = `-- name: UpdateVariant :exec
update product_variants
set updated_at = current_timestamp,
    name = $2,
    product_id = $3,
    type = $4
where id = $1
`

type UpdateVariantParams struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ProductID int64  `json:"product_id"`
	Type      int32  `json:"type"`
}

func (q *Queries) UpdateVariant(ctx context.Context, arg UpdateVariantParams) error {
	_, err := q.db.Exec(ctx, updateVariant,
		arg.ID,
		arg.Name,
		arg.ProductID,
		arg.Type,
	)
	return err
}
