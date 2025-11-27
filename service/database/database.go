/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lorenzougolini/wimf-app/service/models"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CheckIdExistence(barcode string) (bool, error)
	AddItem(productInfo models.ProductInfo, expiration time.Time, addition time.Time) error
	GetItemsByBarcode(barcode string) (bool, []models.Item, error)
	GetItemById(id string) (models.Item, error)
	GetNItemsBy(limit int, orderBy string) ([]models.Item, error)

	GetFridge() ([]models.Item, error)

	DeleteItem(id string) error
	UpdateItem(id string, name string, brand string, date time.Time) error

	IncreaseItemQuantity(barcode string, quantity int) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='items';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `
			CREATE TABLE IF NOT EXISTS items (
				id TEXT NOT NULL PRIMARY KEY,
				barcode TEXT NOT NULL,
				name TEXT NOT NULL,
				brand TEXT NOT NULL,
				quantity INTEGER NOT NULL DEFAULT 1,
				expiration_date TEXT,
				added_at TEXT DEFAULT CURRENT_TIMESTAMP
			);`
		// maybe add unit TEXT field

		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database items structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
