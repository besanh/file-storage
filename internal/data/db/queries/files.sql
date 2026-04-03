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
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
) RETURNING *;
