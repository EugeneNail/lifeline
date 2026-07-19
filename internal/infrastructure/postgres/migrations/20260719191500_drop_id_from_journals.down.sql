ALTER TABLE journals ADD COLUMN id UUID;
CREATE INDEX IF NOT EXISTS journals_account_id_date_idx ON journals (account_id, date);
