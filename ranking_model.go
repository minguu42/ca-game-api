package ca_game_api

import (
	"database/sql"
	"log"
	"net/http"
)

func selectUserRanking(db *sql.DB, w http.ResponseWriter) ([]UserInfo, error) {
	log.Println("INFO START selectUserRanking")
	var users []UserInfo
	const selectSql = `
SELECT UOC.user_id, U.name, SUM(UOC.level * C.power) AS sumPower
FROM user_ownership_characters AS UOC
INNER JOIN users AS U ON UOC.user_id = U.id
INNER JOIN characters AS C ON UOC.character_id = C.id
GROUP BY UOC.user_id
ORDER BY sumPower DESC
LIMIT 3
`
	rows, err := db.Query(selectSql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("ERROR Return 500:", err)
		return nil, err
	}
	for rows.Next() {
		var user UserInfo
		if err := rows.Scan(&user.Id, &user.Name, &user.SumPower); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("ERROR Return 500:", err)
			return nil, err
		}
		users = append(users, user)
	}
	log.Println("INFO END selectUserRanking")
	return users, nil
}
