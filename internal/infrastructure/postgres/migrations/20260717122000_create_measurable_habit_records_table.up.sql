CREATE TABLE measurable_habit_records
(
    measurable_habit_id UUID NOT NULL,
    account_id          UUID NOT NULL,
    date                DATE NOT NULL,
    value               REAL NOT NULL,
    PRIMARY KEY (measurable_habit_id, date)
);
