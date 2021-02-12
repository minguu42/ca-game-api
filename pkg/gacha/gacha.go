package gacha

import "database/sql"

func ApplyGachaResult(db *sql.DB, userId, CharacterId int) error {
	const insertSql = "INSERT INTO gacha_results (user_id, character_id) VALUES (?, ?)"
	if _, err := db.Exec(insertSql, userId, CharacterId); err != nil {
		return err
	}
	const createSql = "INSERT INTO user_ownership_characters (user_id, character_id) VALUES (?, ?)"
	if _, err := db.Exec(createSql, userId, CharacterId); err != nil {
		return err
	}
	return nil
}