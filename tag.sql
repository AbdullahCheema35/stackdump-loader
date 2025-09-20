CREATE TABLE tag (
    id                BIGINT PRIMARY KEY,
    tag_name          TEXT NOT NULL,
    count             INTEGER NOT NULL,
    excerpt_post_id   BIGINT NULL,
    wiki_post_id      BIGINT NULL
);
