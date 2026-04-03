-- name: InsertPhysicalFile :one
INSERT INTO physical_files (
    file_hash,
    size_bytes,
    mime_type,
    storage_path,
    reference_count,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;