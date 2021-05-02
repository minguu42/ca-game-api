package ca_game_api

import (
	"database/sql"
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
	json := strings.NewReader(`{"name":"minguu3"}`)
	request, _ := http.NewRequest("POST", "/user/create", json)
	mux.ServeHTTP(writer, request)
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}
