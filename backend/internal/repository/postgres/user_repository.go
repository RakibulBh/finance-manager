package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
    db *pgxpool.Pool
}


func NewUserRepository(db *pgxpool.Pool) *UserRepository {
    return &UserRepository{db: db}
}

// CreateUser creates a user and returns the created object
func (r *UserRepository) CreateUser(ctx context.Context, email, password string, familyID uuid.UUID) (*models.User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    var user models.User
    query := `INSERT INTO users (email, password_digest, family_id) VALUES ($1, $2, $3) RETURNING id, email, family_id`
    err = r.db.QueryRow(ctx, query, email, string(hashedPassword), familyID).Scan(&user.ID, &user.Email, &user.FamilyID)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// FindByEmail finds a user by email (used for login)
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, string, error) {
    var user models.User
    var passwordHash string
    query := `SELECT id, email, family_id, password_digest FROM users WHERE email = $1`
    err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.FamilyID, &passwordHash)
    if err != nil {
        return nil, "", err
    }
    return &user, passwordHash, nil
}
