package ca_game_api

import (
	"fmt"
)

func selectUserRanking() ([]UserInfo, error) {
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
		return nil, fmt.Errorf("db.Query faild: %w", err)
	}
	for rows.Next() {
		var user UserInfo
		if err := rows.Scan(&user.Id, &user.Name, &user.SumPower); err != nil {
			return nil, fmt.Errorf("rows.Scan faild: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}
