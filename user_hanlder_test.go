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
	name, _ := generateRandomString(8)
	requestBody := strings.NewReader(`{"name":"` + name + `"}`)
	r := httptest.NewRequest("POST", "/user/post", requestBody)
	w := httptest.NewRecorder()

	PostUser(w, r)

	resp := w.Result()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("cannot read test response: %v", w.Code)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Response code is %v", w.Code)
	}
	var response PostUserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Errorf("cannot unmarshal body: %v", body)
	}
	if response.Token == "" {
		t.Errorf("token is none")
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
		t.Errorf("cannot read test response: %v", w.Code)
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
