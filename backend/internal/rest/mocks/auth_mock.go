package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// UserStore is a mock implementation of UserStore for testing
type UserStore struct {
	Users       map[string]*models.User
	Passwords   map[string]string
	CreateError error
	FindError   error
}

func NewUserStore() *UserStore {
	return &UserStore{
		Users:     make(map[string]*models.User),
		Passwords: make(map[string]string),
	}
}

func (m *UserStore) CreateFamily(ctx context.Context, name string) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (m *UserStore) CreateUser(ctx context.Context, email, password string, familyID uuid.UUID) (*models.User, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}

	// Check for duplicate
	if _, exists := m.Users[email]; exists {
		return nil, &DuplicateKeyError{Email: email}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       uuid.New(),
		Email:    email,
		FamilyID: familyID,
		Role:     "admin",
	}

	m.Users[email] = user
	m.Passwords[email] = string(hashedPassword)

	return user, nil
}

func (m *UserStore) FindByEmail(ctx context.Context, email string) (*models.User, string, error) {
	if m.FindError != nil {
		return nil, "", m.FindError
	}

	user, exists := m.Users[email]
	if !exists {
		return nil, "", pgx.ErrNoRows
	}

	password := m.Passwords[email]
	return user, password, nil
}

func (m *UserStore) AddUser(email, hashedPassword string, familyID uuid.UUID) {
	user := &models.User{
		ID:       uuid.New(),
		Email:    email,
		FamilyID: familyID,
		Role:     "admin",
	}
	m.Users[email] = user
	m.Passwords[email] = hashedPassword
}

type DuplicateKeyError struct {
	Email string
}

func (e *DuplicateKeyError) Error() string {
	return "duplicate key value violates unique constraint"
}
