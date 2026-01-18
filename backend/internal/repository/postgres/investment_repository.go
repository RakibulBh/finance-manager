package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)

type InvestmentRepository struct {
	db *pgxpool.Pool
}

func NewInvestmentRepository(db *pgxpool.Pool) *InvestmentRepository {
	return &InvestmentRepository{db: db}
}

func (r *InvestmentRepository) GetOrCreateSecurity(ctx context.Context, ticker, name string) (uuid.UUID, error) {
	var id uuid.UUID
	query := `
		INSERT INTO securities (ticker, name)
		VALUES ($1, $2)
		ON CONFLICT (ticker) DO UPDATE SET ticker = EXCLUDED.ticker
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query, ticker, name).Scan(&id)
	return id, err
}

func (r *InvestmentRepository) CreateTrade(ctx context.Context, entry *models.Entry, trade *models.Trade) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert Trade
	if trade.ID == uuid.Nil {
		trade.ID = uuid.New()
	}
	queryTrade := `
		INSERT INTO trades (id, security_id, qty, price, kind)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, queryTrade, trade.ID, trade.SecurityID, trade.Qty, trade.Price, trade.Kind)
	if err != nil {
		return err
	}

	// 2. Insert Entry
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	entry.EntryableType = "Trade"
	entry.EntryableID = trade.ID

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

	// 3. Update Account Balance
	queryUpdateAcc := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err = tx.Exec(ctx, queryUpdateAcc, entry.Amount, entry.AccountID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *InvestmentRepository) UpdateSecurityPrice(ctx context.Context, ticker string, price float64) error {
	query := `UPDATE securities SET latest_price = $1, last_updated = $2 WHERE ticker = $3`
	_, err := r.db.Exec(ctx, query, price, time.Now(), ticker)
	return err
}

func (r *InvestmentRepository) GetActiveTickers(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT ticker FROM securities`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickers []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		tickers = append(tickers, t)
	}
	return tickers, nil
}
