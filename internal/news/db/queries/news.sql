-- name: CreateNews :one
INSERT INTO news(
    media_id,
    creator_username,
    title,
    content
) VALUES (
    $1, $2, $3, $4
) RETURNING id;

-- name: CreateNewsWithPublishDate :one
INSERT INTO news(
    creator_username,
    title,
    content,
    created_at
) VALUES (
    $1, $2, $3, $4
) RETURNING id;

-- name: GetOneNews :one
SELECT * FROM news WHERE id = $1;

-- name: GetNews :many
SELECT * FROM news WHERE created_at > $4 AND title ILIKE $1 AND content = $2 ORDER BY $3 ASC;

-- name: UpdateNews :one
UPDATE news
SET
    creator_username = COALESCE(sqlc.narg(creator_username), creator_username),
    title = COALESCE(sqlc.narg(title), title),
    content = COALESCE(sqlc.narg(content), content)
WHERE
    id = sqlc.arg(id)
RETURNING id;

-- name: DeleteNews :one
UPDATE news SET deleted_at=now() RETURNING id;