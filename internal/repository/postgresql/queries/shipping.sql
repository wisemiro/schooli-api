-- name: CreateShipping :exec
insert into shipping(
        created_at,
        location,
        user_id,
        order_id,
        status
    )
values(
        current_timestamp,
        ST_GeomFromText($1, 4269),
        @user_id,
        @order_id,
        @status
    );
-- 
-- 
-- name: UpdateShipping :exec
update shipping
set updated_at = current_timestamp,
    location = ST_GeomFromText($1, 4269),
    status = @status
where id = $1;
-- 
-- name: ListShipping :many
select *
from shipping
where status = $1;