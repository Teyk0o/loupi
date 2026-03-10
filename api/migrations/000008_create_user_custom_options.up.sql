-- Custom user options for personalizable lists (symptom types, meal categories, sport types).
-- Each row is one option belonging to a user under a specific category.

CREATE TABLE user_custom_options (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category    VARCHAR(50) NOT NULL,  -- 'symptom_type', 'meal_category', 'sport_type'
    value       VARCHAR(100) NOT NULL, -- machine key, e.g. 'diarrhea', 'homemade', 'running'
    label       VARCHAR(100) NOT NULL, -- display label, e.g. 'Diarrhée', 'Fait maison', 'Course'
    emoji       VARCHAR(10),           -- optional emoji for display
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_custom_options_user_category ON user_custom_options(user_id, category);

-- Unique constraint: one value per category per user
CREATE UNIQUE INDEX idx_user_custom_options_unique ON user_custom_options(user_id, category, value);

-- Function to seed default options for a new user.
CREATE OR REPLACE FUNCTION seed_default_options(p_user_id UUID) RETURNS VOID AS $$
BEGIN
    -- Default symptom types
    INSERT INTO user_custom_options (user_id, category, value, label, sort_order) VALUES
        (p_user_id, 'symptom_type', 'diarrhea',     'Diarrhée',            1),
        (p_user_id, 'symptom_type', 'stomach_ache',  'Maux de ventre',     2),
        (p_user_id, 'symptom_type', 'nausea',        'Nausée',             3),
        (p_user_id, 'symptom_type', 'bloating',      'Ballonnements',      4),
        (p_user_id, 'symptom_type', 'heartburn',     'Brûlures d''estomac', 5),
        (p_user_id, 'symptom_type', 'cramps',        'Crampes',            6),
        (p_user_id, 'symptom_type', 'constipation',  'Constipation',       7),
        (p_user_id, 'symptom_type', 'gas',           'Gaz',                8),
        (p_user_id, 'symptom_type', 'reflux',        'Reflux',             9),
        (p_user_id, 'symptom_type', 'fatigue',       'Fatigue',           10);

    -- Default meal categories
    INSERT INTO user_custom_options (user_id, category, value, label, emoji, sort_order) VALUES
        (p_user_id, 'meal_category', 'homemade',   'Fait maison', '🏠',  1),
        (p_user_id, 'meal_category', 'restaurant', 'Restaurant',  '🍽️', 2),
        (p_user_id, 'meal_category', 'takeout',    'À emporter',  '🥡',  3),
        (p_user_id, 'meal_category', 'snack',      'Collation',   '🍪',  4),
        (p_user_id, 'meal_category', 'fast_food',  'Fast-food',   '🍔',  5),
        (p_user_id, 'meal_category', 'cafeteria',  'Cantine',     '🏫',  6),
        (p_user_id, 'meal_category', 'family',     'En famille',  '👨‍👩‍👧', 7),
        (p_user_id, 'meal_category', 'friends',    'Entre amis',  '👫',  8),
        (p_user_id, 'meal_category', 'other',      'Autre',       '🍴',  9);

    -- Default sport types
    INSERT INTO user_custom_options (user_id, category, value, label, emoji, sort_order) VALUES
        (p_user_id, 'sport_type', 'running',    'Course',       '🏃', 1),
        (p_user_id, 'sport_type', 'cycling',    'Vélo',         '🚴', 2),
        (p_user_id, 'sport_type', 'swimming',   'Natation',     '🏊', 3),
        (p_user_id, 'sport_type', 'gym',        'Musculation',  '🏋️', 4),
        (p_user_id, 'sport_type', 'yoga',       'Yoga',         '🧘', 5),
        (p_user_id, 'sport_type', 'walking',    'Marche',       '🚶', 6),
        (p_user_id, 'sport_type', 'team_sport', 'Sport co.',    '⚽', 7),
        (p_user_id, 'sport_type', 'other',      'Autre',        '🏅', 8);
END;
$$ LANGUAGE plpgsql;

-- Trigger: automatically seed defaults when a new user is created.
CREATE OR REPLACE FUNCTION trigger_seed_user_options() RETURNS TRIGGER AS $$
BEGIN
    PERFORM seed_default_options(NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_seed_user_options
    AFTER INSERT ON users
    FOR EACH ROW
    EXECUTE FUNCTION trigger_seed_user_options();

-- Seed options for all existing users that don't have any yet.
INSERT INTO user_custom_options (user_id, category, value, label, sort_order)
SELECT u.id, vals.category, vals.value, vals.label, vals.sort_order
FROM users u
CROSS JOIN (
    VALUES
        ('symptom_type', 'diarrhea',     'Diarrhée',            1),
        ('symptom_type', 'stomach_ache', 'Maux de ventre',      2),
        ('symptom_type', 'nausea',       'Nausée',              3),
        ('symptom_type', 'bloating',     'Ballonnements',       4),
        ('symptom_type', 'heartburn',    'Brûlures d''estomac', 5),
        ('symptom_type', 'cramps',       'Crampes',             6),
        ('symptom_type', 'constipation', 'Constipation',        7),
        ('symptom_type', 'gas',          'Gaz',                 8),
        ('symptom_type', 'reflux',       'Reflux',              9),
        ('symptom_type', 'fatigue',      'Fatigue',            10)
) AS vals(category, value, label, sort_order)
WHERE NOT EXISTS (
    SELECT 1 FROM user_custom_options uco WHERE uco.user_id = u.id AND uco.category = 'symptom_type'
);

INSERT INTO user_custom_options (user_id, category, value, label, emoji, sort_order)
SELECT u.id, vals.category, vals.value, vals.label, vals.emoji, vals.sort_order
FROM users u
CROSS JOIN (
    VALUES
        ('meal_category', 'homemade',   'Fait maison', '🏠',  1),
        ('meal_category', 'restaurant', 'Restaurant',  '🍽️', 2),
        ('meal_category', 'takeout',    'À emporter',  '🥡',  3),
        ('meal_category', 'snack',      'Collation',   '🍪',  4),
        ('meal_category', 'fast_food',  'Fast-food',   '🍔',  5),
        ('meal_category', 'cafeteria',  'Cantine',     '🏫',  6),
        ('meal_category', 'family',     'En famille',  '👨‍👩‍👧', 7),
        ('meal_category', 'friends',    'Entre amis',  '👫',  8),
        ('meal_category', 'other',      'Autre',       '🍴',  9)
) AS vals(category, value, label, emoji, sort_order)
WHERE NOT EXISTS (
    SELECT 1 FROM user_custom_options uco WHERE uco.user_id = u.id AND uco.category = 'meal_category'
);

INSERT INTO user_custom_options (user_id, category, value, label, emoji, sort_order)
SELECT u.id, vals.category, vals.value, vals.label, vals.emoji, vals.sort_order
FROM users u
CROSS JOIN (
    VALUES
        ('sport_type', 'running',    'Course',       '🏃', 1),
        ('sport_type', 'cycling',    'Vélo',         '🚴', 2),
        ('sport_type', 'swimming',   'Natation',     '🏊', 3),
        ('sport_type', 'gym',        'Musculation',  '🏋️', 4),
        ('sport_type', 'yoga',       'Yoga',         '🧘', 5),
        ('sport_type', 'walking',    'Marche',       '🚶', 6),
        ('sport_type', 'team_sport', 'Sport co.',    '⚽', 7),
        ('sport_type', 'other',      'Autre',        '🏅', 8)
) AS vals(category, value, label, emoji, sort_order)
WHERE NOT EXISTS (
    SELECT 1 FROM user_custom_options uco WHERE uco.user_id = u.id AND uco.category = 'sport_type'
);
