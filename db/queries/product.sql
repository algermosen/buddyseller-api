-- name: CreateProduct :one
INSERT INTO products(name, description, sku, price, stock)
VALUES ($1, $2, $3, $4, $5) RETURNING id;

-- name: ListProducts :many
SELECT * FROM products;

-- name: ListProductsToOrder :many
SELECT id, price, stock FROM products
WHERE id = ANY($1::int[]);

-- name: GetProductById :one
SELECT * FROM products 
WHERE id = $1
LIMIT 1;

-- name: GetProductBySku :one
SELECT * FROM products 
WHERE sku = $1
LIMIT 1;

-- name: UpdateProduct :execrows
UPDATE products
	SET 
		name = $2,
		description = $3,
		sku = $4,
		price = $5,
		stock = $6
	WHERE id = $1;

-- name: UpdateStock :exec
UPDATE products
	SET 
		stock = $2
	WHERE id = $1;

-- name: DeleteProduct :execrows
DELETE FROM products
	WHERE id = $1;

