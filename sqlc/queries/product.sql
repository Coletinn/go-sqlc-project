-- Products queries
-- name: GetProductByID :one
SELECT product_id, name, description, price, sku, category, created_at 
FROM products WHERE product_id = $1;

-- name: GetProductBySKU :one
SELECT product_id, name, description, price, sku, category, created_at 
FROM products WHERE sku = $1;

-- name: CreateProduct :one
INSERT INTO products (name, description, price, sku, category) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING product_id, name, description, price, sku, category, created_at;

-- name: ListProducts :many
SELECT product_id, name, description, price, sku, category, created_at 
FROM products 
ORDER BY created_at DESC;

-- name: ListProductsByCategory :many
SELECT product_id, name, description, price, sku, category, created_at 
FROM products 
WHERE category = $1 
ORDER BY name;

-- name: UpdateProduct :one
UPDATE products 
SET name = $2, description = $3, price = $4, category = $5 
WHERE product_id = $1 
RETURNING product_id, name, description, price, sku, category, created_at;

-- name: DeleteProduct :exec
DELETE FROM products WHERE product_id = $1;