package models

import "time"

type Item struct {
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
