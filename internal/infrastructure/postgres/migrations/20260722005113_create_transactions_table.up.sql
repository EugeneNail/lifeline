CREATE TABLE transactions
(
    id         UUID        NOT NULL,
    money      REAL        NOT NULL,
    date       DATE        NOT NULL,
    account_id UUID        NOT NULL,
    created_at TIMESTAMP   NOT NULL,
    updated_at TIMESTAMP   NOT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS transactions_account_id_date_idx ON transactions (account_id, date);
