CREATE TABLE IF NOT EXISTS `users` (
    id int PRIMARY KEY AUTO_INCREMENT,
    name varchar(255) UNIQUE NOT NULL,
    digest_token varchar(255) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `characters` (
    id int PRIMARY KEY,
    name varchar(255) NOT NULL,
    rarity int NOT NULL
);

CREATE TABLE IF NOT EXISTS `user_ownership_characters` (
    id int PRIMARY KEY AUTO_INCREMENT,
    user_id int NOT NULL,
    character_id int NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `gacha_results` (
    id int PRIMARY KEY AUTO_INCREMENT,
    user_id int NOT NULL,
    character_id int NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE `user_ownership_characters`
    ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `user_ownership_characters`
    ADD FOREIGN KEY (`character_id`) REFERENCES `characters` (`id`);

ALTER TABLE `gacha_results`
    ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `gacha_results`
    ADD FOREIGN KEY (`character_id`) REFERENCES `characters` (`id`);

INSERT INTO characters VALUES (30000001, 'normal_character1', 3);
INSERT INTO characters VALUES (30000002, 'normal_character2', 3);
INSERT INTO characters VALUES (30000003, 'normal_character3', 3);
INSERT INTO characters VALUES (30000004, 'normal_character4', 3);
INSERT INTO characters VALUES (30000005, 'normal_character5', 3);
INSERT INTO characters VALUES (30000006, 'normal_character6', 3);
INSERT INTO characters VALUES (30000007, 'normal_character7', 3);
INSERT INTO characters VALUES (30000008, 'normal_character8', 3);
INSERT INTO characters VALUES (30000009, 'normal_character9', 3);
INSERT INTO characters VALUES (30000010, 'normal_character10', 3);
INSERT INTO characters VALUES (40000001, 'rare_character1', 4);
INSERT INTO characters VALUES (40000002, 'rare_character2', 4);
INSERT INTO characters VALUES (40000003, 'rare_character3', 4);
INSERT INTO characters VALUES (40000004, 'rare_character4', 4);
INSERT INTO characters VALUES (40000005, 'rare_character5', 4);
INSERT INTO characters VALUES (40000006, 'rare_character6', 4);
INSERT INTO characters VALUES (40000007, 'rare_character7', 4);
INSERT INTO characters VALUES (40000008, 'rare_character8', 4);
INSERT INTO characters VALUES (40000009, 'rare_character9', 4);
INSERT INTO characters VALUES (40000010, 'rare_character10', 4);
INSERT INTO characters VALUES (50000001, 'super_rare_character1', 5);
INSERT INTO characters VALUES (50000002, 'super_rare_character2', 5);
INSERT INTO characters VALUES (50000003, 'super_rare_character3', 5);
INSERT INTO characters VALUES (50000004, 'super_rare_character4', 5);
INSERT INTO characters VALUES (50000005, 'super_rare_character5', 5);
INSERT INTO characters VALUES (50000006, 'super_rare_character6', 5);
INSERT INTO characters VALUES (50000007, 'super_rare_character7', 5);
INSERT INTO characters VALUES (50000008, 'super_rare_character8', 5);
INSERT INTO characters VALUES (50000009, 'super_rare_character9', 5);
INSERT INTO characters VALUES (50000010, 'super_rare_character10', 5);
