CREATE TYPE status_types AS ENUM ('pending', 'uploaded', 'ready', 'failed');

CREATE TABLE media (
    id UUID PRIMARY KEY NOT NULL,
    owner_username TEXT NOT NULL,

    purpose TEXT NOT NULL,
    bucket TEXT NOT NULL,
    object_key TEXT NOT NULL,

    declared_content_type TEXT NOT NULL, 
    actual_content_type TEXT,

    status status_types NOT NULL DEFAULT 'pending',

    size_bytes BIGINT,
    checksum TEXT,
    metadata JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    uploaded_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_media_owner ON media(owner_username);
CREATE INDEX idx_media_status ON media(status);
CREATE UNIQUE INDEX idx_media_bucket_objectkey ON media(bucket, object_key);
