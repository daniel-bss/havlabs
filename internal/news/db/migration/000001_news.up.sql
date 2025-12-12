CREATE TABLE news (
    id uuid PRIMARY KEY NOT NULL DEFAULT(uuidv7()),
    creator_username varchar NOT NULL,
    title varchar NOT NULL,
    content varchar NOT NULL,
    created_at timestamp NOT NULL DEFAULT(NOW()),
    updated_at timestamp NOT NULL DEFAULT(NOW()),
    deleted_at timestamp
);

