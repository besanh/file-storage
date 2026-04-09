CREATE TABLE physical_files (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    
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
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    
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
    updated_at TIMESTAMPTZ,

    CONSTRAINT check_folder_integrity CHECK (
        (is_folder = TRUE AND physical_file_id IS NULL) OR 
        (is_folder = FALSE AND physical_file_id IS NOT NULL)
    )
);

-- Crucial Indexes for Performance
CREATE INDEX idx_file_nodes_parent_id ON file_nodes(parent_id);
CREATE INDEX idx_file_nodes_owner_id ON file_nodes(owner_id);
-- Ensure a user cannot have two files with the exact same name in the exact same folder
CREATE UNIQUE INDEX idx_unique_name_per_folder ON file_nodes (owner_id, parent_id, name) WHERE status = 'active';

-- Create the share_links table
CREATE TABLE share_links (
    link_token VARCHAR(32) PRIMARY KEY,       -- e.g., NanoID: "V1StGXR8_Z5jdHi6B-myT"
    
    -- The resource this link points to
    resource_id UUID NOT NULL REFERENCES file_nodes(id) ON DELETE CASCADE,
    resource_type VARCHAR(20) NOT NULL,       -- Explicitly "file" or "folder" for SpiceDB routing
    
    -- Audit & Control
    created_by TEXT NOT NULL, 
    permission_level VARCHAR(20) NOT NULL,    -- e.g., "viewer" or "editor"
    
    -- Lifecycle
    expires_at TIMESTAMP WITH TIME ZONE,      -- NULL means it never expires
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for the owner dashboard (e.g., "Show me all links I created")
CREATE INDEX idx_share_links_created_by ON share_links(created_by);

-- Index for resource management (e.g., "Show me all active links for this specific folder")
CREATE INDEX idx_share_links_resource ON share_links(resource_id);

-- Plans and Subscriptions
CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT UNIQUE NOT NULL,
    storage_quota BIGINT NOT NULL,
    price BIGINT NOT NULL DEFAULT 0,
    discount_price BIGINT NOT NULL DEFAULT 0,
    duration_days INT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ
);

CREATE TABLE user_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id uuid NOT NULL,
    plan_id UUID NOT NULL REFERENCES plans(id),
    started_at TIMESTAMPTZ NOT NULL,
    expired_at TIMESTAMPTZ,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ
);

-- Ensure a user can have only one active subscription at a time
CREATE UNIQUE INDEX uniq_active_sub
ON user_subscriptions(user_id)
WHERE status = 'active';

-- Trigger function to update user storage usage on hard delete of file_nodes
CREATE OR REPLACE FUNCTION trg_file_nodes_hard_delete()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.file_size IS NULL THEN
        RETURN OLD;
    END IF;

    UPDATE users
    SET
        storage_photos_used = GREATEST(
            storage_photos_used - (OLD.file_size * (OLD.file_type = 'photo')::int),
            0
        ),
        storage_video_used = GREATEST(
            storage_video_used - (OLD.file_size * (OLD.file_type = 'video')::int),
            0
        ),
        storage_document_used = GREATEST(
            storage_document_used - (OLD.file_size * (OLD.file_type = 'document')::int),
            0
        ),
        storage_audio_used = GREATEST(
            storage_audio_used - (OLD.file_size * (OLD.file_type = 'audio')::int),
            0
        ),
        storage_compress_used = GREATEST(
            storage_compress_used - (OLD.file_size * (OLD.file_type = 'compress')::int),
            0
        ),
        storage_other_used = GREATEST(
            storage_other_used - (OLD.file_size * (OLD.file_type = 'other')::int),
            0
        ),
        updated_at = now()
    WHERE id = OLD.owner_id;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- Attach trigger to file_nodes table
CREATE TRIGGER trg_file_nodes_after_delete
AFTER DELETE ON file_nodes
FOR EACH ROW
EXECUTE FUNCTION trg_file_nodes_hard_delete();
