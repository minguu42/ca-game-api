package ca_game_api

import (
	"log"
	"net/http"
)

type UserInfo struct {
	Id       string `json:"userID"`
	Name     string `json:"name"`
	SumPower string `json:"sumPower"`
}

type GetUserRankingResponse struct {
	UserRankings []UserInfo `json:"userRankings"`
}

func GetUserRanking(w http.ResponseWriter, r *http.Request) {
	if isStatusMethodInvalid(r, http.MethodGet) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userRankings, err := selectUserRanking()
	if err != nil {
		log.Println("ERROR selectUserRanking error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse := GetUserRankingResponse{
		UserRankings: userRankings,
	}
	if err := encodeResponse(w, jsonResponse); err != nil {
		log.Println("ERROR encodeResponse fail:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
