CREATE TABLE wellness_entries (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date            DATE NOT NULL,
    stress          SMALLINT CHECK (stress BETWEEN 1 AND 5),
    mood            SMALLINT CHECK (mood BETWEEN 1 AND 5),
    energy          SMALLINT CHECK (energy BETWEEN 1 AND 5),
    sleep_hours     REAL,
    sleep_quality   SMALLINT CHECK (sleep_quality BETWEEN 1 AND 5),
    sport           JSONB,
    hydration       SMALLINT,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_wellness_user_date UNIQUE (user_id, date)
);

CREATE INDEX idx_wellness_user_date ON wellness_entries (user_id, date);
