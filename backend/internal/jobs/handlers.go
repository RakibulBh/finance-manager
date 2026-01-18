package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/plaid/plaid-go/v20/plaid"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)

type WorkerServices struct {
	Plaid    PlaidProvider
	DB       ItemStorage
	Ledger   LedgerStorage
	Accounts AccountStorage
}

type PlaidProvider interface {
	DecryptToken(encryptedToken string) (string, error)
	SyncTransactions(ctx context.Context, accessToken string, cursor string) (plaid.TransactionsSyncResponse, error)
}

type ItemStorage interface {
	GetItemsByFamily(ctx context.Context, familyID uuid.UUID) ([]models.PlaidItem, error)
	GetItemByID(ctx context.Context, itemID string) (*models.PlaidItem, error)
	UpdateCursor(ctx context.Context, itemID string, cursor string) error
}

type LedgerStorage interface {
    GetOrCreateMerchant(ctx context.Context, name string, familyID uuid.UUID) (uuid.UUID, error)
	CreateTransaction(ctx context.Context, entry *models.Entry, txDetail *models.Transaction) error
}

type AccountStorage interface {
    GetByPlaidID(ctx context.Context, familyID uuid.UUID, plaidAccountID string) (*models.Account, error)
}

func HandleSyncAccountTask(ctx context.Context, t *asynq.Task, svc *WorkerServices) error {
	var p SyncAccountPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// 1. Fetch Item
	item, err := svc.DB.GetItemByID(ctx, p.ItemID)
	if err != nil {
		return fmt.Errorf("failed to fetch plaid item: %w", err)
	}

	// 2. Decrypt Token
	accessToken, err := svc.Plaid.DecryptToken(item.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to decrypt access token: %w", err)
	}

	// 3. Sync from Plaid
	resp, err := svc.Plaid.SyncTransactions(ctx, accessToken, item.SyncCursor)
	if err != nil {
		return fmt.Errorf("plaid sync failed: %w", err)
	}


	// 4. Process Added Transactions
	for _, plTx := range resp.Added {
		// Find Account
		acc, err := svc.Accounts.GetByPlaidID(ctx, item.FamilyID, plTx.AccountId)
		if err != nil {
			log.Printf("Account not found for Plaid ID %s, skipping transaction %s", plTx.AccountId, plTx.TransactionId)
			continue
		}

		// Prepare Entry
		date, _ := time.Parse("2006-01-02", plTx.Date)

		currency := "USD"
		if plTx.IsoCurrencyCode.Get() != nil {
			currency = *plTx.IsoCurrencyCode.Get()
		}

		entry := &models.Entry{
			AccountID: acc.ID,
			Amount:    plTx.Amount,
			Date:      date,
			Name:      plTx.Name,
			Currency:  currency,
		}


		// Prepare Transaction details
		var merchantID *uuid.UUID
		if plTx.MerchantName.Get() != nil {
			mID, err := svc.Ledger.GetOrCreateMerchant(ctx, *plTx.MerchantName.Get(), item.FamilyID)
			if err == nil {
				merchantID = &mID
			}
		}

		txDetail := &models.Transaction{
			MerchantID: merchantID,
			Kind:       "standard",
		}

		// Save
		if err := svc.Ledger.CreateTransaction(ctx, entry, txDetail); err != nil {
			return fmt.Errorf("failed to save synced transaction %s: %w", plTx.TransactionId, err)
		}
	}

	// 5. Update Cursor
	if err := svc.DB.UpdateCursor(ctx, item.ItemID, resp.NextCursor); err != nil {
		return fmt.Errorf("failed to update cursor: %w", err)
	}

	return nil
}
