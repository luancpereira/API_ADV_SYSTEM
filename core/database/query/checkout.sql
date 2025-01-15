-----------------
---- INSERTS ----
-----------------

-- name: InsertTransaction :one
INSERT INTO "order" (
    description,
    transaction_date,
    transaction_value
) VALUES (
    @description::VARCHAR,
    @transaction_date::TIMESTAMP,
    @transaction_value::FLOAT
) RETURNING id;

-----------------
---- INSERTS ----
-----------------

-----------------
---- SELECTS ----
-----------------

-- name: SelectTransactions :many
SELECT 
    id,
    description,
    transaction_date::TIMESTAMP AS transaction_date,
    transaction_value
FROM
    "order"
WHERE
	(CASE WHEN @transaction_date::VARCHAR <> '' THEN transaction_date::DATE >= @transaction_date::DATE ELSE TRUE END)
    AND (CASE WHEN @transaction_date::VARCHAR <> '' THEN transaction_date::DATE <= @transaction_date::DATE ELSE TRUE END)
LIMIT $1::BIGINT
OFFSET $2::BIGINT;

-- name: SelectTransactionsTotal :one
SELECT 
    count(id) AS total
FROM
    "order"
WHERE
	(CASE WHEN @transaction_date::VARCHAR <> '' THEN transaction_date::DATE >= @transaction_date::DATE ELSE TRUE END)
    AND (CASE WHEN @transaction_date::VARCHAR <> '' THEN transaction_date::DATE <= @transaction_date::DATE ELSE TRUE END);

-- name: SelectTransactionByID :one
SELECT 
    id,
	description,
    transaction_date::TIMESTAMP AS transaction_date,
    transaction_value
FROM 
	"order"
WHERE
	id = @id::BIGINT;
-----------------
---- SELECTS ----
-----------------