package models

import (
	"time"

	"github.com/google/uuid"
)

type Family struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Currency  string    `db:"currency" json:"currency"` // Default currency (USD, GBP)
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type User struct {
	ID             uuid.UUID `db:"id" json:"id"`
	FamilyID       uuid.UUID `db:"family_id" json:"familyId"`
	Email          string    `db:"email" json:"email"`
	PasswordDigest string    `db:"password_digest" json:"-"` // Never return password in JSON
	Role           string    `db:"role" json:"role"`         // "admin" or "member"
}
