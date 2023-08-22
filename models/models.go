package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type News struct {
	ID       int
	Title    string
	ImageURL string
	Link     string
	Source   string
}

func createTableIfNotExists(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS news (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT,
            image_url TEXT,
            link TEXT,
            source TEXT
        )
    `)
	return err
}
