DROP TRIGGER IF EXISTS trg_seed_user_options ON users;
DROP FUNCTION IF EXISTS trigger_seed_user_options();
DROP FUNCTION IF EXISTS seed_default_options(UUID);
DROP TABLE IF EXISTS user_custom_options;
