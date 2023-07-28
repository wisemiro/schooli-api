-- name: ListUsers :many
select *
from users;
-- 
-- 
-- name: GetUser :one
select *
from users
where id = $1;
-- 
-- 
-- name: ListByVerification :many
select *
from users
where is_verified = $1;
-- 
-- 
-- name: CreateUser :exec
insert into users(created_at, email, phone_number, password_hash)
values(
        current_timestamp,
        @email,
        @phone_number,
        @password_hash
    );
-- 
-- 
-- name: UpdateUser :exec
update users
set updated_at = current_timestamp,
    email = @email,
    phone_number = @phone_number
where id = $1;
-- 
--
-- name: DeleteUser :exec
delete from users
where id = $1;
-- 
-- 
-- name: UserByEmail :one
SELECT is_verified,
    created_at,
    email,
    id,
    password_hash,
    phone_number
FROM users
WHERE email = @email
    AND "users"."deleted_at" IS NULL
ORDER BY "users"."id"
LIMIT 1;
-- 
-- 
-- name: CreateDevice :exec
insert into devices(user_id, device)
values(@user_id, @device);
-- 
-- 
-- name: GetOneDevice :one 
select *
from devices
where devices.user_id = $1;
-- 
-- 
-- name: UpdateDevice :exec
update devices
set device = @device
where user_id = $1;