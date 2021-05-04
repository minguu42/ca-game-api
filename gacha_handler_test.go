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

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("io.ReadAll failed: %v", err)
	}
	var body PostGachaDrawResponse
	if err := json.Unmarshal(bytes, &body); err != nil {
		t.Errorf("json.Unmarshal failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Response code should be 200, but %v", resp.StatusCode)
	}
	if len(body.Results) != 3 {
		t.Errorf("Results num should be 3, but %v", len(body.Results))
	}
	if body.Results[0].CharacterId == 0 {
		t.Error(`1st character's characterID does not exist`)
	}
	if body.Results[0].Name == "" {
		t.Error(`1st character's name does not exist`)
	}
	if body.Results[1].CharacterId == 0 {
		t.Error(`2nd character's characterID does not exist`)
	}
	if body.Results[1].Name == "" {
		t.Error(`2nd character's name does not exist`)
	}
	if body.Results[2].CharacterId == 0 {
		t.Error(`3rd character's characterID does not exist`)
	}
	if body.Results[2].Name == "" {
		t.Error(`3rd character's name does not exist`)
	}
}
