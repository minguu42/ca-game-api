package ca_game_api

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostGachaDraw(t *testing.T) {
	reqBody := strings.NewReader(`{"times": 3}`)
	r := httptest.NewRequest("POST", "/gacha/draw", reqBody)
	r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
	w := httptest.NewRecorder()

	PostGachaDraw(w, r)

	resp := w.Result()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("cannot read body: %v", resp.Body)
	}

	if resp.StatusCode != 200 {
		t.Errorf("response code is %v", resp.StatusCode)
	}

	var response PostGachaDrawResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Errorf("cannot unmarshal body: %v", body)
	}
	if len(response.Results) != 3 {
		t.Errorf("results num is 3, but %v", len(response.Results))
	}
	if response.Results[0].CharacterId == "0" {
		t.Error("character id is not 0")
	}
	if response.Results[0].Name == "" {
		t.Error(`character name is not ""`)
	}
}