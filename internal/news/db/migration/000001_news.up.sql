CREATE TABLE news (
    id uuid PRIMARY KEY NOT NULL DEFAULT(uuidv7()),
    media_id uuid NOT NULL,
    creator_username TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at timestamp NOT NULL DEFAULT(NOW()),
    updated_at timestamp NOT NULL DEFAULT(NOW()),
    deleted_at timestamp
);

CREATE INDEX idx_news_id ON news(id);
CREATE INDEX idx_news_media ON news(media_id);
CREATE INDEX idx_news_createdat ON news(created_at);

