package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
	"bale-moderator-bot/config"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite", config.DBFile)
	if err != nil {
		log.Fatal(err)
	}

	_, _ = DB.Exec(`CREATE TABLE IF NOT EXISTS muted (
		user_id INTEGER PRIMARY KEY,
		until INTEGER
	);`)
}
