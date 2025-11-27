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

func (db *appdbimpl) AddItem(product models.ProductInfo, expiration time.Time, addition time.Time) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	newId := id.String()

	query := `
		INSERT INTO items (id, barcode, name, brand, quantity, expiration_date, added_at)
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	_, err = db.c.Exec(query,
		newId,
		product.Barcode,
		product.Name,
		product.Brand,
		1,
		expiration.Format(models.DbTimeLayout),
		addition.Format(models.DbTimeLayout),
	)
	if err != nil {
		return fmt.Errorf("error inserting item %s: %w", product.Barcode, err)
	}
	return nil
}

func (db *appdbimpl) GetItemsByBarcode(barcode string) (bool, []models.Item, error) {
	var items []models.Item

	query := `
		SELECT id, barcode, name, brand, quantity, expiration_date, added_at
		FROM items
		WHERE barcode=?
		ORDER BY expiration_date ASC;
	`

	rows, err := db.c.Query(query, barcode)
	if err != nil {
		return false, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Item
		var exp, add sql.NullString
		if err := rows.Scan(
			&i.Id,
			&i.Barcode,
			&i.Name,
			&i.Brand,
			&i.Quantity,
			&exp,
			&add,
		); err != nil {
			return false, []models.Item{}, nil
		}

		if exp.Valid {
			i.ExpirationDate, _ = time.Parse(models.DbTimeLayout, exp.String)
		}
		if add.Valid {
			i.AdditionDate, _ = time.Parse(models.DbTimeLayout, add.String)
		}

		items = append(items, i)
	}
	if err = rows.Err(); err != nil {
		return false, nil, err
	}
	return len(items) > 0, items, nil
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

func (db *appdbimpl) GetItemById(id string) (models.Item, error) {
	var item models.Item
	var exp, add sql.NullString

	query := `
		SELECT id, barcode, name, brand, quantity, expiration_date, added_at
		FROM items
		WHERE id=?;
	`

	err := db.c.QueryRow(query, id).Scan(
		&item.Id,
		&item.Barcode,
		&item.Name,
		&item.Brand,
		&item.Quantity,
		&exp,
		&add,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Item{}, fmt.Errorf("item with id %s not found", id)
		}
		return models.Item{}, err
	}

	if exp.Valid {
		item.ExpirationDate, _ = time.Parse(models.DbTimeLayout, exp.String)
	}
	if add.Valid {
		item.AdditionDate, _ = time.Parse(models.DbTimeLayout, add.String)
	}

	return item, nil
}

func (db *appdbimpl) IncreaseItemQuantity(barcode string, quantity int) error {
	res, err := db.c.Exec(`
		UPDATE items 
		SET quantity = quantity + ? 
		WHERE id = ?;`,
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

func (db *appdbimpl) DeleteItem(id string) error {
	_, err := db.c.Exec("DELETE FROM items WHERE id=?;", id)
	return err
}

func (db *appdbimpl) UpdateItem(id string, name string, brand string, date time.Time) error {
	query := "UPDATE items SET name=?, brand=?, expiration_date=? WHERE id=?;"
	_, err := db.c.Exec(query, name, brand, date.Format(models.DbTimeLayout), id)
	return err
}
