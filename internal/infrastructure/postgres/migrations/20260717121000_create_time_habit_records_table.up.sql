CREATE TABLE time_habit_records
(
    time_habit_id UUID     NOT NULL,
    account_id    UUID     NOT NULL,
    date          DATE     NOT NULL,
    value         SMALLINT NOT NULL,
    PRIMARY KEY (time_habit_id, date)
);
