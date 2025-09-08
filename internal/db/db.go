package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	// Create complete schema
	schema := `
	CREATE TABLE IF NOT EXISTS vaults (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS prompts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		vault_id INTEGER,
		name TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(vault_id) REFERENCES vaults(id),
		UNIQUE(vault_id, name)
	);

	CREATE TABLE IF NOT EXISTS prompt_versions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		prompt_id INTEGER,
		version INTEGER,
		content TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(prompt_id) REFERENCES prompts(id)
	);

	CREATE TABLE IF NOT EXISTS runs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		prompt_version_id INTEGER,
		provider TEXT,
		params TEXT,
		response TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(prompt_version_id) REFERENCES prompt_versions(id)
	);
	`
	_, err = DB.Exec(schema)
	return err
}
