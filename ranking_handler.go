package ca_game_api

import "net/http"

type GetRankingUserResponse struct {
	Rank1UserId string `json:"rank1UserID"`
	Rank2UserId string `json:"rank2UserID"`
	Rank3UserId string `json:"rank3UserID"`
}

func GetRankingUser(w http.ResponseWriter, r *http.Request) {
  if isStatusMethodInvalid(w, r, http.MethodGet) {
	  return
  }

}