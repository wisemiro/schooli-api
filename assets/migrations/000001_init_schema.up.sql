create table if not exists schema_migrations (
    version bigint not null primary key,
    dirty boolean not null
);
-- 
-- 
--  User
create table if not exists users (
    id bigserial primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    email varchar(100) not null unique,
    password_hash varchar(100) not null,
    phone_number text not null unique,
    is_verified boolean default true
);
-- 
-- indexing
CREATE UNIQUE INDEX users_email_idx ON users (email);
CREATE INDEX users_phone_number_idx ON users (phone_number);
CREATE INDEX users_is_verified_idx ON users (is_verified);
-- 
-- 
-- Devices
create table if not exists devices (
    id bigserial primary key,
    created_at timestamp with time zone default current_timestamp,
    user_id bigint not null constraint fk_user references users,
    device text not null
);
-- 
-- indexing
CREATE UNIQUE INDEX devices_id_idx ON devices (id);
CREATE UNIQUE INDEX devices_device_idx ON devices (device);
-- 
-- 
-- Categories
create table if not exists categories (
    id bigserial primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text not null unique,
    image text not null
);
-- 
-- 
-- indexing
CREATE UNIQUE INDEX categories_name_idx on categories (name);
-- 
-- 
-- Product 
create table if not exists products (
    id bigserial primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    average_rating int default 0,
    name text not null,
    price integer not null,
    discount_price integer default 0,
    sku text not null,
    description text not null,
    stock_count integer not null,
    min_stock_count integer not null,
    category_id bigint not null constraint fk_products_category_id references categories,
    total_ratings integer default 0 not null,
    total_view integer default 0,
    default_image text not null
);
-- 
-- indexing
CREATE UNIQUE INDEX products_name_idx on products (name);
CREATE INDEX products_price_idx on products (price);
CREATE INDEX products_stock_count_idx on products (stock_count);
CREATE INDEX products_discount_price_idx on products (discount_price);
-- 
-- 
--  Product variants
create table if not exists product_variants(
    id bigserial primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    type integer not null,
    name text not null,
    product_id bigint not null constraint fk_product_variants_product_id references products
);
-- 
-- indexing
CREATE UNIQUE INDEX product_variants_name_idx on product_variants (name);
-- 
-- 
-- Ordered products
create table if not exists order_products(
    id bigserial primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    quantity integer not null,
    total_price integer not null,
    device_id text not null,
    product_variant bigint constraint fk_order_products_product_variant references product_variants,
    product_id bigint not null constraint fk_order_products_product_id references products
);
-- Create trigger function
CREATE OR REPLACE FUNCTION update_product_stock() RETURNS TRIGGER AS $$ BEGIN IF TG_OP = 'INSERT' THEN -- decrease product stock when a new order item is created
UPDATE products
SET stock_count = stock_count - NEW.quantity
WHERE id = NEW.product_id;
ELSIF TG_OP = 'UPDATE' THEN -- Calculate the difference in quantity for an updated order item
-- and update the product stock accordingly
UPDATE products
SET stock_count = stock_count + OLD.quantity - NEW.quantity
WHERE id = NEW.product_id;
ELSIF TG_OP = 'DELETE' THEN -- Increase product stock when an order item is deleted
UPDATE products
SET stock_count = stock_count + OLD.quantity
WHERE id = OLD.product_id;
END IF;
RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- Create trigger on order_products table
CREATE TRIGGER update_product_stock_trigger
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON order_products FOR EACH ROW EXECUTE FUNCTION update_product_stock();
-- 
-- Orders
create table if not exists orders(
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    grand_total integer not null,
    serial_number text not null,
    order_products_id bigint not null constraint fk_orders_order_products_id references order_products,
    user_id bigint constraint fk_orders_user_id references users,
    primary key (order_products_id)
);
-- 
-- indexing
CREATE UNIQUE INDEX orders_order_product_id_idx on orders (order_products_id);
-- 
-- 
-- Shipping
create table if not exists shipping(
    id bigserial primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    location geometry(POINT) not null,
    address text,
    apartment text,
    phone_number text not null,
    user_id bigint constraint fk_shipping_user_id references users,
    order_id bigint not null constraint fk_shipping_order_id references orders,
    location geometry(POINT) not null,
    status text default 'pending'
);
-- 
-- indexing
CREATE UNIQUE INDEX shipping_order_id_idx on shipping (order_id);
CREATE INDEX shipping_status_idx on shipping (status);
CREATE INDEX users_location_idx ON users USING GIST (location);
-- TODO: Shipping fee and pick_up? 
-- 
-- 
--  Wishlist 
create table if not exists wishlists(
    id bigserial primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    product_id bigint not null constraint fk_wishlists_product_id references products,
    user_id bigint not null constraint fk_wishlists_user_id references users
);
-- 
-- indexing
CREATE INDEX wishlists_product_id_idx on wishlists (product_id);
CREATE INDEX wishlists_user_id_idx on wishlists (user_id);
-- 
-- 
-- Ratings
create table if not exists product_ratings(
    id bigserial primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    user_id bigint not null constraint fk_product_ratings_user_id references users,
    stars integer not null,
    feedback text default ''::text,
    product_id bigint not null constraint fk_product_ratings_product_id references products
);
-- 
-- indexing
CREATE INDEX product_ratings_user_id_idx on product_ratings (user_id);
-- 
--  Trigger
CREATE OR REPLACE FUNCTION gen_average_trigger_fnc() RETURNS trigger AS $$
DECLARE ratings_avg integer;
BEGIN
SELECT avg(stars)
FROM product_ratings
where product_ratings.product_id = new.product_id into ratings_avg;
update products
set average_rating = ratings_avg
where products.id = new.product_id;
return new;
END $$ LANGUAGE 'plpgsql';
CREATE TRIGGER rating_insert_trigger
AFTER
INSERT
    OR
UPDATE ON product_ratings FOR EACH ROW EXECUTE PROCEDURE gen_average_trigger_fnc();
CREATE OR REPLACE FUNCTION gen_total_trigger_fnc() RETURNS TRIGGER AS $$
DECLARE ratings_total integer;
BEGIN
SELECT count(*) INTO ratings_total
FROM product_ratings
WHERE product_id = NEW.product_id;
UPDATE products
SET total_ratings = ratings_total
WHERE id = NEW.product_id;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER total_ratings_insert_trigger
AFTER
INSERT
    OR
UPDATE ON product_ratings FOR EACH ROW EXECUTE FUNCTION gen_total_trigger_fnc();
-- 
-- product specifications
create table if not exists product_specifications(
    id bigserial primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text not null,
    description text not null,
    product_id bigint not null constraint fk_product_specifications_product_id references products
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- Enable PostGIS (as of 3.0 contains just geometry/geography)
CREATE EXTENSION postgis;
-- enable raster support (for 3+)
CREATE EXTENSION postgis_raster;
-- Enable Topology
CREATE EXTENSION postgis_topology;
-- Enable PostGIS Advanced 3D
-- and other geoprocessing algorithms
-- sfcgal not available with all distributions
CREATE EXTENSION postgis_sfcgal;
-- fuzzy matching needed for Tiger
CREATE EXTENSION fuzzystrmatch;
-- rule based standardizer
CREATE EXTENSION address_standardizer;
-- Enable US Tiger Geocoder
CREATE EXTENSION postgis_tiger_geocoder;