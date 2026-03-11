-- Change symptoms columns from JSONB to TEXT to support encrypted data storage.
ALTER TABLE symptom_entries ALTER COLUMN symptoms TYPE TEXT USING symptoms::TEXT;
ALTER TABLE symptom_entries ALTER COLUMN symptoms SET DEFAULT '';

ALTER TABLE symptom_checkins ALTER COLUMN symptoms TYPE TEXT USING symptoms::TEXT;
ALTER TABLE symptom_checkins ALTER COLUMN symptoms SET DEFAULT '';
