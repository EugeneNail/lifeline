CREATE TABLE mood_records
(
    date       DATE NOT NULL,
    account_id UUID NOT NULL,
    value      INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (date, account_id)
);
