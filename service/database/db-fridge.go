package database

import (
	"database/sql"
	"time"

	"github.com/lorenzougolini/wimf-app/service/models"
)

func (db *appdbimpl) GetFridge() ([]models.Item, error) {
	query := `
		SELECT barcode, name, brand, SUM(quantity) as tot_quantity, MIN(expiration_date) as next_exp, MAX(added_at) as latest_add
		FROM items
		GROUP BY barcode
		ORDER BY next_exp ASC
	`
	rows, err := db.c.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Item
	for rows.Next() {
		var i models.Item
		var nextExp, latestAdd sql.NullString
		if err := rows.Scan(&i.Barcode, &i.Name, &i.Brand, &i.Quantity, &nextExp, &latestAdd); err != nil {
			return nil, err
		}

		if nextExp.Valid {
			i.ExpirationDate, _ = time.Parse(models.DbTimeLayout, nextExp.String)
		}
		if latestAdd.Valid {
			i.AdditionDate, _ = time.Parse(models.DbTimeLayout, latestAdd.String)
		}

		result = append(result, i)
	}
	return result, rows.Err()
}
