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

/*
test1 はGetUser, GetUserRanking, PostGachaDraw, GetCharacter で使用するユーザ. x-token は ceKeMPeYr0eF3K5e4Lfjfe である.
test2 は PutUser で名前変更用ユーザ. 名前はテスト実行時にランダムな文字列に変わる. x-token は yypKkCsMXx2MBBVorFQBsQ である.
test3, test4, test5 は GetUserRanking で使用するユーザ.
*/
INSERT INTO users (name, digest_token) VALUES ('test1', '71a6f9c1007c60601a6d67e7f79d4550602b34ced90cdac86bd340f293bf0247'),
                                              ('test2', '541d9abc4b06e838e471ff564c24585a6ddc5280c9478f2e6e85b2eb7ed979a9'),
                                              ('test3', 'dj9fq2j9u9feq3nfq8fjqf98qfjdb0jqq0db09da38fa3i98qh4vnmz8zqq2ue90'),
                                              ('test4', 'q34avo2q3avj9q28t4nq39vm9uz98qnq984j91oaoj9zu9q3ujvq9832j932q8ud'),
                                              ('test5', 'jdoije928jf9eqj1fnqz9duq921ejf6qwure2qi9jf7qc4xz98qw2urf98j9eqf1');

-- GetUserRanking, GetCharacterList のためのキャラクター. Rank1. test1, Rank2. test3, Rank3. test4 となる.
INSERT INTO user_ownership_characters (user_id, character_id, level, experience) VALUES (1, 50000002, 1, 0),
                                                                                        (1, 40000002, 1, 0),
                                                                                        (1, 50000002, 1, 0),
                                                                                        (3, 50000002, 1, 0),
                                                                                        (4, 40000002, 1, 0),
                                                                                        (5, 30000002, 1, 0);
