package mocks

import (
	"github.com/google/uuid"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)

// DB mocks database behavior for testing
type DB struct {
	Users           map[uuid.UUID]models.User
	Passwords       map[uuid.UUID]string
	EmailToID       map[string]uuid.UUID
	NextID          uuid.UUID
	CreateError     error
	FindError       error
	QueryError      error
	ScanError       error
	ShouldDuplicate bool
}

func NewDB() *DB {
	return &DB{
		Users:     make(map[uuid.UUID]models.User),
		Passwords: make(map[uuid.UUID]string),
		EmailToID: make(map[string]uuid.UUID),
		NextID:    uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	}
}
