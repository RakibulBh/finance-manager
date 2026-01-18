package models

import (
	"time"

	"github.com/google/uuid"
)

type PlaidItem struct {
	ID              uuid.UUID `json:"id"`
	FamilyID        uuid.UUID `json:"familyId"`
	AccessToken     string    `json:"-"` // Encrypted
	ItemID          string    `json:"itemId"`
	InstitutionID   string    `json:"institutionId"`
	InstitutionName string    `json:"institutionName"`
	SyncCursor      string    `json:"-"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
