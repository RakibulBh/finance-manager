package models

import "github.com/google/uuid"

type Account struct {
    ID             uuid.UUID `json:"id"`
    FamilyID       uuid.UUID `json:"familyId"`
    Name           string    `json:"name"`
    Type           string    `json:"type"`           // "depository", "loan", "property", "credit_card"
    Subtype        string    `json:"subtype"`        // "checking", "mortgage"
    Classification string    `json:"classification"` // "asset", "liability"
    Balance        float64   `json:"balance"`        // Raw number
    Currency       string    `json:"currency"`

    // Detailed Attributes (Optional, only filled if applicable)
    PropertyDetails *PropertyDetails `json:"propertyDetails,omitempty"`
    LoanDetails     *LoanDetails     `json:"loanDetails,omitempty"`
}

type PropertyDetails struct {
    Address string `json:"address"`
    Sqft    int    `json:"sqft"`
}

type LoanDetails struct {
    InterestRate float64 `json:"interestRate"`
    TermMonths   int     `json:"termMonths"`
}
