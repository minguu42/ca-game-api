package ca_game_api

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("postgres", "postgres://test:password@localhost:15432/test_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	setupPutCharacterCompose()

	code := m.Run()

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}

func TestPostUser(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		name, err := generateRandomString(8)
		if err != nil {
			t.Errorf("generateRandomString failed: %v", err)
		}
		reqBody := strings.NewReader(`{"name":"` + name + `"}`)
		r := httptest.NewRequest("POST", "/user/post", reqBody)
		w := httptest.NewRecorder()

		PostUser(w, r)

		resp := w.Result()

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("io.ReadAll failed: %v", err)
		}
		var body PostUserResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Errorf("json.Unmarshal failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Status code should be 200, but %v", w.Code)
		}
		if body.Token == "" {
			t.Errorf("Token does not exist")
		}
		if len(body.Token) != 22 {
			t.Errorf("Token length should be 22, but %v", len(body.Token))
		}
	})

	t.Run("Bad request method", func(t *testing.T) {
		r1 := httptest.NewRequest("GET", "/user/create", nil)
		w1 := httptest.NewRecorder()
		r2 := httptest.NewRequest("PUT", "/user/create", nil)
		w2 := httptest.NewRecorder()
		r3 := httptest.NewRequest("DELETE", "/user/create", nil)
		w3 := httptest.NewRecorder()

		PostUser(w1, r1)
		PostUser(w2, r2)
		PostUser(w3, r3)

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

	t.Run("Bad request body", func(t *testing.T) {
		reqBody := strings.NewReader(`{}`)
		r := httptest.NewRequest("POST", "/character/compose", reqBody)
		w := httptest.NewRecorder()

		PostUser(w, r)

		resp := w.Result()

		if resp.StatusCode != 400 {
			t.Errorf("Status code should be 400, but %v", resp.StatusCode)
		}
	})
}

func TestGetUser(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/user/get", nil)
		r.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
		w := httptest.NewRecorder()

		GetUser(w, r)

		resp := w.Result()

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("io.ReadAll failed: %v", err)
		}
		var body GetUserResponse
		if err := json.Unmarshal(bytes, &body); err != nil {
			t.Errorf("json.Unmarshal failed: %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("Status code should be 200, but %v", w.Code)
		}
		if body.Name != "test1" {
			t.Errorf("Name should be test1, but %v", body.Name)
		}
	})

	t.Run("Bad request method", func(t *testing.T) {
		r1 := httptest.NewRequest("POST", "/user/get", nil)
		w1 := httptest.NewRecorder()
		r2 := httptest.NewRequest("PUT", "/user/get", nil)
		w2 := httptest.NewRecorder()
		r3 := httptest.NewRequest("DELETE", "/user/get", nil)
		w3 := httptest.NewRecorder()

		GetUser(w1, r1)
		GetUser(w2, r2)
		GetUser(w3, r3)

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
		r := httptest.NewRequest("GET", "/user/get", nil)
		r.Header.Set("x-token", "xxxxxxxxxxxxxxxxxxxxxx")
		w := httptest.NewRecorder()

		GetUser(w, r)

		resp := w.Result()

		if resp.StatusCode != 403 {
			t.Errorf("Status code should be 403, but %v", resp.StatusCode)
		}
	})
}

func TestPutUser(t *testing.T) {
	name, err := generateRandomString(8)
	if err != nil {
		t.Errorf("generateRandomString failed: %v", err)
	}
	reqBody := strings.NewReader(`{"name":"` + name + `"}`)
	r := httptest.NewRequest("PUT", "/user/updateUser", reqBody)
	r.Header.Set("x-token", "yypKkCsMXx2MBBVorFQBsQ")
	w := httptest.NewRecorder()

	PutUser(w, r)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("Status code should be 200, but %v", resp.StatusCode)
	}
}

func TestGetUserRanking(t *testing.T) {
	r := httptest.NewRequest("GET", "/user/ranking", nil)
	r.Header.Set("x-token", "yypKkCsMXx2MBBVorFQBsQ")
	w := httptest.NewRecorder()

	GetUserRanking(w, r)

	resp := w.Result()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("io.ReadAll failed: %v", err)
	}
	var body GetUserRankingResponse
	if err := json.Unmarshal(bytes, &body); err != nil {
		t.Errorf("json.Unmarshal failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Status code should be 200, but %v", resp.StatusCode)
	}
	if len(body.Users) != 3 {
		t.Errorf("Ranking up to 3rd, but response body include %v users", len(body.Users))
	}
	if body.Users[0].Name != "test1" {
		t.Errorf("1st user name should be test1, but %v", body.Users[0].Name)
	}
	if body.Users[1].Name != "test3" {
		t.Errorf("2nd user name should be test3, but %v", body.Users[1].Name)
	}
	if body.Users[2].Name != "test4" {
		t.Errorf("3rd user name should be test4, but %v", body.Users[2].Name)
	}
	if body.Users[1].SumPower != 130000 {
		t.Errorf("test3's sumPower should be 130000, but %v", body.Users[1].SumPower)
	}
}
