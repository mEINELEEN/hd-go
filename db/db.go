package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite", "./helpdesk.db")
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user'
	);

	CREATE TABLE IF NOT EXISTS tickets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'new',
		admin_response TEXT DEFAULT '',
		user_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблиц: %w", err)
	}

	return nil
}
