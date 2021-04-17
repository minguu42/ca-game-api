package ca_game_api

import "net/http"

type UserInfo struct {
	Id       string `json:"userID"`
	Name     string `json:"name"`
	SumPower string `json:"sumPower"`
}

type GetRankingUserResponse struct {
	UserRankings []UserInfo `json:"userRankings"`
}

func GetRankingUser(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(r, http.MethodGet) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	users, err := selectUserRanking(w)
	if err != nil {
		return
	}

	jsonResponse := GetRankingUserResponse{
		UserRankings: users,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		return
	}
}
