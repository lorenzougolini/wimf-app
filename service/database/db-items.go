package database

import (
	"fmt"

	"github.com/lorenzougolini/wimf-app/service/models"
)

func (db *appdbimpl) AddItem(itemid string) error {
	_, err := db.c.Exec("INSERT INTO items (id) VALUES (?);", itemid)
	if err != nil {
		return fmt.Errorf("error in item addition err: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetLastItems(limit int) ([]models.Item, error) {
	rows, err := db.c.Query(`
		SELECT id FROM items ORDER BY rowid DESC LIMIT ?`,
		limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Item
	for rows.Next() {
		var it models.Item
		err = rows.Scan(&it.ItemID)
		if err != nil {
			return nil, err
		}
		result = append(result, it)
	}
	return result, nil
}

