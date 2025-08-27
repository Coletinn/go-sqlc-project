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