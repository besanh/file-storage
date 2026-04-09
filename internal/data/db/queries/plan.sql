-- name: GetPlan :one
SELECT * FROM plans WHERE id = $1 LIMIT 1;

-- name: GetPlanByName :one
SELECT * FROM plans WHERE name = $1 LIMIT 1;

-- name: ListPlans :many
SELECT * FROM plans ORDER BY price ASC;

-- name: CreatePlan :one
INSERT INTO plans (
    name,
    storage_quota,
    price,
    discount_price,
    duration_days,
    description
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdatePlan :one
UPDATE plans SET
    name = $2,
    storage_quota = $3,
    price = $4,
    discount_price = $5,
    duration_days = $6,
    description = $7,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePlan :exec
DELETE FROM plans WHERE id = $1;
