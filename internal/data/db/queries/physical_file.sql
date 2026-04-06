-- name: InsertPhysicalFile :one
INSERT INTO physical_files (
    file_hash,
    size_bytes,
    mime_type,
    storage_path,
    reference_count,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, now()
) RETURNING *;

-- name: GetPhysicalFileByID :one
SELECT * FROM physical_files WHERE id = $1;

-- name: UpdatePhysicalFileReferenceCount :one
UPDATE physical_files SET reference_count = reference_count + $2 WHERE id = $1 RETURNING *;

-- name: DeletePhysicalFileByID :one
UPDATE physical_files SET deleted_at = now() WHERE id = $1 RETURNING *;

-- name: GetPhysicalFileByHash :one
SELECT * FROM physical_files WHERE file_hash = $1;
