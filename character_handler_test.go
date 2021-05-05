package ca_game_api

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
)

func TestGetCharacterList(t *testing.T) {
	r := httptest.NewRequest("GET", "/character/list", nil)
	r.Header.Set("x-token", "UGjoBQOXIjVHMWT7wpH5Ow")
	w := httptest.NewRecorder()

	GetCharacterList(w, r)

	resp := w.Result()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("io.ReadAll failed: %v", err)
	}
	var body GetCharacterListResponse
	if err := json.Unmarshal(bytes, &body); err != nil {
		t.Errorf("json.Unmarshal failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Response code should be 200, but %v", resp.StatusCode)
	}
	if len(body.Characters) != 2 {
		t.Errorf("number of characters should be 2, but %v", len(body.Characters))
	}

	if body.Characters[0].UserCharacterId == 0 {
		t.Errorf("1st character's userCharacterID does not exist")
	}
	if body.Characters[0].CharacterId != 50000002 {
		t.Errorf("1st character's characterID should be 50000002, but %v", body.Characters[0].CharacterId)
	}
	if body.Characters[0].Name != "super_rare_character2" {
		t.Errorf("1st character's name should be super_rare_character2, but %v", body.Characters[0].Name)
	}
	if body.Characters[0].Level != 1 {
		t.Errorf("1st character's level should be 1, but %v", body.Characters[0].Level)
	}
	if body.Characters[0].Experience != 100 {
		t.Errorf("1st character's experience should be 100, but %v", body.Characters[0].Experience)
	}
	if body.Characters[0].Power != 500 {
		t.Errorf("1st character's power should be 500, but %v", body.Characters[0].Power)
	}

	if body.Characters[1].UserCharacterId == 0 {
		t.Errorf("2nd character's userCharacterID does not exist")
	}
	if body.Characters[1].CharacterId != 30000002 {
		t.Errorf("2nd character's characterID should be 30000002, but %v", body.Characters[0].CharacterId)
	}
	if body.Characters[1].Name != "normal_character2" {
		t.Errorf("2nd character's name should be normal_character2, but %v", body.Characters[0].Name)
	}
	if body.Characters[1].Level != 2 {
		t.Errorf("2nd character's level should be 2, but %v", body.Characters[0].Level)
	}
	if body.Characters[1].Experience != 400 {
		t.Errorf("2nd character's experience should be 400, but %v", body.Characters[0].Experience)
	}
	if body.Characters[1].Power != 400 {
		t.Errorf("2nd character's power should be 400, but %v", body.Characters[0].Power)
	}
}
