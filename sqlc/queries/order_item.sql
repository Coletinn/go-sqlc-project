-- Order Items queries
-- name: GetOrderItemsByOrderID :many
SELECT oi.order_item_id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price, oi.total_price,
       p.name as product_name, p.sku
FROM order_items oi
JOIN products p ON oi.product_id = p.product_id
WHERE oi.order_id = $1
ORDER BY oi.order_item_id;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, unit_price, total_price) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING order_item_id, order_id, product_id, quantity, unit_price, total_price;

-- name: UpdateOrderItem :one
UPDATE order_items 
SET quantity = $3, unit_price = $4, total_price = $5 
WHERE order_item_id = $1 AND order_id = $2 
RETURNING order_item_id, order_id, product_id, quantity, unit_price, total_price;

-- name: DeleteOrderItem :exec
DELETE FROM order_items WHERE order_item_id = $1;