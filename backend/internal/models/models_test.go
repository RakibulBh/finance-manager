package models

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Tests based on Ruby model specifications from maybe/test/models/

// User Model Tests
func TestUser_Model(t *testing.T) {
	t.Run("should create valid user", func(t *testing.T) {
		userID := uuid.New()
		familyID := uuid.New()

		user := User{
			ID:       userID,
			Email:    "test@example.com",
			FamilyID: familyID,
			Role:     "admin",
		}

		assert.Equal(t, userID, user.ID)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, familyID, user.FamilyID)
		assert.Equal(t, "admin", user.Role)
	})

	t.Run("should have valid UUID", func(t *testing.T) {
		userID := uuid.New()

		user := User{
			ID:       userID,
			Email:    "test@example.com",
			FamilyID: uuid.New(),
		}

		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.NotEqual(t, uuid.Nil, user.FamilyID)
	})

	t.Run("should support different roles", func(t *testing.T) {
		roles := []string{"admin", "member"}

		for _, role := range roles {
			user := User{
				ID:       uuid.New(),
				Email:    "test@example.com",
				FamilyID: uuid.New(),
				Role:     role,
			}
			assert.Equal(t, role, user.Role)
		}
	})
}

// Family Model Tests
func TestFamily_Model(t *testing.T) {
	t.Run("should create valid family", func(t *testing.T) {
		familyID := uuid.New()

		family := Family{
			ID:        familyID,
			Name:      "Test Family",
			Currency:  "USD",
			CreatedAt: time.Now(),
		}

		assert.Equal(t, familyID, family.ID)
		assert.Equal(t, "Test Family", family.Name)
		assert.Equal(t, "USD", family.Currency)
		assert.False(t, family.CreatedAt.IsZero())
	})

	t.Run("should support different currencies", func(t *testing.T) {
		currencies := []string{"USD", "GBP", "EUR", "CAD", "JPY", "AUD"}

		for _, currency := range currencies {
			family := Family{
				ID:        uuid.New(),
				Name:      "Test Family",
				Currency:  currency,
				CreatedAt: time.Now(),
			}
			assert.Equal(t, currency, family.Currency)
		}
	})
}

// Account Model Tests
func TestAccount_Model(t *testing.T) {
	t.Run("should create valid account", func(t *testing.T) {
		accountID := uuid.New()
		familyID := uuid.New()

		account := Account{
			ID:            accountID,
			FamilyID:      familyID,
			Name:          "Test Checking",
			Balance:       1000.50,
			Currency:      "USD",
			Type:          "depository",
			Subtype:       "checking",
			Classification: "asset",
		}

		assert.Equal(t, accountID, account.ID)
		assert.Equal(t, familyID, account.FamilyID)
		assert.Equal(t, "Test Checking", account.Name)
		assert.Equal(t, 1000.50, account.Balance)
		assert.Equal(t, "USD", account.Currency)
		assert.Equal(t, "depository", account.Type)
		assert.Equal(t, "checking", account.Subtype)
		assert.Equal(t, "asset", account.Classification)
	})

	t.Run("should support different account classifications", func(t *testing.T) {
		classifications := []string{"asset", "liability"}

		for _, classification := range classifications {
			account := Account{
				ID:            uuid.New(),
				FamilyID:      uuid.New(),
				Name:          "Test Account",
				Classification: classification,
			}
			assert.Equal(t, classification, account.Classification)
		}
	})

	t.Run("should support property details", func(t *testing.T) {
		account := Account{
			ID:            uuid.New(),
			FamilyID:      uuid.New(),
			Name:          "Family Home",
			Type:          "property",
			Classification: "asset",
			PropertyDetails: &PropertyDetails{
				Address: "123 Main St",
				Sqft:    2000,
			},
		}

		assert.NotNil(t, account.PropertyDetails)
		assert.Equal(t, "123 Main St", account.PropertyDetails.Address)
		assert.Equal(t, 2000, account.PropertyDetails.Sqft)
	})

	t.Run("should support loan details", func(t *testing.T) {
		account := Account{
			ID:            uuid.New(),
			FamilyID:      uuid.New(),
			Name:          "Mortgage",
			Type:          "loan",
			Subtype:       "mortgage",
			Classification: "liability",
			LoanDetails: &LoanDetails{
				InterestRate: 4.5,
				TermMonths:   360,
			},
		}

		assert.NotNil(t, account.LoanDetails)
		assert.Equal(t, 4.5, account.LoanDetails.InterestRate)
		assert.Equal(t, 360, account.LoanDetails.TermMonths)
	})

	t.Run("should handle nil optional details", func(t *testing.T) {
		account := Account{
			ID:            uuid.New(),
			FamilyID:      uuid.New(),
			Name:          "Simple Account",
			Classification: "asset",
		}

		assert.Nil(t, account.PropertyDetails)
		assert.Nil(t, account.LoanDetails)
	})
}

// Entry Model Tests
func TestEntry_Model(t *testing.T) {
	t.Run("should create valid entry", func(t *testing.T) {
		entryID := uuid.New()
		accountID := uuid.New()
		entryableID := uuid.New()

		entry := Entry{
			ID:            entryID,
			AccountID:     accountID,
			Amount:        25.50,
			Currency:      "USD",
			Date:          time.Now(),
			Name:          "Coffee Purchase",
			EntryableType: "Transaction",
			EntryableID:   entryableID,
		}

		assert.Equal(t, entryID, entry.ID)
		assert.Equal(t, accountID, entry.AccountID)
		assert.Equal(t, 25.50, entry.Amount)
		assert.Equal(t, "USD", entry.Currency)
		assert.False(t, entry.Date.IsZero())
		assert.Equal(t, "Coffee Purchase", entry.Name)
		assert.Equal(t, "Transaction", entry.EntryableType)
		assert.Equal(t, entryableID, entry.EntryableID)
	})

	t.Run("should support different entryable types", func(t *testing.T) {
		entryableTypes := []string{"Transaction", "Valuation"}

		for _, entryableType := range entryableTypes {
			entry := Entry{
				ID:            uuid.New(),
				AccountID:     uuid.New(),
				Amount:        100,
				Currency:      "USD",
				Date:          time.Now(),
				Name:          "Test Entry",
				EntryableType: entryableType,
				EntryableID:   uuid.New(),
			}
			assert.Equal(t, entryableType, entry.EntryableType)
		}
	})

	t.Run("should support negative amounts", func(t *testing.T) {
		entry := Entry{
			ID:            uuid.New(),
			AccountID:     uuid.New(),
			Amount:        -25.50,
			Currency:      "USD",
			Date:          time.Now(),
			Name:          "Expense",
			EntryableType: "Transaction",
			EntryableID:   uuid.New(),
		}

		assert.Less(t, entry.Amount, 0.0)
	})
}

// Transaction Model Tests
func TestTransaction_Model(t *testing.T) {
	t.Run("should create valid transaction", func(t *testing.T) {
		transactionID := uuid.New()
		categoryID := uuid.New()
		merchantID := uuid.New()

		transaction := Transaction{
			ID:         transactionID,
			CategoryID: &categoryID,
			MerchantID: &merchantID,
			Kind:       "standard",
		}

		assert.Equal(t, transactionID, transaction.ID)
		assert.NotNil(t, transaction.CategoryID)
		assert.NotNil(t, transaction.MerchantID)
		assert.Equal(t, categoryID, *transaction.CategoryID)
		assert.Equal(t, merchantID, *transaction.MerchantID)
		assert.Equal(t, "standard", transaction.Kind)
	})

	t.Run("should support nil category", func(t *testing.T) {
		merchantID := uuid.New()

		transaction := Transaction{
			ID:         uuid.New(),
			CategoryID: nil,
			MerchantID: &merchantID,
			Kind:       "standard",
		}

		assert.Nil(t, transaction.CategoryID)
		assert.NotNil(t, transaction.MerchantID)
	})

	t.Run("should support nil merchant", func(t *testing.T) {
		categoryID := uuid.New()

		transaction := Transaction{
			ID:         uuid.New(),
			CategoryID: &categoryID,
			MerchantID: nil,
			Kind:       "standard",
		}

		assert.NotNil(t, transaction.CategoryID)
		assert.Nil(t, transaction.MerchantID)
	})

	t.Run("should support different transaction kinds", func(t *testing.T) {
		kinds := []string{"standard", "transfer", "reconciliation"}

		for _, kind := range kinds {
			transaction := Transaction{
				ID:   uuid.New(),
				Kind: kind,
			}
			assert.Equal(t, kind, transaction.Kind)
		}
	})
}

// TransactionDetail Model Tests
func TestTransactionDetail_Model(t *testing.T) {
	t.Run("should create valid transaction detail", func(t *testing.T) {
		detail := TransactionDetail{
			Entry: Entry{
				ID:        uuid.New(),
				AccountID: uuid.New(),
				Amount:    25.50,
				Currency:  "USD",
				Date:      time.Now(),
				Name:      "Coffee",
			},
			CategoryName: "Food & Drink",
			MerchantName: "Starbucks",
			Kind:         "standard",
		}

		assert.Equal(t, "Coffee", detail.Name)
		assert.Equal(t, 25.50, detail.Amount)
		assert.Equal(t, "Food & Drink", detail.CategoryName)
		assert.Equal(t, "Starbucks", detail.MerchantName)
		assert.Equal(t, "standard", detail.Kind)
	})

	t.Run("should embed Entry fields", func(t *testing.T) {
		entry := Entry{
			ID:        uuid.New(),
			AccountID: uuid.New(),
			Amount:    100,
			Currency:  "USD",
			Date:      time.Now(),
			Name:      "Test",
		}

		detail := TransactionDetail{
			Entry:         entry,
			CategoryName:  "Category",
			MerchantName:  "Merchant",
			Kind:          "standard",
		}

		assert.Equal(t, entry.ID, detail.ID)
		assert.Equal(t, entry.AccountID, detail.AccountID)
		assert.Equal(t, entry.Amount, detail.Amount)
		assert.Equal(t, entry.Currency, detail.Currency)
		assert.Equal(t, entry.Date, detail.Date)
		assert.Equal(t, entry.Name, detail.Name)
	})
}

// PropertyDetails Model Tests
func TestPropertyDetails_Model(t *testing.T) {
	t.Run("should create valid property details", func(t *testing.T) {
		details := PropertyDetails{
			Address: "123 Main St, Anytown, USA",
			Sqft:    2000,
		}

		assert.Equal(t, "123 Main St, Anytown, USA", details.Address)
		assert.Equal(t, 2000, details.Sqft)
	})
}

// LoanDetails Model Tests
func TestLoanDetails_Model(t *testing.T) {
	t.Run("should create valid loan details", func(t *testing.T) {
		details := LoanDetails{
			InterestRate: 4.5,
			TermMonths:   360,
		}

		assert.Equal(t, 4.5, details.InterestRate)
		assert.Equal(t, 360, details.TermMonths)
	})
}

// JSON Serialization Tests
func TestModel_JSONTags(t *testing.T) {
	t.Run("user should serialize correctly", func(t *testing.T) {
		user := User{
			ID:             uuid.New(),
			Email:          "test@example.com",
			FamilyID:       uuid.New(),
			Role:           "admin",
			PasswordDigest: "secret",
		}

		// Password digest should not be in JSON (tag is "-")
		assert.Equal(t, "id", getJSONTag(&user, "ID"))
		assert.Equal(t, "email", getJSONTag(&user, "Email"))
		assert.Equal(t, "familyId", getJSONTag(&user, "FamilyID"))
		assert.Equal(t, "role", getJSONTag(&user, "Role"))
	})

	t.Run("account should serialize correctly", func(t *testing.T) {
		account := Account{
			ID:            uuid.New(),
			FamilyID:      uuid.New(),
			Name:          "Test",
			Type:          "depository",
			Subtype:       "checking",
			Classification: "asset",
			Balance:       100,
			Currency:      "USD",
		}

		assert.Equal(t, "id", getJSONTag(&account, "ID"))
		assert.Equal(t, "familyId", getJSONTag(&account, "FamilyID"))
		assert.Equal(t, "type", getJSONTag(&account, "Type"))
		assert.Equal(t, "subtype", getJSONTag(&account, "Subtype"))
		assert.Equal(t, "classification", getJSONTag(&account, "Classification"))
	})
}

// Helper function to check JSON tags
func getJSONTag(v interface{}, fieldName string) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	f, ok := t.FieldByName(fieldName)
	if !ok {
		return ""
	}
	tag := f.Tag.Get("json")
	if tag == "" {
		return ""
	}
	parts := strings.Split(tag, ",")
	return parts[0]
}


// Currency Support Tests
func TestModel_CurrencySupport(t *testing.T) {
	t.Run("should support common currencies", func(t *testing.T) {
		currencies := []string{
			"USD", // US Dollar
			"GBP", // British Pound
			"EUR", // Euro
			"JPY", // Japanese Yen
			"CAD", // Canadian Dollar
			"AUD", // Australian Dollar
			"CHF", // Swiss Franc
			"CNY", // Chinese Yuan
		}

		for _, currency := range currencies {
			entry := Entry{
				ID:       uuid.New(),
				Amount:   100,
				Currency: currency,
			}
			assert.Equal(t, currency, entry.Currency)
		}
	})
}

// UUID Validation Tests
func TestModel_UUIDValidation(t *testing.T) {
	t.Run("should have valid UUIDs", func(t *testing.T) {
		user := User{
			ID:       uuid.New(),
			Email:    "test@example.com",
			FamilyID: uuid.New(),
		}

		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.NotEqual(t, uuid.Nil, user.FamilyID)
		assert.True(t, user.ID.Version() == uuid.Version(4))
	})
}

// Time Zone Tests
func TestModel_TimeHandling(t *testing.T) {
	t.Run("should handle time correctly", func(t *testing.T) {
		now := time.Now().UTC()

		family := Family{
			ID:        uuid.New(),
			Name:      "Test",
			Currency:  "USD",
			CreatedAt: now,
		}

		assert.Equal(t, now, family.CreatedAt)
	})

	t.Run("should handle entry dates", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

		entry := Entry{
			ID:       uuid.New(),
			Amount:   100,
			Currency: "USD",
			Date:     date,
		}

		assert.Equal(t, 2024, entry.Date.Year())
		assert.Equal(t, time.January, entry.Date.Month())
		assert.Equal(t, 15, entry.Date.Day())
	})
}

// Edge Cases Tests
func TestModel_EdgeCases(t *testing.T) {
	t.Run("should handle zero balance", func(t *testing.T) {
		account := Account{
			ID:       uuid.New(),
			Name:     "Empty Account",
			Balance:  0,
			Currency: "USD",
		}

		assert.Equal(t, 0.0, account.Balance)
	})

	t.Run("should handle negative balances", func(t *testing.T) {
		account := Account{
			ID:       uuid.New(),
			Name:     "Credit Card",
			Balance:  -500.50,
			Currency: "USD",
		}

		assert.Less(t, account.Balance, 0.0)
	})

	t.Run("should handle very large amounts", func(t *testing.T) {
		account := Account{
			ID:       uuid.New(),
			Name:     "Property",
			Balance:  1000000000,
			Currency: "USD",
		}

		assert.Greater(t, account.Balance, 999999999.0)
	})

	t.Run("should handle very small amounts", func(t *testing.T) {
		entry := Entry{
			ID:       uuid.New(),
			Amount:   0.01,
			Currency: "USD",
		}

		assert.Equal(t, 0.01, entry.Amount)
	})
}
