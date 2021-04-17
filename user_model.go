package ca_game_api

import (
	"log"
	"net/http"
)

func insertUser(name string, token string, w http.ResponseWriter) error {
	log.Println("INFO START insertUser")
	const createSql = `INSERT INTO users (name, digest_token) VALUES ($1, $2);`
	digestToken := HashToken(token)
	if _, err := db.Exec(createSql, name, digestToken); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Return 403:", err)
		return err
	}
	log.Println("INFO END insertUser")
	return nil
}

func selectUserName(token string, w http.ResponseWriter) (string, error) {
	log.Println("INFO START selectUserName")
	const selectSql = `SELECT name FROM users WHERE digest_token = $1`
	digestToken := HashToken(token)

	var name string
	row := db.QueryRow(selectSql, digestToken)
	if err := row.Scan(&name); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("ERROR Return 401: x-token is invalid")
		return "", err
	}
	log.Println("INFO END selectUserName")
	return name, nil
}

func selectUserId(token string, w http.ResponseWriter) (int, error) {
	const selectSql = `SELECT id FROM users WHERE digest_token = $1`
	digestToken := HashToken(token)
	row := db.QueryRow(selectSql, digestToken)
	var id int
	if err := row.Scan(&id); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("ERROR Return 401: x-token is invalid")
		return 0, err
	}
	return id, nil
}

func selectUserIdByUserCharacterId(userCharacterId int, w http.ResponseWriter) (int, error) {
	const selectSql = `SELECT user_id FROM user_ownership_characters WHERE id = $1`
	row := db.QueryRow(selectSql, userCharacterId)
	var id int
	if err := row.Scan(&id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Return 400:", err)
		return 0, err
	}
	return id, nil
}

func updateUser(token, newName string, w http.ResponseWriter) error {
	log.Println("INFO START updateUser")
	const updateSql = `UPDATE users SET name = $1 WHERE digest_token = $2`
	digestToken := HashToken(token)
	if _, err := db.Exec(updateSql, newName, digestToken); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR Return 403:", err)
		return err
	}
	log.Println("INFO END updateUser")
	return nil
}
