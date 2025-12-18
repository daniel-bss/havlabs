-- name: CreateUpload :one
INSERT INTO media (
  id,
  owner_username,

  purpose,
  bucket,
  object_key,

  declared_content_type
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING id;


-- name: ConfirmGoodUpload :one
UPDATE media 
SET 
    bucket = $2,

    actual_content_type = $3,

    status = 'ready',

    size_bytes = $4,
    checksum = $5,
    metadata = $6,

    uploaded_at = NOW(),
    updated_at = NOW()
WHERE
    id = $1
RETURNING id;

-- name: UpdateUploadStatus :one
UPDATE media 
SET 
    status = $2,
    updated_at = NOW()
WHERE
    id = $1
RETURNING id;

-- name: GetMediaById :one
SELECT id, purpose FROM media WHERE id=$1;
