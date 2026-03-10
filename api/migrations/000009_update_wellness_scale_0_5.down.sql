-- Revert wellness metric constraints back to 1-5 scale.

ALTER TABLE wellness_entries DROP CONSTRAINT IF EXISTS wellness_entries_stress_check;
ALTER TABLE wellness_entries ADD CONSTRAINT wellness_entries_stress_check CHECK (stress BETWEEN 1 AND 5);

ALTER TABLE wellness_entries DROP CONSTRAINT IF EXISTS wellness_entries_mood_check;
ALTER TABLE wellness_entries ADD CONSTRAINT wellness_entries_mood_check CHECK (mood BETWEEN 1 AND 5);

ALTER TABLE wellness_entries DROP CONSTRAINT IF EXISTS wellness_entries_energy_check;
ALTER TABLE wellness_entries ADD CONSTRAINT wellness_entries_energy_check CHECK (energy BETWEEN 1 AND 5);

ALTER TABLE wellness_entries DROP CONSTRAINT IF EXISTS wellness_entries_sleep_quality_check;
ALTER TABLE wellness_entries ADD CONSTRAINT wellness_entries_sleep_quality_check CHECK (sleep_quality BETWEEN 1 AND 5);
