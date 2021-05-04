package ca_game_api

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
)

func TestGetCharacterList(t *testing.T) {
	r := httptest.NewRequest("GET", "/character/list", nil)
	r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
	w := httptest.NewRecorder()

	GetCharacterList(w, r)

	resp := w.Result()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("cannot read body: %v", resp.Body)
	}

	if resp.StatusCode != 200 {
		t.Errorf("response code is %v", resp.StatusCode)
	}

	var response GetCharacterListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Errorf("cannot unmarshal body: %v", body)
	}
	if response.Characters[0].Name == "" {
		t.Errorf("character name is not %v", response.Characters[0].Name)
	}
}
