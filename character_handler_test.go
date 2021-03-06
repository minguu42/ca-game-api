package ca_game_api

import (
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestGetCharacterList(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
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
			t.Errorf("Status code should be 200, but %v", resp.StatusCode)
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
	})

	t.Run("Bad request method", func(t *testing.T) {
		r1 := httptest.NewRequest("POST", "/character/list", nil)
		w1 := httptest.NewRecorder()
		r2 := httptest.NewRequest("PUT", "/character/list", nil)
		w2 := httptest.NewRecorder()
		r3 := httptest.NewRequest("DELETE", "/character/list", nil)
		w3 := httptest.NewRecorder()

		GetCharacterList(w1, r1)
		GetCharacterList(w2, r2)
		GetCharacterList(w3, r3)

		resp1 := w1.Result()
		resp2 := w2.Result()
		resp3 := w3.Result()

		if resp1.StatusCode != 405 {
			t.Errorf("Status code should be 405, but %v", resp1.StatusCode)
		}
		if resp2.StatusCode != 405 {
			t.Errorf("Status code should be 405, but %v", resp2.StatusCode)
		}
		if resp3.StatusCode != 405 {
			t.Errorf("Status code should be 405, but %v", resp3.StatusCode)
		}
	})

	t.Run("Bad request parameters", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/character/list", nil)
		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxxx")
		w := httptest.NewRecorder()

		GetCharacterList(w, r)

		resp := w.Result()

		if resp.StatusCode != 400 {
			t.Errorf("Status code should be 400, but %v", resp.StatusCode)
		}
	})
}

var materialUserCharacterId int

func setupPutCharacterCompose() {
	if err := db.QueryRow(`
INSERT INTO user_characters (user_id, character_id, experience)
VALUES (1, 30000002, 100)
RETURNING id
`).Scan(&materialUserCharacterId); err != nil {
		log.Println("setupPutCharacterCompose failed: ", err)
	}
}

func TestPutCharacterCompose(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		reqBody := strings.NewReader(`
{
  "baseUserCharacterID": 1,
  "materialUserCharacterID":` + strconv.Itoa(materialUserCharacterId) + `
}
`)
		r := httptest.NewRequest("PUT", "/character/compose", reqBody)
		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
		w := httptest.NewRecorder()

		PutCharacterCompose(w, r)

		resp := w.Result()

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("io.ReadAll failed: %v", err)
		}
		var body PutCharacterComposeResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Errorf("json.Unmarshal failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Status code should be 200, but %v", resp.StatusCode)
		}
		if body.UserCharacterId != 1 {
			t.Errorf("userCharacterID should be 1, but %v", body.UserCharacterId)
		}
		if body.CharacterId != 50000002 {
			t.Errorf("characterID should be 50000002, but %v", body.CharacterId)
		}
		if body.Name != "super_rare_character2" {
			t.Errorf("name should be super_rare_character2, but %v", body.Name)
		}
		if body.Level <= 0 {
			t.Errorf("level should be positive number, but %v", body.Level)
		}
		if body.Power == 0 {
			t.Errorf("power does not exist")
		}
	})

	t.Run("Bad request method", func(t *testing.T) {
		r1 := httptest.NewRequest("GET", "/character/compose", nil)
		w1 := httptest.NewRecorder()
		r2 := httptest.NewRequest("POST", "/character/compose", nil)
		w2 := httptest.NewRecorder()
		r3 := httptest.NewRequest("DELETE", "/character/compose", nil)
		w3 := httptest.NewRecorder()

		PutCharacterCompose(w1, r1)
		PutCharacterCompose(w2, r2)
		PutCharacterCompose(w3, r3)

		resp1 := w1.Result()
		resp2 := w2.Result()
		resp3 := w3.Result()

		if resp1.StatusCode != 405 {
			t.Errorf("Status code should be 405, but %v", resp1.StatusCode)
		}
		if resp2.StatusCode != 405 {
			t.Errorf("Status code should be 405, but %v", resp2.StatusCode)
		}
		if resp3.StatusCode != 405 {
			t.Errorf("Status code should be 405, but %v", resp3.StatusCode)
		}
	})

	t.Run("Bad request parameters", func(t *testing.T) {
		reqBody := strings.NewReader(`
{
  "baseUserCharacterID": 1,
  "materialUserCharacterID":` + strconv.Itoa(materialUserCharacterId) + `
}
`)
		r := httptest.NewRequest("PUT", "/character/compose", reqBody)
		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxx")
		w := httptest.NewRecorder()

		PutCharacterCompose(w, r)

		resp := w.Result()

		if resp.StatusCode != 403 {
			t.Errorf("Status code should be 403, but %v", resp.StatusCode)
		}
	})

	t.Run("Bad request body", func(t *testing.T) {
		reqBody := strings.NewReader(`
{
  "baseUserCharacterID": 1
}
`)
		r := httptest.NewRequest("PUT", "/character/compose", reqBody)
		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
		w := httptest.NewRecorder()

		PutCharacterCompose(w, r)

		resp := w.Result()

		if resp.StatusCode != 400 {
			t.Errorf("Status code should be 400, but %v", resp.StatusCode)
		}
	})
}
