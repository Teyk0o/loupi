-- Revert symptoms columns back to JSONB (only works if data is valid JSON).
ALTER TABLE symptom_entries ALTER COLUMN symptoms TYPE JSONB USING symptoms::JSONB;
ALTER TABLE symptom_entries ALTER COLUMN symptoms SET DEFAULT '[]';

ALTER TABLE symptom_checkins ALTER COLUMN symptoms TYPE JSONB USING symptoms::JSONB;
ALTER TABLE symptom_checkins ALTER COLUMN symptoms SET DEFAULT '[]';
