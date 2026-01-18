package models

import "github.com/google/uuid"

type Category struct {
	ID       uuid.UUID `json:"id"`
	FamilyID uuid.UUID `json:"familyId"`
	Name     string    `json:"name"`
	Color    string    `json:"color"`
}

type Merchant struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
