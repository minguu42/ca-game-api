package ca_game_api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

var mux *http.ServeMux
var writer *httptest.ResponseRecorder

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	var err error
	mux = http.NewServeMux()
	mux.HandleFunc("/user/create", PostUser)
	mux.HandleFunc("/user/get", GetUser)
	writer = httptest.NewRecorder()
	db, err = sql.Open("postgres", "postgres://test:password@localhost:15432/ca_game_api_db_test?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func tearDown() {
	if err := db.Close(); err != nil {
		log.Println(err)
	}
}

func TestPostUser(t *testing.T) {
	name, _ := generateRandomString(8)
	requestBody := strings.NewReader(`{"name":"` + name + `"}`)
	request, _ := http.NewRequest("POST", "/user/create", requestBody)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	var response PostUserResponse
	if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
		log.Println(err)
	}
	if response.Token == "" {
		t.Errorf("token is none")
	}
}

func TestGetUser(t *testing.T) {
	writer = httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/user/get", nil)
	request.Header.Set("x-token", "ceKeMPeYr0eF3K5e4Lfjfe")
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	var response GetUserResponse
	if err := json.Unmarshal(writer.Body.Bytes(), &response); err != nil {
		log.Println("json.Unmarshal failed:", err)
	}
	if response.Name != "test user" {
		t.Errorf("user name is %v", response.Name)
	}
}
