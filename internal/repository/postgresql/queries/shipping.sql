-- name: CreateShipping :exec
insert into shipping(
        created_at,
        location,
        address,
        apartment,
        phone_number,
        user_id,
        order_id,
        status
    )
values(
        current_timestamp,
        @location,
        @address,
        @apartment,
        @phone_number,
        @user_id,
        @order_id,
        @status
    );
-- 
-- 
-- name: UpdateShipping :exec
update shipping
set updated_at = current_timestamp,
    phone_number = @phone_number,
    status = @status
where id = $1;
-- 
-- TODO
-- name: ListShipping :many
select *
from shipping
where status = $1;