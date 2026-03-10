CREATE TABLE symptom_checkins (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meal_id     UUID NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    delay_hours SMALLINT NOT NULL CHECK (delay_hours IN (6, 8, 12)),
    symptoms    JSONB NOT NULL DEFAULT '[]',
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_symptom_checkins_meal ON symptom_checkins (meal_id);
