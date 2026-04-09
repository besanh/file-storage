-- name: GetUserSubscription :one
SELECT * FROM user_subscriptions WHERE user_id = $1 AND status = 'active' LIMIT 1;

-- name: GetSubscription :one
SELECT * FROM user_subscriptions WHERE id = $1 LIMIT 1;

-- name: ListUserSubscriptions :many
SELECT * FROM user_subscriptions WHERE user_id = $1 ORDER BY created_at DESC;

-- name: CreateSubscription :one
INSERT INTO user_subscriptions (
    user_id,
    plan_id,
    started_at,
    expired_at,
    status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateSubscriptionStatus :one
UPDATE user_subscriptions SET
    status = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateSubscriptionExpiration :one
UPDATE user_subscriptions SET
    expired_at = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSubscription :exec
DELETE FROM user_subscriptions WHERE id = $1;