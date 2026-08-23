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
		return fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	return createTables()
}

func createTables() error {
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
	);

	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ticket_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		text TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (ticket_id) REFERENCES tickets(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблиц: %w", err)
	}

	DB.Exec("ALTER TABLE tickets ADD COLUMN admin_response TEXT DEFAULT '';")

	return nil
}
