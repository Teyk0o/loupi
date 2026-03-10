CREATE TABLE symptom_entries (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    symptoms    JSONB NOT NULL DEFAULT '[]',
    notes       TEXT,
    entry_time  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_symptom_entries_user_date ON symptom_entries (user_id, entry_time);
