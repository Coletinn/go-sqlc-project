-- Store Inventory queries
-- name: GetInventoryByStore :many
SELECT si.inventory_id, si.store_id, si.product_id, si.quantity, si.last_updated,
       p.name as product_name, p.price, p.sku 
FROM store_inventory si
JOIN products p ON si.product_id = p.product_id
WHERE si.store_id = $1
ORDER BY p.name;

-- name: GetInventoryByStoreAndProduct :one
SELECT inventory_id, store_id, product_id, quantity, last_updated 
FROM store_inventory 
WHERE store_id = $1 AND product_id = $2;

-- name: CreateInventoryItem :one
INSERT INTO store_inventory (store_id, product_id, quantity) 
VALUES ($1, $2, $3) 
RETURNING inventory_id, store_id, product_id, quantity, last_updated;

-- name: UpdateInventoryQuantity :one
UPDATE store_inventory 
SET quantity = $3, last_updated = CURRENT_TIMESTAMP 
WHERE store_id = $1 AND product_id = $2 
RETURNING inventory_id, store_id, product_id, quantity, last_updated;

-- name: DeleteInventoryItem :exec
DELETE FROM store_inventory WHERE store_id = $1 AND product_id = $2;

-- name: GetProductsAvailableInStore :many
SELECT p.product_id, p.name, p.description, p.price, p.sku, p.category, si.quantity
FROM products p
JOIN store_inventory si ON p.product_id = si.product_id
WHERE si.store_id = $1 AND si.quantity > 0
ORDER BY p.name;