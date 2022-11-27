package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const dbFileName string = "db.sqlite3"

func GetDB() (dbSession *sql.DB) {
	dbSession, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	createTables(dbSession)
	return
}

func createTables(dbSession *sql.DB) (err error) {
	stmt, err := dbSession.Prepare(`
		CREATE TABLE IF NOT EXISTS agents (id TEXT NOT NULL UNIQUE, created_at TEXT NOT NULL, updated_at TEXT)
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		log.Fatal(err)
	}

	stmt, err = dbSession.Prepare(`
		CREATE TABLE IF NOT EXISTS 
			commands (id TEXT NOT NULL UNIQUE, agent_id TEXT, command TEXT, result TEXT, created_at TEXT NOT NULL, updated_at TEXT)
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		log.Fatal(err)
	}
	return
}
