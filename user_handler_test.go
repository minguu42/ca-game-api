package ca_game_api

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostUser(t *testing.T) {
	setupTestDB(t)

	t.Run("Success", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name": "new user"}`)
		r := httptest.NewRequest("POST", "/user/post", reqBody)
		w := httptest.NewRecorder()

		PostUser(w, r)

		var body PostUserResponse
		resp := generateTestResponse(t, w, &body)

		if resp.StatusCode != 201 {
			t.Errorf("Status code should be 201, but %v", w.Code)
		}
		if body.Token == "" {
			t.Errorf("Token does not exist")
		}
		if len(body.Token) != 22 {
			t.Errorf("Token length should be 22, but %v", len(body.Token))
		}
	})

	t.Run("Exception with bad request method", func(t *testing.T) {
		r1 := httptest.NewRequest("GET", "/user/create", nil)
		w1 := httptest.NewRecorder()
		r2 := httptest.NewRequest("PUT", "/user/create", nil)
		w2 := httptest.NewRecorder()
		r3 := httptest.NewRequest("DELETE", "/user/create", nil)
		w3 := httptest.NewRecorder()

		PostUser(w1, r1)
		PostUser(w2, r2)
		PostUser(w3, r3)

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

	t.Run("Exception with name already exists", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name": "test user 1"}`)
		r := httptest.NewRequest("POST", "/user/create", reqBody)
		w := httptest.NewRecorder()

		PostUser(w, r)

		resp := generateTestResponse(t, w, nil)

		if resp.StatusCode != 409 {
			t.Errorf("Status code should be 409, but %v", resp.StatusCode)
		}
	})

	shutdownTestDB(t)
}

//func TestGetUser(t *testing.T) {
//	t.Run("OK", func(t *testing.T) {
//		r := httptest.NewRequest("GET", "/user/get", nil)
//		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
//		w := httptest.NewRecorder()
//
//		GetUser(w, r)
//
//		resp := w.Result()
//
//		bytes, err := io.ReadAll(resp.Body)
//		if err != nil {
//			t.Errorf("io.ReadAll failed: %v", err)
//		}
//		var body GetUserResponse
//		if err := json.Unmarshal(bytes, &body); err != nil {
//			t.Errorf("json.Unmarshal failed: %v", err)
//		}
//
//		if resp.StatusCode != 200 {
//			t.Errorf("Status code should be 200, but %v", w.Code)
//		}
//		if body.Name != "test1" {
//			t.Errorf("Name should be test1, but %v", body.Name)
//		}
//	})
//
//	t.Run("Bad request method", func(t *testing.T) {
//		r1 := httptest.NewRequest("POST", "/user/get", nil)
//		w1 := httptest.NewRecorder()
//		r2 := httptest.NewRequest("PUT", "/user/get", nil)
//		w2 := httptest.NewRecorder()
//		r3 := httptest.NewRequest("DELETE", "/user/get", nil)
//		w3 := httptest.NewRecorder()
//
//		GetUser(w1, r1)
//		GetUser(w2, r2)
//		GetUser(w3, r3)
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
//		r := httptest.NewRequest("GET", "/user/get", nil)
//		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxx")
//		w := httptest.NewRecorder()
//
//		GetUser(w, r)
//
//		resp := w.Result()
//
//		if resp.StatusCode != 403 {
//			t.Errorf("Status code should be 403, but %v", resp.StatusCode)
//		}
//	})
//}
//
//func TestPutUser(t *testing.T) {
//	t.Run("OK", func(t *testing.T) {
//		name, err := generateRandomString(8)
//		if err != nil {
//			t.Errorf("generateRandomString failed: %v", err)
//		}
//		reqBody := strings.NewReader(`{"name":"` + name + `"}`)
//		r := httptest.NewRequest("PUT", "/user/update", reqBody)
//		r.Header.Set("x-token", "yypKkCsMXx2MBBVorFQBsQ")
//		w := httptest.NewRecorder()
//
//		PutUser(w, r)
//
//		resp := w.Result()
//
//		if resp.StatusCode != 200 {
//			t.Errorf("Status code should be 200, but %v", resp.StatusCode)
//		}
//	})
//
//	t.Run("Bad request method", func(t *testing.T) {
//		r1 := httptest.NewRequest("GET", "/user/update", nil)
//		w1 := httptest.NewRecorder()
//		r2 := httptest.NewRequest("POST", "/user/update", nil)
//		w2 := httptest.NewRecorder()
//		r3 := httptest.NewRequest("DELETE", "/user/update", nil)
//		w3 := httptest.NewRecorder()
//
//		PutUser(w1, r1)
//		PutUser(w2, r2)
//		PutUser(w3, r3)
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
//		r := httptest.NewRequest("PUT", "/user/update", nil)
//		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxx")
//		w := httptest.NewRecorder()
//
//		PutUser(w, r)
//
//		resp := w.Result()
//
//		if resp.StatusCode != 400 {
//			t.Errorf("Status code should be 403, but %v", resp.StatusCode)
//		}
//	})
//
//	t.Run("Bad request body", func(t *testing.T) {
//		reqBody := strings.NewReader(`{}`)
//		r := httptest.NewRequest("PUT", "/user/update", reqBody)
//		r.Header.Set("x-token", "yypKkCsMXx2MBBVorFQBsQ")
//		w := httptest.NewRecorder()
//
//		PutUser(w, r)
//
//		resp := w.Result()
//
//		if resp.StatusCode != 400 {
//			t.Errorf("Status code should be 400, but %v", resp.StatusCode)
//		}
//	})
//}
//
//func TestGetUserRanking(t *testing.T) {
//	t.Run("OK", func(t *testing.T) {
//		r := httptest.NewRequest("GET", "/user/ranking", nil)
//		r.Header.Set("x-token", "yypKkCsMXx2MBBVorFQBsQ")
//		w := httptest.NewRecorder()
//
//		GetUserRanking(w, r)
//
//		resp := w.Result()
//
//		bytes, err := io.ReadAll(resp.Body)
//		if err != nil {
//			t.Errorf("io.ReadAll failed: %v", err)
//		}
//		var body GetUserRankingResponse
//		if err := json.Unmarshal(bytes, &body); err != nil {
//			t.Errorf("json.Unmarshal failed: %v", err)
//		}
//
//		if resp.StatusCode != 200 {
//			t.Errorf("Status code should be 200, but %v", resp.StatusCode)
//		}
//		if len(body.Users) != 3 {
//			t.Errorf("Ranking up to 3rd, but response body include %v users", len(body.Users))
//		}
//		if body.Users[0].Name != "test1" {
//			t.Errorf("1st user name should be test1, but %v", body.Users[0].Name)
//		}
//		if body.Users[1].Name != "test3" {
//			t.Errorf("2nd user name should be test3, but %v", body.Users[1].Name)
//		}
//		if body.Users[2].Name != "test4" {
//			t.Errorf("3rd user name should be test4, but %v", body.Users[2].Name)
//		}
//		if body.Users[1].SumPower != 130000 {
//			t.Errorf("test3's sumPower should be 130000, but %v", body.Users[1].SumPower)
//		}
//	})
//
//	t.Run("Bad request method", func(t *testing.T) {
//		r1 := httptest.NewRequest("POST", "/user/ranking", nil)
//		w1 := httptest.NewRecorder()
//		r2 := httptest.NewRequest("PUT", "/user/ranking", nil)
//		w2 := httptest.NewRecorder()
//		r3 := httptest.NewRequest("DELETE", "/user/ranking", nil)
//		w3 := httptest.NewRecorder()
//
//		GetUserRanking(w1, r1)
//		GetUserRanking(w2, r2)
//		GetUserRanking(w3, r3)
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
//		r := httptest.NewRequest("GET", "/user/ranking", nil)
//		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxx")
//		w := httptest.NewRecorder()
//
//		GetUserRanking(w, r)
//
//		resp := w.Result()
//
//		if resp.StatusCode != 403 {
//			t.Errorf("Status code should be 403, but %v", resp.StatusCode)
//		}
//	})
//}
