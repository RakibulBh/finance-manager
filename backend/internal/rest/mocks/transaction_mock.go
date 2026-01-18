package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)

// TransactionStore is a mock implementation of TransactionStore for testing
type TransactionStore struct {
	Merchants     map[string]uuid.UUID
	Transactions  []models.Entry
	CreateError   error
	MerchantError error
	TransferError error
}

func NewTransactionStore() *TransactionStore {
	return &TransactionStore{
		Merchants:    make(map[string]uuid.UUID),
		Transactions: []models.Entry{},
	}
}

func (m *TransactionStore) GetOrCreateMerchant(ctx context.Context, name string, familyID uuid.UUID) (uuid.UUID, error) {
	if m.MerchantError != nil {
		return uuid.Nil, m.MerchantError
	}

	if id, exists := m.Merchants[name]; exists {
		return id, nil
	}

	id := uuid.New()
	m.Merchants[name] = id
	return id, nil
}

func (m *TransactionStore) CreateTransaction(ctx context.Context, entry *models.Entry, txDetail *models.Transaction) error {
	if m.CreateError != nil {
		return m.CreateError
	}

	entry.ID = uuid.New()
	m.Transactions = append(m.Transactions, *entry)
	return nil
}

func (m *TransactionStore) CreateTransfer(ctx context.Context, fromEntry, toEntry *models.Entry) error {
	if m.TransferError != nil {
		return m.TransferError
	}

	fromEntry.ID = uuid.New()
	toEntry.ID = uuid.New()
	m.Transactions = append(m.Transactions, *fromEntry, *toEntry)
	return nil
}
