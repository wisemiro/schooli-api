-- name: CreateCategory :exec
insert into categories(created_at, name, image)
values(current_timestamp, @name, @image);
-- 
-- 
-- name: UpdateCategory :exec
update categories
set updated_at = current_timestamp,
    name = @name,
    image = @image
where id = $1;
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
-- name: ListOrderProducts :many
select p.id,
    p.created_at,
    p.updated_at,
    p.name,
    p.price,
    p.discount_price,
    p.sku,
    p.description,
    p.stock_count,
    p.category_id,
    p.default_image,
    pv.id,
    pv.name,
    order_products.id,
    order_products.device_id,
    order_products.total_price,
    order_products.quantity
from order_products
    left join product_variants pv on pv.id = order_products.product_variant
    left join products p on p.id = order_products.product_id;
-- 
-- 
-- name: GetOrderProduct :one
select p.id,
    p.created_at,
    p.updated_at,
    p.name,
    p.price,
    p.discount_price,
    p.sku,
    p.description,
    p.stock_count,
    p.category_id,
    p.default_image,
    pv.id,
    pv.name,
    order_products.id,
    order_products.total_price,
    order_products.device_id,
    order_products.quantity
from order_products
    left join product_variants pv on pv.id = order_products.product_variant
    left join products p on p.id = order_products.product_id
where order_products.id = $1;
-- 
-- 
-- name: CreateOrderProduct :exec
insert into order_products(
        created_at,
        quantity,
        total_price,
        product_variant,
        product_id,
        device_id
    )
values(
        current_timestamp,
        @quantity,
        @total_price,
        @product_variant,
        @product_id,
        @device_id
    );
-- 
-- 
-- name: UpdateOrderProduct :exec
update order_products
set updated_at = current_timestamp,
    quantity = @quantity,
    total_price = @total_price,
    product_variant = @product_variant
where id = $1;
-- 
-- 
-- name: DeleteOrderProduct :exec
delete from order_products
where id = $1;
-- 
-- 
-- name: CreateOrder :exec
insert into orders(
        created_at,
        grand_total,
        serial_number,
        order_products_id,
        user_id
    )
values(
        current_timestamp,
        @grand_total,
        @serial_number,
        @order_products_id,
        @user_id
    );
-- 
-- 
-- name: UpdateOrder :exec
update orders
set updated_at = current_timestamp,
    grand_total = @grand_total
where order_products_id = $1;
-- 
-- 
-- name: ListOrders :many
select orders.created_at,
    orders.grand_total,
    orders.serial_number,
    orders.order_products_id,
    op.quantity,
    op.total_price,
    op.device_id,
    op.product_variant,
    p.id,
    p.name,
    p.price,
    p.discount_price,
    p.sku,
    p.default_image,
    u.id,
    u.email,
    u.phone_number,
    pv.name
from orders
    left join order_products op on op.id = orders.order_products_id
    left join products p on p.id = op.product_id
    left join product_variants pv on pv.id = op.product_variant
    left join users u on u.id = orders.user_id;
-- 
-- 
-- name: ListUserOrders :many
select orders.created_at,
    orders.grand_total,
    orders.serial_number,
    orders.order_products_id,
    op.quantity,
    op.total_price,
    op.device_id,
    op.product_variant,
    p.id,
    p.name,
    p.price,
    p.discount_price,
    p.sku,
    p.default_image,
    u.id,
    u.email,
    u.phone_number,
    pv.name
from orders
    left join order_products op on op.id = orders.order_products_id
    left join products p on p.id = op.product_id
    left join product_variants pv on pv.id = op.product_variant
    left join users u on u.id = orders.user_id
where orders.user_id = $1;
-- 
-- 
-- name: DeleteOrder :exec
delete from orders
where order_products_id = $1;
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