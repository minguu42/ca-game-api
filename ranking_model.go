package ca_game_api

import (
	"fmt"
	"log"
)

func selectUserRanking() ([]UserInfo, error) {
	log.Println("INFO START selectUserRanking")
	var users []UserInfo
	const selectSql = `
SELECT U.id, U.name, SUM(UOC.level * C.power) AS sumPower
FROM user_ownership_characters AS UOC
INNER JOIN users AS U ON UOC.user_id = U.id
INNER JOIN characters AS C ON UOC.character_id = C.id
GROUP BY U.id
ORDER BY sumPower DESC
LIMIT 3
`
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, fmt.Errorf("query fail: %w", err)
	}
	for rows.Next() {
		var user UserInfo
		if err := rows.Scan(&user.Id, &user.Name, &user.SumPower); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	log.Println("INFO END selectUserRanking")
	return users, nil
}
