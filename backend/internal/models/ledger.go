package models

import (
	"time"

	"github.com/google/uuid"
)

// Entry represents the 'entries' table - the core ledger
type Entry struct {
	ID            uuid.UUID `db:"id" json:"id"`
	AccountID     uuid.UUID `db:"account_id" json:"accountId"`
	Amount        float64   `db:"amount" json:"amount"`
	Currency      string    `db:"currency" json:"currency"`
	Date          time.Time `db:"date" json:"date"`
	Name          string    `db:"name" json:"name"` // Description: "Starbucks"

	// Polymorphic Fields
	EntryableType string    `db:"entryable_type" json:"entryableType"` // "Transaction", "Valuation"
	EntryableID   uuid.UUID `db:"entryable_id" json:"entryableId"`
}

// Transaction represents the 'transactions' table - specific info
type Transaction struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	CategoryID *uuid.UUID `db:"category_id" json:"categoryId"`
	MerchantID *uuid.UUID `db:"merchant_id" json:"merchantId"`
	Kind       string     `db:"kind" json:"kind"` // "standard", "transfer"
}

// API Response model: Combining them into one usable struct
type TransactionDetail struct {
	Entry           // Embed the Entry fields
	CategoryName string `json:"categoryName"`
	MerchantName string `json:"merchantName"`
	Kind         string `json:"kind"`
}
