-- Orders queries
-- name: GetOrderByID :one
SELECT order_id, user_id, store_id, total_amount, status, delivery_address, order_date 
FROM orders WHERE order_id = $1;

-- name: CreateOrder :one
INSERT INTO orders (user_id, store_id, total_amount, delivery_address) 
VALUES ($1, $2, $3, $4) 
RETURNING order_id, user_id, store_id, total_amount, status, delivery_address, order_date;

-- name: ListOrdersByUser :many
SELECT o.order_id, o.user_id, o.store_id, o.total_amount, o.status, o.delivery_address, o.order_date,
       s.name as store_name
FROM orders o
JOIN stores s ON o.store_id = s.store_id
WHERE o.user_id = $1 
ORDER BY o.order_date DESC;

-- name: ListOrdersByStore :many
SELECT o.order_id, o.user_id, o.store_id, o.total_amount, o.status, o.delivery_address, o.order_date,
       u.name as user_name, u.email as user_email
FROM orders o
JOIN users u ON o.user_id = u.id
WHERE o.store_id = $1 
ORDER BY o.order_date DESC;

-- name: UpdateOrderStatus :one
UPDATE orders 
SET status = $2 
WHERE order_id = $1 
RETURNING order_id, user_id, store_id, total_amount, status, delivery_address, order_date;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE order_id = $1;