-- name: CreateShareLink :one
INSERT INTO share_links (
    link_token,
    resource_id,
    resource_type,
    created_by,
    permission_level,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetShareLinkByToken :one
SELECT * FROM share_links
WHERE link_token = $1 LIMIT 1;

-- name: RevokeShareLink :exec
DELETE FROM share_links
WHERE link_token = $1;

-- name: ListShareLinksByCreator :many
SELECT * FROM share_links
WHERE created_by = $1
LIMIT $2 OFFSET $3;

-- name: ListShareLinksByResource :many
SELECT * FROM share_links
WHERE resource_id = $1
LIMIT $2 OFFSET $3;

-- name: UpdateShareLink :exec
UPDATE share_links
SET permission_level = $2
WHERE link_token = $1;

-- name: GetUserShareLinkCount :one
SELECT COUNT(*)::bigint FROM share_links
WHERE created_by = $1;
