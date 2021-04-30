CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    name varchar(255) UNIQUE NOT NULL,
    digest_token varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS characters (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    rarity int NOT NULL,
    power int NOT NULL,
    calorie int NOT NULL
);

CREATE TABLE IF NOT EXISTS user_ownership_characters (
    id serial PRIMARY KEY,
    user_id int REFERENCES users,
    character_id int REFERENCES characters,
    level int NOT NULL,
    experience int NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gacha_results (
    id serial PRIMARY KEY,
    user_id int REFERENCES users,
    character_id int REFERENCES characters,
    level int NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_on_users
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TRIGGER set_timestamp_on_user_ownership_characters
    BEFORE UPDATE ON user_ownership_characters
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
