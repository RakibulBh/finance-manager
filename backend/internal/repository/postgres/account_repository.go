package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)


type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, acc *models.Account) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert Account
	queryAcc := `
		INSERT INTO accounts (family_id, name, balance, currency, subtype, classification)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err = tx.QueryRow(ctx, queryAcc,
		acc.FamilyID, acc.Name, acc.Balance, acc.Currency, acc.Subtype, acc.Classification,
	).Scan(&acc.ID)
	if err != nil {
		return err
	}

	// 2. If initial balance > 0, create a Valuation Entry
	if acc.Balance != 0 {
		// Valuation entry id
		valuationID := uuid.New()

		// For simplicity, we just insert the valuation ID into the transactions table
		// if we are following the structure where valuation is a type of entryable.
		// However, for now, let's just insert into entries directly with a generated ID
		// as 'Valuation' type.

		queryEntry := `
			INSERT INTO entries (account_id, amount, date, currency, name, entryable_type, entryable_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = tx.Exec(ctx, queryEntry,
			acc.ID, acc.Balance, time.Now(), acc.Currency, "Initial Balance", "Valuation", valuationID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}


func (r *AccountRepository) ListByFamilyID(ctx context.Context, familyID uuid.UUID) ([]models.Account, error) {
	query := `
		SELECT id, family_id, name, balance, currency, subtype, classification
		FROM accounts
		WHERE family_id = $1 AND status = 'active'
	`
	rows, err := r.db.Query(ctx, query, familyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var acc models.Account
		err := rows.Scan(
			&acc.ID, &acc.FamilyID, &acc.Name, &acc.Balance, &acc.Currency, &acc.Subtype, &acc.Classification,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func (r *AccountRepository) GetNetWorth(ctx context.Context, familyID uuid.UUID) (float64, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN classification = 'asset' THEN balance ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN classification = 'liability' THEN balance ELSE 0 END), 0)
		FROM accounts
		WHERE family_id = $1 AND status = 'active'
	`
	var netWorth float64
	err := r.db.QueryRow(ctx, query, familyID).Scan(&netWorth)
	return netWorth, err
}
