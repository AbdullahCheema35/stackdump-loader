DROP TABLE IF EXISTS comment;

CREATE TABLE comment (
    id                BIGINT PRIMARY KEY,
    post_id           BIGINT NOT NULL,
    score             INTEGER NOT NULL DEFAULT 0,
    text              TEXT NULL,
    creation_date     TIMESTAMP NOT NULL,
    user_display_name TEXT NULL,
    user_id           BIGINT NULL,
    content_license   TEXT NULL
);
