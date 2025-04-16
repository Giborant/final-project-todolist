package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var (
	db *sql.DB
)

const (
	defaultPath = "./pkg/db"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT DEFAULT "",
    repeat VARCHAR(50) DEFAULT ""
);

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func getPath(dbFile string) string {
	pathStr := os.Getenv("TODO_DBFILE")
	if pathStr == "" {
		return filepath.Join(defaultPath, dbFile)
	}

	if filepath.Ext(pathStr) != "" {
		return pathStr
	}

	return filepath.Join(pathStr, dbFile)
}

func Init(dbFile string) error {
	dbPath := getPath(dbFile)

	_, err := os.Stat(dbPath)
	install := os.IsNotExist(err)

	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	if install {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания схемы БД: %w", err)
		}
		log.Printf("Создана новая БД в %s\n", dbPath)
	} else {
		log.Printf("Используется существующая БД из %s\n", dbPath)
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
