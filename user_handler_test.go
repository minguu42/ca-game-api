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
	db, err = sql.Open("postgres", "postgres://test:password@localhost:15432/ca_game_api_db_test?sslmode=disable")
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
		t.Errorf("Response code should be 200, but %v", w.Code)
	}
	if body.Token == "" {
		t.Errorf("Token does not exist")
	}
	if len(body.Token) != 22 {
		t.Errorf("Token length should be 22, but %v", len(body.Token))
	}
}

func TestGetUser(t *testing.T) {
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
		t.Errorf("Response code should be 200, but %v", w.Code)
	}
	if body.Name != "test1" {
		t.Errorf("Name should be test1, but %v", body.Name)
	}
}

func TestPutUser(t *testing.T) {
	name, err := generateRandomString(8)
	if err != nil {
		t.Errorf("generateRandomString failed: %v", err)
	}
	reqBody := strings.NewReader(`{"name":"` + name + `"}`)
	r := httptest.NewRequest("PUT", "/user/update", reqBody)
	r.Header.Set("x-token", "yypKkCsMXx2MBBVorFQBsQ")
	w := httptest.NewRecorder()

	PutUser(w, r)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("Response code should be 200, but %v", resp.StatusCode)
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
		t.Errorf("Response code should be 200, but %v", resp.StatusCode)
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
	if body.Users[1].SumPower != 900 {
		t.Errorf("test3's sumPower should be 900, but %v", body.Users[1].SumPower)
	}
}
