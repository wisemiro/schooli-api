-- name: CreateCategory :exec
insert into categories(created_at, name, image)
values(current_timestamp, @name, @image);
-- 
-- 
-- name: UpdateCategory :exec
update categories
SET updated_at = current_timestamp,
    name = COALESCE(@name::text, name),
    image = CASE
        WHEN @image::text IS NOT NULL THEN @image::text
        ELSE image
    END
WHERE id = @id;
-- 
-- 
-- name: ListCategories :many
select *
from categories;
-- 
-- 
-- name: GetCategory :one
select *
from categories
where id = $1;
-- 
-- 
-- name: DeleteCategory :exec
delete from categories
where id = $1;
-- 
-- 
-- name: CreateProduct :exec
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
        @name,
        @price,
        @discount_price,
        @sku,
        @description,
        @stock_count,
        @category_id,
        @default_image
    );
-- 
-- 
-- name: UpdateProduct :exec
update products
set updated_at = current_timestamp,
    name = @name,
    price = @price,
    discount_price = @discount_price,
    sku = @sku,
    description = @description,
    stock_count = @stock_count,
    category_id = @category_id,
    default_image = @default_image
where id = $1;
-- 
-- 
-- name: ListProducts :many
select *
from products
    left join categories c on c.id = products.id;
-- 
-- 
-- name: DiscountedProducts :many
SELECT *
from products
where products.discount_price > 0;
-- name: ProductsByCategory :many
select *
from products
where category_id = $1
    and products.id > $2
order by products.created_at desc
limit 50;
-- 
-- 
-- name: OneProduct :one
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
WHERE products.id = $1;
-- 
-- 
-- name: DeleteProduct :exec
delete from products
where id = $1;
-- 
-- 
-- name: CreateWishlist :exec
insert into wishlists(created_at, product_id, user_id)
values(
        current_timestamp,
        @product_id,
        @user_id
    );
-- 
-- 
-- name: DeleteWishlist :exec
delete from wishlists
where id = $1;
-- 
-- 
-- name: ListWishlist :many
select *
from wishlists
    left join products p on p.id = product_id
where user_id = $1;
-- 
-- 
-- name: GetWishlist :one
select *
from wishlists
where user_id = $1
    and product_id = $2;
-- 
-- 
-- name: CreateVariant :exec
insert into product_variants(created_at, name, product_id, type)
values(
        current_timestamp,
        @name,
        @product_id,
        @type
    );
-- 
-- 
-- name: UpdateVariant :exec
update product_variants
set updated_at = current_timestamp,
    name = @name,
    product_id = @product_id,
    type = @type
where id = $1;
-- 
-- 
-- name: DeleteVariant :exec
delete from product_variants
where id = $1;
-- 
-- 
-- name: ListVariants :many
select *
from product_variants
where product_id = $1;
-- 
-- 
-- name: GetVariant :one
select *
from product_variants
where id = $1;
-- 
-- 
-- name: CreateProductSpec :exec
insert into product_specifications(
        created_at,
        name,
        description,
        product_id
    )
values(
        current_timestamp,
        @name,
        @description,
        @product_id
    );
-- 
-- 
-- name: UpdateProductSpec :exec
update product_specifications
set updated_at = current_timestamp,
    name = @name,
    description = @description,
    product_id = @product_id
where id = $1;
-- 
-- 
-- name: DeleteProductSpec :exec
delete from product_specifications
where id = $1;
-- 
-- 
-- name: ProductSpecifications :many
select *
from product_specifications
where product_id = $1;
-- 
-- TODO: Delete stale order_products