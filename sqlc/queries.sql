-- Users queries
-- name: GetUserByID :one
SELECT id, name, email, phone, created_at FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, phone) 
VALUES ($1, $2, $3) 
RETURNING id, name, email, phone, created_at;

-- name: ListUsers :many
SELECT id, name, email, phone, created_at FROM users 
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users 
SET name = $2, email = $3, phone = $4 
WHERE id = $1 
RETURNING id, name, email, phone, created_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

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

-- Complex queries
-- name: GetOrderWithItems :many
SELECT 
    o.order_id, o.user_id, o.store_id, o.total_amount, o.status, o.delivery_address, o.order_date,
    u.name as user_name, u.email as user_email,
    s.name as store_name, s.address as store_address,
    oi.order_item_id, oi.product_id, oi.quantity, oi.unit_price, oi.total_price,
    p.name as product_name, p.sku as product_sku
FROM orders o
JOIN users u ON o.user_id = u.id
JOIN stores s ON o.store_id = s.store_id
LEFT JOIN order_items oi ON o.order_id = oi.order_id
LEFT JOIN products p ON oi.product_id = p.product_id
WHERE o.order_id = $1
ORDER BY oi.order_item_id;

-- name: GetTopSellingProducts :many
SELECT 
    p.product_id, p.name, p.price, p.sku, p.category,
    SUM(oi.quantity) as total_sold,
    COUNT(DISTINCT oi.order_id) as total_orders
FROM products p
JOIN order_items oi ON p.product_id = oi.product_id
JOIN orders o ON oi.order_id = o.order_id
WHERE o.status IN ('confirmed', 'processing', 'shipped', 'delivered')
GROUP BY p.product_id, p.name, p.price, p.sku, p.category
ORDER BY total_sold DESC
LIMIT $1;

-- name: GetUserOrderSummary :many
SELECT 
    u.id as user_id, u.name, u.email,
    COUNT(o.order_id) as total_orders,
    SUM(o.total_amount) as total_spent,
    MAX(o.order_date) as last_order_date
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name, u.email
ORDER BY total_spent DESC NULLS LAST;