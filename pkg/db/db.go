package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const achema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);`

func Init(dbFile string) error {
	var install bool

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть БД: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("БД недоступна: %w", err)
	}

	if install {
		_, err = DB.Exec(achema)
		if err != nil {
			return fmt.Errorf("ошибка при создании БД: %w", err)
		}
	}
	return nil
}
