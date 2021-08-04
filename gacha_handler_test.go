package ca_game_api

//import (
//	"encoding/json"
//	"io"
//	"net/http/httptest"
//	"strings"
//	"testing"
//)
//
//func TestPostGachaDraw(t *testing.T) {
//	t.Run("OK", func(t *testing.T) {
//		reqBody := strings.NewReader(`{"times": 3}`)
//		r := httptest.NewRequest("POST", "/gacha/draw", reqBody)
//		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
//		w := httptest.NewRecorder()
//
//		PostGachaDraw(w, r)
//
//		resp := w.Result()
//
//		bytes, err := io.ReadAll(resp.Body)
//		if err != nil {
//			t.Errorf("io.ReadAll failed: %v", err)
//		}
//		var body PostGachaDrawResponse
//		if err := json.Unmarshal(bytes, &body); err != nil {
//			t.Errorf("json.Unmarshal failed: %v", err)
//		}
//
//		if resp.StatusCode != 200 {
//			t.Errorf("Status code should be 200, but %v", resp.StatusCode)
//		}
//		if len(body.Results) != 3 {
//			t.Errorf("Results num should be 3, but %v", len(body.Results))
//		}
//		if body.Results[0].CharacterId == 0 {
//			t.Error(`1st character's characterID does not exist`)
//		}
//		if body.Results[0].Name == "" {
//			t.Error(`1st character's name does not exist`)
//		}
//		if body.Results[1].CharacterId == 0 {
//			t.Error(`2nd character's characterID does not exist`)
//		}
//		if body.Results[1].Name == "" {
//			t.Error(`2nd character's name does not exist`)
//		}
//		if body.Results[2].CharacterId == 0 {
//			t.Error(`3rd character's characterID does not exist`)
//		}
//		if body.Results[2].Name == "" {
//			t.Error(`3rd character's name does not exist`)
//		}
//	})
//
//	t.Run("Bad request method", func(t *testing.T) {
//		r1 := httptest.NewRequest("GET", "/gacha/draw", nil)
//		w1 := httptest.NewRecorder()
//		r2 := httptest.NewRequest("PUT", "/gacha/draw", nil)
//		w2 := httptest.NewRecorder()
//		r3 := httptest.NewRequest("DELETE", "/gacha/draw", nil)
//		w3 := httptest.NewRecorder()
//
//		PostGachaDraw(w1, r1)
//		PostGachaDraw(w2, r2)
//		PostGachaDraw(w3, r3)
//
//		resp1 := w1.Result()
//		resp2 := w2.Result()
//		resp3 := w3.Result()
//
//		if resp1.StatusCode != 405 {
//			t.Errorf("Status code should be 405, but %v", resp1.StatusCode)
//		}
//		if resp2.StatusCode != 405 {
//			t.Errorf("Status code should be 405, but %v", resp2.StatusCode)
//		}
//		if resp3.StatusCode != 405 {
//			t.Errorf("Status code should be 405, but %v", resp3.StatusCode)
//		}
//	})
//
//	t.Run("Bad request parameters", func(t *testing.T) {
//		reqBody := strings.NewReader(`{"times": 3}`)
//		r := httptest.NewRequest("POST", "/gacha/draw", reqBody)
//		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxxx")
//		w := httptest.NewRecorder()
//
//		PostGachaDraw(w, r)
//
//		resp := w.Result()
//
//		if resp.StatusCode != 403 {
//			t.Errorf("Status code should be 403, but %v", resp.StatusCode)
//		}
//	})
//
//	t.Run("Bad request body", func(t *testing.T) {
//		reqBody := strings.NewReader(`{ "time": 5 }`)
//		r := httptest.NewRequest("POST", "/gacha/draw", reqBody)
//		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
//		w := httptest.NewRecorder()
//
//		PostGachaDraw(w, r)
//
//		resp := w.Result()
//
//		if resp.StatusCode != 400 {
//			t.Errorf("Status code should be 400, but %v", resp.StatusCode)
//		}
//	})
//}
