CREATE TABLE physical_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Cryptographic hash (e.g., SHA-256) of the file content.
    -- If two users upload the exact same video, you only store it once!
    file_hash VARCHAR(64) UNIQUE NOT NULL, 
    
    size_bytes BIGINT NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    
    -- The actual location where the file is saved (e.g., "s3://bucket/hash.pdf")
    storage_path TEXT NOT NULL, 

    -- Optional: Track how many logical files point to this physical file
    -- If this count drops to 0, you can safely delete the file from S3
    reference_count INT DEFAULT 1,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE file_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- The user who owns this specific file or folder
    owner_id UUID NOT NULL,
    
    -- Self-referencing Foreign Key for nested folders. 
    -- NULL means it is in the root directory.
    parent_id UUID REFERENCES file_nodes(id) ON DELETE CASCADE,
    
    -- Link to the actual data. 
    -- NULL means this node is a Folder, not a File.
    physical_file_id UUID REFERENCES physical_files(id),
    
    name VARCHAR(255) NOT NULL,
    is_folder BOOLEAN NOT NULL DEFAULT FALSE,

    file_hash TEXT,
    file_size BIGINT,
    file_type TEXT,
    file_ext TEXT,
    file_mime_type TEXT,
    file_video_resolution TEXT,
    recent_accessed_at TIMESTAMPTZ,
    favorite BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMPTZ,     -- Trash
    deleted_by TEXT,
    
    -- Status tracking for your Trash and Restore features
    status VARCHAR(50) DEFAULT 'active', -- 'active', 'trashed'
    
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ
);

-- Crucial Indexes for Performance
CREATE INDEX idx_file_nodes_parent_id ON file_nodes(parent_id);
CREATE INDEX idx_file_nodes_owner_id ON file_nodes(owner_id);
-- Ensure a user cannot have two files with the exact same name in the exact same folder
CREATE UNIQUE INDEX idx_unique_name_per_folder ON file_nodes (owner_id, parent_id, name) WHERE status = 'active';