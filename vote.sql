DROP TABLE IF EXISTS vote;

CREATE TABLE vote (
    id             BIGINT PRIMARY KEY,
    post_id        BIGINT NOT NULL,
    vote_type_id   SMALLINT NOT NULL,
    user_id        BIGINT NULL,
    creation_date  DATE NOT NULL,
    bounty_amount  INTEGER NULL
);
