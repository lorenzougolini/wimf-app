package store

import (
	"database/sql"
	"errors"
	"fmt"
)

type AppDatabase interface {
	AddItem(itemId string, itemName string) error
}

type appdb struct {
	db *sql.DB
}

func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	var tableName string
	err := db.QueryRow(
		`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)

	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE IF NOT EXISTS items (
			itemId TEXT NOT NULL PRIMARY KEY
			itemName TEXT NOT NULL UNIQUE
			);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating db table: %w", err)
		}
	}

	return &appdb{
		db: db,
	}, nil
}
