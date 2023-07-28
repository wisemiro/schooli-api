
-- name: ListOrderProducts :many
select p.id,
    p.created_at,
    p.updated_at,
    p.name,
    p.price,
    p.discount_price,
    p.sku,
    p.description,
    p.category_id,
    p.default_image,
    order_products.id,
    order_products.total_price,
    order_products.quantity,
        (
        SELECT json_agg(
            json_build_object(
                'id', v.id,
                'name', v.name,
                'type', v.type
            )
        )
        FROM product_variants v
        WHERE v.id = ANY(order_products.product_variants)
    ) AS product_variants
from order_products
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
    p.category_id,
    p.default_image,
    order_products.id,
    order_products.total_price,
    order_products.quantity,
        (
        SELECT json_agg(
            json_build_object(
                'id', v.id,
                'name', v.name,
                'type', v.type
            )
        )
        FROM product_variants v
        WHERE v.id = ANY(order_products.product_variants)
    ) AS product_variants
from order_products
    left join products p on p.id = order_products.product_id
where order_products.id = $1;
-- 
-- 
-- name: CreateOrderProduct :exec
insert into order_products(
        created_at,
        quantity,
        total_price,
        product_variants,
        product_id,
        order_id
    )
values(
        current_timestamp,
        @quantity,
        @total_price,
        @product_variants,
        @product_id,
        @order_id
    );
-- 
-- 
-- name: UpdateOrderProduct :exec
update order_products
set updated_at = current_timestamp,
    quantity = @quantity,
    total_price = @total_price,
    product_variants = @product_variants
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
        shipping_address,
        user_id
    )
values(
        current_timestamp,
        @grand_total,
        @serial_number,
        @shipping_address,
        @user_id
    );
-- 
-- 
-- name: UpdateOrder :exec
update orders
set updated_at = current_timestamp,
    grand_total = @grand_total,
    confirmed = @confirmed
where id = $1;
-- 
-- 
-- name: ListOrders :many
select
  p.id,
  p.created_at,
  p.updated_at,
  p.name,
  p.price,
  p.discount_price,
  p.sku,
  p.description,
  p.category_id,
  p.default_image,
  order_products.id,
  order_products.total_price,
  order_products.quantity,
  o.id,
  o.created_at,
  o.grand_total,
  o.serial_number,
  o.confirmed,
  u.id,
  u.email,
  u.phone_number,
  s.id,
  st_x(s.location) as latitude,
  st_y(s.location) as longitude,
  (
    SELECT
      json_agg(
        json_build_object('id', v.id, 'name', v.name, 'type', v.type)
      )
    FROM
      product_variants v
    WHERE
      v.id = ANY(order_products.product_variants)
  ) AS product_variants
from
  order_products
  left join products p on p.id = order_products.product_id
  left join users u on u.id = order_products.user_id
  left join shipping s on s.id = orders.shipping_address
  left join orders o on o.id = order_products.order_id;
-- 
-- 
-- name: ListUserOrders :many
select
  p.id,
  p.created_at,
  p.updated_at,
  p.name,
  p.price,
  p.discount_price,
  p.sku,
  p.description,
  p.category_id,
  p.default_image,
  order_products.id,
  order_products.total_price,
  order_products.quantity,
  o.id,
  o.created_at,
  o.grand_total,
  o.serial_number,
  o.confirmed,
  u.id,
  u.email,
  u.phone_number,
  s.id,
  st_x(s.location) as latitude,
  st_y(s.location) as longitude,
  (
    SELECT
      json_agg(
        json_build_object('id', v.id, 'name', v.name, 'type', v.type)
      )
    FROM
      product_variants v
    WHERE
      v.id = ANY(order_products.product_variants)
  ) AS product_variants
from
  order_products
  left join products p on p.id = order_products.product_id
  left join users u on u.id = order_products.user_id
  left join shipping s on s.id = orders.shipping_address
  left join orders o on o.id = order_products.order_id
where orders.user_id = $1;
-- 
-- 
-- name: DeleteOrder :exec
delete from orders
where id = $1;