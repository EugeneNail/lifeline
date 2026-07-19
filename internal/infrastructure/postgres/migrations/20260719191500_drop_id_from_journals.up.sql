ALTER TABLE journals DROP CONSTRAINT IF EXISTS entries_pkey;
ALTER TABLE journals DROP CONSTRAINT IF EXISTS journals_pkey;
ALTER TABLE journals DROP CONSTRAINT IF EXISTS entries_account_id_date_key;
ALTER TABLE journals DROP COLUMN id;
CREATE INDEX IF NOT EXISTS journals_account_id_date_idx ON journals (account_id, date);
