-- name: CreateProduct :one
INSERT INTO products(name, description, sku, price, stock)
VALUES ($1, $2, $3, $4, $5) RETURNING id;

-- name: ListProducts :many
SELECT * FROM products;

-- name: ListProductPrices :many
SELECT id, price FROM products
WHERE id = sqlc.slice('ids');

-- name: GetProductById :one
SELECT * FROM products 
WHERE id = $1
LIMIT 1;

-- name: GetProductBySku :one
SELECT * FROM products 
WHERE sku = $1
LIMIT 1;

-- name: UpdateProduct :exec
UPDATE products
	SET 
		name = $2,
		description = $3,
		sku = $4,
		price = $5,
		stock = $6
	WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
	WHERE id = $1;