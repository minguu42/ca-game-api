CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    name varchar(255) UNIQUE NOT NULL,
    digest_token varchar(255) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS characters (
    id int PRIMARY KEY,
    name varchar(255) NOT NULL,
    rarity int NOT NULL,
    base_power int NOT NULL,
    calorie int NOT NULL
);

CREATE TABLE IF NOT EXISTS user_characters (
    id serial PRIMARY KEY,
    user_id int REFERENCES users,
    character_id int REFERENCES characters,
    experience int NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS gacha_results (
    id serial PRIMARY KEY,
    user_id int REFERENCES users,
    character_id int REFERENCES characters,
    experience int NOT NULL,
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

CREATE TRIGGER set_timestamp_on_user_characters
    BEFORE UPDATE ON user_characters
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

INSERT INTO characters VALUES (30000001, 'normal_character1', 3, 1, 10000),
                              (30000002, 'normal_character2', 3, 200, 500),
                              (30000003, 'normal_character3', 3, 180, 500),
                              (30000004, 'normal_character4', 3, 195, 500),
                              (30000005, 'normal_character5', 3, 200, 500),
                              (30000006, 'normal_character6', 3, 205, 500),
                              (30000007, 'normal_character7', 3, 200, 500),
                              (30000008, 'normal_character8', 3, 230, 500),
                              (30000009, 'normal_character9', 3, 210, 500),
                              (30000010, 'normal_character10', 3, 200, 500),
                              (40000001, 'rare_character1', 4, 1, 20000),
                              (40000002, 'rare_character2', 4, 300, 1000),
                              (40000003, 'rare_character3', 4, 320, 1000),
                              (40000004, 'rare_character4', 4, 310, 1000),
                              (40000005, 'rare_character5', 4, 320, 1000),
                              (40000006, 'rare_character6', 4, 330, 1000),
                              (40000007, 'rare_character7', 4, 300, 1000),
                              (40000008, 'rare_character8', 4, 310, 1000),
                              (40000009, 'rare_character9', 4, 315, 1000),
                              (40000010, 'rare_character10', 4, 350, 1000),
                              (50000001, 'super_rare_character1', 5, 1, 50000),
                              (50000002, 'super_rare_character2', 5, 500, 2000),
                              (50000003, 'super_rare_character3', 5, 540, 2000),
                              (50000004, 'super_rare_character4', 5, 560, 2000),
                              (50000005, 'super_rare_character5', 5, 550, 2000),
                              (50000006, 'super_rare_character6', 5, 540, 2000),
                              (50000007, 'super_rare_character7', 5, 510, 2000),
                              (50000008, 'super_rare_character8', 5, 550, 2000),
                              (50000009, 'super_rare_character9', 5, 560, 2000),
                              (50000010, 'super_rare_character10', 5, 600, 2000);
