CREATE TABLE entries
(
    id         UUID PRIMARY KEY,
    date       DATE      NOT NULL,
    mood       SMALLINT  NOT NULL,
    note       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    account_id UUID      NOT NULL,
    UNIQUE (account_id, date)
);
