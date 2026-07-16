CREATE TABLE measurable_habits
(
    id          UUID PRIMARY KEY,
    label       VARCHAR(32)      NOT NULL,
    icon        SMALLINT         NOT NULL,
    step        REAL             NOT NULL,
    unit        VARCHAR(8)       NOT NULL,
    created_at  TIMESTAMP        NOT NULL,
    updated_at  TIMESTAMP        NOT NULL,
    archived_at TIMESTAMP,
    deleted_at  TIMESTAMP,
    account_id  UUID             NOT NULL
);
