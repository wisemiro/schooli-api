-- name: CreateRating :exec
insert into product_ratings (
        created_at,
        user_id,
        stars,
        feedback,
        product_id
    )
values(
        current_timestamp,
        @user_id,
        @stars,
        @feedback,
        @product_id
    );
-- 
-- 
-- name: UpdateRating :exec
update product_ratings
set updated_at = current_timestamp,
    stars = @stars,
    feedback = @feedback
where id = $1;
-- 
-- 
-- name: ListProductRatings :many
select product_ratings.id,
    product_ratings.created_at,
    product_ratings.stars,
    product_ratings.feedback,
    product_ratings.user_id,
    product_ratings.product_id,
    u.email,
    p.name
from product_ratings
    left join users u on u.id = product_ratings.user_id
    left join products p on p.id = product_ratings.product_id
where product_id = $1;
-- 
-- 
-- name: GetRating :one
select product_ratings.id,
    product_ratings.created_at,
    product_ratings.stars,
    product_ratings.feedback,
    product_ratings.user_id,
    product_ratings.product_id,
    u.email,
    p.name
from product_ratings
    left join users u on u.id = product_ratings.user_id
    left join products p on p.id = product_ratings.product_id
where product_ratings.id = $1;
-- 
-- 
-- name: DeleteRatings :exec
delete from product_ratings
where id = $1;