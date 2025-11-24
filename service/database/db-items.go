package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lorenzougolini/wimf-app/service/models"
)

func (db *appdbimpl) CheckIdExistence(barcode string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM items WHERE barcode=?)", barcode).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("Check existence error: %w", err)
	}
	return exists, nil
}

func (db *appdbimpl) AddItem(product models.ProductInfo) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	newId := id.String()
	now := time.Now()
	exp := now.AddDate(0, 0, 14)
	query := `
		INSERT INTO items (id, barcode, name, brand, quantity, expiration_date, added_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.c.Exec(query,
		newId,
		product.Barcode,
		product.Name,
		product.Brand,
		1,
		exp.Format(models.DbTimeLayout),
		now.Format(models.DbTimeLayout),
	)
	if err != nil {
		return fmt.Errorf("error inserting item %s: %w", product.Barcode, err)
	}
	return nil
}

func (db *appdbimpl) GetItemByBarcode(barcode string) (bool, models.Item, error) {
	var item models.Item
	var expiration sql.NullString
	var added string

	err := db.c.QueryRow(`
		SELECT barcode, name, brand, quantity, expiration_date, added_at
		FROM items
		WHERE barcode=?
		LIMIT 1;`,
		barcode).Scan(
		&item.Barcode,
		&item.Name,
		&item.Brand,
		&item.Quantity,
		&expiration,
		&added,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, models.Item{}, nil
		}
		return false, models.Item{}, err
	}

	parsedAdd, err := time.Parse(models.DbTimeLayout, added)
	if err == nil {
		item.AdditionDate = parsedAdd
	}

	if expiration.Valid {
		parsedExp, err := time.Parse(models.DbTimeLayout, expiration.String)
		if err == nil {
			item.ExpirationDate = parsedExp
		}
	}
	return true, item, nil
}

func (db *appdbimpl) GetNItemsBy(limit int, orderBy string) ([]models.Item, error) {
	var orderByClause string
	var latestAdd string
	var nextExpirationDate string

	switch orderBy {
	case "latest":
		orderByClause = "latest_date DESC"
	case "expiring":
		orderByClause = "next_expiration_date ASC"
	default:
		return nil, fmt.Errorf("unsupported sort mode: %s", orderBy)
	}

	query := fmt.Sprintf(`
		SELECT barcode, name, brand, SUM(quantity) as tot_quantity, MIN(expiration_date) as next_expiration_date, MAX(added_at) as latest_date
		FROM items
		GROUP BY barcode
		ORDER BY %s
		LIMIT ?;`,
		orderByClause)
	rows, err := db.c.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.Item, 0, limit)
	for rows.Next() {
		var item models.Item
		err = rows.Scan(
			&item.Barcode,
			&item.Name,
			&item.Brand,
			&item.Quantity,
			&nextExpirationDate,
			&latestAdd,
		)
		if err != nil {
			return nil, err
		}

		parsedAdd, err := time.Parse(models.DbTimeLayout, latestAdd)
		if err == nil {
			item.AdditionDate = parsedAdd
		} else {
			fmt.Printf("ERROR parsing AdditionDate: '%s' with layout '%s' -> %v\n", latestAdd, models.DbTimeLayout, err)
		}
		parsedExp, err := time.Parse(models.DbTimeLayout, nextExpirationDate)
		if err == nil {
			item.ExpirationDate = parsedExp
		} else {
			fmt.Printf("ERROR parsing AdditionDate: '%s' with layout '%s' -> %v\n", latestAdd, models.DbTimeLayout, err)
		}

		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (db *appdbimpl) IncreaseItemQuantity(barcode string, quantity int) error {
	res, err := db.c.Exec(`
		UPDATE items 
		SET quantity = quantity + ? 
		WHERE id = ?`,
		quantity, barcode)
	if err != nil {
		return fmt.Errorf("error updating quantity: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil || affected == 0 {
		return fmt.Errorf("error checking affected rows: %w", err)
	}

	return nil
}
