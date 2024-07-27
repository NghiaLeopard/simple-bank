-- name: CreateTransfer :one
INSERT INTO transfer (
  amount
) VALUES (
  $1
)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfer
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfer
LIMIT $1
OFFSET $2;


-- name: UpdateTransfer :one
UPDATE transfer
SET amount = $1
WHERE id = $2
RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM transfer
WHERE id = $1;


