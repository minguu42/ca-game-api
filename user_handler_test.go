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

	code := m.Run()
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func TestPostUser(t *testing.T) {
	name, err := generateRandomString(8)
	if err != nil {
		t.Errorf("cannot generate random string")
	}
	reqBody := strings.NewReader(`{"name":"` + name + `"}`)
	r := httptest.NewRequest("POST", "/user/post", reqBody)
	w := httptest.NewRecorder()

	PostUser(w, r)

	resp := w.Result()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("cannot read response body: %v", resp.Body)
	}
	var body PostUserResponse
	if err := json.Unmarshal(bytes, &body); err != nil {
		t.Errorf("cannot unmarshal bytes")
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("cannot read body: %v", resp.Body)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Response code is %v", w.Code)
	}

	var response GetUserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Errorf("cannot unmarshal body: %v", body)
	}
	if response.Name != "test user" {
		t.Errorf("user name is %v", response.Name)
	}
}

func TestPutUser(t *testing.T) {
	name, _ := generateRandomString(8)
	reqBody := strings.NewReader(`{"name":"` + name + `"}`)
	r := httptest.NewRequest("PUT", "/user/update", reqBody)
	r.Header.Set("x-token", "yypKkCsMXx2MBBVorFQBsQ")
	w := httptest.NewRecorder()

	PutUser(w, r)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("response code is %v", resp.StatusCode)
	}
}

func TestGetUserRanking(t *testing.T) {
	r := httptest.NewRequest("GET", "/user/ranking", nil)
	r.Header.Set("x-token", "yypKkCsMXx2MBBVorFQBsQ")
	w := httptest.NewRecorder()

	GetUserRanking(w, r)

	resp := w.Result()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("cannot read body: %v", resp.Body)
	}

	if resp.StatusCode != 200 {
		t.Errorf("response code is %v", resp.StatusCode)
	}

	var response GetUserRankingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Errorf("cannot unmarshal body: %v", body)
	}
	if len(response.Users) != 3 {
		t.Errorf("ranking up to 3rd, but %v", len(response.Users))
	}
	if response.Users[0].Name != "test user" {
		t.Errorf("rank1 user is not test user, %v", response.Users[0].Name)
	}
	if response.Users[1].Name != "test user3" {
		t.Errorf("rank2 user is not test user2, %v", response.Users[1].Name)
	}
	if response.Users[2].Name != "test user4" {
		t.Errorf("rank3 user is not test user3, %v", response.Users[2].Name)
	}
}
