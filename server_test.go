package ca_game_api

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	setupPutCharacterCompose()

	m.Run()
}

func setupTestDB() {
	file, err := os.ReadFile(filepath.Join(".", "build", "setup.sql"))
	if err != nil {
		log.Fatal("os.ReadFile failed:", err)
	}

	queries := strings.Split(string(file), ";")

	for _, query := range queries {
		if query == "" {
			break
		}

		_, err := db.Exec(query)
		if err != nil {
			log.Fatal("db.Exec failed:", err)
		}
	}

	/*
		test1 はGetUser, GetUserRanking, PostGachaDraw で使用するユーザ. x-token は ceKeMPeYr0eF3K5e4Lfjfe である.
		test2 は PutUser で名前変更用ユーザ. 名前はテスト実行時にランダムな文字列に変わる. x-token は yypKkCsMXx2MBBVorFQBsQ である.
		test3 は GetCharacterList で使用するユーザ. x-token は UGjoBQOXIjVHMWT7wpH5Ow である.
	*/
	const createTestUser = `
INSERT INTO users (name, digest_token) VALUES ('test1', '71a6f9c1007c60601a6d67e7f79d4550602b34ced90cdac86bd340f293bf0247'),
                                              ('test2', '541d9abc4b06e838e471ff564c24585a6ddc5280c9478f2e6e85b2eb7ed979a9'),
                                              ('test3', '3a80da5cb241be83d0275219c728c9e40cb8f17a433d776dbdb51741a7b49bce'),
                                              ('test4', 'q34avo2q3avj9q28t4nq39vm9uz98qnq984j91oaoj9zu9q3ujvq9832j932q8ud'),
                                              ('test5', 'jdoije928jf9eqj1fnqz9duq921ejf6qwure2qi9jf7qc4xz98qw2urf98j9eqf1');
`
	if _, err := db.Exec(createTestUser); err != nil {
		log.Fatal("db.Exec createTestUser failed:", err)
	}
}

func shutdownTestDB() {
	const dropUserCharactersTable = `DROP TABLE IF EXISTS user_characters`
	const dropGachaResultsTable = `DROP TABLE IF EXISTS gacha_results`
	const dropCharactersTable = `DROP TABLE IF EXISTS characters`
	const dropUsersTable = `DROP TABLE IF EXISTS users`

	if _, err := db.Exec(dropUserCharactersTable); err != nil {
		log.Fatal("db.Exec dropUserCharactersTable failed:", err)
	}
	if _, err := db.Exec(dropGachaResultsTable); err != nil {
		log.Fatal("db.Exec dropGachaResultsTable failed:", err)
	}
	if _, err := db.Exec(dropCharactersTable); err != nil {
		log.Fatal("db.Exec dropCharactersTable failed:", err)
	}
	if _, err := db.Exec(dropUsersTable); err != nil {
		log.Fatal("db.Exec dropUsersTable failed:", err)
	}
}
