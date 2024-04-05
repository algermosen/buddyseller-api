-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: ListOrders :many
SELECT * FROM orders;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2
WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders(total_amount, tax, client_name, client_email, client_phone, note)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;

-- name: CreateOrderItems :copyfrom
INSERT INTO order_items(unit_price, product_id, quantity, order_id)
VALUES ($1, $2, $3, $4);
