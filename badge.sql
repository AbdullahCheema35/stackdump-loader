DROP TABLE IF EXISTS badge;

CREATE TABLE badge (
    id        BIGINT PRIMARY KEY,
    user_id   BIGINT NOT NULL,
    name      TEXT NOT NULL,
    date      TIMESTAMP NOT NULL,
    class     SMALLINT NOT NULL CHECK (class IN (1, 2, 3)),
    tag_based BOOLEAN NOT NULL DEFAULT FALSE
);
