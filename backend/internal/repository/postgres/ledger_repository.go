package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)

type LedgerRepository struct {
	db *pgxpool.Pool
}

func NewLedgerRepository(db *pgxpool.Pool) *LedgerRepository {
	return &LedgerRepository{db: db}
}

func (r *LedgerRepository) GetOrCreateMerchant(ctx context.Context, name string, familyID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	query := `
		INSERT INTO merchants (name, family_id)
		VALUES ($1, $2)
		ON CONFLICT (name, family_id) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query, name, familyID).Scan(&id)
	return id, err
}

func (r *LedgerRepository) CreateTransaction(ctx context.Context, entry *models.Entry, txDetail *models.Transaction) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert Transaction Metadata
	if txDetail.ID == uuid.Nil {
		txDetail.ID = uuid.New()
	}
	queryTx := `
		INSERT INTO transactions (id, category_id, merchant_id, kind)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(ctx, queryTx, txDetail.ID, txDetail.CategoryID, txDetail.MerchantID, txDetail.Kind)
	if err != nil {
		return err
	}

	// 2. Insert Entry
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	entry.EntryableType = "Transaction"
	entry.EntryableID = txDetail.ID

	queryEntry := `
		INSERT INTO entries (id, account_id, amount, date, currency, name, entryable_type, entryable_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.Exec(ctx, queryEntry,
		entry.ID, entry.AccountID, entry.Amount, entry.Date, entry.Currency, entry.Name, entry.EntryableType, entry.EntryableID,
	)
	if err != nil {
		return err
	}

	// 3. Update Account Balance (Optional if we use triggers, but let's do it manually for now as per doc tip)
	queryUpdateAcc := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err = tx.Exec(ctx, queryUpdateAcc, entry.Amount, entry.AccountID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *LedgerRepository) CreateTransfer(ctx context.Context, fromEntry, toEntry *models.Entry) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Create one Transaction row for the transfer
	txID := uuid.New()
	queryTx := `INSERT INTO transactions (id, kind) VALUES ($1, 'transfer')`
	_, err = tx.Exec(ctx, queryTx, txID)
	if err != nil {
		return err
	}

	entries := []*models.Entry{fromEntry, toEntry}
	for _, e := range entries {
		if e.ID == uuid.Nil {
			e.ID = uuid.New()
		}
		e.EntryableType = "Transaction"
		e.EntryableID = txID

		queryEntry := `
			INSERT INTO entries (id, account_id, amount, date, currency, name, entryable_type, entryable_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err = tx.Exec(ctx, queryEntry,
			e.ID, e.AccountID, e.Amount, e.Date, e.Currency, e.Name, e.EntryableType, e.EntryableID,
		)
		if err != nil {
			return err
		}

		// Update Account Balance
		queryUpdateAcc := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
		_, err = tx.Exec(ctx, queryUpdateAcc, e.Amount, e.AccountID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
