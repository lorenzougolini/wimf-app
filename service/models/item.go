package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Item struct {
	Id             uuid.UUID
	Barcode        string
	Name           string
	Brand          string
	Quantity       int
	ExpirationDate time.Time
	AdditionDate   time.Time
}

type HomeItems struct {
	RecentItems   []Item
	ExpiringItems []Item
}
