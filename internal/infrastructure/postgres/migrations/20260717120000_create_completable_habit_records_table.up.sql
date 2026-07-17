CREATE TABLE completable_habit_records
(
    completable_habit_id UUID    NOT NULL,
    account_id           UUID    NOT NULL,
    date                 DATE    NOT NULL,
    value                BOOLEAN NOT NULL,
    PRIMARY KEY (completable_habit_id, date)
);
