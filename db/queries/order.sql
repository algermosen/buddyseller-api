-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: ListOrders :many
SELECT * FROM orders;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2::order_status,
    shipped = CASE
        WHEN $2::order_status = 'shipped' THEN NOW () ELSE shipped END,
    delivered = CASE
        WHEN $2::order_status = 'delivered' THEN NOW () ELSE delivered END,
    cancelled = CASE
        WHEN $2::order_status = 'cancelled' THEN NOW () ELSE cancelled END
WHERE id = $1;

-- name: CancelOrder :exec
UPDATE orders
SET status = 'cancelled',
    cancelled = NOW (),
    cancellation_reason = $2
WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders(total_amount, tax, user_id, client_name, client_email, client_phone, note)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- name: CreateOrderItems :copyfrom
INSERT INTO order_items(unit_price, product_id, quantity, order_id)
VALUES ($1, $2, $3, $4);

-- name: GetOrderItemsDetail :many
SELECT p.name, oi.unit_price, oi.quantity FROM order_items oi
JOIN products p ON oi.product_id = p.id
WHERE oi.order_id = $1;