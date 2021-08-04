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
	t.Run("Exception by bad request method", func(t *testing.T) {
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
	t.Run("Exception by name already exists", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name": "test1"}`)
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

func TestGetUser(t *testing.T) {
	setupTestDB(t)
	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/user/get", nil)
		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
		w := httptest.NewRecorder()

		GetUser(w, r)

		var body GetUserResponse
		resp := generateTestResponse(t, w, &body)

		if resp.StatusCode != 200 {
			t.Errorf("Status code should be 200, but %v", resp.StatusCode)
		}
		if body.Name != "test1" {
			t.Errorf("Name should be test1, but %v", body.Name)
		}
	})
	t.Run("Exception by bad x-token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/user/get", nil)
		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxx")
		w := httptest.NewRecorder()

		GetUser(w, r)

		resp := generateTestResponse(t, w, nil)

		if resp.StatusCode != 401 {
			t.Errorf("Status code should be 401, but %v", resp.StatusCode)
		}
	})
	t.Run("Exception by bad request method", func(t *testing.T) {
		r1 := httptest.NewRequest("POST", "/user/get", nil)
		w1 := httptest.NewRecorder()
		r2 := httptest.NewRequest("PUT", "/user/get", nil)
		w2 := httptest.NewRecorder()
		r3 := httptest.NewRequest("DELETE", "/user/get", nil)
		w3 := httptest.NewRecorder()

		GetUser(w1, r1)
		GetUser(w2, r2)
		GetUser(w3, r3)

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

func TestPutUser(t *testing.T) {
	setupTestDB(t)
	t.Run("Success", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name": "new user"}`)
		r := httptest.NewRequest("PUT", "/user/update", reqBody)
		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
		w := httptest.NewRecorder()

		PutUser(w, r)

		resp := w.Result()

		if resp.StatusCode != 200 {
			t.Errorf("Status code should be 200, but %v", resp.StatusCode)
		}
	})
	t.Run("Exception by bad request body", func(t *testing.T) {
		reqBody := strings.NewReader(``)
		r := httptest.NewRequest("PUT", "/user/update", reqBody)
		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
		w := httptest.NewRecorder()

		PutUser(w, r)

		resp := generateTestResponse(t, w, nil)

		if resp.StatusCode != 400 {
			t.Errorf("Status code should be 400, but %v", resp.StatusCode)
		}
	})
	t.Run("Exception by bad x-token", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name": "new user"}`)
		r := httptest.NewRequest("PUT", "/user/update", reqBody)
		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxx")
		w := httptest.NewRecorder()

		PutUser(w, r)

		resp := generateTestResponse(t, w, nil)

		if resp.StatusCode != 401 {
			t.Errorf("Status code should be 401, but %v", resp.StatusCode)
		}
	})
	t.Run("Exception by bad request method", func(t *testing.T) {
		r1 := httptest.NewRequest("GET", "/user/update", nil)
		w1 := httptest.NewRecorder()
		r2 := httptest.NewRequest("POST", "/user/update", nil)
		w2 := httptest.NewRecorder()
		r3 := httptest.NewRequest("DELETE", "/user/update", nil)
		w3 := httptest.NewRecorder()

		PutUser(w1, r1)
		PutUser(w2, r2)
		PutUser(w3, r3)

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
	t.Run("Exception by name already exists", func(t *testing.T) {
		reqBody := strings.NewReader(`{"name": "test2"}`)
		r := httptest.NewRequest("PUT", "/user/update", reqBody)
		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
		w := httptest.NewRecorder()

		PutUser(w, r)

		resp := generateTestResponse(t, w, nil)

		if resp.StatusCode != 409 {
			t.Errorf("Status code should be 409, but %v", resp.StatusCode)
		}
	})
	shutdownTestDB(t)
}

func setupTestGetUserRanking(tb testing.TB) {
	const createUserCharacters = `
INSERT INTO user_characters (user_id, character_id, experience) VALUES (1, 50000002, 100),
                                                                       (2, 40000002, 100),
                                                                       (3, 30000002, 100);
`
	if _, err := db.Exec(createUserCharacters); err != nil {
		tb.Fatal("setupTestGetUserRanking failed:", err)
	}
}

func TestGetUserRanking(t *testing.T) {
	setupTestDB(t)
	setupTestGetUserRanking(t)
	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/user/ranking", nil)
		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
		w := httptest.NewRecorder()

		GetUserRanking(w, r)

		var body GetUserRankingResponse
		resp := generateTestResponse(t, w, &body)

		if resp.StatusCode != 200 {
			t.Errorf("Status code should be 200, but %v", resp.StatusCode)
		}
		if len(body.Users) != 3 {
			t.Errorf("Ranking up to 3rd, but response body include %v users", len(body.Users))
		}

		test1 := body.Users[0]
		test2 := body.Users[1]
		test3 := body.Users[2]
		if test1.Name != "test1" {
			t.Errorf("1st user name should be test1, but %v", test1.Name)
		}
		if test1.SumPower != 500 {
			t.Errorf("1st user sumPower should be 500, but %v", test1.SumPower)
		}
		if test2.Name != "test2" {
			t.Errorf("2st user name should be test2, but %v", test2.Name)
		}
		if test2.SumPower != 300 {
			t.Errorf("2st user sumPower should be 300, but %v", test2.SumPower)
		}
		if test3.Name != "test3" {
			t.Errorf("3st user name should be test3, but %v", test3.Name)
		}
		if test3.SumPower != 200 {
			t.Errorf("3st user sumPower should be 200, but %v", test3.SumPower)
		}
	})
	t.Run("Exception by bad x-token", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/user/ranking", nil)
		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxx")
		w := httptest.NewRecorder()

		GetUserRanking(w, r)

		resp := w.Result()

		if resp.StatusCode != 401 {
			t.Errorf("Status code should be 401, but %v", resp.StatusCode)
		}
	})
	t.Run("Exception by bad request method", func(t *testing.T) {
		r1 := httptest.NewRequest("POST", "/user/ranking", nil)
		w1 := httptest.NewRecorder()
		r2 := httptest.NewRequest("PUT", "/user/ranking", nil)
		w2 := httptest.NewRecorder()
		r3 := httptest.NewRequest("DELETE", "/user/ranking", nil)
		w3 := httptest.NewRecorder()

		GetUserRanking(w1, r1)
		GetUserRanking(w2, r2)
		GetUserRanking(w3, r3)

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
