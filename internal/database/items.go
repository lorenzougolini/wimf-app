package store

import (
	"fmt"
)

func (db *appdb) AddItem(itemId string, itemName string) error {
	_, err := db.db.Exec(`insert into items (itemId, itemName) values (?, ?);`, itemId, itemName)
	if err != nil {
		return fmt.Errorf("error in item addition: %w", err)
	}
	return nil
}
