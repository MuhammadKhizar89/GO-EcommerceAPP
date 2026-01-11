-- name: GetProducts :many
SELECT *
FROM products;

-- name: FindProductByID :one
SELECT *
FROM products
WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (customer_id) VALUES ($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateProduct :one
INSERT INTO products (name, price, quantity) VALUES ($1, $2, $3) RETURNING *;

-- name: GetOrdersByCustomerID :many
SELECT id, created_at
FROM orders
WHERE customer_id = $1
ORDER BY created_at DESC;

-- name: GetOrderItemsByOrderID :many
SELECT id, order_id, product_id, quantity, price
FROM order_items
WHERE order_id = $1;

-- name: FindProductsByIDs :many
SELECT id, name, price, quantity
FROM products
WHERE id = ANY($1::int[]);
