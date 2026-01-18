package models

import "github.com/google/uuid"

type Account struct {
	ID              uuid.UUID `db:"id" json:"id"`
	FamilyID        uuid.UUID `db:"family_id" json:"familyId"`
	Name            string    `db:"name" json:"name"`
	Type            string    `db:"subtype" json:"type"` // "checking", "savings", "credit_card"
	Balance         float64   `db:"balance" json:"balance"`
	Currency        string    `db:"currency" json:"currency"`
	Status          string    `db:"status" json:"status"` // "active", "archived"

	// Linking to Plaid
	PlaidAccountID *uuid.UUID `db:"plaid_account_id" json:"plaidAccountId,omitempty"`
}
