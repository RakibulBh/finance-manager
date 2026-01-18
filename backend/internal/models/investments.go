package models

import (
	"time"

	"github.com/google/uuid"
)

type Security struct {
	ID          uuid.UUID `json:"id"`
	Ticker      string    `json:"ticker"`
	Name        string    `json:"name"`
	LatestPrice float64   `json:"latestPrice"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type Trade struct {
	ID         uuid.UUID `json:"id"`
	SecurityID uuid.UUID `json:"securityId"`
	Qty        float64   `json:"qty"`
	Price      float64   `json:"price"`
	Kind       string    `json:"kind"` // "buy", "sell"
}
