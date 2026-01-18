package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)

// AccountStore is a mock implementation of AccountStore for testing
type AccountStore struct {
	Accounts      map[uuid.UUID][]models.Account
	NetWorth      map[uuid.UUID]float64
	CreateError   error
	ListError     error
	NetWorthError error
}

func NewAccountStore() *AccountStore {
	return &AccountStore{
		Accounts: make(map[uuid.UUID][]models.Account),
		NetWorth: make(map[uuid.UUID]float64),
	}
}

func (m *AccountStore) Create(ctx context.Context, acc *models.Account) error {
	if m.CreateError != nil {
		return m.CreateError
	}

	acc.ID = uuid.New()
	if m.Accounts[acc.FamilyID] == nil {
		m.Accounts[acc.FamilyID] = []models.Account{}
	}
	m.Accounts[acc.FamilyID] = append(m.Accounts[acc.FamilyID], *acc)

	// Update net worth
	if acc.Classification == "asset" {
		m.NetWorth[acc.FamilyID] += acc.Balance
	} else {
		m.NetWorth[acc.FamilyID] -= acc.Balance
	}

	return nil
}

func (m *AccountStore) ListByFamilyID(ctx context.Context, familyID uuid.UUID) ([]models.Account, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}

	return m.Accounts[familyID], nil
}

func (m *AccountStore) GetNetWorth(ctx context.Context, familyID uuid.UUID) (float64, error) {
	if m.NetWorthError != nil {
		return 0, m.NetWorthError
	}

	return m.NetWorth[familyID], nil
}

func (m *AccountStore) AddAccount(familyID uuid.UUID, acc models.Account) {
	if m.Accounts[familyID] == nil {
		m.Accounts[familyID] = []models.Account{}
	}
	m.Accounts[familyID] = append(m.Accounts[familyID], acc)

	if acc.Classification == "asset" {
		m.NetWorth[familyID] += acc.Balance
	} else {
		m.NetWorth[familyID] -= acc.Balance
	}
}
