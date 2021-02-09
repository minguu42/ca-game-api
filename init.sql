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
