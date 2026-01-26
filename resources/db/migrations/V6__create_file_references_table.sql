-- File References Table
-- V6: Create file_references table for tracking uploaded files

CREATE TABLE IF NOT EXISTS file_references (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    storage_key VARCHAR(500) NOT NULL UNIQUE,
    content_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    storage_provider VARCHAR(50) NOT NULL DEFAULT 'S3',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_file_references_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_file_references_user_id ON file_references(user_id);
CREATE INDEX IF NOT EXISTS idx_file_references_storage_key ON file_references(storage_key);
CREATE INDEX IF NOT EXISTS idx_file_references_file_type ON file_references(file_type);

-- Comments for documentation
COMMENT ON TABLE file_references IS 'Tracks all uploaded files with their storage locations';
COMMENT ON COLUMN file_references.storage_key IS 'Unique key/path in the storage provider (S3 key, file path, etc)';
COMMENT ON COLUMN file_references.file_type IS 'Type of file: PROFILE_IMAGE, DOCUMENT, ATTACHMENT, etc';
COMMENT ON COLUMN file_references.storage_provider IS 'Storage provider used: S3, LOCAL, GCS, etc';
