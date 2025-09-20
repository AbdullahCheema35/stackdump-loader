DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id                BIGINT PRIMARY KEY,
    reputation        INTEGER NOT NULL,
    creation_date     TIMESTAMP NOT NULL,
    display_name      TEXT NOT NULL,
    last_access_date  TIMESTAMP NULL,
    website_url       TEXT NULL,
    location          TEXT NULL,
    about_me          TEXT NULL,
    views             INTEGER NOT NULL DEFAULT 0,
    up_votes          INTEGER NOT NULL DEFAULT 0,
    down_votes        INTEGER NOT NULL DEFAULT 0,
    profile_image_url TEXT NULL,
    account_id        BIGINT NULL
);
