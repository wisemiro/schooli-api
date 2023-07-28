-- name: CreateShipping :exec
insert into shipping(
        created_at,
        location,
        user_id
    )
values(
        current_timestamp,
        ST_GeomFromText($1, 4269),
        @user_id
    );
-- 
-- 
-- name: UpdateShipping :exec
update shipping
set updated_at = current_timestamp,
    location = ST_GeomFromText($1, 4269)
where shipping.id = $2;
-- 
-- 
-- name: ListShipping :many
select shipping.id,
    shipping.created_at,
    shipping.updated_at,
    st_x(shipping.location) as latitude,
    st_y(shipping.location) as longitude,
    u.id,
    u.email,
    u.phone_number
from shipping
    left join users u on u.id = shipping.user_id;
-- 
-- 
-- name: ListUserShipping :many
select shipping.id,
    shipping.created_at,
    shipping.updated_at,
    st_x(shipping.location) as latitude,
    st_y(shipping.location) as longitude,
    u.id,
    u.email,
    u.phone_number
from shipping
    left join users u on u.id = shipping.user_id
where shipping.user_id = $1;