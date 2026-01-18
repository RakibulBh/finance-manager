package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rakibulbh/ai-finance-manager/internal/models"
)

type TransactionStore interface {
	GetOrCreateMerchant(ctx context.Context, name string) (uuid.UUID, error)
	CreateTransaction(ctx context.Context, entry *models.Entry, txDetail *models.Transaction) error
	CreateTransfer(ctx context.Context, fromEntry, toEntry *models.Entry) error
}

type TransactionHandler struct {
	repo TransactionStore
}

func NewTransactionHandler(repo TransactionStore) *TransactionHandler {
	return &TransactionHandler{repo: repo}
}

type CreateTransactionRequest struct {
	AccountID    uuid.UUID `json:"account_id"`
	Amount       float64   `json:"amount"`
	Date         time.Time `json:"date"`
	Name         string    `json:"name"`
	CategoryID   uuid.UUID `json:"category_id"`
	MerchantName string    `json:"merchant_name"`
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 1. Get/Create Merchant
	var merchantID uuid.UUID
	if req.MerchantName != "" {
		id, err := h.repo.GetOrCreateMerchant(r.Context(), req.MerchantName)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Failed to handle merchant")
			return
		}
		merchantID = id
	}

	// 2. Prepare Entry & Transaction
	entry := &models.Entry{
		AccountID: req.AccountID,
		Amount:    req.Amount,
		Date:      req.Date,
		Name:      req.Name,
		Currency:  "USD", // Default or from account? Let's assume USD for now
	}
	if entry.Date.IsZero() {
		entry.Date = time.Now()
	}

	txDetail := &models.Transaction{
		CategoryID: &req.CategoryID,
		MerchantID: &merchantID,
		Kind:       "standard",
	}

	// 3. Create in DB
	if err := h.repo.CreateTransaction(r.Context(), entry, txDetail); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	sendJSON(w, http.StatusCreated, entry)
}

type CreateTransferRequest struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        float64   `json:"amount"`
	Date          time.Time `json:"date"`
	Name          string    `json:"name"`
}

func (h *TransactionHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	var req CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	fromEntry := &models.Entry{
		AccountID: req.FromAccountID,
		Amount:    -req.Amount,
		Date:      req.Date,
		Name:      req.Name,
		Currency:  "USD",
	}

	toEntry := &models.Entry{
		AccountID: req.ToAccountID,
		Amount:    req.Amount,
		Date:      req.Date,
		Name:      req.Name,
		Currency:  "USD",
	}

	if err := h.repo.CreateTransfer(r.Context(), fromEntry, toEntry); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create transfer")
		return
	}

	sendJSON(w, http.StatusCreated, map[string]string{"message": "Transfer successful"})
}
