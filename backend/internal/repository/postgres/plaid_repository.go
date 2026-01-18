package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)

type PlaidRepository struct {
	db *pgxpool.Pool
}

func NewPlaidRepository(db *pgxpool.Pool) *PlaidRepository {
	return &PlaidRepository{db: db}
}

func (r *PlaidRepository) SaveItem(ctx context.Context, item *models.PlaidItem) error {
	query := `
		INSERT INTO plaid_items (family_id, access_token, item_id, institution_id, institution_name, sync_cursor)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (item_id) DO UPDATE SET
			access_token = EXCLUDED.access_token,
			sync_cursor = EXCLUDED.sync_cursor,
			updated_at = NOW()
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		item.FamilyID, item.AccessToken, item.ItemID, item.InstitutionID, item.InstitutionName, item.SyncCursor,
	).Scan(&item.ID)
}

func (r *PlaidRepository) GetItemsByFamily(ctx context.Context, familyID uuid.UUID) ([]models.PlaidItem, error) {
	query := `
		SELECT id, family_id, access_token, item_id, institution_id, institution_name, sync_cursor, status, created_at, updated_at
		FROM plaid_items
		WHERE family_id = $1
	`
	rows, err := r.db.Query(ctx, query, familyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.PlaidItem
	for rows.Next() {
		var item models.PlaidItem
		err := rows.Scan(
			&item.ID, &item.FamilyID, &item.AccessToken, &item.ItemID,
			&item.InstitutionID, &item.InstitutionName, &item.SyncCursor,
			&item.Status, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
