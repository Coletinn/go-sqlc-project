-- Stores queries
-- name: GetStoreByID :one
SELECT store_id, name, address, phone, created_at FROM stores WHERE store_id = $1;

-- name: CreateStore :one
INSERT INTO stores (name, address, phone) 
VALUES ($1, $2, $3) 
RETURNING store_id, name, address, phone, created_at;

-- name: ListStores :many
SELECT store_id, name, address, phone, created_at FROM stores 
ORDER BY created_at DESC;

-- name: UpdateStore :one
UPDATE stores 
SET name = $2, address = $3, phone = $4 
WHERE store_id = $1 
RETURNING store_id, name, address, phone, created_at;

-- name: DeleteStore :exec
DELETE FROM stores WHERE store_id = $1;