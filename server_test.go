package ca_game_api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("postgres", "postgres://test:password@localhost:15432/db_test?sslmode=disable")
	if err != nil {
		log.Fatal("open testDB failed:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("close testDB failed:", err)
		}
	}()

	m.Run()
}

func setupTestDB(tb testing.TB) {
	file := filepath.Join(".", "build", "setup.sql")
	cmd := exec.Command("psql", "-U", "test", "-h", "localhost", "-p", "15432", "-d", "db_test", "-a", "-f", file)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run failed.\n Command Output: %+v\n %+v, %v", stdout.String(), stderr.String(), err)
	}

	/*
		test user 1 の x-token は ceKeMPeYr0eF3K5e4Lfjfe
		test user 2 の x-token は yypKkCsMXx2MBBVorFQBsQ
		test user 3 の x-token は UGjoBQOXIjVHMWT7wpH5Ow
	*/
	const createTestUser = `
INSERT INTO users (name, digest_token) VALUES ('test1', '71a6f9c1007c60601a6d67e7f79d4550602b34ced90cdac86bd340f293bf0247'),
                                              ('test2', '541d9abc4b06e838e471ff564c24585a6ddc5280c9478f2e6e85b2eb7ed979a9'),
                                              ('test3', '3a80da5cb241be83d0275219c728c9e40cb8f17a433d776dbdb51741a7b49bce');
`
	if _, err := db.Exec(createTestUser); err != nil {
		tb.Fatal("db.Exec createTestUser failed:", err)
	}
}

func shutdownTestDB(tb testing.TB) {
	const dropUserCharactersTable = `DROP TABLE IF EXISTS user_characters;`
	const dropGachaResultsTable = `DROP TABLE IF EXISTS gacha_results;`
	const dropCharactersTable = `DROP TABLE IF EXISTS characters;`
	const dropUsersTable = `DROP TABLE IF EXISTS users;`

	if _, err := db.Exec(dropUserCharactersTable); err != nil {
		tb.Fatal("db.Exec dropUserCharactersTable failed:", err)
	}
	if _, err := db.Exec(dropGachaResultsTable); err != nil {
		tb.Fatal("db.Exec dropGachaResultsTable failed:", err)
	}
	if _, err := db.Exec(dropCharactersTable); err != nil {
		tb.Fatal("db.Exec dropCharactersTable failed:", err)
	}
	if _, err := db.Exec(dropUsersTable); err != nil {
		tb.Fatal("db.Exec dropUsersTable failed:", err)
	}
}

func generateTestResponse(tb testing.TB, w *httptest.ResponseRecorder, body interface{}) *http.Response {
	resp := w.Result()

	if body != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			tb.Error("io.ReadAll failed:", err)
		}
		if err := json.Unmarshal(bodyBytes, body); err != nil {
			tb.Error("json.Unmarshal failed:", err)
		}
	}

	return resp
}
