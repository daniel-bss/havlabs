CREATE TYPE status_types AS ENUM ('pending', 'uploaded', 'ready', 'failed');

CREATE TABLE media (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,

    bucket TEXT NOT NULL,
    object_key TEXT NOT NULL,
    content_type TEXT NOT NULL,

    status status_types NOT NULL DEFAULT 'pending',
    size_bytes BIGINT,

    checksum TEXT,
    metadata JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    uploaded_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_media_owneruser ON media(user_id);
CREATE INDEX idx_media_status ON media(status);
CREATE UNIQUE INDEX idx_media_bucket_objectkey ON media(bucket, object_key);
