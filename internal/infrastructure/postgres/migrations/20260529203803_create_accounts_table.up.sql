CREATE TABLE accounts
(
    id         UUID PRIMARY KEY,
    email      VARCHAR(200) NOT NULL UNIQUE,
    password   VARCHAR(150) NOT NULL,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL
);