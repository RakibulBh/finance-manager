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

func (r *PlaidRepository) GetItemByID(ctx context.Context, itemID string) (*models.PlaidItem, error) {
	query := `
		SELECT id, family_id, access_token, item_id, institution_id, institution_name, sync_cursor, status, created_at, updated_at
		FROM plaid_items
		WHERE item_id = $1
	`
	var item models.PlaidItem
	err := r.db.QueryRow(ctx, query, itemID).Scan(
		&item.ID, &item.FamilyID, &item.AccessToken, &item.ItemID,
		&item.InstitutionID, &item.InstitutionName, &item.SyncCursor,
		&item.Status, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PlaidRepository) UpdateCursor(ctx context.Context, itemID string, cursor string) error {
	query := `UPDATE plaid_items SET sync_cursor = $1, updated_at = NOW() WHERE item_id = $2`
	_, err := r.db.Exec(ctx, query, cursor, itemID)
	return err
}
