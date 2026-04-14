-- name: InsertFile :one
INSERT INTO file_nodes (
    owner_id,
    parent_id,
    physical_file_id,
    name,
    is_folder,
    file_hash,
    file_size,
    file_type,
    file_ext,
    file_mime_type,
    file_video_resolution,
    recent_accessed_at,
    favorite,
    status,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, now()
) RETURNING *;

-- name: GetUserStorageUsed :one
SELECT COALESCE(SUM(file_size), 0)::bigint AS total_used
FROM file_nodes
WHERE owner_id = $1
  AND is_folder = FALSE
  AND status = 'active';

-- name: GetUserStorageBreakdown :one
SELECT 
    COALESCE(SUM(CASE WHEN file_type = 'photo' THEN file_size ELSE 0 END), 0)::bigint AS photos,
    COALESCE(SUM(CASE WHEN file_type = 'video' THEN file_size ELSE 0 END), 0)::bigint AS videos,
    COALESCE(SUM(CASE WHEN file_type = 'document' THEN file_size ELSE 0 END), 0)::bigint AS documents,
    COALESCE(SUM(CASE WHEN file_type = 'audio' THEN file_size ELSE 0 END), 0)::bigint AS audio,
    COALESCE(SUM(CASE WHEN file_type = 'compress' THEN file_size ELSE 0 END), 0)::bigint AS compress,
    COALESCE(SUM(CASE WHEN file_type = 'other' THEN file_size ELSE 0 END), 0)::bigint AS other
FROM file_nodes
WHERE owner_id = $1
  AND is_folder = FALSE
  AND status = 'active';

-- name: GetRecentFiles :many
SELECT *
FROM file_nodes
WHERE owner_id = $1
AND status = 'active'
ORDER BY recent_accessed_at DESC NULLS LAST, created_at DESC
LIMIT $2;

-- name: GetFile :one
SELECT *
FROM file_nodes
WHERE id = $1
AND status = 'active';