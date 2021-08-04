package ca_game_api

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestGetCharacterList(tb testing.TB) {
	const createUserCharacters = `
INSERT INTO user_characters (user_id, character_id, experience) VALUES (1, 30000001, 100),
                                                                       (1, 30000002, 100),
                                                                       (1, 30000003, 900);
`
	if _, err := db.Exec(createUserCharacters); err != nil {
		tb.Fatal("setupTestGetCharacterList failed:", err)
	}
}

func TestGetCharacterList(t *testing.T) {
	setupTestDB(t)
	setupTestGetCharacterList(t)
	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/character/list", nil)
		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
		w := httptest.NewRecorder()

		GetCharacterList(w, r)

		var body GetCharacterListResponse
		resp := generateTestResponse(t, w, &body)

		if resp.StatusCode != 200 {
			t.Errorf("Status code should be 200, but %v", resp.StatusCode)
		}
		if len(body.Characters) != 3 {
			t.Errorf("number of characters should be 3, but %v", len(body.Characters))
		}

		character1 := body.Characters[0]
		if character1.UserCharacterId != 1 {
			t.Errorf("1st character's userCharacterID should be 1, but %v", character1.UserCharacterId)
		}
		if character1.CharacterId != 30000001 {
			t.Errorf("1st character's characterID should be 30000001, but %v", character1.CharacterId)
		}
		if character1.Name != "normal_character1" {
			t.Errorf("1st character's name should be normal_character1, but %v", character1.Name)
		}
		if character1.Level != 1 {
			t.Errorf("1st character's level should be 1, but %v", character1.Level)
		}
		if character1.Experience != 100 {
			t.Errorf("1st character's experience should be 100, but %v", character1.Experience)
		}
		if character1.Power != 1 {
			t.Errorf("1st character's power should be 1, but %v", character1.Power)
		}

		character2 := body.Characters[1]
		if character2.UserCharacterId != 2 {
			t.Errorf("2nd character's userCharacterID should be 2, but %v", character2.UserCharacterId)
		}
		if character2.CharacterId != 30000002 {
			t.Errorf("2nd character's characterID should be 30000002, but %v", character2.CharacterId)
		}
		if character2.Name != "normal_character2" {
			t.Errorf("2nd character's name should be normal_character2, but %v", character2.Name)
		}
		if character2.Level != 1 {
			t.Errorf("2nd character's level should be 1, but %v", character2.Level)
		}
		if character2.Experience != 100 {
			t.Errorf("2nd character's experience should be 100, but %v", character2.Experience)
		}
		if character2.Power != 200 {
			t.Errorf("2nd character's power should be 200, but %v", character2.Power)
		}

		character3 := body.Characters[2]
		if character3.UserCharacterId != 3 {
			t.Errorf("3rd character's userCharacterID should be 3, but %v", character3.UserCharacterId)
		}
		if character3.CharacterId != 30000003 {
			t.Errorf("3rd character's characterID should be 30000003, but %v", character3.CharacterId)
		}
		if character3.Name != "normal_character3" {
			t.Errorf("3rd character's name should be normal_character3, but %v", character3.Name)
		}
		if character3.Level != 3 {
			t.Errorf("3rd character's level should be 3, but %v", character3.Level)
		}
		if character3.Experience != 900 {
			t.Errorf("3rd character's experience should be 900, but %v", character3.Experience)
		}
		if character3.Power != 540 {
			t.Errorf("3rd character's power should be 540, but %v", character3.Power)
		}
	})
	t.Run("Exception by bad x-token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/character/list", nil)
		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxxx")
		w := httptest.NewRecorder()

		GetCharacterList(w, r)

		resp := w.Result()

		if resp.StatusCode != 401 {
			t.Errorf("Status code should be 401, but %v", resp.StatusCode)
		}
	})
	t.Run("Exception by bad request method", func(t *testing.T) {
		r1 := httptest.NewRequest("POST", "/character/list", nil)
		w1 := httptest.NewRecorder()
		r2 := httptest.NewRequest("PUT", "/character/list", nil)
		w2 := httptest.NewRecorder()
		r3 := httptest.NewRequest("DELETE", "/character/list", nil)
		w3 := httptest.NewRecorder()

		GetCharacterList(w1, r1)
		GetCharacterList(w2, r2)
		GetCharacterList(w3, r3)

		resp1 := generateTestResponse(t, w1, nil)
		resp2 := generateTestResponse(t, w2, nil)
		resp3 := generateTestResponse(t, w3, nil)

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
	shutdownTestDB(t)
}

func setupTestPutCharacterCompose(tb testing.TB) {
	const createUserCharacters = `
INSERT INTO user_characters (user_id, character_id, experience) VALUES (1, 30000001, 100),
                                                                       (1, 30000002, 100);
`
	if _, err := db.Exec(createUserCharacters); err != nil {
		tb.Fatal("setupTestGetCharacterList failed:", err)
	}
}

func TestPutCharacterCompose(t *testing.T) {
	setupTestDB(t)
	setupTestPutCharacterCompose(t)
	t.Run("Success", func(t *testing.T) {
		reqBody := strings.NewReader(`
{
 "baseUserCharacterID": 2,
 "materialUserCharacterID": 1 
}
`)
		r := httptest.NewRequest("PUT", "/character/compose", reqBody)
		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
		w := httptest.NewRecorder()

		PutCharacterCompose(w, r)

		var body PutCharacterComposeResponse
		resp := generateTestResponse(t, w, &body)

		if resp.StatusCode != 200 {
			t.Errorf("Status code should be 200, but %v", resp.StatusCode)
		}
		if body.UserCharacterId != 2 {
			t.Errorf("userCharacterID should be 2, but %v", body.UserCharacterId)
		}
		if body.CharacterId != 30000002 {
			t.Errorf("characterID should be 30000002, but %v", body.CharacterId)
		}
		if body.Name != "normal_character2" {
			t.Errorf("name should be normal_character2, but %v", body.Name)
		}
		if body.Level == 1 {
			t.Errorf("level should be up, but %v", body.Level)
		}
		if body.Power == 200 {
			t.Errorf("power should be up, but %v", body.Power)
		}
	})
	t.Run("Exception by bad request body", func(t *testing.T) {
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
	t.Run("Exception by bad x-token", func(t *testing.T) {
		reqBody := strings.NewReader(`
{
 "baseUserCharacterID": 1,
 "materialUserCharacterID": 2
}
`)
		r := httptest.NewRequest("PUT", "/character/compose", reqBody)
		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxx")
		w := httptest.NewRecorder()

		PutCharacterCompose(w, r)

		resp := generateTestResponse(t, w, nil)

		if resp.StatusCode != 401 {
			t.Errorf("Status code should be 401, but %v", resp.StatusCode)
		}
	})
	t.Run("Exception by bad request method", func(t *testing.T) {
		r1 := httptest.NewRequest("GET", "/character/compose", nil)
		w1 := httptest.NewRecorder()
		r2 := httptest.NewRequest("POST", "/character/compose", nil)
		w2 := httptest.NewRecorder()
		r3 := httptest.NewRequest("DELETE", "/character/compose", nil)
		w3 := httptest.NewRecorder()

		PutCharacterCompose(w1, r1)
		PutCharacterCompose(w2, r2)
		PutCharacterCompose(w3, r3)

		resp1 := generateTestResponse(t, w1, nil)
		resp2 := generateTestResponse(t, w2, nil)
		resp3 := generateTestResponse(t, w3, nil)

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
	shutdownTestDB(t)
}
