-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: ListOrders :many
SELECT * FROM orders;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders(status, total_amount, tax, client_name, client_email, client_phone, note)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;

-- name: CreateOrderItems :exec
INSERT INTO order_items(unit_price, product_id, quantity, order_id)
SELECT p.price, p.id, $2, $3
FROM products p
WHERE p.id in ($1);
