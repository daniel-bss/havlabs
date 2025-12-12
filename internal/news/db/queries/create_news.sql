-- name: CreateNews :one
INSERT INTO news(
    creator_username,
    title,
    content
) VALUES (
    $1, $2, $3
) RETURNING id;

-- name: GetOneNews :one
SELECT * FROM news WHERE id = $1;

-- name: GetAllNews :one
SELECT * FROM news;

-- name: UpdateNews :one
UPDATE news
SET
    creator_username = COALESCE(sqlc.narg(creator_username), creator_username),
    title = COALESCE(sqlc.narg(title), title),
    content = COALESCE(sqlc.narg(content), content)
WHERE
    id = sqlc.arg(id)
RETURNING id;